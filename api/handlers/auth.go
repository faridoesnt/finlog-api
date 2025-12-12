package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"finlog-api/api/models/responses"
	"finlog-api/api/services/auth"
)

// AuthLogin returns tokens in a flat response for the mobile client.
func AuthLogin(c *fiber.Ctx) error {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return responses.BadRequest(err)
	}
	access, refresh, user, err := app.Services.Auth.Login(context.Background(), body.Email, body.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrEmailNotVerified()):
			return responses.UnAuthorized(fmt.Errorf("email_not_verified"))
		default:
			return responses.UnAuthorized(err)
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  access,
		"refresh_token": refresh,
		"email":         user.Email,
	})
}

// Register creates a new user and issues tokens.
func Register(c *fiber.Ctx) error {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return responses.BadRequest(err)
	}
	_, err := app.Services.Auth.Register(context.Background(), body.Email, body.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials()) || errors.Is(err, auth.ErrInvalidInput()):
			return responses.BadRequest(err)
		case errors.Is(err, auth.ErrEmailExists()):
			return responses.Conflict(err)
		default:
			return responses.BadRequest(err)
		}
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Verification email sent",
	})
}

func VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return responses.BadRequest(errors.New("verification token is required"))
	}
	if _, err := app.Services.Auth.VerifyEmail(context.Background(), token); err != nil {
		switch {
		case errors.Is(err, auth.ErrVerificationTokenInvalid()):
			return responses.BadRequest(err)
		case errors.Is(err, auth.ErrVerificationTokenExpired()):
			return responses.BadRequest(err)
		case errors.Is(err, auth.ErrEmailAlreadyVerified()):
			// already verified, continue to redirect
		default:
			return responses.BadRequest(err)
		}
	}
	redirectURL := getVerificationRedirect(c.Get(fiber.HeaderUserAgent))
	return c.Redirect(redirectURL, fiber.StatusFound)
}

func ResendVerification(c *fiber.Ctx) error {
	type req struct {
		Email string `json:"email"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return responses.BadRequest(err)
	}
	if err := app.Services.Auth.ResendVerification(context.Background(), body.Email); err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound()):
			return responses.NotFound(err)
		case errors.Is(err, auth.ErrEmailAlreadyVerified()):
			return responses.BadRequest(err)
		default:
			return responses.InternalServerError(err)
		}
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Email aktivasi berhasil dikirim ulang",
	})
}

func getVerificationRedirect(agent string) string {
	lower := strings.ToLower(agent)
	if strings.Contains(lower, "android") || strings.Contains(lower, "iphone") || strings.Contains(lower, "ipad") {
		return "finlog://login?verified=true"
	}
	return "https://finlog.app/activated"
}

// Refresh renews access/refresh tokens using the provided refresh token.
func Refresh(c *fiber.Ctx) error {
	type req struct {
		RefreshToken string `json:"refresh_token"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return responses.BadRequest(err)
	}
	access, refresh, user, err := app.Services.Auth.Refresh(context.Background(), body.RefreshToken)
	if err != nil {
		return responses.UnAuthorized(err)
	}
	return c.JSON(fiber.Map{
		"access_token":  access,
		"refresh_token": refresh,
		"email":         user.Email,
	})
}

// Logout ends the session (stateless JWT, no server state is kept).
func Logout(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	if err := app.Services.Auth.Logout(context.Background(), userID); err != nil {
		return responses.BadRequest(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
