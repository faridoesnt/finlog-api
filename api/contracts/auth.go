package contracts

import (
	"context"
	"time"

	"finlog-api/api/entities"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByID(ctx context.Context, id int64) (*entities.User, error)
	CreateUser(ctx context.Context, user *entities.User) (int64, error)
	UpdateVerificationToken(ctx context.Context, userID int64, token *string, expiresAt *time.Time) error
	FindByVerificationToken(ctx context.Context, token string) (*entities.User, error)
	MarkUserAsVerified(ctx context.Context, userID int64) error
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, string, *entities.User, error)
	Register(ctx context.Context, email, password string) (*entities.User, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, *entities.User, error)
	Logout(ctx context.Context, userID int64) error
	VerifyEmail(ctx context.Context, token string) (*entities.User, error)
	ResendVerification(ctx context.Context, email string) error
}
