package repository

import (
	db "UserAgeAPI/db/sqlc/generated"
	"context"
)

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) *UserRepository {
	return &UserRepository{
		queries: queries,
	}
}

func (r *UserRepository) CreateUser(
	ctx context.Context,
	params db.CreateUserParams,
) (db.User, error) {

	return r.queries.CreateUser(ctx, params)
}

func (r *UserRepository) GetUser(
	ctx context.Context,
	id int32,
) (db.User, error) {
	return r.queries.GetUser(ctx, id)

}

func (r *UserRepository) ListUsers(
	ctx context.Context,
) ([]db.User, error) {

	return r.queries.ListUsers(ctx)
}

func (r *UserRepository) UpdateUser(
	ctx context.Context,
	params db.UpdateUserParams,
) (db.User, error) {

	return r.queries.UpdateUser(ctx, params)
}

func (r *UserRepository) DeleteUser(
	ctx context.Context,
	id int32,
) (int32, error) {
	return r.queries.DeleteUser(ctx, id)
}
func (r *UserRepository) UserExists(
	ctx context.Context,
	id int32,
) (bool, error) {
	return r.queries.UserExists(ctx, id)
}
