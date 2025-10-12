package handlers

import (
	"context"
	"encoding/json"
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

	// run generation
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	tasks, err := services.GeneratePlan(ctx, req.Goal)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "generation_failed", "detail": err.Error()})
	}

	// upsert user based on auth_sub
	authSub := c.Locals("auth_sub")
	userID := authSub.(string)
	// ensure user exists
	_, err = findOrCreateUser(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "user_upsert_failed", "detail": err.Error()})
	}

	// save plan
	id := uuid.NewString()
	planJson, _ := json.Marshal(tasks)
	_, err = db.Pool.Exec(context.Background(),
		"INSERT INTO plans (id, user_id, title, goal, plan_json, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,now(),now())",
		id, userID, req.Title, req.Goal, planJson,
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "save_failed", "detail": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"id": id, "plan": tasks})
}

func findOrCreateUser(sub string) (string, error) {
	var id string
	err := db.Pool.QueryRow(context.Background(), "SELECT id FROM users WHERE auth0_id=$1", sub).Scan(&id)
	if err == nil {
		return id, nil
	}
	// create
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
