package service

import (
	db "UserAgeAPI/db/sqlc/generated"
	customErrors "UserAgeAPI/internal/errors"
	"UserAgeAPI/internal/models"
	"UserAgeAPI/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type UserService struct {
	repo   *repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo *repository.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

func (s *UserService) GetUser(
	ctx context.Context,
	id int32,
) (models.UserResponse, error) {

	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserResponse{}, customErrors.ErrUserNotFound
		}
		s.logger.Error(
			"failed to get user",
			zap.Error(err),
		)
		return models.UserResponse{}, err
	}

	return mapToUserResponse(user), nil
}

func (s *UserService) CreateUser(
	ctx context.Context,
	req models.CreateUserRequest,
) (models.UserWriteResponse, error) {

	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return models.UserWriteResponse{}, customErrors.ErrInvalidDate
	}

	user, err := s.repo.CreateUser(ctx, db.CreateUserParams{
		Name: req.Name,
		Dob: pgtype.Date{
			Time:  dob,
			Valid: true,
		},
	})

	if err != nil {
		s.logger.Error(
			"failed to create user",
			zap.Error(err),
		)
		return models.UserWriteResponse{}, err
	}

	s.logger.Info(
		"user created",
		zap.Int32("id", user.ID),
		zap.String("name", user.Name),
	)

	return mapToUserWriteResponse(user), nil
}

func (s *UserService) UpdateUser(
	ctx context.Context,
	id int32,
	req models.UpdateUserRequest,
) (models.UserWriteResponse, error) {

	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return models.UserWriteResponse{}, customErrors.ErrInvalidDate
	}

	user, err := s.repo.UpdateUser(ctx, db.UpdateUserParams{
		ID:   id,
		Name: req.Name,
		Dob: pgtype.Date{
			Time:  dob,
			Valid: true,
		},
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserWriteResponse{}, customErrors.ErrUserNotFound
		}
		s.logger.Error(
			"failed to update user",
			zap.Error(err),
		)
		return models.UserWriteResponse{}, err
	}

	s.logger.Info(
		"user updated",
		zap.Int32("id", user.ID),
		zap.String("name", user.Name),
	)
	return mapToUserWriteResponse(user), nil
}

func (s *UserService) DeleteUser(
	ctx context.Context,
	id int32,
) error {

	_, err := s.repo.DeleteUser(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return customErrors.ErrUserNotFound
		}

		s.logger.Error(
			"failed to delete user",
			zap.Error(err),
		)
		return err
	}
	s.logger.Info(
		"user deleted",
		zap.Int32("id", id),
	)

	return nil
}
func (s *UserService) ListUsers(
	ctx context.Context,
	page int32,
	limit int32,
) ([]models.UserResponse, error) {

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	users, err := s.repo.ListUsersPaginated(
		ctx,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"failed to list users",
			zap.Error(err),
		)
		return nil, err
	}

	var response []models.UserResponse

	for _, user := range users {
		response = append(response, mapToUserResponse(user))
	}

	return response, nil
}
func mapToUserResponse(user db.User) models.UserResponse {
	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format("2006-01-02"),
		Age:  CalculateAge(user.Dob.Time, time.Now()),
	}
}
func mapToUserWriteResponse(user db.User) models.UserWriteResponse {
	return models.UserWriteResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format("2006-01-02"),
	}
}
func CalculateAge(dob time.Time, today time.Time) int {
	age := today.Year() - dob.Year()

	if today.Month() < dob.Month() ||
		(today.Month() == dob.Month() && today.Day() < dob.Day()) {
		age--
	}

	return age
}
