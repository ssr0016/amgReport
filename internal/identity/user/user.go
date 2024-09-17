package user

import "context"

type Service interface {
	CreateUser(ctx context.Context, cmd *CreateUserCommand) error
	UpdateUser(ctx context.Context, cmd *UpdateUserCommand) error
	GetByUserID(ctx context.Context, id int64) (*User, error)
	DeleteUser(ctx context.Context, id int64) error
	SearchUser(ctx context.Context, query *SearchUserQuery) (*SearchUserResult, error)
	GetUserByEmail(ctx context.Context, cmd *LoginUserCommand) (string, error)

	RegisterDefaultUser(ctx context.Context, cmd *RegisterUserCommand) error
}
