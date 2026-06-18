package dao

import (
	"os"
	"testing"

	"backend/domain/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("no se pudo abrir la BD de test: " + err.Error())
	}
	db.AutoMigrate(&models.User{}, &models.Event{}, &models.Ticket{})
	DB = db
	os.Exit(m.Run())
}

func seedEvento(cupo int) models.Event {
	e := models.Event{
		Titulo: "Evento DAO Test", Categoria: "TestCat",
		Lugar: "Córdoba", Precio: 10, CupoMaximo: cupo, CupoDispon: cupo,
	}
	DB.Create(&e)
	return e
}

func seedUsuario(email, dni string) models.User {
	u := models.User{
		Email: email, Password: "hash",
		Nombre: "T", Apellido: "U", Rol: "cliente", DNI: dni,
	}
	DB.Create(&u)
	return u
}

// ── Eventos ───────────────────────────────────────────────────────

func TestGetAllEvents_DBNula(t *testing.T) {
	saved := DB
	DB = nil
	_, err := GetAllEvents("")
	DB = saved
	if err == nil {
		t.Error("se esperaba error con DB nula")
	}
}

func TestGetEventByID_DBNula(t *testing.T) {
	saved := DB
	DB = nil
	_, err := GetEventByID(1)
	DB = saved
	if err == nil {
		t.Error("se esperaba error con DB nula")
	}
}

func TestGetAllEvents_SinFiltro(t *testing.T) {
	seedEvento(10)
	events, err := GetAllEvents("")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) == 0 {
		t.Error("se esperaba al menos un evento")
	}
}

func TestGetAllEvents_ConCategoria(t *testing.T) {
	DB.Create(&models.Event{
		Titulo: "Especial", Categoria: "Unica",
		Lugar: "X", Precio: 1, CupoMaximo: 1, CupoDispon: 1,
	})
	events, err := GetAllEvents("Unica")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) == 0 {
		t.Error("se esperaba un evento de categoría Unica")
	}
}

func TestGetAllEvents_CategoriaInexistente(t *testing.T) {
	events, err := GetAllEvents("NoExiste999")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 0 {
		t.Error("no debería haber eventos de categoría inexistente")
	}
}

func TestGetEventByID_Existente(t *testing.T) {
	e := seedEvento(5)
	found, err := GetEventByID(e.ID)
	if err != nil {
		t.Fatal(err)
	}
	if found.ID != e.ID {
		t.Errorf("ID incorrecto: esperado %d, obtenido %d", e.ID, found.ID)
	}
}

func TestGetEventByID_Inexistente(t *testing.T) {
	_, err := GetEventByID(999999)
	if err == nil {
		t.Error("se esperaba error para ID inexistente")
	}
	if err.Error() != "evento no encontrado" {
		t.Errorf("mensaje de error incorrecto: %s", err.Error())
	}
}

// ── Cupo ──────────────────────────────────────────────────────────

func TestDecrementCupo_ConCupo(t *testing.T) {
	e := seedEvento(5)
	if err := DecrementCupo(DB, e.ID); err != nil {
		t.Fatalf("DecrementCupo falló: %v", err)
	}
	var updated models.Event
	DB.First(&updated, e.ID)
	if updated.CupoDispon != 4 {
		t.Errorf("cupo esperado 4, obtenido %d", updated.CupoDispon)
	}
}

func TestDecrementCupo_SinCupo(t *testing.T) {
	e := seedEvento(0)
	err := DecrementCupo(DB, e.ID)
	if err == nil {
		t.Error("se esperaba error cuando cupo_disponible = 0")
	}
	if err.Error() != "sin cupo disponible" {
		t.Errorf("mensaje incorrecto: %s", err.Error())
	}
}

func TestIncrementCupo(t *testing.T) {
	e := seedEvento(3)
	if err := IncrementCupo(DB, e.ID); err != nil {
		t.Fatal(err)
	}
	var updated models.Event
	DB.First(&updated, e.ID)
	if updated.CupoDispon != 4 {
		t.Errorf("cupo esperado 4, obtenido %d", updated.CupoDispon)
	}
}

// ── Tickets ───────────────────────────────────────────────────────

func TestCreateTicket(t *testing.T) {
	e := seedEvento(10)
	u := seedUsuario("create@dao.test", "10010010")
	ticket := &models.Ticket{UserID: u.ID, EventID: e.ID, Estado: "activo"}
	if err := CreateTicket(DB, ticket); err != nil {
		t.Fatal(err)
	}
	if ticket.ID == 0 {
		t.Error("ticket sin ID luego de crearlo")
	}
}

