package models

type Role int64

const (
	AdminRole Role = iota
	UserRole
)

type User struct {
	ID       int64
	Email    string
	Password string
	Role     Role
}
