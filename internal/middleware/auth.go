package middleware

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"

	"github.com/KILLERGTG01/smart-task-planner-be/internal/config"
)

type AuthMiddleware struct {
	jwks   *keyfunc.JWKS
	mu     sync.RWMutex
	config *config.Config
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		config: cfg,
	}
}

func (a *AuthMiddleware) ensureJWKS(ctx context.Context) error {
	a.mu.RLock()
	if a.jwks != nil {
		a.mu.RUnlock()
		return nil
	}
	a.mu.RUnlock()

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.jwks != nil {
		return nil
	}

	issuer := strings.TrimSuffix(a.config.Auth0Issuer, "/")
	jwksURL := issuer + "/.well-known/jwks.json"

	options := keyfunc.Options{
		Ctx:             ctx,
		RefreshInterval: time.Hour,
		RefreshTimeout:  10 * time.Second,
		RefreshErrorHandler: func(err error) {
			zap.L().Error("JWKS refresh failed", zap.Error(err))
		},
	}

	k, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		zap.L().Error("Failed to initialize JWKS", zap.Error(err), zap.String("jwks_url", jwksURL))
		return err
	}

	a.jwks = k
	zap.L().Info("JWKS initialized successfully")
	return nil
}

func (a *AuthMiddleware) Cleanup() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.jwks != nil {
		a.jwks.EndBackground()
		a.jwks = nil
	}
}

func (a *AuthMiddleware) AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
		defer cancel()

		if err := a.ensureJWKS(ctx); err != nil {
			zap.L().Error("JWKS initialization failed", zap.Error(err))
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "authentication_service_unavailable",
			})
		}

		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		tokenStr := parts[1]
		if len(tokenStr) == 0 {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		a.mu.RLock()
		jwks := a.jwks
		a.mu.RUnlock()

		token, err := jwt.Parse(tokenStr, jwks.Keyfunc)

		if err != nil || !token.Valid {
			zap.L().Warn("Invalid token provided", zap.Error(err))
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		// Validate issuer
		if iss, ok := claims["iss"].(string); !ok || iss != a.config.Auth0Issuer {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		// Validate audience
		if aud, ok := claims["aud"].(string); !ok || aud != a.config.Auth0Aud {
			// Also check if audience is an array
			if audArray, ok := claims["aud"].([]interface{}); ok {
				found := false
				for _, audItem := range audArray {
					if audStr, ok := audItem.(string); ok && audStr == a.config.Auth0Aud {
						found = true
						break
					}
				}
				if !found {
					return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
						"error": "unauthorized",
					})
				}
			} else {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "unauthorized",
				})
			}
		}

		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "token_expired",
				})
			}
		} else {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		if sub, ok := claims["sub"].(string); ok && sub != "" {
			c.Locals("auth_sub", sub)
		} else {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		c.Locals("auth_claims", claims)
		return c.Next()
	}
}
