package auth

import (
	"context"
	"time"

	"finlog-api/api/contracts"
	"finlog-api/api/datasources"
	"finlog-api/api/entities"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	app  *contracts.App
	stmt Statement
}

type Statement struct {
	findByEmail             *sqlx.Stmt
	findByID                *sqlx.Stmt
	findByVerificationToken *sqlx.Stmt
	insertUser              *sqlx.Stmt
	updateVerificationToken *sqlx.Stmt
	markUserVerified        *sqlx.Stmt
}

func initRepository(app *contracts.App) contracts.AuthRepository {
	stmts := Statement{
		findByEmail:             datasources.Prepare(app.Ds.ReaderDB, findByEmail),
		findByID:                datasources.Prepare(app.Ds.ReaderDB, findByID),
		findByVerificationToken: datasources.Prepare(app.Ds.ReaderDB, findByVerificationToken),
		insertUser:              datasources.Prepare(app.Ds.WriterDB, insertUser),
		updateVerificationToken: datasources.Prepare(app.Ds.WriterDB, updateVerificationToken),
		markUserVerified:        datasources.Prepare(app.Ds.WriterDB, markUserVerified),
	}

	r := Repository{
		app:  app,
		stmt: stmts,
	}

	return &r
}

// FindByEmail returns user by email.
func (r *Repository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	user := new(entities.User)
	err := r.stmt.findByEmail.GetContext(ctx, user, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindByID returns user by ID.
func (r *Repository) FindByID(ctx context.Context, id int64) (*entities.User, error) {
	user := new(entities.User)
	err := r.stmt.findByID.GetContext(ctx, user, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser persists a new user and returns its id.
func (r *Repository) CreateUser(ctx context.Context, user *entities.User) (int64, error) {
	res, err := r.stmt.insertUser.ExecContext(
		ctx,
		user.Email,
		user.Name,
		user.Role,
		user.Password,
		user.IsVerified,
		user.VerificationToken,
		user.VerificationExpiresAt,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repository) FindByVerificationToken(ctx context.Context, token string) (*entities.User, error) {
	user := new(entities.User)
	if err := r.stmt.findByVerificationToken.GetContext(ctx, user, token); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) UpdateVerificationToken(ctx context.Context, userID int64, token *string, expiresAt *time.Time) error {
	_, err := r.stmt.updateVerificationToken.ExecContext(ctx, token, expiresAt, userID)
	return err
}

func (r *Repository) MarkUserAsVerified(ctx context.Context, userID int64) error {
	_, err := r.stmt.markUserVerified.ExecContext(ctx, userID)
	return err
}
