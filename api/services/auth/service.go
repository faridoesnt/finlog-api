package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	// "net"
	"net/mail"
	"net/smtp"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"finlog-api/api/constants"
	"finlog-api/api/contracts"
	"finlog-api/api/entities"
	"finlog-api/api/seeds"
	"finlog-api/api/services/category"
)

const (
	verificationTTL = 24 * time.Hour
	emailSubject    = "Aktivasi Akun FinLog"
)

var (
	errInvalidCredentials       = errors.New("invalid credentials")
	errEmailExists              = errors.New("email already registered")
	errInvalidInput             = errors.New("invalid credentials")
	errEmailNotVerified         = errors.New("email not verified")
	errVerificationTokenInvalid = errors.New("invalid verification token")
	errVerificationTokenExpired = errors.New("verification token expired")
	errEmailAlreadyVerified     = errors.New("email already verified")
	errUserNotFound             = errors.New("user not found")
)

func ErrInvalidCredentials() error       { return errInvalidCredentials }
func ErrEmailExists() error              { return errEmailExists }
func ErrInvalidInput() error             { return errInvalidInput }
func ErrEmailNotVerified() error         { return errEmailNotVerified }
func ErrVerificationTokenInvalid() error { return errVerificationTokenInvalid }
func ErrVerificationTokenExpired() error { return errVerificationTokenExpired }
func ErrEmailAlreadyVerified() error     { return errEmailAlreadyVerified }
func ErrUserNotFound() error             { return errUserNotFound }

// Service endpoints require the application context.
type Service struct {
	app     *contracts.App
	repo    contracts.AuthRepository
	catRepo contracts.CategoryRepository
}

func Init(app *contracts.App) contracts.AuthService {
	repo := initRepository(app)

	return &Service{
		app:     app,
		repo:    repo,
		catRepo: category.NewRepository(app),
	}
}

// Login authenticates user.
func (s *Service) Login(ctx context.Context, email, password string) (string, string, *entities.User, error) {
	email = normalizeEmail(email)
	if err := validateCredentials(email, password); err != nil {
		return "", "", nil, errInvalidCredentials
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", nil, errInvalidCredentials
		}
		return "", "", nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", "", nil, errInvalidCredentials
	}
	if !user.IsVerified {
		return "", "", nil, errEmailNotVerified
	}

	return s.issueTokens(user)
}

// Register creates a new user account.
func (s *Service) Register(ctx context.Context, email, password string) (*entities.User, error) {
	email = normalizeEmail(email)
	if err := validateCredentials(email, password); err != nil {
		return nil, errInvalidInput
	}
	if _, err := s.repo.FindByEmail(ctx, email); err == nil {
		return nil, errEmailExists
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	rawToken, hashedToken, expiresAt, err := s.prepareVerificationToken()
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Email:                 email,
		Name:                  defaultName(email),
		Role:                  "user",
		Password:              string(hashedPassword),
		IsVerified:            false,
		VerificationToken:     &hashedToken,
		VerificationExpiresAt: &expiresAt,
	}

	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = id
	s.seedDefaultCategories(ctx, user.ID)

	if err := s.sendVerificationEmail(user, rawToken); err != nil {
		return nil, err
	}

	return user, nil
}

// Refresh renews tokens using a refresh token.
func (s *Service) Refresh(ctx context.Context, refreshToken string) (string, string, *entities.User, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return "", "", nil, errInvalidCredentials
	}
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalidCredentials
		}
		return []byte(s.app.Config[constants.REFRESH_SECRET]), nil
	})
	if err != nil || !token.Valid {
		return "", "", nil, errInvalidCredentials
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", nil, errInvalidCredentials
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		return "", "", nil, errInvalidCredentials
	}
	user, err := s.repo.FindByID(ctx, int64(sub))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", nil, errInvalidCredentials
		}
		return "", "", nil, err
	}
	return s.issueTokens(user)
}

// Logout keeps stateless JWT flow but allows future revocation hook.
func (s *Service) Logout(ctx context.Context, userID int64) error {
	_ = userID
	return nil
}

// VerifyEmail completes account activation.
func (s *Service) VerifyEmail(ctx context.Context, token string) (*entities.User, error) {
	hashed := hashToken(token)
	user, err := s.repo.FindByVerificationToken(ctx, hashed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errVerificationTokenInvalid
		}
		return nil, err
	}
	if user.IsVerified {
		return nil, errEmailAlreadyVerified
	}
	if user.VerificationExpiresAt == nil || user.VerificationExpiresAt.Before(time.Now()) {
		return nil, errVerificationTokenExpired
	}
	if err := s.repo.MarkUserAsVerified(ctx, user.ID); err != nil {
		return nil, err
	}
	user.IsVerified = true
	user.VerificationToken = nil
	user.VerificationExpiresAt = nil
	return user, nil
}