func TestGetTicketByID_Existente(t *testing.T) {
	e := seedEvento(10)
	u := seedUsuario("getbyid@dao.test", "20020020")
	ticket := &models.Ticket{UserID: u.ID, EventID: e.ID, Estado: "activo"}
	DB.Create(ticket)

	found, err := GetTicketByID(ticket.ID)
	if err != nil {
		t.Fatal(err)
	}
	if found.ID != ticket.ID {
		t.Errorf("ID incorrecto: esperado %d, obtenido %d", ticket.ID, found.ID)
	}
}

func TestGetTicketByID_Inexistente(t *testing.T) {
	_, err := GetTicketByID(999999)
	if err == nil {
		t.Error("se esperaba error para ID inexistente")
	}
	if err.Error() != "ticket no encontrado" {
		t.Errorf("mensaje de error incorrecto: %s", err.Error())
	}
}

func TestGetTicketsByUserID(t *testing.T) {
	u := seedUsuario("mytickets@dao.test", "30030030")
	tickets, err := GetTicketsByUserID(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	_ = tickets
}

func TestUpdateTicketEstado(t *testing.T) {
	e := seedEvento(10)
	u := seedUsuario("estado@dao.test", "40040040")
	ticket := &models.Ticket{UserID: u.ID, EventID: e.ID, Estado: "activo"}
	DB.Create(ticket)

	if err := UpdateTicketEstado(DB, ticket.ID, "cancelado"); err != nil {
		t.Fatal(err)
	}
	var updated models.Ticket
	DB.First(&updated, ticket.ID)
	if updated.Estado != "cancelado" {
		t.Errorf("estado esperado 'cancelado', obtenido '%s'", updated.Estado)
	}
}

func TestTransferTicketOwner(t *testing.T) {
	e := seedEvento(10)
	owner := seedUsuario("owner@dao.test", "50050050")
	newOwner := seedUsuario("newowner@dao.test", "60060060")
	ticket := &models.Ticket{UserID: owner.ID, EventID: e.ID, Estado: "activo"}
	DB.Create(ticket)

	if err := TransferTicketOwner(DB, ticket.ID, newOwner.ID); err != nil {
		t.Fatal(err)
	}
	var updated models.Ticket
	DB.First(&updated, ticket.ID)
	if updated.UserID != newOwner.ID {
		t.Errorf("propietario no transferido: esperado %d, obtenido %d", newOwner.ID, updated.UserID)
	}
}

func TestGetUserByDNI_Existente(t *testing.T) {
	u := seedUsuario("dniok@dao.test", "70070070")
	found, err := GetUserByDNI(u.DNI)
	if err != nil {
		t.Fatal(err)
	}
	if found.DNI != u.DNI {
		t.Errorf("DNI incorrecto: esperado %s, obtenido %s", u.DNI, found.DNI)
	}
}

func TestGetUserByDNI_Inexistente(t *testing.T) {
	_, err := GetUserByDNI("00000000")
	if err == nil {
		t.Error("se esperaba error para DNI inexistente")
	}
	if err.Error() != "usuario con ese DNI no encontrado" {
		t.Errorf("mensaje de error incorrecto: %s", err.Error())
	}
}

func TestGetTicketsByUserID_ConTickets(t *testing.T) {
	e := seedEvento(10)
	u := seedUsuario("contickets@dao.test", "80080080")
	ticket := &models.Ticket{UserID: u.ID, EventID: e.ID, Estado: "activo"}
	DB.Create(ticket)

	tickets, err := GetTicketsByUserID(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(tickets) == 0 {
		t.Error("se esperaba al menos un ticket para el usuario")
	}
}

func TestUpdateEvent(t *testing.T) {
	e := seedEvento(10)
	fields := map[string]interface{}{"titulo": "Actualizado", "precio": 999.0}
	if err := UpdateEvent(e.ID, fields); err != nil {
		t.Fatal(err)
	}
	updated, _ := GetEventByID(e.ID)
	if updated.Titulo != "Actualizado" {
		t.Errorf("título no actualizado: obtenido %s", updated.Titulo)
	}
}

func TestDeleteEvent(t *testing.T) {
	e := seedEvento(5)
	if err := DeleteEvent(e.ID); err != nil {
		t.Fatal(err)
	}
	_, err := GetEventByID(e.ID)
	if err == nil {
		t.Error("el evento debería haberse eliminado")
	}
}
