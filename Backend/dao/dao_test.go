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
	db.AutoMigrate(&models.User{}, &models.Venue{}, &models.Event{}, &models.Seat{}, &models.Ticket{})
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
	if err := DecrementCupo(DB, e.ID, 1); err != nil {
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
	err := DecrementCupo(DB, e.ID, 1)
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

func TestGetTicketsByUserID_DBNula(t *testing.T) {
	saved := DB
	DB = nil
	_, err := GetTicketsByUserID(1)
	DB = saved
	if err == nil {
		t.Error("se esperaba error con DB nula")
	}
}

func TestGetTicketByID_DBNula(t *testing.T) {
	saved := DB
	DB = nil
	_, err := GetTicketByID(1)
	DB = saved
	if err == nil {
		t.Error("se esperaba error con DB nula")
	}
}

func TestGetUserByDNI_DBNula(t *testing.T) {
	saved := DB
	DB = nil
	_, err := GetUserByDNI("12345678")
	DB = saved
	if err == nil {
		t.Error("se esperaba error con DB nula")
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

// ── dao.CreateEvent ───────────────────────────────────────────────

func TestCreateEventDAO_DBNula(t *testing.T) {
	saved := DB
	DB = nil
	err := CreateEvent(&models.Event{})
	DB = saved
	if err == nil {
		t.Error("se esperaba error con DB nula")
	}
}

func TestCreateEventDAO_Exitoso(t *testing.T) {
	e := &models.Event{Titulo: "DAO Direct", Categoria: "Test", Lugar: "L", Precio: 10, CupoMaximo: 5, CupoDispon: 5}
	if err := CreateEvent(e); err != nil {
		t.Fatal(err)
	}
	if e.ID == 0 {
		t.Error("evento sin ID luego de crearlo")
	}
}

// ── dao.GetAllVenues / GetVenueByID / CreateVenue / UpdateVenue / DeleteVenue ──

func seedVenueDAO() models.Venue {
	v := models.Venue{Nombre: "Venue DAO", Direccion: "Dir DAO", Filas: 2, ColumnasPorFila: 3, Capacidad: 6}
	DB.Create(&v)
	return v
}

func TestGetAllVenues_DBNula(t *testing.T) {
	saved := DB
	DB = nil
	_, err := GetAllVenues()
	DB = saved
	if err == nil {
		t.Error("se esperaba error con DB nula")
	}
}

func TestGetAllVenues_OK(t *testing.T) {
	seedVenueDAO()
	venues, err := GetAllVenues()
	if err != nil {
		t.Fatal(err)
	}
	if len(venues) == 0 {
		t.Error("se esperaba al menos un venue")
	}
}

func TestGetVenueByID_DBNula(t *testing.T) {
	saved := DB
	DB = nil
	_, err := GetVenueByID(1)
	DB = saved
	if err == nil {
		t.Error("se esperaba error con DB nula")
	}
}

func TestGetVenueByID_Existente(t *testing.T) {
	v := seedVenueDAO()
	found, err := GetVenueByID(v.ID)
	if err != nil {
		t.Fatal(err)
	}
	if found.ID != v.ID {
		t.Errorf("ID incorrecto: esperado %d, obtenido %d", v.ID, found.ID)
	}
}

func TestGetVenueByID_Inexistente(t *testing.T) {
	_, err := GetVenueByID(99999)
	if err == nil {
		t.Error("se esperaba error para ID inexistente")
	}
}

func TestCreateVenue_DBNula(t *testing.T) {
	saved := DB
	DB = nil
	err := CreateVenue(&models.Venue{})
	DB = saved
	if err == nil {
		t.Error("se esperaba error con DB nula")
	}
}

func TestCreateVenue_DAO(t *testing.T) {
	v := &models.Venue{Nombre: "V", Direccion: "D", Filas: 1, ColumnasPorFila: 1, Capacidad: 1}
	if err := CreateVenue(v); err != nil {
		t.Fatal(err)
	}
	if v.ID == 0 {
		t.Error("venue sin ID")
	}
}

func TestUpdateVenue_DAO(t *testing.T) {
	v := seedVenueDAO()
	fields := map[string]interface{}{"nombre": "Venue Actualizado DAO"}
	if err := UpdateVenue(v.ID, fields); err != nil {
		t.Fatal(err)
	}
	updated, _ := GetVenueByID(v.ID)
	if updated.Nombre != "Venue Actualizado DAO" {
		t.Errorf("nombre no actualizado: %s", updated.Nombre)
	}
}

func TestDeleteVenue_DAO(t *testing.T) {
	v := seedVenueDAO()
	if err := DeleteVenue(v.ID); err != nil {
		t.Fatal(err)
	}
}

// ── dao.Seat ──────────────────────────────────────────────────────

func seedSeatDAO(eventID uint) models.Seat {
	s := models.Seat{EventID: eventID, Fila: "A", Numero: 99, Ocupado: false}
	DB.Create(&s)
	return s
}

func TestGetSeatsByEventID_DBNula(t *testing.T) {
	saved := DB
	DB = nil
	_, err := GetSeatsByEventID(1)
	DB = saved
	if err == nil {
		t.Error("se esperaba error con DB nula")
	}
}

func TestGetSeatsByEventID_OK(t *testing.T) {
	e := seedEvento(5)
	seedSeatDAO(e.ID)
	seats, err := GetSeatsByEventID(e.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(seats) == 0 {
		t.Error("se esperaba al menos un asiento")
	}
}

func TestGetSeatByID_DBNula(t *testing.T) {
	saved := DB
	DB = nil
	_, err := GetSeatByID(1)
	DB = saved
	if err == nil {
		t.Error("se esperaba error con DB nula")
	}
}

func TestGetSeatByID_Existente(t *testing.T) {
	e := seedEvento(5)
	s := seedSeatDAO(e.ID)
	found, err := GetSeatByID(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if found.ID != s.ID {
		t.Errorf("ID incorrecto: esperado %d, obtenido %d", s.ID, found.ID)
	}
}

func TestGetSeatByID_Inexistente(t *testing.T) {
	_, err := GetSeatByID(99999)
	if err == nil {
		t.Error("se esperaba error para ID inexistente")
	}
}

func TestOccupyAndFreeSeats(t *testing.T) {
	e := seedEvento(5)
	s := models.Seat{EventID: e.ID, Fila: "Z", Numero: 1, Ocupado: false}
	DB.Create(&s)

	if err := OccupySeats(DB, []uint{s.ID}); err != nil {
		t.Fatal(err)
	}
	var afterOccupy models.Seat
	DB.First(&afterOccupy, s.ID)
	if !afterOccupy.Ocupado {
		t.Error("asiento debería estar ocupado después de OccupySeats")
	}

	if err := FreeSeats(DB, []uint{s.ID}); err != nil {
		t.Fatal(err)
	}
	var afterFree models.Seat
	DB.First(&afterFree, s.ID)
	if afterFree.Ocupado {
		t.Error("asiento debería estar libre después de FreeSeats")
	}
}

func TestFreeSeatByTicketSeatID_DAO(t *testing.T) {
	e := seedEvento(5)
	s := models.Seat{EventID: e.ID, Fila: "Y", Numero: 1, Ocupado: true}
	DB.Create(&s)

	if err := FreeSeatByTicketSeatID(DB, s.ID); err != nil {
		t.Fatal(err)
	}
	var updated models.Seat
	DB.First(&updated, s.ID)
	if updated.Ocupado {
		t.Error("asiento debería estar libre")
	}
}

func TestCreateSeatsForEvent_OK(t *testing.T) {
	e := seedEvento(10)
	if err := CreateSeatsForEvent(DB, e.ID, []SectorDef{{Nombre: "Test", Capacidad: 12}}); err != nil {
		t.Fatal(err)
	}
	seats, _ := GetSeatsByEventID(e.ID)
	if len(seats) < 12 {
		t.Errorf("cantidad de asientos: esperado >=12, obtenido %d", len(seats))
	}
}

func TestGetActiveTicketCount(t *testing.T) {
	e := seedEvento(10)
	u := seedUsuario("counttickets@dao.test", "90190190")
	DB.Create(&models.Ticket{UserID: u.ID, EventID: e.ID, Estado: "activo"})
	DB.Create(&models.Ticket{UserID: u.ID, EventID: e.ID, Estado: "activo"})
	DB.Create(&models.Ticket{UserID: u.ID, EventID: e.ID, Estado: "cancelado"})

	count, err := GetActiveTicketCountByEventAndUser(u.ID, e.ID)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("esperado 2 tickets activos, obtenido %d", count)
	}
}

func TestGetEventReport_Existente(t *testing.T) {
	e := seedEvento(10)
	u := seedUsuario("report@dao.test", "90290290")
	DB.Create(&models.Ticket{UserID: u.ID, EventID: e.ID, Estado: "activo"})

	report, err := GetEventReport(e.ID)
	if err != nil {
		t.Fatal(err)
	}
	if report["entradas_vendidas"].(int64) != 1 {
		t.Errorf("esperado 1 vendida, obtenido %v", report["entradas_vendidas"])
	}
}

func TestGetEventReport_Inexistente(t *testing.T) {
	_, err := GetEventReport(99999)
	if err == nil {
		t.Error("se esperaba error para evento inexistente")
	}
}