// ResendVerification sends a fresh token to the user email.
func (s *Service) ResendVerification(ctx context.Context, email string) error {
	email = normalizeEmail(email)
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errUserNotFound
		}
		return err
	}
	if user.IsVerified {
		return errEmailAlreadyVerified
	}

	rawToken, hashedToken, expiresAt, err := s.prepareVerificationToken()
	if err != nil {
		return err
	}
	if err := s.repo.UpdateVerificationToken(ctx, user.ID, &hashedToken, &expiresAt); err != nil {
		return err
	}
	user.VerificationToken = &hashedToken
	user.VerificationExpiresAt = &expiresAt

	return s.sendVerificationEmail(user, rawToken)
}

func (s *Service) issueTokens(user *entities.User) (string, string, *entities.User, error) {
	if user.Role == "" {
		user.Role = "user"
	}

	access, err := s.generateToken(user, []byte(s.app.Config[constants.JWT_SECRET]), mustDuration(s.app.Config[constants.JWT_TTL], time.Hour))
	if err != nil {
		return "", "", nil, err
	}
	refresh, err := s.generateToken(user, []byte(s.app.Config[constants.REFRESH_SECRET]), mustDuration(s.app.Config[constants.REFRESH_TTL], 7*24*time.Hour))
	if err != nil {
		return "", "", nil, err
	}
	safeUser := *user
	safeUser.Password = ""
	safeUser.VerificationToken = nil
	safeUser.VerificationExpiresAt = nil
	return access, refresh, &safeUser, nil
}

func (s *Service) generateToken(user *entities.User, secret []byte, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
		"exp":   time.Now().Add(ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func mustDuration(v string, fallback time.Duration) time.Duration {
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func validateCredentials(email, password string) error {
	if email == "" || password == "" {
		return errInvalidCredentials
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errInvalidCredentials
	}
	if len(password) < 6 {
		return errInvalidCredentials
	}
	return nil
}

func defaultName(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return email
}

func (s *Service) seedDefaultCategories(ctx context.Context, userID int64) {
	for _, seed := range seeds.DefaultCategories() {
		category := &entities.Category{
			UserID:    userID,
			Name:      seed.Name,
			IsExpense: seed.IsExpense,
			IconKey:   seed.IconKey,
		}
		if _, err := s.catRepo.Create(ctx, category); err != nil {
			s.app.Logger.Err(err).
				Str("category", seed.Name).
				Msg("failed to seed default category")
		}
	}
}

func (s *Service) prepareVerificationToken() (string, string, time.Time, error) {
	raw, err := generateRandomToken()
	if err != nil {
		return "", "", time.Time{}, err
	}
	hashed := hashToken(raw)
	expiresAt := time.Now().UTC().Add(verificationTTL)
	return raw, hashed, expiresAt, nil
}

func generateRandomToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func hashToken(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func (s *Service) sendVerificationEmail(user *entities.User, token string) error {
	host := strings.TrimSpace(s.app.Config[constants.SMTPHost])
	if host == "" {
		return errors.New("smtp host is not configured")
	}
	port := strings.TrimSpace(s.app.Config[constants.SMTPPort])
	if port == "" {
		port = "587"
	}
	username := strings.TrimSpace(s.app.Config[constants.SMTPUsername])
	password := strings.TrimSpace(s.app.Config[constants.SMTPPassword])
	from := strings.TrimSpace(s.app.Config[constants.SMTPFrom])
	if from == "" {
		return errors.New("smtp from address is not configured")
	}
	fromName := strings.TrimSpace(s.app.Config[constants.SMTPFromName])
	if fromName == "" {
		fromName = "FinLog"
	}

	_ = smtp.Auth(nil)
	if username != "" && password != "" {
		// auth = smtp.PlainAuth("", username, password, host)
	}

	activationURL := s.buildActivationURL(token)
	body := buildEmailBody(defaultName(user.Email), activationURL)

	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", fromName, from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", user.Email))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", emailSubject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	// addr := net.JoinHostPort(host, port)
	// recipients := []string{user.Email}
	// if err := smtp.SendMail(addr, auth, from, recipients, []byte(msg.String())); err != nil {
	// 	return fmt.Errorf("failed to send verification email: %w", err)
	// }
	return nil
}

func (s *Service) buildActivationURL(token string) string {
	baseURL := strings.TrimRight(strings.TrimSpace(s.app.Config[constants.APIBaseURL]), "/")
	if baseURL == "" {
		baseURL = "https://api.finlog.app"
	}
	escaped := url.QueryEscape(token)
	return fmt.Sprintf("%s/auth/verify?token=%s", baseURL, escaped)
}

func buildEmailBody(name, activationURL string) string {
	return fmt.Sprintf(`
<p>Hai %s,</p>
<p>Terima kasih telah mendaftar FinLog. Klik tombol di bawah ini untuk mengaktifkan akun Anda. Link ini berlaku 24 jam dan hanya bisa dipakai sekali.</p>
<p><a href="%s" style="background:#2e7d32;color:#fff;padding:12px 20px;border-radius:8px;text-decoration:none;display:inline-block;">Aktifkan Akun</a></p>
<p>Terima kasih,<br/>Tim FinLog</p>
`, name, activationURL)
}
