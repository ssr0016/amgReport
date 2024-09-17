package user

import (
	"amg/internal/api/errors"
	util "amg/pkg/util/password"
	"amg/pkg/util/validation"
	"strings"
)

var (
	ErrInvalidEmail       = errors.New("user.invalid-email", "Invalid email")
	ErrInvalidID          = errors.New("user.invalid-id", "Invalid id")
	ErrUserAlreadyExists  = errors.New("user.already-exists", "User already exists")
	ErrUserNotFound       = errors.New("user.not-found", "User not found")
	ErrInvalidPassword    = errors.New("user.invalid-password", "Invalid password")
	ErrInvalidFirstName   = errors.New("user.invalid-first-name", "Invalid first name")
	ErrInvalidLastName    = errors.New("user.invalid-last-name", "Invalid last name")
	ErrInvalidAddress     = errors.New("user.invalid-address", "Invalid address")
	ErrInvalidPhoneNumber = errors.New("user.invalid-phone-number", "Invalid phone number")
	ErrInvalidDateOfBirth = errors.New("user.invalid-date-of-birth", "Invalid date of birth")
	ErrEmailAlreadyExists = errors.New("user.email-already-exists", "Email already exists")
	ErrorInvalidRole      = errors.New("user.invalid-role", "Invalid role")
	ErrInvalidStatus      = errors.New("user.invalid-status", "Invalid status")
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID           int64  `db:"id" json:"id"`
	FirstName    string `db:"first_name" json:"first_name"`
	LastName     string `db:"last_name" json:"last_name"`
	Email        string `db:"email" json:"email"`
	PasswordHash string `db:"password_hash" json:"-"`
	Address      string `db:"address" json:"address"`
	PhoneNumber  string `db:"phone_number" json:"phone_number"`
	DateOfBirth  string `db:"date_of_birth" json:"date_of_birth"`
	Role         string `db:"role" json:"role"`
	CreatedAt    string `db:"created_at" json:"created_at"`
	UpdatedAt    string `db:"updated_at" json:"updated_at"`
}

var validRoles = map[string]bool{
	RoleUser:  true,
	RoleAdmin: true,
}

func IsValidRole(role string) bool {
	return validRoles[role]
}

type CreateUserCommand struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password_hash"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	DateOfBirth string `json:"date_of_birth"`
	Role        string `json:"role"`
}

type UpdateUserCommand struct {
	ID          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	DateOfBirth string `json:"date_of_birth"`
	Role        string `json:"role"`
}

type SearchUserQuery struct {
	FirstName   string `query:"first_name"`
	LastName    string `query:"last_name"`
	Email       string `query:"email"`
	Address     string `query:"address"`
	PhoneNumber string `query:"phone_number"`
	DateOfBirth string `query:"date_of_birth"`
	Role        string `query:"role"`
	Page        int    `query:"page"`
	PerPage     int    `query:"per_page"`
}

type SearchUserResult struct {
	TotalCount int64   `json:"total_count"`
	Users      []*User `json:"users"`
	Page       int     `json:"page"`
	PerPage    int     `json:"per_page"`
}

type RegisterUserCommand struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password_hash"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	DateOfBirth string `json:"date_of_birth"`
}

type LoginUserCommand struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogoutUserCommand struct {
	Token string `json:"token"`
}

func (cmd *CreateUserCommand) Validate() error {
	if len(cmd.FirstName) == 0 || len(cmd.FirstName) <= 2 {
		return ErrInvalidFirstName
	}
	if len(cmd.LastName) == 0 || len(cmd.LastName) <= 2 {
		return ErrInvalidLastName
	}
	if len(cmd.Email) == 0 || !validation.IsValidEmail(cmd.Email) {
		return ErrInvalidEmail
	}
	if len(cmd.Password) == 0 || !util.IsValidPassword(cmd.Password) {
		return ErrInvalidPassword
	}
	if len(cmd.Address) == 0 {
		return ErrInvalidAddress
	}
	if len(cmd.PhoneNumber) == 0 {
		return ErrInvalidPhoneNumber
	}
	if len(cmd.DateOfBirth) == 0 {
		return ErrInvalidDateOfBirth
	}
	if !IsValidRole(cmd.Role) {
		return ErrorInvalidRole
	}

	return nil
}

func (cmd *UpdateUserCommand) Validate() error {
	if cmd.ID == 0 {
		return ErrUserNotFound
	}
	if len(strings.TrimSpace(cmd.FirstName)) == 0 || len(cmd.FirstName) <= 2 {
		return ErrInvalidFirstName
	}
	if len(strings.TrimSpace(cmd.LastName)) == 0 || len(cmd.LastName) <= 2 {
		return ErrInvalidLastName
	}
	if len(strings.TrimSpace(cmd.Address)) == 0 {
		return ErrInvalidAddress
	}
	if len(cmd.PhoneNumber) == 0 || !validation.IsValidPhoneNumber(cmd.PhoneNumber) {
		return ErrInvalidPhoneNumber
	}
	if len(cmd.Email) > 0 && !validation.IsValidEmail(cmd.Email) {
		return ErrInvalidEmail
	}
	if !IsValidRole(cmd.Role) {
		return ErrorInvalidRole
	}
	return nil
}

func (cmd *RegisterUserCommand) Validate() error {
	if len(cmd.FirstName) == 0 || len(cmd.FirstName) <= 2 {
		return ErrInvalidFirstName
	}
	if len(cmd.LastName) == 0 || len(cmd.LastName) <= 2 {
		return ErrInvalidLastName
	}
	if len(cmd.Email) == 0 || !validation.IsValidEmail(cmd.Email) {
		return ErrInvalidEmail
	}
	if len(cmd.Password) == 0 || !util.IsValidPassword(cmd.Password) {
		return ErrInvalidPassword
	}
	if len(cmd.Address) == 0 {
		return ErrInvalidAddress
	}
	if len(cmd.PhoneNumber) == 0 {
		return ErrInvalidPhoneNumber
	}
	if len(cmd.DateOfBirth) == 0 {
		return ErrInvalidDateOfBirth
	}
	return nil
}

// Validation for LoginUserCommand
func (cmd *LoginUserCommand) Validate() error {
	if len(cmd.Email) == 0 || !validation.IsValidEmail(cmd.Email) {
		return ErrInvalidEmail
	}
	if len(cmd.Password) == 0 || !util.IsValidPassword(cmd.Password) {
		return ErrInvalidPassword
	}
	return nil
}
