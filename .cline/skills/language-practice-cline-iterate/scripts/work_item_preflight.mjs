#!/usr/bin/env node

import { readFile, readdir } from "node:fs/promises";
import { basename, join, relative, resolve } from "node:path";
import process from "node:process";
import { execFile } from "node:child_process";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);
const WORK_ITEM_KEY = /^(WI|GH)-\d{4}$/u;
const BRANCH_NAME = /^(wi|gh)-\d{4}-[a-z0-9]+(?:-[a-z0-9]+)*$/u;
const STEP_FILE = /^(WI|GH)-(\d{4})-(\d+)-.+\.md$/u;
const COMMIT_SUBJECT = /^(WI|GH)-(\d{4}): (\d{2,}): (feat|bug|doc|test|refactor|perf|build|ci|chore)\(([^)]+)\): (.+)$/u;

function addError(errors, code, message) {
  errors.push({ code, message });
}

function parseValue(value) {
  const trimmed = value.trim();
  if (trimmed === "[]") return [];
  if (trimmed.startsWith("[")) {
    try {
      const parsed = JSON.parse(trimmed);
      return Array.isArray(parsed) ? parsed.map(String) : parsed;
    } catch {
      return trimmed;
    }
  }
  if ((trimmed.startsWith("\"") && trimmed.endsWith("\""))
    || (trimmed.startsWith("'") && trimmed.endsWith("'"))) {
    return trimmed.slice(1, -1);
  }
  return trimmed;
}

export function parseFrontmatter(text) {
  const lines = text.split(/\r?\n/u);
  if (lines[0] !== "---") return {};
  const metadata = {};
  for (const line of lines.slice(1)) {
    if (line === "---") break;
    const match = line.match(/^([A-Za-z][A-Za-z0-9_-]*):\s*(.*)$/u);
    if (match) metadata[match[1]] = parseValue(match[2]);
  }
  return metadata;
}

