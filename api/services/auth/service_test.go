package auth

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"finlog-api/api/contracts"
	"finlog-api/api/entities"

	"golang.org/x/crypto/bcrypt"
)

type fakeRepo struct {
	usersByEmail map[string]*entities.User
	usersByID    map[int64]*entities.User
	nextID       int64
}

func newFakeRepo(user *entities.User) *fakeRepo {
	usersByEmail := map[string]*entities.User{}
	usersByID := map[int64]*entities.User{}
	if user != nil {
		usersByEmail[user.Email] = user
		usersByID[user.ID] = user
	}
	return &fakeRepo{
		usersByEmail: usersByEmail,
		usersByID:    usersByID,
		nextID:       1,
	}
}

func (f *fakeRepo) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	if u, ok := f.usersByEmail[email]; ok {
		return u, nil
	}
	return nil, sql.ErrNoRows
}

func (f *fakeRepo) FindByID(ctx context.Context, id int64) (*entities.User, error) {
	if u, ok := f.usersByID[id]; ok {
		return u, nil
	}
	return nil, sql.ErrNoRows
}

func (f *fakeRepo) CreateUser(ctx context.Context, user *entities.User) (int64, error) {
	f.nextID++
	user.ID = f.nextID
	f.usersByEmail[user.Email] = user
	f.usersByID[user.ID] = user
	return user.ID, nil
}

func (f *fakeRepo) FindByVerificationToken(ctx context.Context, token string) (*entities.User, error) {
	return nil, sql.ErrNoRows
}

func (f *fakeRepo) UpdateVerificationToken(ctx context.Context, userID int64, token *string, expiresAt *time.Time) error {
	return nil
}

func (f *fakeRepo) MarkUserAsVerified(ctx context.Context, userID int64) error {
	return nil
}

func TestLoginSuccess(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	repo := newFakeRepo(&entities.User{
		ID:         1,
		Email:      "user@example.com",
		Name:       "user",
		Role:       "user",
		Password:   string(hashed),
		IsVerified: true,
	})
	svc := &Service{
		app: &contracts.App{Config: map[string]string{
			"JWT_SECRET":     "secret",
			"REFRESH_SECRET": "refresh",
			"JWT_TTL":        "1h",
			"REFRESH_TTL":    "24h",
		}},
		repo: repo,
	}

	access, refresh, user, err := svc.Login(context.Background(), "user@example.com", "secret123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if access == "" || refresh == "" {
		t.Fatalf("tokens should not be empty")
	}
	if user.Password != "" {
		t.Fatalf("password should be sanitized")
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	repo := newFakeRepo(&entities.User{
		ID:         1,
		Email:      "user@example.com",
		Name:       "user",
		Role:       "user",
		Password:   string(hashed),
		IsVerified: true,
	})
	svc := &Service{
		app: &contracts.App{Config: map[string]string{
			"JWT_SECRET":     "secret",
			"REFRESH_SECRET": "refresh",
			"JWT_TTL":        "1h",
			"REFRESH_TTL":    "24h",
		}},
		repo: repo,
	}

	_, _, _, err := svc.Login(context.Background(), "user@example.com", "wrong")
	if err == nil || !errors.Is(err, errInvalidCredentials) {
		t.Fatalf("expected invalid credentials error, got %v", err)
	}
}

func TestMustDurationFallback(t *testing.T) {
	d := mustDuration("bad", time.Hour)
	if d != time.Hour {
		t.Fatalf("expected fallback duration, got %v", d)
	}
}
