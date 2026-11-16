// Package main is the manual starting point for the Google ADK learning track.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"google.golang.org/adk/v2/model/openaimodel"
)

const defaultOpenRouterBaseURL = "https://openrouter.ai/api/v1"

func main() {
	if err := run(context.Background(), os.Getenv, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, getenv func(string) string, stdout io.Writer) error {
	apiKey := strings.TrimSpace(getenv("OPENROUTER_API_KEY"))
	if apiKey == "" {
		return fmt.Errorf("OPENROUTER_API_KEY is required")
	}

	modelName := strings.TrimSpace(getenv("OPENROUTER_MODEL"))
	if modelName == "" {
		return fmt.Errorf("OPENROUTER_MODEL is required")
	}

	baseURL := strings.TrimSpace(getenv("OPENROUTER_BASE_URL"))
	if baseURL == "" {
		baseURL = defaultOpenRouterBaseURL
	}

	llm, err := openaimodel.NewModel(ctx, modelName, &openaimodel.ClientConfig{
		APIKey:  apiKey,
		BaseURL: baseURL,
		API:     openaimodel.APIChatCompletions,
	})
	if err != nil {
		return fmt.Errorf("configure OpenRouter model: %w", err)
	}

	_, err = fmt.Fprintf(stdout, "ADK model configured: %s\nEndpoint: %s\nNo model request was sent.\n", llm.Name(), baseURL)
	return err
}
