package services

import (
	"os"
	"testing"
	"time"

	"backend/dao"
	domain "backend/domain/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("no se pudo abrir la BD de test: " + err.Error())
	}
	db.AutoMigrate(&domain.User{}, &domain.Event{}, &domain.Ticket{})
	dao.DB = db
	os.Setenv("JWT_SECRET", "secret-de-test-123")
	os.Exit(m.Run())
}

func seedEvent(db *gorm.DB) domain.Event {
	e := domain.Event{
		Titulo:      "Evento Test",
		Descripcion: "Descripcion de prueba",
		Categoria:   "Recitales",
		Lugar:       "Buenos Aires",
		Precio:      100.0,
		CupoMaximo:  50,
		CupoDispon:  50,
	}
	db.Create(&e)
	return e
}

func seedUser(db *gorm.DB, email, dni string) domain.User {
	u := domain.User{
		Email:    email,
		Password: hashPassword("password123"),
		Nombre:   "Test",
		Apellido: "User",
		Rol:      "cliente",
		DNI:      dni,
	}
	db.Create(&u)
	return u
}

// ── Eventos ──────────────────────────────────────────────────────

func TestGetAllEvents_SinCategoria(t *testing.T) {
	seedEvent(dao.DB)
	events, err := GetAllEvents("")
	if err != nil {
		t.Fatalf("GetAllEvents falló: %v", err)
	}
	if len(events) == 0 {
		t.Error("Se esperaba al menos un evento")
	}
}

func TestGetAllEvents_ConCategoria(t *testing.T) {
	events, err := GetAllEvents("Recitales")
	if err != nil {
		t.Fatalf("GetAllEvents con categoría falló: %v", err)
	}
	t.Logf("Eventos en categoría Recitales: %d", len(events))
}

func TestGetAllEvents_CategoriaInexistente(t *testing.T) {
	events, err := GetAllEvents("Deportes")
	if err != nil {
		t.Fatalf("GetAllEvents con categoría vacía falló: %v", err)
	}
	if len(events) != 0 {
		t.Error("No debería haber eventos de Deportes")
	}
}

func TestGetEventByID_Existente(t *testing.T) {
	e := seedEvent(dao.DB)
	_, err := GetEventByID(e.ID)
	if err != nil {
		t.Fatalf("GetEventByID falló para ID existente: %v", err)
	}
}

func TestGetEventByID_Inexistente(t *testing.T) {
	_, err := GetEventByID(99999)
	if err == nil {
		t.Error("Se esperaba error para ID inexistente")
	}
}

// ── Auth ─────────────────────────────────────────────────────────

func TestRegister_Exitoso(t *testing.T) {
	input := RegisterInput{
		Email:    "nuevo@test.com",
		Password: "password123",
		Nombre:   "Juan",
		Apellido: "Perez",
		DNI:      "11111111",
	}
	user, err := Register(input)
	if err != nil {
		t.Fatalf("Register falló: %v", err)
	}
	if user.Email != input.Email {
		t.Errorf("Email incorrecto: %s", user.Email)
	}
}

func TestRegister_EmailDuplicado(t *testing.T) {
	input := RegisterInput{
		Email:    "duplicado@test.com",
		Password: "password123",
		Nombre:   "Test",
		Apellido: "User",
		DNI:      "22222222",
	}
	Register(input)
	_, err := Register(input)
	if err == nil {
		t.Error("Se esperaba error por email duplicado")
	}
}

func TestLogin_CredencialesInvalidas(t *testing.T) {
	input := LoginInput{
		Email:    "noexiste@test.com",
		Password: "password123",
	}
	_, _, err := Login(input)
	if err == nil {
		t.Error("Se esperaba error con email inexistente")
	}
}

// ── JWT (función pura) ────────────────────────────────────────────

func TestGenerateJWT_ConSecret(t *testing.T) {
	user := domain.User{ID: 1, Email: "simon@test.com", Rol: "cliente"}
	token, err := generateJWT(user)
	if err != nil {
		t.Fatalf("generateJWT falló: %v", err)
	}
	if token == "" {
		t.Error("Token vacío")
	}
}

func TestGenerateJWT_SinSecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")
	user := domain.User{ID: 2, Email: "admin@test.com", Rol: "admin"}
	token, err := generateJWT(user)
	if err != nil {
		t.Fatalf("generateJWT falló sin JWT_SECRET: %v", err)
	}
	if token == "" {
		t.Error("Token vacío sin secret")
	}
	os.Setenv("JWT_SECRET", "secret-de-test-123")
}

// ── Tickets ───────────────────────────────────────────────────────

