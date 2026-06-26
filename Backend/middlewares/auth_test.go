package middlewares

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func makeToken(t *testing.T, userID float64, exp time.Duration) string {
	t.Helper()
	secret := "test-secret"
	os.Setenv("JWT_SECRET", secret)

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   "test@test.com",
		"rol":     "cliente",
		"exp":     time.Now().Add(exp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("no se pudo crear el token de prueba: %v", err)
	}
	return signed
}

func runMiddleware(token string) int {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", AuthMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Code
}

func TestAuthMiddleware_TokenValido(t *testing.T) {
	token := makeToken(t, 1, 24*time.Hour)
	if got := runMiddleware(token); got != http.StatusOK {
		t.Errorf("esperaba 200, obtuve %d", got)
	}
}

func TestAuthMiddleware_SinToken(t *testing.T) {
	if got := runMiddleware(""); got != http.StatusUnauthorized {
		t.Errorf("esperaba 401, obtuve %d", got)
	}
}

func TestAuthMiddleware_TokenExpirado(t *testing.T) {
	token := makeToken(t, 1, -1*time.Hour)
	if got := runMiddleware(token); got != http.StatusUnauthorized {
		t.Errorf("esperaba 401, obtuve %d", got)
	}
}

func TestAuthMiddleware_TokenFalsificado(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	claims := jwt.MapClaims{"user_id": float64(1), "exp": time.Now().Add(time.Hour).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte("clave-distinta"))

	if got := runMiddleware(signed); got != http.StatusUnauthorized {
		t.Errorf("esperaba 401, obtuve %d", got)
	}
}

// ── AdminMiddleware ───────────────────────────────────────────────

func runAdminMiddleware(rol string, setRol bool) int {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/admin", func(c *gin.Context) {
		if setRol {
			c.Set("rol", rol)
		}
		c.Next()
	}, AdminMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Code
}

func TestAdminMiddleware_RolAdmin(t *testing.T) {
	if got := runAdminMiddleware("admin", true); got != http.StatusOK {
		t.Errorf("esperaba 200 para rol admin, obtuve %d", got)
	}
}

func TestAdminMiddleware_RolAdministrador(t *testing.T) {
	if got := runAdminMiddleware("Administrador", true); got != http.StatusOK {
		t.Errorf("esperaba 200 para rol Administrador, obtuve %d", got)
	}
}

func TestAdminMiddleware_RolCliente(t *testing.T) {
	if got := runAdminMiddleware("Cliente", true); got != http.StatusForbidden {
		t.Errorf("esperaba 403 para rol cliente, obtuve %d", got)
	}
}

func TestAdminMiddleware_SinRol(t *testing.T) {
	if got := runAdminMiddleware("", false); got != http.StatusForbidden {
		t.Errorf("esperaba 403 sin rol seteado, obtuve %d", got)
	}
}

func TestAuthMiddleware_FormatoInvalido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", AuthMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token sinprefijobearer")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, obtuve %d", w.Code)
	}
}

func TestAuthMiddleware_AlgoritmoNoHMAC(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	claims := jwt.MapClaims{
		"user_id": float64(1),
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	if got := runMiddleware(signed); got != http.StatusUnauthorized {
		t.Errorf("esperaba 401 con algoritmo RSA, obtuve %d", got)
	}
}

func TestAuthMiddleware_UserIDInvalido(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	claims := jwt.MapClaims{
		"user_id": "no-soy-un-numero",
		"email":   "test@test.com",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte("test-secret"))
	if got := runMiddleware(signed); got != http.StatusUnauthorized {
		t.Errorf("esperaba 401 con user_id inválido, obtuve %d", got)
	}
}

func TestAuthMiddleware_SecretDefault(t *testing.T) {
	os.Unsetenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", "test-secret")
	claims := jwt.MapClaims{
		"user_id": float64(1),
		"email":   "test@test.com",
		"rol":     "cliente",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte("secret-cambiar-en-produccion"))
	if got := runMiddleware(signed); got != http.StatusOK {
		t.Errorf("esperaba 200 con secret por defecto, obtuve %d", got)
	}
}
