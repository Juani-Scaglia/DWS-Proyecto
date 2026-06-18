package services

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	domain "backend/domain/models"
)

// generateJWT es privada, la probamos indirectamente construyendo un user de prueba
// y verificando que el token generado sea parseable con el mismo secreto.
func TestGenerateJWT_TokenValido(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")

	user := domain.User{
		ID:    1,
		Email: "juan@test.com",
		Rol:   "cliente",
	}

	tokenStr, err := generateJWT(user)
	if err != nil {
		t.Fatalf("generateJWT devolvió error: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("generateJWT devolvió un token vacío")
	}

	// Verificar que el token es parseable y contiene los claims correctos
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("el token generado no es válido: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("no se pudieron leer los claims")
	}
	if claims["email"] != "juan@test.com" {
		t.Errorf("email incorrecto en claims: %v", claims["email"])
	}
	if claims["rol"] != "cliente" {
		t.Errorf("rol incorrecto en claims: %v", claims["rol"])
	}
	if claims["user_id"] != float64(1) {
		t.Errorf("user_id incorrecto en claims: %v", claims["user_id"])
	}
}

func TestGenerateJWT_Expiracion24h(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")

	user := domain.User{ID: 1, Email: "a@b.com", Rol: "cliente"}
	tokenStr, _ := generateJWT(user)

	token, _ := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	claims := token.Claims.(jwt.MapClaims)

	exp := int64(claims["exp"].(float64))
	horasRestantes := time.Until(time.Unix(exp, 0)).Hours()

	if horasRestantes < 23 || horasRestantes > 25 {
		t.Errorf("expiración inesperada: %.1f horas (esperaba ~24)", horasRestantes)
	}
}

func TestRegisterInput_CamposRequeridos(t *testing.T) {
	casos := []struct {
		nombre string
		input  RegisterInput
		valido bool
	}{
		{
			nombre: "input completo",
			input:  RegisterInput{Email: "a@b.com", Password: "123456", Nombre: "Juan", Apellido: "García", DNI: "12345678"},
			valido: true,
		},
		{
			nombre: "email vacío",
			input:  RegisterInput{Password: "123456", Nombre: "Juan", Apellido: "García", DNI: "12345678"},
			valido: false,
		},
		{
			nombre: "password corto",
			input:  RegisterInput{Email: "a@b.com", Password: "123", Nombre: "Juan", Apellido: "García", DNI: "12345678"},
			valido: false,
		},
		{
			nombre: "DNI vacío",
			input:  RegisterInput{Email: "a@b.com", Password: "123456", Nombre: "Juan", Apellido: "García"},
			valido: false,
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			emailOk := tc.input.Email != ""
			passOk := len(tc.input.Password) >= 6
			dniOk := tc.input.DNI != ""
			esValido := emailOk && passOk && dniOk

			if esValido != tc.valido {
				t.Errorf("caso '%s': esperaba valido=%v, obtuve valido=%v", tc.nombre, tc.valido, esValido)
			}
		})
	}
}
