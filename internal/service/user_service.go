package service

import (
	db "UserAgeAPI/db/sqlc/generated"
	"UserAgeAPI/internal/models"
	"UserAgeAPI/internal/repository"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetUser(
	ctx context.Context,
	id int32,
) (models.UserResponse, error) {

	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return models.UserResponse{}, err
	}

	response := models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format("2006-01-02"),
		Age:  CalculateAge(user.Dob.Time),
	}

	return response, nil
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

func (s *UserService) CreateUser(
	ctx context.Context,
	req models.CreateUserRequest,
) (models.UserResponse, error) {

	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.CreateUser(ctx, db.CreateUserParams{
		Name: req.Name,
		Dob: pgtype.Date{
			Time:  dob,
			Valid: true,
		},
	})

	if err != nil {
		return models.UserResponse{}, err
	}

	response := models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format("2006-01-02"),
		Age:  CalculateAge(user.Dob.Time),
	}

	return response, nil
}

func (s *UserService) ListUsers(
	ctx context.Context,
) ([]models.UserResponse, error) {

	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	var response []models.UserResponse

	for _, user := range users {
		response = append(response, models.UserResponse{
			ID:   user.ID,
			Name: user.Name,
			Dob:  user.Dob.Time.Format("2006-01-02"),
			Age:  CalculateAge(user.Dob.Time),
		})
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
		return models.UserResponse{}, err
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
		return models.UserResponse{}, err
	}

	response := models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format("2006-01-02"),
		Age:  CalculateAge(user.Dob.Time),
	}

	return response, nil
}

func (s *UserService) DeleteUser(
	ctx context.Context,
	id int32,
) error {

	return s.repo.DeleteUser(ctx, id)
}