func TestPurchaseTicket_Exitoso(t *testing.T) {
	e := seedEvent(dao.DB)
	u := seedUser(dao.DB, "comprador@test.com", "33333333")

	ticket, err := PurchaseTicket(u.ID, e.ID)
	if err != nil {
		t.Fatalf("PurchaseTicket falló: %v", err)
	}
	if ticket.UserID != u.ID {
		t.Errorf("UserID incorrecto: %d", ticket.UserID)
	}
}

func TestPurchaseTicket_SinCupo(t *testing.T) {
	e := domain.Event{
		Titulo:     "Evento Sin Cupo",
		Categoria:  "Teatro",
		Lugar:      "Rosario",
		Precio:     50.0,
		CupoMaximo: 1,
		CupoDispon: 0,
	}
	dao.DB.Create(&e)
	u := seedUser(dao.DB, "sinsuerte@test.com", "44444444")

	_, err := PurchaseTicket(u.ID, e.ID)
	if err == nil {
		t.Error("Se esperaba error por falta de cupo")
	}
}

func TestGetMyTickets_Exitoso(t *testing.T) {
	u := seedUser(dao.DB, "myticketsuser@test.com", "55555555")
	tickets, err := GetMyTickets(u.ID)
	if err != nil {
		t.Fatalf("GetMyTickets falló: %v", err)
	}
	t.Logf("Tickets del usuario: %d", len(tickets))
}

func TestCancelTicket_NoAutorizado(t *testing.T) {
	e := seedEvent(dao.DB)
	u := seedUser(dao.DB, "owner@test.com", "66666666")
	ticket, _ := PurchaseTicket(u.ID, e.ID)

	err := CancelTicket(9999, ticket.ID)
	if err == nil {
		t.Error("Se esperaba error de no autorizado")
	}
}

func TestCancelTicket_Exitoso(t *testing.T) {
	e := seedEvent(dao.DB)
	u := seedUser(dao.DB, "canceler@test.com", "77777777")
	ticket, _ := PurchaseTicket(u.ID, e.ID)

	err := CancelTicket(u.ID, ticket.ID)
	if err != nil {
		t.Fatalf("CancelTicket falló: %v", err)
	}
}

// ── Admin — Eventos ───────────────────────────────────────────────

func TestCreateEvent_Exitoso(t *testing.T) {
	input := EventInput{
		Titulo:     "Evento Admin",
		Categoria:  "Rock",
		Lugar:      "Córdoba",
		Precio:     100,
		CupoMaximo: 50,
		Fecha:      time.Now().Add(24 * time.Hour),
	}
	event, err := CreateEvent(input)
	if err != nil {
		t.Fatal(err)
	}
	if event.CupoDispon != 50 {
		t.Errorf("cupo disponible debe igualar cupo máximo al crear: obtenido %d", event.CupoDispon)
	}
}

func TestUpdateEvent_Exitoso(t *testing.T) {
	e := seedEvent(dao.DB)
	input := EventInput{
		Titulo:     "Modificado",
		Categoria:  "Pop",
		Lugar:      "Rosario",
		Precio:     200,
		CupoMaximo: 30,
		Fecha:      time.Now().Add(48 * time.Hour),
	}
	updated, err := UpdateEvent(e.ID, input)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Titulo != "Modificado" {
		t.Errorf("título no actualizado: obtenido %s", updated.Titulo)
	}
}

func TestUpdateEvent_Inexistente(t *testing.T) {
	input := EventInput{
		Titulo:     "X",
		Categoria:  "X",
		Lugar:      "X",
		Precio:     1,
		CupoMaximo: 1,
		Fecha:      time.Now().Add(24 * time.Hour),
	}
	_, err := UpdateEvent(99999, input)
	if err == nil {
		t.Error("se esperaba error para evento inexistente")
	}
}

func TestDeleteEvent_Exitoso(t *testing.T) {
	e := seedEvent(dao.DB)
	if err := DeleteEvent(e.ID); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEvent_Inexistente(t *testing.T) {
	if err := DeleteEvent(99999); err == nil {
		t.Error("se esperaba error para evento inexistente")
	}
}

func TestTransferTicket_AutoTransferencia(t *testing.T) {
	e := seedEvent(dao.DB)
	u := seedUser(dao.DB, "autotransfer@test.com", "12312312")
	ticket, _ := PurchaseTicket(u.ID, e.ID)

	err := TransferTicket(u.ID, ticket.ID, u.DNI)
	if err == nil {
		t.Error("no se debería poder transferir un ticket a uno mismo")
	}
}
