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
) (models.UserResponse, error) {

	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return models.UserResponse{}, customErrors.ErrInvalidDate
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
		return models.UserResponse{}, err
	}

	s.logger.Info(
		"user created",
		zap.Int32("id", user.ID),
		zap.String("name", user.Name),
	)

	return mapToUserResponse(user), nil
}

func (s *UserService) ListUsers(
	ctx context.Context,
) ([]models.UserResponse, error) {

	users, err := s.repo.ListUsers(ctx)
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
func (s *UserService) UpdateUser(
	ctx context.Context,
	id int32,
	req models.UpdateUserRequest,
) (models.UserResponse, error) {

	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return models.UserResponse{}, customErrors.ErrInvalidDate
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
			return models.UserResponse{}, customErrors.ErrUserNotFound
		}
		s.logger.Error(
			"failed to update user",
			zap.Error(err),
		)
		return models.UserResponse{}, err
	}

	s.logger.Info(
		"user updated",
		zap.Int32("id", user.ID),
		zap.String("name", user.Name),
	)
	return mapToUserResponse(user), nil
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
func mapToUserResponse(user db.User) models.UserResponse {
	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format("2006-01-02"),
		Age:  CalculateAge(user.Dob.Time),
	}
}
func CalculateAge(dob time.Time) int {
	today := time.Now()

	age := today.Year() - dob.Year()

	if today.Month() < dob.Month() ||
		(today.Month() == dob.Month() && today.Day() < dob.Day()) {
		age--
	}

	return age
}
