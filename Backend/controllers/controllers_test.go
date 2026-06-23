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
	// Migramos también la entidad Seat requerida en las transacciones de compra
	db.AutoMigrate(&domain.User{}, &domain.Event{}, &domain.Ticket{}, &domain.Seat{}, &domain.Venue{})
	dao.DB = db
	os.Setenv("JWT_SECRET", "test-secret-controllers")

	// Creamos un Venue por defecto ID=1 para que no fallen las FK/relaciones obligatorias
	dao.DB.Create(&domain.Venue{Nombre: "Estadio Test", Direccion: "Calle Falsa 123", Filas: 5, ColumnasPorFila: 10, Capacidad: 50})

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

	admin := api.Group("/admin")
	admin.Use(middlewares.AuthMiddleware())
	admin.Use(middlewares.AdminMiddleware())
	admin.POST("/events", controllers.CreateEventAdmin)
	admin.PUT("/events/:id", controllers.UpdateEventAdmin)
	admin.DELETE("/events/:id", controllers.DeleteEventAdmin)
	admin.POST("/venues", controllers.CreateVenueAdmin)
	admin.PUT("/venues/:id", controllers.UpdateVenueAdmin)
	admin.DELETE("/venues/:id", controllers.DeleteVenueAdmin)
	admin.GET("/events/:id/report", controllers.GetOccupationReportAdmin)

	api.GET("/venues", controllers.GetVenues)
	api.GET("/venues/:id", controllers.GetVenueByID)
	api.GET("/events/:id/seats", controllers.GetEventSeats)

	os.Exit(m.Run())
}

