package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"backend/dao"
	domain "backend/domain/models"

	"github.com/golang-jwt/jwt/v5"
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

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func Register(input RegisterInput) (*domain.User, error) {
	if dao.DB == nil {
		return nil, errors.New("base de datos no inicializada")
	}
	var existing domain.User
	if err := dao.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		return nil, errors.New("el email ya esta registrado")
	}

	user := domain.User{
		Email:    input.Email,
		Password: hashPassword(input.Password),
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
	if dao.DB == nil {
		return "", nil, errors.New("base de datos no inicializada")
	}
	var user domain.User
	if err := dao.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("credenciales invalidas")
		}
		return "", nil, err
	}

	if user.Password != hashPassword(input.Password) {
		return "", nil, errors.New("credenciales invalidas")
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
		"dni":     user.DNI,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
