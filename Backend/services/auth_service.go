package services

import (
	"errors"
	"os"
	"time"

	"backend/dao"
	domain "backend/domain/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nombre   string `json:"nombre" binding:"required"`
	Apellido string `json:"apellido" binding:"required"`
	DNI      string `json:"dni" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(input RegisterInput) (*domain.User, error) {
	var existing domain.User
	if err := dao.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		return nil, errors.New("el email ya está registrado")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		Email:    input.Email,
		Password: string(hashed),
		Nombre:   input.Nombre,
		Apellido: input.Apellido,
		Rol:      "cliente",
		DNI:      input.DNI,
	}

	if err := dao.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func Login(input LoginInput) (string, *domain.User, error) {
	var user domain.User
	if err := dao.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("credenciales inválidas")
		}
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", nil, errors.New("credenciales inválidas")
	}

	token, err := generateJWT(user)
	if err != nil {
		return "", nil, err
	}
	return token, &user, nil
}

func generateJWT(user domain.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret-cambiar-en-produccion"
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"rol":     user.Rol,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
