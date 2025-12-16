package handlers

import (
	"context"
	"errors"
	"fmt"

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
	redirectURL := getVerificationRedirect()
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

func getVerificationRedirect() string {
	return "https://api.finlog.asia/activated"
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

func ActivatedHandler(c *fiber.Ctx) error {
	return c.Type("html").SendString(`
		<!DOCTYPE html>
		<html lang="id">
		<head>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		<title>FinLog</title>

		<script>
			window.onload = function () {
			window.location.href = "finlog://login?verified=true";
			setTimeout(function () {
				document.getElementById("openApp").style.display = "block";
			}, 1500);
			};
		</script>

		<style>
			body {
			font-family: system-ui, -apple-system, BlinkMacSystemFont;
			text-align: center;
			padding: 40px;
			}
			a {
			display: inline-block;
			margin-top: 20px;
			padding: 12px 20px;
			background: #111;
			color: #fff;
			text-decoration: none;
			border-radius: 8px;
			}
		</style>
		</head>

		<body>
		<h2>Email berhasil diverifikasi</h2>
		<p>Membuka aplikasi FinLog…</p>

		<a id="openApp" href="finlog://login?verified=true" style="display:none">
			Buka Aplikasi FinLog
		</a>
		</body>
		</html>
	`)
}