export function declaredCommitSubject(text) {
  const sectionStart = text.search(/^## Commit subject\s*$/mu);
  if (sectionStart < 0) return null;
  const section = text.slice(sectionStart);
  const fence = section.match(/```text\s*\n([^\n]+)\n```/u);
  return fence?.[1]?.trim() ?? null;
}

function stepRecord(name, text) {
  const fileMatch = basename(name).match(STEP_FILE);
  const metadata = parseFrontmatter(text);
  const step = String(metadata.step ?? fileMatch?.[3] ?? "");
  const dependencies = Array.isArray(metadata.depends_on)
    ? metadata.depends_on.map(String)
    : [];
  return {
    name: basename(name),
    text,
    metadata,
    key: fileMatch ? `${fileMatch[1]}-${fileMatch[2]}` : String(metadata.id ?? ""),
    step,
    stepNumber: Number.parseInt(step, 10),
    target: typeof metadata.target === "string" ? metadata.target : "",
    dependencies,
    status: String(metadata.status ?? "planned"),
    subject: declaredCommitSubject(text),
  };
}

function expectedRecordPath(root, workItemId, directoryName) {
  if (workItemId.startsWith("WI-")) return join(root, ".docs", "work-items", workItemId);
  return join(root, ".docs", "work-items", directoryName ?? workItemId);
}

function validateSubject(subject, key, step) {
  if (!subject) return "missing ## Commit subject declaration";
  const match = subject.match(COMMIT_SUBJECT);
  if (!match) return `malformed commit subject: ${subject}`;
  const expectedStep = String(step).padStart(2, "0");
  if (`${match[1]}-${match[2]}` !== key) return `subject work-item key does not match ${key}`;
  if (match[3] !== expectedStep) return `subject step ${match[3]} does not match ${expectedStep}`;
  if (!match[6].trim() || /^(update|changes|fix)$/iu.test(match[6].trim())) {
    return "subject description is empty or generic";
  }
  return null;
}

function hasCycle(records) {
  const states = new Map();
  const byStep = new Map(records.map((record) => [record.step, record]));
  function visit(step) {
    const state = states.get(step);
    if (state === "active") return true;
    if (state === "done") return false;
    states.set(step, "active");
    for (const dependency of byStep.get(step)?.dependencies ?? []) {
      if (byStep.has(dependency) && visit(dependency)) return true;
    }
    states.set(step, "done");
    return false;
  }
  return records.some((record) => visit(record.step));
}

function dependencyClosure(records, selectedStep) {
  if (selectedStep === undefined) return new Set();
  const byStep = new Map(records.map((record) => [record.step, record]));
  const required = new Set();
  const visit = (step) => {
    if (required.has(step)) return;
    required.add(step);
    for (const dependency of byStep.get(step)?.dependencies ?? []) visit(dependency);
  };
  visit(String(selectedStep).padStart(2, "0"));
  return required;
}

export function validateWorkItem(input) {
  const errors = [];
  const workItemId = String(input.workItemId ?? "");
  const readme = parseFrontmatter(String(input.readmeText ?? ""));
  const records = (input.stepFiles ?? []).map((file) => stepRecord(file.name, file.text));
  const expectedRoot = input.repositoryRoot ? resolve(input.repositoryRoot) : null;
  const requiredDependencies = dependencyClosure(records, input.selectedStep);

  if (input.selectedStep === undefined || input.selectedStep === null || String(input.selectedStep).trim() === "") {
    addError(errors, "SELECTED_STEP_REQUIRED", "preflight requires one explicitly selected step");
  }

  if (!WORK_ITEM_KEY.test(workItemId)) addError(errors, "INVALID_WORK_ITEM", `invalid work-item key: ${workItemId}`);
  if (readme.id !== workItemId) addError(errors, "README_ID_MISMATCH", "README identity does not match the selected work item");
  if (expectedRoot && readme.repository && resolve(String(readme.repository)) !== expectedRoot) {
    addError(errors, "README_ROOT_MISMATCH", "README repository does not match the resolved Git root");
  }
  if (!BRANCH_NAME.test(String(input.branchName ?? ""))) addError(errors, "INVALID_BRANCH", "branch must be wi-NNNN-<slug> or gh-NNNN-<slug>");
  const branchKey = String(input.branchName ?? "").slice(0, 7).toUpperCase();
  if (branchKey !== workItemId) addError(errors, "BRANCH_KEY_MISMATCH", "topic branch key does not match the work item");
  if (readme.branch && readme.branch !== input.branchName) addError(errors, "README_BRANCH_MISMATCH", "README branch does not match the current branch");

  const rootCommitMatch = String(input.readmeText ?? "").match(/^[-*]\s*\*\*Root record commit:\*\*\s*`([^`]+)`/mu);
  if (!rootCommitMatch) addError(errors, "ROOT_COMMIT_MISSING", "README does not declare a root record commit");
  if (input.rootDocumentationCommitted !== true) addError(errors, "ROOT_DOC_UNCOMMITTED", "root work-item documentation is not committed");
  if (input.implementationStarted === true && input.rootDocumentationCommitted !== true) {
    addError(errors, "PREMATURE_IMPLEMENTATION", "implementation started before root documentation was committed");
  }
  if (input.actualRepositoryRoot && expectedRoot && resolve(input.actualRepositoryRoot) !== expectedRoot) {
    addError(errors, "REPOSITORY_ROOT_MISMATCH", "resolved repository root differs from the requested root");
  }

  if (records.length === 0) addError(errors, "STEPS_MISSING", "work-item steps directory has no step records");
  const seenSteps = new Set();
  for (const record of records) {
    if (!record.metadata.id || record.key !== workItemId) addError(errors, "STEP_ID_MISMATCH", `${record.name} does not match ${workItemId}`);
    if (!Number.isInteger(record.stepNumber) || record.stepNumber < 1) addError(errors, "STEP_INVALID", `${record.name} has an invalid step identifier`);
    if (seenSteps.has(record.step)) addError(errors, "STEP_DUPLICATE", `duplicate step ${record.step}`);
    seenSteps.add(record.step);
    if (!record.target) addError(errors, "STEP_TARGET_MISSING", `${record.name} does not declare a target`);
    const subjectError = validateSubject(record.subject, workItemId, record.stepNumber);
    if (subjectError) addError(errors, "COMMIT_SUBJECT_INVALID", `${record.name}: ${subjectError}`);
    for (const dependency of record.dependencies) {
      const dependencyRecord = records.find((candidate) => candidate.step === dependency);
      if (!dependencyRecord) addError(errors, "DEPENDENCY_MISSING", `${record.name} depends on missing step ${dependency}`);
      else if (dependencyRecord.stepNumber >= record.stepNumber) addError(errors, "DEPENDENCY_ORDER", `${record.name} depends on later step ${dependency}`);
      else if (requiredDependencies.has(record.step) && !["complete", "proven"].includes(dependencyRecord.status)) addError(errors, "DEPENDENCY_UNPROVEN", `${record.name} depends on unproven step ${dependency}`);
    }
  }
  if (hasCycle(records)) addError(errors, "DEPENDENCY_CYCLE", "step dependency graph contains a cycle");

  if (input.selectedStep !== undefined) {
    const selected = records.find((record) => record.step === String(input.selectedStep).padStart(2, "0") || record.step === String(input.selectedStep));
    if (!selected) addError(errors, "SELECTED_STEP_MISSING", `selected step ${input.selectedStep} is not owned by this work item`);
    else if (input.selectedTarget !== undefined && selected.target !== input.selectedTarget) addError(errors, "SELECTED_TARGET_MISMATCH", "selected target is outside the selected step's ownership");
  }
  if (input.actualCommitSubjects) {
    for (const [step, actual] of Object.entries(input.actualCommitSubjects)) {
      const record = records.find((candidate) => candidate.step === step);
      if (record && actual !== record.subject) addError(errors, "COMMIT_SUBJECT_MISMATCH", `${record.name} actual commit subject does not match its declaration`);
    }
  }
  for (const path of input.worktreePaths ?? []) addError(errors, "UNRELATED_WORKTREE_CHANGE", `unrelated worktree change: ${path}`);

  return {
    ok: errors.length === 0,
    errors,
    workItemId,
    recordDirectory: expectedRoot ? expectedRecordPath(expectedRoot, workItemId, input.recordDirectoryName) : null,
    steps: records.map(({ name, step, target, dependencies, status, subject }) => ({ name, step, target, dependencies, status, subject })),
  };
}

async function git(root, args) {
  const result = await execFileAsync("git", ["-C", root, ...args], { maxBuffer: 1024 * 1024 });
  return result.stdout.trim();
}

async function findRecordDirectory(root, workItemId) {
  const workItems = join(root, ".docs", "work-items");
  const entries = await readdir(workItems, { withFileTypes: true });
  const candidates = entries
    .filter((entry) => entry.isDirectory() && (entry.name === workItemId || entry.name.startsWith(`${workItemId}-`)))
    .map((entry) => entry.name)
    .sort();
  if (candidates.length !== 1) throw new Error(`expected one authoritative record directory for ${workItemId}, found ${candidates.length}`);
  return join(workItems, candidates[0]);
}

async function runCli() {
  const args = process.argv.slice(2);
  const value = (name) => {
    const index = args.indexOf(name);
    return index < 0 ? undefined : args[index + 1];
  };
  const root = resolve(value("--root") ?? (await git(process.cwd(), ["rev-parse", "--show-toplevel"])));
  const workItemId = value("--work-item");
  if (!workItemId) throw new Error("--work-item is required");
  const recordDirectory = await findRecordDirectory(root, workItemId);
  const readmeText = await readFile(join(recordDirectory, "README.md"), "utf8");
  const stepDirectory = join(recordDirectory, "steps");
  const stepNames = (await readdir(stepDirectory)).filter((name) => name.endsWith(".md")).sort();
  const stepFiles = await Promise.all(stepNames.map(async (name) => ({ name, text: await readFile(join(stepDirectory, name), "utf8") })));
  const branchName = await git(root, ["branch", "--show-current"]);
  const rootCommitSubject = readmeText.match(/^[-*]\s*\*\*Root record commit:\*\*\s*`([^`]+)`/mu)?.[1];
  let rootDocumentationCommitted = false;
  if (rootCommitSubject) {
    try {
      const recordRelativePath = relative(root, recordDirectory);
      const history = await git(root, ["log", "--format=%s", "--all", "--", recordRelativePath]);
      rootDocumentationCommitted = history.split(/\r?\n/u).includes(rootCommitSubject);
    } catch {
      rootDocumentationCommitted = false;
    }
  }
  const status = await git(root, ["status", "--porcelain=v1", "--untracked-files=all"]);
  const result = validateWorkItem({
    workItemId,
    readmeText,
    stepFiles,
    branchName,
    repositoryRoot: root,
    actualRepositoryRoot: root,
    rootDocumentationCommitted,
    selectedStep: value("--step"),
    selectedTarget: value("--target"),
    worktreePaths: status ? status.split(/\r?\n/u).filter(Boolean).map((line) => line.slice(3)) : [],
    recordDirectoryName: relative(join(root, ".docs", "work-items"), recordDirectory),
  });
  process.stdout.write(`${JSON.stringify(result, null, 2)}\n`);
  if (!result.ok) process.exitCode = 1;
}

if (import.meta.url === `file://${process.argv[1]}`) {
  runCli().catch((error) => {
    console.error(error instanceof Error ? error.message : String(error));
    process.exitCode = 1;
  });
}
