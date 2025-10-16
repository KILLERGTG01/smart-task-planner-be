package handlers

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/KILLERGTG01/smart-task-planner-be/internal/config"
	"github.com/KILLERGTG01/smart-task-planner-be/internal/db"
)

type Auth0TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type Auth0UserResponse struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type Auth0ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func LoginHandler(c *fiber.Ctx) error {
	provider := c.Query("provider", "google")
	if provider != "google" && provider != "github" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_provider"})
	}

	cfg := c.Locals("config").(*config.Config)
	state := generateState()

	authURL := fmt.Sprintf("https://%s/authorize?"+
		"response_type=code&"+
		"client_id=%s&"+
		"redirect_uri=%s&"+
		"scope=openid profile email&"+
		"state=%s&"+
		"connection=%s",
		cfg.Auth0Domain,
		cfg.Auth0ClientID,
		url.QueryEscape(cfg.Auth0RedirectURI),
		state,
		getAuth0Connection(provider))

	return c.JSON(fiber.Map{
		"auth_url": authURL,
		"state":    state,
	})
}

func CallbackHandler(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	if errorParam != "" {
		errorDesc := c.Query("error_description")
		zap.L().Error("Auth0 callback error", zap.String("error", errorParam), zap.String("description", errorDesc))
		return c.Redirect(fmt.Sprintf("%s/auth/error?error=%s", c.Locals("config").(*config.Config).FrontendURL, errorParam))
	}

	if code == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "authorization_code_required"})
	}

	if state == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "state_parameter_required"})
	}

	cfg := c.Locals("config").(*config.Config)

	tokenResp, err := exchangeCodeForToken(cfg, code)
	if err != nil {
		zap.L().Error("Failed to exchange code for token", zap.Error(err))
		return c.Redirect(fmt.Sprintf("%s/auth/error?error=token_exchange_failed", cfg.FrontendURL))
	}

	userInfo, err := getUserInfoFromAuth0(cfg, tokenResp.AccessToken)
	if err != nil {
		zap.L().Error("Failed to get user info from Auth0", zap.Error(err))
		return c.Redirect(fmt.Sprintf("%s/auth/error?error=user_info_failed", cfg.FrontendURL))
	}

	user, err := findOrCreateUserFromAuth0(userInfo.Sub, userInfo.Email, userInfo.Name)
	if err != nil {
		zap.L().Error("Failed to create/find user", zap.Error(err))
		return c.Redirect(fmt.Sprintf("%s/auth/error?error=user_creation_failed", cfg.FrontendURL))
	}

	successURL := fmt.Sprintf("%s/auth/success?token=%s&user=%s",
		cfg.FrontendURL,
		url.QueryEscape(tokenResp.AccessToken),
		url.QueryEscape(fmt.Sprintf(`{"id":"%s","email":"%s","name":"%s"}`, user.ID, user.Email, user.Name)))

	return c.Redirect(successURL)
}

func RefreshTokenHandler(c *fiber.Ctx) error {
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_request_body"})
	}

	if req.RefreshToken == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "refresh_token_required"})
	}

	cfg := c.Locals("config").(*config.Config)

	tokenResp, err := refreshAuth0Token(cfg, req.RefreshToken)
	if err != nil {
		zap.L().Error("Token refresh failed", zap.Error(err))
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid_refresh_token"})
	}

	return c.JSON(fiber.Map{
		"access_token": tokenResp.AccessToken,
		"token_type":   tokenResp.TokenType,
		"expires_in":   tokenResp.ExpiresIn,
	})
}

func LogoutHandler(c *fiber.Ctx) error {
	cfg := c.Locals("config").(*config.Config)

	logoutURL := fmt.Sprintf("https://%s/v2/logout?client_id=%s&returnTo=%s",
		cfg.Auth0Domain,
		cfg.Auth0ClientID,
		url.QueryEscape(cfg.FrontendURL))

	return c.JSON(fiber.Map{
		"logout_url": logoutURL,
	})
}

func UserProfileHandler(c *fiber.Ctx) error {
	authSub := c.Locals("auth_sub")
	if authSub == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthenticated"})
	}

	userID := authSub.(string)
	var user db.User
	err := db.Pool.QueryRow(context.Background(),
		"SELECT id, auth0_id, email, name, created_at FROM users WHERE auth0_id=$1",
		userID).Scan(&user.ID, &user.Auth0ID, &user.Email, &user.Name, &user.CreatedAt)

	if err != nil {
		zap.L().Error("Failed to get user profile", zap.Error(err))
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "user_not_found"})
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":         user.ID,
			"email":      user.Email,
			"name":       user.Name,
			"created_at": user.CreatedAt,
		},
	})
}

func exchangeCodeForToken(cfg *config.Config, code string) (*Auth0TokenResponse, error) {
	url := fmt.Sprintf("https://%s/oauth/token", cfg.Auth0Domain)

	payload := map[string]interface{}{
		"grant_type":    "authorization_code",
		"client_id":     cfg.Auth0ClientID,
		"client_secret": cfg.Auth0ClientSecret,
		"code":          code,
		"redirect_uri":  cfg.Auth0RedirectURI,
	}

	return makeAuth0TokenRequest(url, payload)
}

func refreshAuth0Token(cfg *config.Config, refreshToken string) (*Auth0TokenResponse, error) {
	url := fmt.Sprintf("https://%s/oauth/token", cfg.Auth0Domain)

	payload := map[string]interface{}{
		"grant_type":    "refresh_token",
		"client_id":     cfg.Auth0ClientID,
		"client_secret": cfg.Auth0ClientSecret,
		"refresh_token": refreshToken,
	}

	return makeAuth0TokenRequest(url, payload)
}

func makeAuth0TokenRequest(url string, payload map[string]interface{}) (*Auth0TokenResponse, error) {
	jsonData, _ := json.Marshal(payload)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errorResp Auth0ErrorResponse
		json.Unmarshal(body, &errorResp)
		return nil, fmt.Errorf("auth0 error: %s - %s", errorResp.Error, errorResp.ErrorDescription)
	}

	var tokenResp Auth0TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func getUserInfoFromAuth0(cfg *config.Config, accessToken string) (*Auth0UserResponse, error) {
	url := fmt.Sprintf("https://%s/userinfo", cfg.Auth0Domain)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info from Auth0")
	}

	body, _ := io.ReadAll(resp.Body)
	var userResp Auth0UserResponse
	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, err
	}

	return &userResp, nil
}

func findOrCreateUserFromAuth0(auth0ID, email, name string) (*db.User, error) {
	var user db.User
	err := db.Pool.QueryRow(context.Background(),
		"SELECT id, auth0_id, email, name, created_at FROM users WHERE auth0_id=$1",
		auth0ID).Scan(&user.ID, &user.Auth0ID, &user.Email, &user.Name, &user.CreatedAt)

	if err == nil {
		return &user, nil
	}

	user = db.User{
		ID:      auth0ID,
		Auth0ID: auth0ID,
		Email:   email,
		Name:    name,
	}

	_, err = db.Pool.Exec(context.Background(),
		"INSERT INTO users (id, auth0_id, email, name, created_at) VALUES ($1,$2,$3,$4,now())",
		user.ID, user.Auth0ID, user.Email, user.Name)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func getAuth0Connection(provider string) string {
	switch provider {
	case "google":
		return "google-oauth2"
	case "github":
		return "github"
	default:
		return "google-oauth2"
	}
}
