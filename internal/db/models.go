package db

import "time"

type User struct {
	ID        string    `json:"id" validate:"required,uuid4"`
	Auth0ID   string    `json:"auth0_id" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Name      string    `json:"name" validate:"required,min=1,max=100"`
	CreatedAt time.Time `json:"created_at"`
}

type Plan struct {
	ID        string      `json:"id" validate:"required,uuid4"`
	UserID    string      `json:"user_id" validate:"required,uuid4"`
	Title     string      `json:"title" validate:"required,min=1,max=200"`
	Goal      string      `json:"goal" validate:"required,min=1,max=1000"`
	PlanJSON  interface{} `json:"plan_json" validate:"required"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
