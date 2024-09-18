package userimpl

import (
	"amg/config"
	"amg/internal/db"
	"amg/internal/identity/user"
	"amg/pkg/util/jwt"
	util "amg/pkg/util/password"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type service struct {
	store       *store
	cfg         *config.Config
	log         *zap.Logger
	db          db.DB
	redisClient *redis.Client
}

func NewService(db db.DB, cfg *config.Config) *service {
	return &service{
		store:       NewStore(db),
		cfg:         cfg,
		db:          db,
		redisClient: cfg.RedisClient,
		log:         zap.L().Named("user.service"),
	}
}

func (s *service) CreateUser(ctx context.Context, cmd *user.CreateUserCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		result, err := s.store.userTaken(ctx, 0, cmd.Email)
		if err != nil {
			return err
		}

		if len(result) > 0 {
			return user.ErrUserAlreadyExists
		}

		passwordHash, err := util.HashPassword(cmd.Password)
		if err != nil {
			return err
		}

		cmd.Password = passwordHash

		err = s.store.create(ctx, cmd)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) GetByUserID(ctx context.Context, id int64) (*user.User, error) {
	result, err := s.store.getUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, user.ErrUserNotFound
	}

	return result, nil
}

func (s *service) UpdateUser(ctx context.Context, cmd *user.UpdateUserCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		result, err := s.store.userTaken(ctx, cmd.ID, cmd.Email)
		if err != nil {
			return err
		}

		if len(result) == 0 {
			return user.ErrUserNotFound
		}
		if len(result) > 1 || (len(result) == 1 && result[0].ID != cmd.ID) {
			return user.ErrUserAlreadyExists
		}

		err = s.store.update(ctx, cmd)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) SearchUser(ctx context.Context, query *user.SearchUserQuery) (*user.SearchUserResult, error) {
	if query.Page <= 0 {
		query.Page = s.cfg.Pagination.Page
	}

	if query.PerPage <= 0 {
		query.PerPage = s.cfg.Pagination.PageLimit
	}

	result, err := s.store.search(ctx, query)
	if err != nil {
		return nil, err
	}

	result.PerPage = query.PerPage
	result.Page = query.Page

	return result, nil
}

func (s *service) DeleteUser(ctx context.Context, id int64) error {
	result, err := s.store.getUserByID(ctx, id)
	if err != nil {
		return err
	}

	if result == nil {
		return user.ErrUserNotFound
	}

	err = s.store.delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) RegisterDefaultUser(ctx context.Context, cmd *user.RegisterUserCommand) error {
	role := "user"

	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		result, err := s.store.userTaken(ctx, 0, cmd.Email)
		if err != nil {
			return err
		}

		if len(result) > 0 {
			return user.ErrUserAlreadyExists
		}

		passwordHash, err := util.HashPassword(cmd.Password)
		if err != nil {
			return err
		}

		cmd.Password = passwordHash

		err = s.store.registerDefaultUser(ctx, cmd, role)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) GetUserByEmail(ctx context.Context, cmd *user.LoginUserCommand) (string, error) {
	result, err := s.store.getUserByEmail(ctx, cmd.Email)
	if err != nil {
		return "", err
	}

	if result == nil {
		return "", user.ErrUserNotFound
	}

	err = util.CheckPasswordHash(result.PasswordHash, cmd.Password)
	if err != nil {
		return "", user.ErrInvalidPassword
	}

	token, err := jwt.GenerateToken(result.Email, result.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

// InvalidateToken adds the token to a blacklist using Redis
func (s *service) InvalidateToken(ctx context.Context, token string) error {
	expiration := 24 * time.Hour
	err := s.redisClient.Set(ctx, token, "blacklisted", expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

// IsTokenBlacklisted checks if the token is blacklisted in Redis
func (s *service) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	_, err := s.redisClient.Get(ctx, token).Result()
	if err == redis.Nil {
		// Token is no found in blacklist
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
