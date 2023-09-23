package user

import (
	"context"

	userUsecase "github.com/f4mk/travel/backend/travel-api/internal/app/usecase/user"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Storer struct {
	repo *sqlx.DB
	log  *zerolog.Logger
}

func NewStorer(l *zerolog.Logger, r *sqlx.DB) *Storer {
	return &Storer{repo: r, log: l}
}

func (s *Storer) QueryAll(ctx context.Context) ([]userUsecase.User, error) {
	users := []StorerUser{}
	q := `SELECT * from users;`
	if err := s.repo.SelectContext(ctx, &users, q); err != nil {
		return []userUsecase.User{}, err
	}
	var res []userUsecase.User
	for _, user := range res {
		res = append(res, userUsecase.User{
			ID:           user.ID,
			Name:         user.Name,
			Email:        user.Email,
			IsActive:     user.IsActive,
			IsDeleted:    user.IsDeleted,
			TokenVersion: user.TokenVersion,
			Roles:        user.Roles,
			PasswordHash: user.PasswordHash,
			DateCreated:  user.DateCreated,
			DateUpdated:  user.DateUpdated,
		})
	}
	return res, nil
}

func (s *Storer) Create(ctx context.Context, u userUsecase.User) error {
	user := StorerUser{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		IsActive:     u.IsActive,
		IsDeleted:    u.IsDeleted,
		TokenVersion: u.TokenVersion,
		Roles:        u.Roles,
		PasswordHash: u.PasswordHash,
		DateCreated:  u.DateCreated,
		DateUpdated:  u.DateUpdated,
	}
	q := `INSERT INTO users(user_id, name, email, is_active, is_deleted, token_version, roles, password_hash, date_created, date_updated) 
			VALUES(:user_id, :name, :email, :is_active, :is_deleted, :token_version, :roles, :password_hash, :date_created, :date_updated);`
	_, err := s.repo.NamedExecContext(ctx, q, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storer) QueryByID(ctx context.Context, userID string) (userUsecase.User, error) {
	user := StorerUser{}
	q := `SELECT * from users WHERE user_id = $1;`
	if err := s.repo.GetContext(ctx, &user, q, userID); err != nil {
		return userUsecase.User{}, err
	}
	res := userUsecase.User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		IsActive:     user.IsActive,
		IsDeleted:    user.IsDeleted,
		TokenVersion: user.TokenVersion,
		Roles:        user.Roles,
		PasswordHash: user.PasswordHash,
		DateCreated:  user.DateCreated,
		DateUpdated:  user.DateUpdated,
	}
	return res, nil
}

func (s *Storer) QueryByEmail(ctx context.Context, email string) (userUsecase.User, error) {
	user := StorerUser{}
	q := `SELECT * from users WHERE email = $1;`
	if err := s.repo.GetContext(ctx, &user, q, email); err != nil {
		return userUsecase.User{}, err
	}
	res := userUsecase.User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		IsActive:     user.IsActive,
		IsDeleted:    user.IsDeleted,
		TokenVersion: user.TokenVersion,
		Roles:        user.Roles,
		PasswordHash: user.PasswordHash,
		DateCreated:  user.DateCreated,
		DateUpdated:  user.DateUpdated,
	}
	return res, nil
}

func (s *Storer) QueryTokenByEmail(ctx context.Context, email string) (userUsecase.VerifyToken, error) {
	token := StorerToken{}
	q := `SELECT * from verify_tokens WHERE email = $1;`
	if err := s.repo.GetContext(ctx, &token, q, email); err != nil {
		return userUsecase.VerifyToken{}, err
	}
	res := userUsecase.VerifyToken{
		TokenID:   token.TokenID,
		UserID:    token.UserID,
		Email:     token.Email,
		ExpiresAt: token.ExpiresAt,
		IssuedAt:  token.IssuedAt,
	}
	return res, nil
}

func (s *Storer) Update(ctx context.Context, u userUsecase.User) error {
	user := StorerUser{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		IsActive:     u.IsActive,
		IsDeleted:    u.IsDeleted,
		TokenVersion: u.TokenVersion,
		Roles:        u.Roles,
		PasswordHash: u.PasswordHash,
		DateCreated:  u.DateCreated,
		DateUpdated:  u.DateUpdated,
	}
	q := `UPDATE users SET 
					name = :name, email = :email, 
					roles = :roles, password_hash = :password_hash, 
					date_updated = :date_updated
				WHERE user_id = :user_id;`
	_, err := s.repo.NamedExecContext(ctx, q, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storer) Verify(ctx context.Context, u userUsecase.User) error {
	user := StorerUser{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		IsActive:     u.IsActive,
		IsDeleted:    u.IsDeleted,
		TokenVersion: u.TokenVersion,
		Roles:        u.Roles,
		PasswordHash: u.PasswordHash,
		DateCreated:  u.DateCreated,
		DateUpdated:  u.DateUpdated,
	}
	q := `UPDATE users SET 
					is_active = :is_active, 
					date_updated = :date_updated
				WHERE user_id = :user_id;`
	_, err := s.repo.NamedExecContext(ctx, q, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storer) Delete(ctx context.Context, u userUsecase.User) error {
	user := StorerUser{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		IsActive:     u.IsActive,
		IsDeleted:    u.IsDeleted,
		TokenVersion: u.TokenVersion,
		Roles:        u.Roles,
		PasswordHash: u.PasswordHash,
		DateCreated:  u.DateCreated,
		DateUpdated:  u.DateUpdated,
	}
	q := `UPDATE users SET 
					is_active = :is_active, 
					is_deleted = :is_deleted, 
					date_updated = :date_updated,
					token_version = :token_version
				WHERE user_id = :user_id;`
	_, err := s.repo.NamedExecContext(ctx, q, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storer) StoreVerifyToken(ctx context.Context, vt userUsecase.VerifyToken) error {
	q := `INSERT INTO verify_tokens (token_id, user_id, email, expires_at, issued_at)
	VALUES (:token_id, :user_id, :email, :expires_at, :issued_at);
	`
	_, err := s.repo.NamedExecContext(ctx, q, vt)
	return err
}

func (s *Storer) DeleteVerifyTokensByUserID(ctx context.Context, uID string) error {
	q := `DELETE from verify_tokens WHERE user_id = $1;`
	_, err := s.repo.ExecContext(ctx, q, uID)
	return err
}
