package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"backend/controllers"
	"backend/dao"
	domain "backend/domain/models"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var testRouter *gin.Engine

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("no se pudo abrir la BD de test: " + err.Error())
	}
	db.AutoMigrate(&domain.User{}, &domain.Event{}, &domain.Ticket{})
	dao.DB = db
	os.Setenv("JWT_SECRET", "test-secret-controllers")

	testRouter = gin.New()
	api := testRouter.Group("/api")
	api.POST("/auth/register", controllers.RegisterUser)
	api.POST("/auth/login", controllers.LoginUser)
	api.GET("/events", controllers.GetEvents)
	api.GET("/events/:id", controllers.GetEventByID)

	protected := api.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	protected.POST("/tickets/purchase", controllers.PurchaseTicket)
	protected.GET("/tickets/my-tickets", controllers.GetMyTickets)
	protected.POST("/tickets/:id/cancel", controllers.CancelTicket)
	protected.POST("/tickets/:id/transfer", controllers.TransferTicket)

	os.Exit(m.Run())
}

func tokenParaUsuario(userID uint) string {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"user_id": float64(userID),
		"email":   "test@test.com",
		"rol":     "cliente",
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func crearUsuario(email, dni string) domain.User {
	u := domain.User{
		Email: email, Password: "hash",
		Nombre: "Test", Apellido: "User", Rol: "cliente", DNI: dni,
	}
	dao.DB.Create(&u)
	return u
}

func crearEvento(cupo int) domain.Event {
	e := domain.Event{
		Titulo: "Evento Test", Categoria: "Ctrl",
		Lugar: "BsAs", Precio: 100, CupoMaximo: cupo, CupoDispon: cupo,
	}
	dao.DB.Create(&e)
	return e
}

// ── Protección de endpoints (todos los del rol cliente requieren JWT) ──

func TestEndpointsProtegidos_SinToken_Retorna401(t *testing.T) {
	endpoints := []struct{ method, path, body string }{
		{"POST", "/api/tickets/purchase", `{"event_id": 1}`},
		{"GET", "/api/tickets/my-tickets", ""},
		{"POST", "/api/tickets/1/cancel", ""},
		{"POST", "/api/tickets/1/transfer", `{"dni": "12345678"}`},
	}
	for _, ep := range endpoints {
		req, _ := http.NewRequest(ep.method, ep.path, strings.NewReader(ep.body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("%s %s sin token: esperado 401, obtenido %d", ep.method, ep.path, w.Code)
		}
	}
}

// ── Auth — Register ───────────────────────────────────────────────

func TestRegister_Exitoso(t *testing.T) {
	body := `{"email":"nuevo@ctrl.test","password":"pass123","nombre":"N","apellido":"A","dni":"11122233"}`
	req, _ := http.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("esperado 201, obtenido %d — body: %s", w.Code, w.Body.String())
	}
}

func TestRegister_InputInvalido(t *testing.T) {
	body := `{"email":"no-es-email","password":"123"}` // email inválido, sin nombre/apellido/dni
	req, _ := http.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestRegister_EmailDuplicado(t *testing.T) {
	body := `{"email":"dup@ctrl.test","password":"pass123","nombre":"D","apellido":"U","dni":"22233344"}`
	req, _ := http.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(httptest.NewRecorder(), req)

	req2, _ := http.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req2)
	if w.Code != http.StatusConflict {
		t.Errorf("esperado 409 por email duplicado, obtenido %d", w.Code)
	}
}

// ── Auth — Login ──────────────────────────────────────────────────

func TestLogin_Exitoso(t *testing.T) {
	reg := `{"email":"login@ctrl.test","password":"pass123","nombre":"L","apellido":"U","dni":"33344455"}`
	req, _ := http.NewRequest("POST", "/api/auth/register", strings.NewReader(reg))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(httptest.NewRecorder(), req)

	body := `{"email":"login@ctrl.test","password":"pass123"}`
	req2, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req2)
	if w.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtenido %d — body: %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] == "" {
		t.Error("token ausente en la respuesta del login")
	}
}

func TestLogin_CredencialesInvalidas(t *testing.T) {
	body := `{"email":"noexiste@ctrl.test","password":"pass123"}`
	req, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401, obtenido %d", w.Code)
	}
}

// ── Eventos (públicos) ────────────────────────────────────────────

func TestGetEvents_Publico(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/events", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d", w.Code)
	}
}

func TestGetEventByID_Inexistente(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/events/999999", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtenido %d", w.Code)
	}
}

// ── Tickets (protegidos) ──────────────────────────────────────────

func TestPurchaseTicket_InputInvalido(t *testing.T) {
	u := crearUsuario("purchase@ctrl.test", "44455566")
	token := tokenParaUsuario(u.ID)

	body := `{}` // falta event_id
	req, _ := http.NewRequest("POST", "/api/tickets/purchase", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400 por input inválido, obtenido %d", w.Code)
	}
}

func TestPurchaseTicket_SinCupo(t *testing.T) {
	u := crearUsuario("sincupo@ctrl.test", "55566677")
	e := crearEvento(0)
	token := tokenParaUsuario(u.ID)

	body := fmt.Sprintf(`{"event_id":%d}`, e.ID)
	req, _ := http.NewRequest("POST", "/api/tickets/purchase", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("esperado 409 sin cupo, obtenido %d", w.Code)
	}
}

func TestGetMyTickets_Autenticado(t *testing.T) {
	u := crearUsuario("mytickets@ctrl.test", "66677788")
	token := tokenParaUsuario(u.ID)

	req, _ := http.NewRequest("GET", "/api/tickets/my-tickets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d", w.Code)
	}
}

func TestCancelTicket_IDInvalido(t *testing.T) {
	u := crearUsuario("cancelid@ctrl.test", "77788899")
	token := tokenParaUsuario(u.ID)

	req, _ := http.NewRequest("POST", "/api/tickets/abc/cancel", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400 por ID inválido, obtenido %d", w.Code)
	}
}

func TestCancelTicket_NoAutorizado(t *testing.T) {
	e := crearEvento(10)
	owner := crearUsuario("ownercancel@ctrl.test", "88899900")
	otro := crearUsuario("otro@ctrl.test", "99900011")

	ticket := &domain.Ticket{UserID: owner.ID, EventID: e.ID, Estado: "activo"}
	dao.DB.Create(ticket)

	token := tokenParaUsuario(otro.ID)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/tickets/%d/cancel", ticket.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("esperado 403 no autorizado, obtenido %d", w.Code)
	}
}

func TestTransferTicket_InputInvalido(t *testing.T) {
	u := crearUsuario("transferinput@ctrl.test", "10111213")
	token := tokenParaUsuario(u.ID)

	body := `{}` // falta dni
	req, _ := http.NewRequest("POST", "/api/tickets/1/transfer", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400 por input inválido, obtenido %d", w.Code)
	}
}
