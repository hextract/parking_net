package domain

type UserRole string

const (
	UserRoleDriver UserRole = "driver"
	UserRoleOwner  UserRole = "owner"
	UserRoleAdmin  UserRole = "admin"
)

type User struct {
	ID         string
	Login      string
	Email      string
	Role       UserRole
	TelegramID int64
}

func (u *User) IsDriver() bool {
	return u.Role == UserRoleDriver
}

func (u *User) IsOwner() bool {
	return u.Role == UserRoleOwner
}

func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

