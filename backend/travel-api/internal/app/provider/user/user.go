package user

import (
	"context"

	userUsecase "github.com/f4mk/travel/backend/travel-api/internal/app/usecase/user"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Repo struct {
	repo *sqlx.DB
	log  *zerolog.Logger
}

func NewRepo(l *zerolog.Logger, r *sqlx.DB) *Repo {
	return &Repo{repo: r, log: l}
}

func (r *Repo) QueryAll(ctx context.Context) ([]userUsecase.User, error) {
	res := []userUsecase.User{}
	q := `SELECT * from users`
	if err := r.repo.SelectContext(ctx, &res, q); err != nil {
		return []userUsecase.User{}, err
	}
	return res, nil
}

func (r *Repo) Create(ctx context.Context, u userUsecase.User) error {
	q := `INSERT INTO users(user_id, name, email, is_active, is_deleted, token_version, roles, password_hash, date_created, date_updated) 
			VALUES(:user_id, :name, :email, :is_active, :is_deleted, :token_version, :roles, :password_hash, :date_created, :date_updated);`
	_, err := r.repo.NamedExecContext(ctx, q, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) QueryByID(ctx context.Context, userID string) (userUsecase.User, error) {
	res := userUsecase.User{}
	q := `SELECT * from users WHERE user_id = $1`
	if err := r.repo.GetContext(ctx, &res, q, userID); err != nil {
		return userUsecase.User{}, err
	}
	return res, nil
}

func (r *Repo) QueryByEmail(ctx context.Context, email string) (userUsecase.User, error) {
	res := userUsecase.User{}
	q := `SELECT * from users WHERE email = $1`
	if err := r.repo.GetContext(ctx, &res, q, email); err != nil {
		return userUsecase.User{}, err
	}
	return res, nil
}

func (r *Repo) QueryTokenByEmail(ctx context.Context, email string) (userUsecase.VerifyToken, error) {
	res := userUsecase.VerifyToken{}
	q := `SELECT * from verify_tokens WHERE email = $1`
	if err := r.repo.GetContext(ctx, &res, q, email); err != nil {
		return userUsecase.VerifyToken{}, err
	}
	return res, nil
}

func (r *Repo) Update(ctx context.Context, u userUsecase.User) error {
	q := `UPDATE users SET 
					name = :name, email = :email, 
					roles = :roles, password_hash = :password_hash, 
					date_updated = :date_updated
				WHERE user_id = :user_id;`
	_, err := r.repo.NamedExecContext(ctx, q, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) Verify(ctx context.Context, u userUsecase.User) error {
	q := `UPDATE users SET 
					is_active = :is_active, 
					date_updated = :date_updated
				WHERE user_id = :user_id;`
	_, err := r.repo.NamedExecContext(ctx, q, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) Delete(ctx context.Context, user userUsecase.User) error {
	q := `UPDATE users SET 
					is_active = :is_active, 
					is_deleted = :is_deleted, 
					date_updated = :date_updated,
					token_version = :token_version
				WHERE user_id = :user_id`
	_, err := r.repo.NamedExecContext(ctx, q, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) StoreVerifyToken(ctx context.Context, vt userUsecase.VerifyToken) error {
	q := `INSERT INTO verify_tokens (token_id, user_id, email, expires_at, issued_at)
	VALUES (:token_id, :user_id, :email, :expires_at, :issued_at)
	`
	_, err := r.repo.NamedExecContext(ctx, q, vt)
	return err
}

func (r *Repo) DeleteVerifyTokensByUserID(ctx context.Context, uID string) error {
	q := `DELETE from verify_tokens WHERE user_id = $1;`
	_, err := r.repo.ExecContext(ctx, q, uID)
	return err
}
