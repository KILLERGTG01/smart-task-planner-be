package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/KILLERGTG01/smart-task-planner-be/internal/db"
	"github.com/KILLERGTG01/smart-task-planner-be/internal/services"
)

type generateReq struct {
	Goal  string `json:"goal"`
	Title string `json:"title"`
}

func GenerateHandler(c *fiber.Ctx) error {
	var req generateReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_body"})
	}
	if req.Goal == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "goal_required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	tasks, err := services.GeneratePlan(ctx, req.Goal)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "generation_failed", "detail": err.Error()})
	}

	response := fiber.Map{"plan": tasks}

	authSub := c.Locals("auth_sub")
	if authSub != nil {
		userID := authSub.(string)
		_, err = findOrCreateUser(userID)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "user_upsert_failed", "detail": err.Error()})
		}

		id := uuid.NewString()
		planJson, _ := json.Marshal(tasks)
		_, err = db.Pool.Exec(context.Background(),
			"INSERT INTO plans (id, user_id, title, goal, plan_json, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,now(),now())",
			id, userID, req.Title, req.Goal, planJson,
		)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "save_failed", "detail": err.Error()})
		}

		response["id"] = id
		response["saved"] = true
	} else {
		response["saved"] = false
		response["message"] = "Plan generated but not saved. Login to save your plans."
	}

	return c.Status(http.StatusOK).JSON(response)
}

func findOrCreateUser(sub string) (string, error) {
	var id string
	err := db.Pool.QueryRow(context.Background(), "SELECT id FROM users WHERE auth0_id=$1", sub).Scan(&id)
	if err == nil {
		return id, nil
	}
	id = sub
	_, err = db.Pool.Exec(context.Background(), "INSERT INTO users (id, auth0_id, created_at) VALUES ($1,$2,now())", id, sub)
	if err != nil {
		return "", err
	}
	return id, nil
}

func HistoryHandler(c *fiber.Ctx) error {
	authSub := c.Locals("auth_sub")
	if authSub == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthenticated"})
	}
	sub := authSub.(string)
	rows, err := db.Pool.Query(context.Background(), "SELECT id, title, goal, plan_json, created_at FROM plans WHERE user_id=$1 ORDER BY created_at DESC LIMIT 100", sub)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db_query_failed", "detail": err.Error()})
	}
	defer rows.Close()
	var res []map[string]interface{}
	for rows.Next() {
		var id, title, goal string
		var planJson []byte
		var createdAt time.Time
		if err := rows.Scan(&id, &title, &goal, &planJson, &createdAt); err != nil {
			continue
		}
		var plan interface{}
		_ = json.Unmarshal(planJson, &plan)
		res = append(res, map[string]interface{}{
			"id":        id,
			"title":     title,
			"goal":      goal,
			"plan":      plan,
			"createdAt": createdAt,
		})
	}
	return c.JSON(fiber.Map{"plans": res})
}

func GenerateStreamHandler(c *fiber.Ctx) error {
	var req generateReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_body"})
	}
	if req.Goal == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "goal_required"})
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "Cache-Control")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		writeSSE := func(event, data string) {
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, data)
			w.Flush()
		}

		writeSSE("status", `{"message": "Starting plan generation..."}`)

		tasks, err := services.GeneratePlan(ctx, req.Goal)
		if err != nil {
			writeSSE("error", fmt.Sprintf(`{"error": "generation_failed", "detail": "%s"}`, err.Error()))
			return
		}

		writeSSE("progress", `{"message": "Plan generated successfully!"}`)

		planData, _ := json.Marshal(fiber.Map{"plan": tasks})
		writeSSE("plan", string(planData))

		authSub := c.Locals("auth_sub")
		if authSub != nil {
			userID := authSub.(string)
			_, err = findOrCreateUser(userID)
			if err != nil {
				writeSSE("warning", fmt.Sprintf(`{"message": "Plan generated but not saved: %s"}`, err.Error()))
				writeSSE("complete", `{"saved": false}`)
				return
			}

			id := uuid.NewString()
			planJson, _ := json.Marshal(tasks)
			_, err = db.Pool.Exec(context.Background(),
				"INSERT INTO plans (id, user_id, title, goal, plan_json, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,now(),now())",
				id, userID, req.Title, req.Goal, planJson,
			)
			if err != nil {
				writeSSE("warning", fmt.Sprintf(`{"message": "Plan generated but not saved: %s"}`, err.Error()))
				writeSSE("complete", `{"saved": false}`)
				return
			}

			writeSSE("saved", fmt.Sprintf(`{"id": "%s", "message": "Plan saved successfully!"}`, id))
			writeSSE("complete", `{"saved": true}`)
		} else {
			writeSSE("info", `{"message": "Plan generated but not saved. Login to save your plans."}`)
			writeSSE("complete", `{"saved": false}`)
		}
	})

	return nil
}
