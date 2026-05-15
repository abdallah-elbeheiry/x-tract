package controllers

import (
	"context"
	"errors"
	"net/http"
	"x-tract/auth"
	"x-tract/data"
	"x-tract/models"

	"github.com/gin-gonic/gin"
)

// LoginStore describes the credential lookup behavior needed by the auth controller.
type LoginStore interface {
	Authenticate(ctx context.Context, email string, password string) (*models.User, error)
}

// AuthController exposes authentication endpoints.
type AuthController struct {
	store   LoginStore
	manager *auth.Manager
}

// NewAuthController builds an AuthController.
func NewAuthController(store LoginStore, manager *auth.Manager) *AuthController {
	return &AuthController{store: store, manager: manager}
}

// Login handles POST /auth/login.
func (ctl *AuthController) Login(c *gin.Context) {
	input, ok := bindJSON[models.LoginRequest](c)
	if !ok {
		return
	}

	user, err := ctl.store.Authenticate(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		if errors.Is(err, data.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		writeError(c, err)
		return
	}

	token, expiresAt, err := ctl.manager.Issue(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue token"})
		return
	}

	respondWithData(c, http.StatusOK, models.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt,
		User:        user,
	})
}
