package user

import (
	"log"
	"net/http"

	"github.com/reyimanuel/template/internal/infrastructures/pkg/errs"
	"github.com/reyimanuel/template/internal/infrastructures/pkg/helpers"
	"github.com/reyimanuel/template/internal/infrastructures/pkg/token"
	"github.com/reyimanuel/template/internal/migration"
)

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) Login(payload *LoginRequest) (*Response, error) {
	user, err := s.Repo.GetByEmail(payload.Email)
	if err != nil {
		return nil, errs.Unauthorized("Email atau Password Salah")
	}

	if !helpers.CheckPasswordHash(payload.Password, user.Password) {
		return nil, errs.Unauthorized("Email atau Password Salah")
	}

	access, err := token.GenerateToken(user.ID, user.Email)
	if err != nil {
		log.Printf("error saat membuat token: %v", err)
		return nil, errs.InternalServerError("Gagal membuat akses token")
	}

	refresh, err := token.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("error saat membuat refresh token: %v", err)
		return nil, errs.InternalServerError("Gagal membuat refresh token")
	}

	return &Response{
		StatusCode: http.StatusOK,
		Message:    "Login Berhasil",
		Data: TokenResponse{
			AccessToken:  access,
			RefreshToken: refresh,
		},
	}, nil
}

func (s *Service) Register(payload *RegisterRequest) (*Response, error) {
	// Check if email already exists
	existingUser, err := s.Repo.GetByEmail(payload.Email)
	if err == nil && existingUser != nil {
		return nil, errs.BadRequest("Email sudah terdaftar")
	}

	// Create new user
	hashedPassword, err := helpers.HashPassword(payload.Password)
	if err != nil {
		log.Printf("error saat hash password: %v", err)
		return nil, errs.InternalServerError("Gagal membuat akun")
	}

	user := &migration.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashedPassword,
	}

	if err := s.Repo.CreateAccount(user); err != nil {
		log.Printf("error saat membuat user: %v", err)
		return nil, errs.InternalServerError("Gagal membuat akun")
	}
	return &Response{
		StatusCode: http.StatusCreated,
		Message:    "Akun berhasil dibuat",
		Data:       nil,
	}, nil
}
