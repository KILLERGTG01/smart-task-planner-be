package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Task represents a planner task returned by LLM
type Task struct {
	Task         string   `json:"task"`
	DurationDays int      `json:"duration_days"`
	DependsOn    []string `json:"depends_on"`
}

func GeneratePlan(ctx context.Context, goal string) ([]Task, error) {
	base := strings.TrimSuffix(os.Getenv("GEMINI_BASE_URL"), "/")
	key := os.Getenv("GEMINI_API_KEY")
	if base == "" || key == "" {
		return nil, errors.New("gemini config not set")
	}

	// Build a prompt request for Google's Generative API (simple)
	// NOTE: Adapt this payload to the exact API shape of KILLERGTG01r Gemini endpoint.
	payload := map[string]interface{}{
		"prompt": map[string]interface{}{
			"text": "You are an expert task planner. Return ONLY valid JSON array of tasks. Each task: {\\\"task\\\", \\\"duration_days\\\", \\\"depends_on\\\"}. Goal: " + goal,
		},
		"max_output_tokens": 800,
	}
	body, _ := json.Marshal(payload)
	// Example endpoint: ${GEMINI_BASE_URL}/v1/models/gemini-1.5-pro:generateText?key=API_KEY
	url := base + "/v1/models/gemini-2.5-pro:generateText?key=" + key

	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	text := string(b)

	// Extract JSON array from response text
	jsonStart := strings.Index(text, "[")
	jsonEnd := strings.LastIndex(text, "]")
	if jsonStart == -1 || jsonEnd == -1 || jsonEnd <= jsonStart {
		return nil, errors.New("failed to extract json array from model response")
	}
	arrText := text[jsonStart : jsonEnd+1]
	var tasks []Task
	if err := json.Unmarshal([]byte(arrText), &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}
