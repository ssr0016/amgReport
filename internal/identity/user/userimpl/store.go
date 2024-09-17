package userimpl

import (
	"amg/internal/db"
	"amg/internal/identity/user"
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type store struct {
	db     db.DB
	logger *zap.Logger
}

func NewStore(db db.DB) *store {
	return &store{
		db:     db,
		logger: zap.L().Named("user.store"),
	}
}

func (s *store) create(ctx context.Context, cmd *user.CreateUserCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		rawSQL := `
		INSERT INTO users (
			first_name,
			last_name,
			email,
			password_hash,
			address,
			phone_number,
			date_of_birth,
			role
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8 
		) RETURNING id
	`

		var id int64

		err := tx.QueryRow(
			ctx,
			rawSQL,
			cmd.FirstName,
			cmd.LastName,
			cmd.Email,
			cmd.Password,
			&cmd.Address,
			&cmd.PhoneNumber,
			&cmd.DateOfBirth,
			cmd.Role,
		).Scan(&id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *store) userTaken(ctx context.Context, id int64, email string) ([]*user.User, error) {
	var result []*user.User

	rawSQL := `
	SELECT
		id,
		first_name,
		last_name,
		email,
		password_hash,
		address,
		phone_number,
		date_of_birth,
		role
	FROM
		users
	WHERE
		id = $1 OR
		email = $2
	`

	err := s.db.Select(ctx, &result, rawSQL, id, email)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *store) getUserByID(ctx context.Context, id int64) (*user.User, error) {
	var result user.User

	rawSQL := `
	SELECT
		id,
		first_name,
		last_name,
		email,
		password_hash,
		address,
		phone_number,
		date_of_birth,
		role,
		created_at,
		updated_at
	FROM
		users
	WHERE
		id = $1
	`

	err := s.db.Get(ctx, &result, rawSQL, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, nil
		}
	}

	return &result, nil
}

func (s *store) update(ctx context.Context, cmd *user.UpdateUserCommand) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		rawSQL := `
		UPDATE
			users
		SET
			first_name = $1,
			last_name = $2,
			email = $3,
			address = $4,
			phone_number = $5,
			date_of_birth = $6,
			role = $7
		WHERE
			id = $8
		`

		_, err := tx.Exec(
			ctx,
			rawSQL,
			cmd.FirstName,
			cmd.LastName,
			cmd.Email,
			cmd.Address,
			cmd.PhoneNumber,
			cmd.DateOfBirth,
			cmd.Role,
			cmd.ID,
		)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *store) search(ctx context.Context, query *user.SearchUserQuery) (*user.SearchUserResult, error) {
	var (
		result = user.SearchUserResult{
			Users: make([]*user.User, 0),
		}
		sql            bytes.Buffer
		whereCondition = make([]string, 0)
		whereParams    = make([]interface{}, 0)
		paramIndex     = 1
	)

	sql.WriteString(`
	SELECT
		id,
		first_name,
		last_name,
		email,
		address,
		phone_number,
		date_of_birth,
		role,
		created_at,
		updated_at
	FROM
		users
	
	`)

	if len(query.FirstName) > 0 {
		whereCondition = append(whereCondition, fmt.Sprintf("first_name ILIKE $%d", paramIndex))
		whereParams = append(whereParams, "%"+query.FirstName+"%")
		paramIndex++
	}

	if len(query.LastName) > 0 {
		whereCondition = append(whereCondition, fmt.Sprintf("last_name ILIKE $%d", paramIndex))
		whereParams = append(whereParams, "%"+query.LastName+"%")
		paramIndex++
	}

	if len(query.Email) > 0 {
		whereCondition = append(whereCondition, fmt.Sprintf("email ILIKE $%d", paramIndex))
		whereParams = append(whereParams, "%"+query.Email+"%")
		paramIndex++
	}

	if len(query.Address) > 0 {
		whereCondition = append(whereCondition, fmt.Sprintf("address ILIKE $%d", paramIndex))
		whereParams = append(whereParams, "%"+query.Address+"%")
		paramIndex++
	}

	if len(query.PhoneNumber) > 0 {
		whereCondition = append(whereCondition, fmt.Sprintf("phone_number ILIKE $%d", paramIndex))
		whereParams = append(whereParams, "%"+query.PhoneNumber+"%")
		paramIndex++
	}

	if len(query.DateOfBirth) > 0 {
		whereCondition = append(whereCondition, fmt.Sprintf("date_of_birth = $%d", paramIndex))
		whereParams = append(whereParams, query.DateOfBirth)
		paramIndex++
	}

	if len(whereCondition) > 0 {
		sql.WriteString(" WHERE " + strings.Join(whereCondition, " AND "))
	}

	sql.WriteString(" ORDER BY created_at DESC")

	if query.PerPage > 0 {
		offset := query.PerPage * (query.Page - 1)
		sql.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1))
		whereParams = append(whereParams, query.PerPage, offset)
	}

	count, err := s.getCount(ctx, sql, whereParams)
	if err != nil {
		return nil, err
	}

	err = s.db.Select(ctx, &result.Users, sql.String(), whereParams...)
	if err != nil {
		return nil, err
	}

	result.TotalCount = count

	return &result, nil

}

func (s *store) getCount(ctx context.Context, sql bytes.Buffer, whereParams []interface{}) (int64, error) {
	var count int64

	rawSQL := "SELECT COUNT(*) FROM (" + sql.String() + ") as t1"

	err := s.db.Get(ctx, &count, rawSQL, whereParams...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *store) delete(ctx context.Context, id int64) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context, tx db.Tx) error {
		rawSQL := `
			DELETE 
			FROM
				users
			WHERE
				id = $1		
		`

		_, err := tx.Exec(ctx, rawSQL, id)
		if err != nil {
			return err
		}

		return nil
	})
}
