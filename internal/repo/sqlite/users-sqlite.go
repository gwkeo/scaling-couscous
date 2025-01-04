package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int64
	Email    string
	Password string
}

type UsersRepoSqlite struct {
	db *sql.DB
}

func NewUsersRepo(ctx context.Context, path string) (*UsersRepoSqlite, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	stmt, err := db.PrepareContext(ctx, `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(255) UNIQUE NOT NULL,
	    password VARCHAR(255) NOT NULL);
	
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare create table: %w", err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create table: %w", err)
	}

	return &UsersRepoSqlite{db: db}, nil
}

func (r *UsersRepoSqlite) CreateUser(ctx context.Context, user *User) (int64, error) {
	errLocation := "create user"
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO users (email, password) VALUES (?, ?);")
	if err != nil {
		return 0, fmt.Errorf("%s: failed to prepare statement: %w", errLocation, err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	res, err := stmt.ExecContext(ctx, user.Email, user.Password)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to execute create user: %w", errLocation, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last user id: %w", errLocation, err)
	}

	return id, nil
}

func (r *UsersRepoSqlite) User(ctx context.Context, id int64) (*User, error) {
	errLocation := "get user by id"
	user := &User{}

	stmt, err := r.db.PrepareContext(ctx, "SELECT * FROM users WHERE id = ?;")
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare statement: %w", errLocation, err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	err = stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("user not found")
		}
		return nil, fmt.Errorf("%s: failed to query user by id: %w", errLocation, err)
	}

	return user, nil
}

func (r *UsersRepoSqlite) UserByEmail(ctx context.Context, email string) (*User, error) {
	errLocation := "get user by email"
	user := User{}

	stmt, err := r.db.PrepareContext(ctx, "SELECT * FROM users WHERE email = ?;")
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare statement: %w", errLocation, err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	rows, err := stmt.QueryContext(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %w", errLocation, err)
	}

	if rows.Next() {
		if err = rows.Scan(&user.Email, &user.Password); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", errLocation, err)
		}
	}

	return &user, nil
}

func (r *UsersRepoSqlite) UpdateUser(ctx context.Context, user *User) error {
	errLocation := "update user"
	stmt, err := r.db.PrepareContext(ctx, "UPDATE users SET email = ?, password = ? WHERE id = ?;")
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", errLocation, err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	_, err = stmt.ExecContext(ctx, user.Email, user.Password, user.ID)
	if err != nil {
		return fmt.Errorf("%s: failed to execute query: %w", errLocation, err)
	}

	return nil
}
func (r *UsersRepoSqlite) DeleteUser(ctx context.Context, id int64) error {
	errLocation := "delete user"

	stmt, err := r.db.PrepareContext(ctx, "DELETE FROM users WHERE id = ?;")
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", errLocation, err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: failed to execute query: %w", errLocation, err)
	}

	return nil
}