func tokenParaUsuario(userID uint) string {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"user_id": float64(userID),
		"email":   "test@test.com",
		"rol":     "Cliente", // Rol Capitalizado según Middlewares estándar
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func crearUsuario(email, dni string) domain.User {
	u := domain.User{
		Email: email, Password: "hash",
		Nombre: "Test", Apellido: "User", Rol: "Cliente", DNI: dni,
	}
	dao.DB.Create(&u)
	return u
}

func tokenParaAdmin() string {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"user_id": float64(9999),
		"email":   "admin@test.com",
		"rol":     "admin",
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func crearEvento(cupo int) domain.Event {
	e := domain.Event{
		Titulo: "Evento Test", Categoria: "Ctrl",
		Fecha: time.Now().Add(48 * time.Hour), CupoDispon: cupo,
		VenueID: 1, // Vinculado al establecimiento reglamentario
	}
	dao.DB.Create(&e)
	return e
}

// ── Protección de endpoints ──

func TestEndpointsProtegidos_SinToken_Retorna401(t *testing.T) {
	endpoints := []struct{ method, path, body string }{
		{"POST", "/api/tickets/purchase", `{"event_id": 1, "seat_ids": [1]}`},
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

// ── Auth — Register ──

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
	body := `{"email":"no-es-email","password":"123"}`
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

// ── Auth — Login ──

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
	
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	token, _ := resp["token"].(string)
	if token == "" {
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

// ── Eventos (públicos) ──

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

// ── Tickets (protegidos) ──

func TestPurchaseTicket_InputInvalido(t *testing.T) {
	u := crearUsuario("purchase@ctrl.test", "44455566")
	token := tokenParaUsuario(u.ID)

	body := `{}` // falta seat_ids y event_id
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
	e := crearEvento(0) // Evento con aforo cero
	token := tokenParaUsuario(u.ID)

	body := fmt.Sprintf(`{"event_id":%d, "seat_ids":[999]}`, e.ID)
	req, _ := http.NewRequest("POST", "/api/tickets/purchase", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest && w.Code != http.StatusConflict {
		t.Errorf("esperado error 400 o 409 por falta de asientos válidos, obtenido %d", w.Code)
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

// ── Admin endpoints ──

func TestCreateEventAdmin_SinRolAdmin(t *testing.T) {
	u := crearUsuario("clienteadmin@ctrl.test", "20212223")
	token := tokenParaUsuario(u.ID)
	body := `{"titulo":"T","categoria":"C","venue_id":1,"cupo_maximo":10,"fecha":"2027-01-01T20:00:00Z"}`
	req, _ := http.NewRequest("POST", "/api/admin/events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("esperado 403, obtenido %d", w.Code)
	}
}

func TestCreateEventAdmin_InputInvalido(t *testing.T) {
	token := tokenParaAdmin()
	body := `{}`
	req, _ := http.NewRequest("POST", "/api/admin/events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestCreateEventAdmin_Exitoso(t *testing.T) {
	token := tokenParaAdmin()
	body := `{"titulo":"Evento Admin","categoria":"Rock","venue_id":1,"precio":100,"fecha":"2027-06-01T20:00:00Z"}`
	req, _ := http.NewRequest("POST", "/api/admin/events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("esperado 201, obtenido %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateEventAdmin_Exitoso(t *testing.T) {
	token := tokenParaAdmin()
	createBody := `{"titulo":"Para Editar","categoria":"Jazz","venue_id":1,"precio":100,"fecha":"2027-07-01T20:00:00Z"}`
	req, _ := http.NewRequest("POST", "/api/admin/events", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var event map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &event)
	id := fmt.Sprintf("%v", int(event["id"].(float64)))

	updateBody := `{"titulo":"Editado","categoria":"Jazz","venue_id":1,"precio":100,"fecha":"2027-07-01T20:00:00Z"}`
	req2, _ := http.NewRequest("PUT", "/api/admin/events/"+id, strings.NewReader(updateBody))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	testRouter.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d — body: %s", w2.Code, w2.Body.String())
	}
}

func TestDeleteEventAdmin_Inexistente(t *testing.T) {
	token := tokenParaAdmin()
	req, _ := http.NewRequest("DELETE", "/api/admin/events/999999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtenido %d", w.Code)
	}
}

func TestTransferTicket_InputInvalido(t *testing.T) {
	u := crearUsuario("transferinput@ctrl.test", "10111213")
	token := tokenParaUsuario(u.ID)

	body := `{}`
	req, _ := http.NewRequest("POST", "/api/tickets/1/transfer", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400 por input inválido, obtenido %d", w.Code)
	}
}

// ── Login adicional ──

func TestLogin_PasswordIncorrecto(t *testing.T) {
	reg := `{"email":"wrongpwd@ctrl.test","password":"pass123","nombre":"W","apellido":"P","dni":"98765432"}`
	req, _ := http.NewRequest("POST", "/api/auth/register", strings.NewReader(reg))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(httptest.NewRecorder(), req)

	body := `{"email":"wrongpwd@ctrl.test","password":"contraseniaincorrecta"}`
	req2, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req2)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401 por password incorrecto, obtenido %d", w.Code)
	}
}

// ── Tickets adicionales ──

func TestCancelTicket_TicketNoEncontrado(t *testing.T) {
	u := crearUsuario("cancelnotfound@ctrl.test", "50505050")
	token := tokenParaUsuario(u.ID)
	req, _ := http.NewRequest("POST", "/api/tickets/999999/cancel", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtenido %d", w.Code)
	}
}

func TestTransferTicket_IDInvalido(t *testing.T) {
	u := crearUsuario("transferidinvalid@ctrl.test", "60606060")
	token := tokenParaUsuario(u.ID)
	body := `{"dni":"12345678"}`
	req, _ := http.NewRequest("POST", "/api/tickets/abc/transfer", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestTransferTicket_TicketNoEncontrado(t *testing.T) {
	u := crearUsuario("transfernotfound@ctrl.test", "70707070")
	token := tokenParaUsuario(u.ID)
	body := `{"dni":"99999999"}`
	req, _ := http.NewRequest("POST", "/api/tickets/999999/transfer", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtenido %d", w.Code)
	}
}

func TestTransferTicket_NoAutorizado(t *testing.T) {
	owner := crearUsuario("transferowner2@ctrl.test", "80808080")
	otro := crearUsuario("transferotro2@ctrl.test", "81818181")
	e := crearEvento(10)
	ticket := &domain.Ticket{UserID: owner.ID, EventID: e.ID, Estado: "activo"}
	dao.DB.Create(ticket)

	body := `{"dni":"99999999"}`
	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/tickets/%d/transfer", ticket.ID), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenParaUsuario(otro.ID))
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("esperado 403, obtenido %d", w.Code)
	}
}

func TestTransferTicket_Exitoso(t *testing.T) {
	owner := crearUsuario("transfersuccess1@ctrl.test", "90909090")
	target := crearUsuario("transfersuccess2@ctrl.test", "91919191")
	e := crearEvento(10)
	ticket := &domain.Ticket{UserID: owner.ID, EventID: e.ID, Estado: "activo"}
	dao.DB.Create(ticket)

	body := fmt.Sprintf(`{"dni":"%s"}`, target.DNI)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/tickets/%d/transfer", ticket.ID), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenParaUsuario(owner.ID))
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d — body: %s", w.Code, w.Body.String())
	}
}

// ── UpdateEvent / DeleteEvent adicionales ──

func TestUpdateEventAdmin_IDInvalido(t *testing.T) {
	token := tokenParaAdmin()
	body := `{"titulo":"X","categoria":"X","venue_id":1,"precio":10,"fecha":"2027-07-01T20:00:00Z"}`
	req, _ := http.NewRequest("PUT", "/api/admin/events/abc", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestUpdateEventAdmin_EventoInexistente(t *testing.T) {
	token := tokenParaAdmin()
	body := `{"titulo":"X","categoria":"X","venue_id":1,"precio":10,"fecha":"2027-07-01T20:00:00Z"}`
	req, _ := http.NewRequest("PUT", "/api/admin/events/99999", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtenido %d — body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteEventAdmin_Exitoso(t *testing.T) {
	token := tokenParaAdmin()
	createBody := `{"titulo":"Para Borrar","categoria":"Pop","venue_id":1,"precio":50,"fecha":"2027-08-01T20:00:00Z"}`
	req, _ := http.NewRequest("POST", "/api/admin/events", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var event map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &event)
	id := fmt.Sprintf("%v", int(event["id"].(float64)))

	req2, _ := http.NewRequest("DELETE", "/api/admin/events/"+id, nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	testRouter.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d — body: %s", w2.Code, w2.Body.String())
	}
}

// ── Venues (público) ──

func TestGetVenues_OK(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/venues", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d", w.Code)
	}
}

func TestGetVenueByID_Existente(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/venues/1", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d", w.Code)
	}
}

func TestGetVenueByID_IDInvalido(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/venues/abc", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestGetVenueByID_Inexistente(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/venues/99999", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtenido %d", w.Code)
	}
}

// ── Venues (admin) ──

func crearVenueAPI(t *testing.T) uint {
	t.Helper()
	token := tokenParaAdmin()
	body := `{"nombre":"Venue Temp","direccion":"Calle Test 1","cap_platea_norte":10}`
	req, _ := http.NewRequest("POST", "/api/admin/venues", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	var v map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &v)
	return uint(int(v["id"].(float64)))
}

func TestCreateVenueAdmin_InputInvalido(t *testing.T) {
	token := tokenParaAdmin()
	req, _ := http.NewRequest("POST", "/api/admin/venues", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestCreateVenueAdmin_Exitoso(t *testing.T) {
	token := tokenParaAdmin()
	body := `{"nombre":"Nuevo Venue","direccion":"Av. Test 100","cap_platea_norte":50}`
	req, _ := http.NewRequest("POST", "/api/admin/venues", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("esperado 201, obtenido %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateVenueAdmin_IDInvalido(t *testing.T) {
	token := tokenParaAdmin()
	body := `{"nombre":"X","direccion":"X","cap_platea_norte":1}`
	req, _ := http.NewRequest("PUT", "/api/admin/venues/abc", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestUpdateVenueAdmin_Inexistente(t *testing.T) {
	token := tokenParaAdmin()
	body := `{"nombre":"X","direccion":"X","cap_platea_norte":1}`
	req, _ := http.NewRequest("PUT", "/api/admin/venues/99999", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtenido %d", w.Code)
	}
}

func TestUpdateVenueAdmin_Exitoso(t *testing.T) {
	id := crearVenueAPI(t)
	token := tokenParaAdmin()
	body := `{"nombre":"Venue Actualizado","direccion":"Nueva Dir","cap_platea_norte":24}`
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/admin/venues/%d", id), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d — body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteVenueAdmin_IDInvalido(t *testing.T) {
	token := tokenParaAdmin()
	req, _ := http.NewRequest("DELETE", "/api/admin/venues/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestDeleteVenueAdmin_Inexistente(t *testing.T) {
	token := tokenParaAdmin()
	req, _ := http.NewRequest("DELETE", "/api/admin/venues/99999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtenido %d", w.Code)
	}
}

func TestDeleteVenueAdmin_Exitoso(t *testing.T) {
	id := crearVenueAPI(t)
	token := tokenParaAdmin()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/admin/venues/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d — body: %s", w.Code, w.Body.String())
	}
}

// ── Asientos ──

func TestGetEventSeats_IDInvalido(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/events/abc/seats", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestGetEventSeats_OK(t *testing.T) {
	e := crearEvento(10)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/events/%d/seats", e.ID), nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d", w.Code)
	}
}

// ── Reporte (admin) ──

func TestGetOccupationReport_IDInvalido(t *testing.T) {
	token := tokenParaAdmin()
	req, _ := http.NewRequest("GET", "/api/admin/events/abc/report", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestGetOccupationReport_EventoInexistente(t *testing.T) {
	token := tokenParaAdmin()
	req, _ := http.NewRequest("GET", "/api/admin/events/99999/report", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtenido %d", w.Code)
	}
}

func TestGetOccupationReport_OK(t *testing.T) {
	e := crearEvento(50)
	token := tokenParaAdmin()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/admin/events/%d/report", e.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d — body: %s", w.Code, w.Body.String())
	}
}

// ── Paths adicionales de controllers ──

func TestGetEventByID_IDInvalido_Ctrl(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/events/abc", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestGetEventByID_Existente_Ctrl(t *testing.T) {
	e := crearEvento(10)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/events/%d", e.ID), nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d", w.Code)
	}
}

func TestUpdateEventAdmin_InputInvalido(t *testing.T) {
	token := tokenParaAdmin()
	e := crearEvento(10)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/admin/events/%d", e.ID), strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestDeleteEventAdmin_IDInvalido(t *testing.T) {
	token := tokenParaAdmin()
	req, _ := http.NewRequest("DELETE", "/api/admin/events/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestCreateEventAdmin_VenueInexistente(t *testing.T) {
	token := tokenParaAdmin()
	body := `{"titulo":"T","categoria":"C","venue_id":99999,"precio":10,"fecha":"2027-07-01T20:00:00Z"}`
	req, _ := http.NewRequest("POST", "/api/admin/events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtenido %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateVenueAdmin_InputInvalido(t *testing.T) {
	id := crearVenueAPI(t)
	token := tokenParaAdmin()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/admin/venues/%d", id), strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", w.Code)
	}
}

func TestGetEvents_DBNula(t *testing.T) {
	saved := dao.DB
	dao.DB = nil
	defer func() { dao.DB = saved }()
	req, _ := http.NewRequest("GET", "/api/events", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtenido %d", w.Code)
	}
}

func TestGetVenues_DBNula(t *testing.T) {
	saved := dao.DB
	dao.DB = nil
	defer func() { dao.DB = saved }()
	req, _ := http.NewRequest("GET", "/api/venues", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtenido %d", w.Code)
	}
}

func TestGetMyTickets_DBNula(t *testing.T) {
	u := crearUsuario("dbnulatickets@ctrl.test", "11122200")
	token := tokenParaUsuario(u.ID)
	saved := dao.DB
	dao.DB = nil
	defer func() { dao.DB = saved }()
	req, _ := http.NewRequest("GET", "/api/tickets/my-tickets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtenido %d", w.Code)
	}
}

func TestGetOccupationReport_VenueSinAsociar(t *testing.T) {
	e := domain.Event{
		Titulo: "Sin Venue Ctrl", Categoria: "Test",
		Lugar: "X", Precio: 10, CupoMaximo: 10, CupoDispon: 10,
		VenueID: 88888,
	}
	dao.DB.Create(&e)
	token := tokenParaAdmin()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/admin/events/%d/report", e.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtenido %d — body: %s", w.Code, w.Body.String())
	}
}