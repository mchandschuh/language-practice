import test from "node:test";
import assert from "node:assert/strict";
import { validateWorkItem } from "../scripts/work_item_preflight.mjs";

const root = "/tmp/language-practice";
const subject = "WI-0001: 01: feat(repo): implement selected step";

function readme({ id = "WI-0001", branch = "wi-0001-demo", rootCommit = "abc123" } = {}) {
  return `---\nid: ${id}\nrepository: ${root}\nbranch: ${branch}\n---\n\n- **Root record commit:** \`${rootCommit}\`\n`;
}

function step({ id = "WI-0001", step = "01", depends = "[]", status = "planned", target = ".cline/skills/example/SKILL.md", declaredSubject = subject } = {}) {
  return {
    name: `${id}-${step}-example.md`,
    text: `---\nid: ${id}\nstep: "${step}"\nstatus: ${status}\ndepends_on: ${depends}\ntarget: ${target}\n---\n\n## Commit subject\n\n\`\`\`text\n${declaredSubject}\n\`\`\`\n`,
  };
}

function valid(overrides = {}) {
  return {
    workItemId: "WI-0001",
    readmeText: readme(),
    stepFiles: [step()],
    branchName: "wi-0001-demo",
    repositoryRoot: root,
    actualRepositoryRoot: root,
    rootDocumentationCommitted: true,
    selectedStep: "01",
    selectedTarget: ".cline/skills/example/SKILL.md",
    worktreePaths: [],
    ...overrides,
  };
}

test("accepts a valid local work item and canonical branch", () => {
  assert.equal(validateWorkItem(valid()).ok, true);
});

test("rejects a work-item preflight without an explicitly selected step", () => {
  const input = valid();
  delete input.selectedStep;
  assert.ok(validateWorkItem(input).errors.some((error) => error.code === "SELECTED_STEP_REQUIRED"));
});

test("accepts a GitHub work item with its slugged record directory contract", () => {
  const result = validateWorkItem(valid({
    workItemId: "GH-0042",
    readmeText: readme({ id: "GH-0042", branch: "gh-0042-auth-timeout" }),
    branchName: "gh-0042-auth-timeout",
    stepFiles: [step({ id: "GH-0042", declaredSubject: "GH-0042: 01: bug(repo): fix auth timeout" })],
    selectedTarget: ".cline/skills/example/SKILL.md",
  }));
  assert.equal(result.ok, true);
});

test("rejects a missing work item README identity", () => {
  const result = validateWorkItem(valid({ readmeText: "" }));
  assert.ok(result.errors.some((error) => error.code === "README_ID_MISMATCH"));
});

test("rejects missing steps", () => {
  const result = validateWorkItem(valid({ stepFiles: [] }));
  assert.ok(result.errors.some((error) => error.code === "STEPS_MISSING"));
});

test("rejects invalid dependency references and ordering", () => {
  const result = validateWorkItem(valid({
    stepFiles: [step({ step: "01", status: "complete" }), step({ step: "02", depends: '["99"]', declaredSubject: "WI-0001: 02: feat(repo): implement second step" })],
    selectedStep: "02",
    selectedTarget: ".cline/skills/example/SKILL.md",
  }));
  assert.ok(result.errors.some((error) => error.code === "DEPENDENCY_MISSING"));
});

test("rejects a dependency on a later step", () => {
  const result = validateWorkItem(valid({
    stepFiles: [step({ step: "01", depends: '["02"]' }), step({ step: "02", declaredSubject: "WI-0001: 02: feat(repo): implement second step" })],
  }));
  assert.ok(result.errors.some((error) => error.code === "DEPENDENCY_ORDER"));
});

test("rejects a cyclic dependency graph", () => {
  const result = validateWorkItem(valid({
    stepFiles: [step({ step: "01", depends: '["02"]' }), step({ step: "02", depends: '["01"]', declaredSubject: "WI-0001: 02: feat(repo): implement second step" })],
  }));
  assert.ok(result.errors.some((error) => error.code === "DEPENDENCY_CYCLE"));
});

test("rejects an uncommitted root record and premature implementation", () => {
  const result = validateWorkItem(valid({ rootDocumentationCommitted: false, implementationStarted: true }));
  assert.ok(result.errors.some((error) => error.code === "ROOT_DOC_UNCOMMITTED"));
  assert.ok(result.errors.some((error) => error.code === "PREMATURE_IMPLEMENTATION"));
});

test("rejects a mismatched implementation subject", () => {
  const result = validateWorkItem(valid({ actualCommitSubjects: { "01": "WI-0001: 01: feat(repo): wrong subject" } }));
  assert.ok(result.errors.some((error) => error.code === "COMMIT_SUBJECT_MISMATCH"));
});

test("rejects invalid and mismatched branches", () => {
  const result = validateWorkItem(valid({ branchName: "docs/WI-0001-demo" }));
  assert.ok(result.errors.some((error) => error.code === "INVALID_BRANCH"));
});

test("rejects a selected target outside step ownership", () => {
  const result = validateWorkItem(valid({ selectedTarget: ".cline/skills/other/SKILL.md" }));
  assert.ok(result.errors.some((error) => error.code === "SELECTED_TARGET_MISMATCH"));
});

test("rejects unrelated worktree changes", () => {
  const result = validateWorkItem(valid({ worktreePaths: ["unrelated.txt"] }));
  assert.ok(result.errors.some((error) => error.code === "UNRELATED_WORKTREE_CHANGE"));
});