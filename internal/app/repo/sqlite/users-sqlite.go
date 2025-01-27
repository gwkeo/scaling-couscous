package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gwkeo/scaling-couscous/internal/app/repo"
	"github.com/gwkeo/scaling-couscous/internal/app/repo/models"
	_ "github.com/mattn/go-sqlite3"
)

type UsersRepoSqlite struct {
	db *sql.DB
}

func NewUsersRepo(ctx context.Context, path string) (*UsersRepoSqlite, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, errors.New("unable to connect to sqlite")
	}

	stmt, err := db.PrepareContext(ctx, `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(255) UNIQUE NOT NULL,
	    password VARCHAR(255) NOT NULL);
	
	`)
	if err != nil {
		return nil, errors.New("failed to create table")
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return nil, errors.New("failed to insert user")
	}

	return &UsersRepoSqlite{db: db}, nil
}

func (r *UsersRepoSqlite) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO users (email, password) VALUES (?, ?);")
	if err != nil {
		return 0, repo.ErrFailedToPrepareStmt
	}

	defer func() {
		_ = stmt.Close()
	}()

	res, err := stmt.ExecContext(ctx, user.Email, user.Password)
	if err != nil {
		return 0, repo.ErrFailedToInsertUser
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, repo.ErrFailedToGetLastId
	}

	return id, nil
}

func (r *UsersRepoSqlite) User(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}

	stmt, err := r.db.PrepareContext(ctx, "SELECT * FROM users WHERE id = ?;")
	if err != nil {
		return nil, repo.ErrFailedToPrepareStmt
	}
	defer func() {
		_ = stmt.Close()
	}()

	err = stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, repo.ErrUserNotFound
		}
		return nil, repo.ErrFailedToGetUser
	}

	return user, nil
}

func (r *UsersRepoSqlite) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := models.User{}

	stmt, err := r.db.PrepareContext(ctx, "SELECT * FROM users WHERE email = ?;")
	if err != nil {
		return nil, repo.ErrFailedToPrepareStmt
	}
	defer func() {
		_ = stmt.Close()
	}()

	err = stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrUserNotFound
		}
		return nil, repo.ErrFailedToGetUserByEmail
	}

	return &user, nil
}

func (r *UsersRepoSqlite) UpdateUser(ctx context.Context, user *models.User) error {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE users SET email = ?, password = ? WHERE id = ?;")
	if err != nil {
		return repo.ErrFailedToPrepareStmt
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	_, err = stmt.ExecContext(ctx, user.Email, user.Password, user.ID)
	if err != nil {
		return repo.ErrFailedToUpdateUser
	}

	return nil
}
func (r *UsersRepoSqlite) DeleteUser(ctx context.Context, id int64) error {
	stmt, err := r.db.PrepareContext(ctx, "DELETE FROM users WHERE id = ?;")
	if err != nil {
		return repo.ErrFailedToPrepareStmt
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return repo.ErrFailedToDeleteUser
	}

	return nil
}
