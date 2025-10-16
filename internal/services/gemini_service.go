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

	payload := map[string]interface{}{
		"prompt": map[string]interface{}{
			"text": "You are an expert task planner. Return ONLY valid JSON array of tasks. Each task: {\\\"task\\\", \\\"duration_days\\\", \\\"depends_on\\\"}. Goal: " + goal,
		},
		"max_output_tokens": 800,
	}
	body, _ := json.Marshal(payload)
	url := base + "/v1/models/gemini-2.5-flash-lite:generateText?key=" + key

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
