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
	db.AutoMigrate(&domain.User{}, &domain.Venue{}, &domain.Event{}, &domain.Seat{}, &domain.Ticket{})
	dao.DB = db
	os.Setenv("JWT_SECRET", "secret-de-test-123")
	os.Exit(m.Run())
}

func seedVenue(db *gorm.DB) domain.Venue {
	v := domain.Venue{
		Nombre:          "Venue Test",
		Direccion:       "Calle Falsa 123",
		Filas:           5,
		ColumnasPorFila: 10,
		Capacidad:       50,
	}
	db.Create(&v)
	return v
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

func seedSeat(db *gorm.DB, eventID uint) domain.Seat {
	s := domain.Seat{
		EventID: eventID,
		Fila:    "A",
		Numero:  1,
		Ocupado: false,
	}
	db.Create(&s)
	return s
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
	s := seedSeat(dao.DB, e.ID)
	u := seedUser(dao.DB, "comprador@test.com", "33333333")

	tickets, err := PurchaseTickets(u.ID, e.ID, []uint{s.ID})
	if err != nil {
		t.Fatalf("PurchaseTickets falló: %v", err)
	}
	if tickets[0].UserID != u.ID {
		t.Errorf("UserID incorrecto: %d", tickets[0].UserID)
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
	s := seedSeat(dao.DB, e.ID)
	u := seedUser(dao.DB, "sinsuerte@test.com", "44444444")

	_, err := PurchaseTickets(u.ID, e.ID, []uint{s.ID})
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
	s := seedSeat(dao.DB, e.ID)
	u := seedUser(dao.DB, "owner@test.com", "66666666")
	tickets, _ := PurchaseTickets(u.ID, e.ID, []uint{s.ID})

	err := CancelTicket(9999, tickets[0].ID)
	if err == nil {
		t.Error("Se esperaba error de no autorizado")
	}
}

func TestCancelTicket_Exitoso(t *testing.T) {
	e := seedEvent(dao.DB)
	s := seedSeat(dao.DB, e.ID)
	u := seedUser(dao.DB, "canceler@test.com", "77777777")
	tickets, _ := PurchaseTickets(u.ID, e.ID, []uint{s.ID})

	err := CancelTicket(u.ID, tickets[0].ID)
	if err != nil {
		t.Fatalf("CancelTicket falló: %v", err)
	}
}

// ── Admin — Eventos ───────────────────────────────────────────────

func TestCreateEvent_Exitoso(t *testing.T) {
	v := seedVenue(dao.DB)
	input := EventInput{
		Titulo:    "Evento Admin",
		Categoria: "Rock",
		Fecha:     time.Now().Add(24 * time.Hour),
		Precio:    100,
		VenueID:   v.ID,
	}
	event, err := CreateEvent(input)
	if err != nil {
		t.Fatal(err)
	}
	if event.CupoDispon != v.Capacidad {
		t.Errorf("cupo disponible debe igualar capacidad del venue: obtenido %d", event.CupoDispon)
	}
}

func TestUpdateEvent_Exitoso(t *testing.T) {
	v := seedVenue(dao.DB)
	e := seedEvent(dao.DB)
	input := EventInput{
		Titulo:    "Modificado",
		Categoria: "Pop",
		Fecha:     time.Now().Add(48 * time.Hour),
		Precio:    200,
		VenueID:   v.ID,
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
	v := seedVenue(dao.DB)
	input := EventInput{
		Titulo:    "X",
		Categoria: "X",
		Fecha:     time.Now().Add(24 * time.Hour),
		Precio:    1,
		VenueID:   v.ID,
	}
	_, err := UpdateEvent(99999, input)
	if err == nil {
		t.Error("se esperaba error para evento inexistente")
	}
}

// ── Venue service ─────────────────────────────────────────────────

func TestGetAllVenues_OK(t *testing.T) {
	_, err := GetAllVenues()
	if err != nil {
		t.Fatalf("GetAllVenues falló: %v", err)
	}
}

func TestGetVenueByID_Existente(t *testing.T) {
	v := seedVenue(dao.DB)
	_, err := GetVenueByID(v.ID)
	if err != nil {
		t.Fatalf("GetVenueByID falló para ID existente: %v", err)
	}
}

func TestGetVenueByID_Inexistente(t *testing.T) {
	_, err := GetVenueByID(99999)
	if err == nil {
		t.Error("se esperaba error para ID inexistente")
	}
}

func TestCreateVenue_Exitoso(t *testing.T) {
	input := VenueInput{
		Nombre:          "Estadio Test Service",
		Direccion:       "Calle Test 123",
		Filas:           4,
		ColumnasPorFila: 5,
	}
	venue, err := CreateVenue(input)
	if err != nil {
		t.Fatalf("CreateVenue falló: %v", err)
	}
	if venue.Capacidad != 20 {
		t.Errorf("capacidad incorrecta: %d", venue.Capacidad)
	}
}

func TestUpdateVenue_Exitoso(t *testing.T) {
	v := seedVenue(dao.DB)
	input := VenueInput{
		Nombre:          "Venue Actualizado",
		Direccion:       "Nueva Dirección 456",
		Filas:           3,
		ColumnasPorFila: 6,
	}
	updated, err := UpdateVenue(v.ID, input)
	if err != nil {
		t.Fatalf("UpdateVenue falló: %v", err)
	}
	if updated.Nombre != "Venue Actualizado" {
		t.Errorf("nombre no actualizado: %s", updated.Nombre)
	}
}

func TestUpdateVenue_Inexistente(t *testing.T) {
	input := VenueInput{Nombre: "X", Direccion: "X", Filas: 1, ColumnasPorFila: 1}
	_, err := UpdateVenue(99999, input)
	if err == nil {
		t.Error("se esperaba error para venue inexistente")
	}
}

func TestUpdateVenue_ConEventos(t *testing.T) {
	v := seedVenue(dao.DB)
	dao.DB.Create(&domain.Event{
		Titulo: "Evento Venue Linked", Categoria: "Test",
		Lugar: "Somewhere", Precio: 10, CupoMaximo: 10, CupoDispon: 10,
		VenueID: v.ID,
	})
	input := VenueInput{Nombre: "X", Direccion: "X", Filas: 1, ColumnasPorFila: 1}
	_, err := UpdateVenue(v.ID, input)
	if err == nil {
		t.Error("se esperaba error al actualizar venue con eventos asociados")
	}
}

func TestDeleteVenue_Exitoso(t *testing.T) {
	v := seedVenue(dao.DB)
	if err := DeleteVenue(v.ID); err != nil {
		t.Fatalf("DeleteVenue falló: %v", err)
	}
}

func TestDeleteVenue_Inexistente(t *testing.T) {
	if err := DeleteVenue(99999); err == nil {
		t.Error("se esperaba error para venue inexistente")
	}
}

func TestDeleteVenue_ConEventos(t *testing.T) {
	v := seedVenue(dao.DB)
	dao.DB.Create(&domain.Event{
		Titulo: "Evento No Borrar Venue", Categoria: "Test",
		Lugar: "Somewhere", Precio: 10, CupoMaximo: 10, CupoDispon: 10,
		VenueID: v.ID,
	})
	if err := DeleteVenue(v.ID); err == nil {
		t.Error("se esperaba error al eliminar venue con eventos asociados")
	}
}

// ── Reporte ───────────────────────────────────────────────────────

func TestGetOccupationReport_Exitoso(t *testing.T) {
	v := seedVenue(dao.DB)
	e := domain.Event{
		Titulo: "Evento Reporte", Categoria: "Test",
		Lugar: "Somewhere", Precio: 10, CupoMaximo: 50, CupoDispon: 50,
		VenueID: v.ID,
	}
	dao.DB.Create(&e)

	report, err := GetOccupationReport(e.ID)
	if err != nil {
		t.Fatalf("GetOccupationReport falló: %v", err)
	}
	if report.EventID != e.ID {
		t.Errorf("EventID incorrecto: %d", report.EventID)
	}
}

func TestGetOccupationReport_Inexistente(t *testing.T) {
	_, err := GetOccupationReport(99999)
	if err == nil {
		t.Error("se esperaba error para evento inexistente")
	}
}

// ── Asientos ──────────────────────────────────────────────────────

func TestGetSeatsByEventID_OK(t *testing.T) {
	v := seedVenue(dao.DB)
	event, err := CreateEvent(EventInput{
		Titulo:    "Evento Con Asientos",
		Categoria: "Test",
		Fecha:     time.Now().Add(24 * time.Hour),
		Precio:    10,
		VenueID:   v.ID,
	})
	if err != nil {
		t.Fatalf("CreateEvent falló: %v", err)
	}
	seats, err := GetSeatsByEventID(event.ID)
	if err != nil {
		t.Fatalf("GetSeatsByEventID falló: %v", err)
	}
	expected := v.Filas * v.ColumnasPorFila
	if len(seats) != expected {
		t.Errorf("cantidad de asientos: esperado %d, obtenido %d", expected, len(seats))
	}
}

// ── Login paths adicionales ───────────────────────────────────────

func TestLogin_Exitoso(t *testing.T) {
	Register(RegisterInput{
		Email: "loginok@svc.test", Password: "pass123",
		Nombre: "L", Apellido: "T", DNI: "99887766",
	})
	_, _, err := Login(LoginInput{Email: "loginok@svc.test", Password: "pass123"})
	if err != nil {
		t.Fatalf("Login exitoso falló: %v", err)
	}
}

func TestLogin_PasswordIncorrecto(t *testing.T) {
	Register(RegisterInput{
		Email: "wrongpwd@svc.test", Password: "pass123",
		Nombre: "W", Apellido: "P", DNI: "11223344",
	})
	_, _, err := Login(LoginInput{Email: "wrongpwd@svc.test", Password: "contraseniaincorrecta"})
	if err == nil {
		t.Error("se esperaba error por password incorrecto")
	}
}

// ── Tickets edge cases ─────────────────────────────────────────────

func TestCancelTicket_TicketNoActivo(t *testing.T) {
	e := seedEvent(dao.DB)
	u := seedUser(dao.DB, "cancelinactivo@test.com", "55551111")
	s := seedSeat(dao.DB, e.ID)
	tickets, _ := PurchaseTickets(u.ID, e.ID, []uint{s.ID})
	CancelTicket(u.ID, tickets[0].ID)
	err := CancelTicket(u.ID, tickets[0].ID)
	if err == nil {
		t.Error("se esperaba error al cancelar ticket ya cancelado")
	}
}

func TestPurchaseTicket_LimiteExcedido(t *testing.T) {
	e := seedEvent(dao.DB)
	dao.DB.Model(&domain.Event{}).Where("id = ?", e.ID).Update("cupo_disponible", 20)
	u := seedUser(dao.DB, "limiteuser@test.com", "11112222")

	var seatIDs []uint
	for i := 1; i <= 11; i++ {
		s := domain.Seat{EventID: e.ID, Fila: "X", Numero: i, Ocupado: false}
		dao.DB.Create(&s)
		seatIDs = append(seatIDs, s.ID)
	}
	_, err := PurchaseTickets(u.ID, e.ID, seatIDs)
	if err == nil {
		t.Error("se esperaba error por límite de tickets excedido")
	}
}

func TestPurchaseTicket_AsientoNoExiste(t *testing.T) {
	e := seedEvent(dao.DB)
	u := seedUser(dao.DB, "noasiento@test.com", "22221111")
	_, err := PurchaseTickets(u.ID, e.ID, []uint{999999})
	if err == nil {
		t.Error("se esperaba error por asiento inexistente")
	}
}

func TestCreateEvent_VenueInexistente(t *testing.T) {
	input := EventInput{
		Titulo: "T", Categoria: "C",
		Fecha:   time.Now().Add(24 * time.Hour),
		Precio:  10,
		VenueID: 99999,
	}
	_, err := CreateEvent(input)
	if err == nil {
		t.Error("se esperaba error para venue inexistente")
	}
}

func TestUpdateEvent_VenueInexistente(t *testing.T) {
	e := seedEvent(dao.DB)
	input := EventInput{
		Titulo: "T", Categoria: "C",
		Fecha:   time.Now().Add(24 * time.Hour),
		Precio:  10,
		VenueID: 99999,
	}
	_, err := UpdateEvent(e.ID, input)
	if err == nil {
		t.Error("se esperaba error para venue inexistente en update")
	}
}

func TestRegister_DNIDuplicado(t *testing.T) {
	Register(RegisterInput{
		Email: "dniA@svc.test", Password: "pass123",
		Nombre: "A", Apellido: "B", DNI: "duplicadodni",
	})
	_, err := Register(RegisterInput{
		Email: "dniB@svc.test", Password: "pass123",
		Nombre: "C", Apellido: "D", DNI: "duplicadodni",
	})
	if err == nil {
		t.Error("se esperaba error por DNI duplicado")
	}
}

func TestRegister_DBNula(t *testing.T) {
	saved := dao.DB
	dao.DB = nil
	_, err := Register(RegisterInput{
		Email: "dbnula@svc.test", Password: "pass123",
		Nombre: "X", Apellido: "Y", DNI: "00000001",
	})
	dao.DB = saved
	if err == nil {
		t.Error("se esperaba error con DB nula")
	}
}

func TestPurchaseTicket_AsientoOcupado(t *testing.T) {
	e := seedEvent(dao.DB)
	u1 := seedUser(dao.DB, "compradorA@test.com", "33331111")
	u2 := seedUser(dao.DB, "compradorB@test.com", "33332222")
	s := seedSeat(dao.DB, e.ID)

	PurchaseTickets(u1.ID, e.ID, []uint{s.ID})

	_, err := PurchaseTickets(u2.ID, e.ID, []uint{s.ID})
	if err == nil {
		t.Error("se esperaba error por asiento ya ocupado")
	}
}

func TestPurchaseTicket_AsientoDeOtroEvento(t *testing.T) {
	e1 := seedEvent(dao.DB)
	e2 := seedEvent(dao.DB)
	s := seedSeat(dao.DB, e1.ID)
	u := seedUser(dao.DB, "wrongevent@test.com", "44441111")

	_, err := PurchaseTickets(u.ID, e2.ID, []uint{s.ID})
	if err == nil {
		t.Error("se esperaba error por asiento perteneciente a otro evento")
	}
}

func TestCancelTicket_SinSeat(t *testing.T) {
	e := seedEvent(dao.DB)
	u := seedUser(dao.DB, "cancelnoseat@test.com", "55551234")
	ticket := domain.Ticket{UserID: u.ID, EventID: e.ID, Estado: "activo"}
	dao.DB.Create(&ticket)

	err := CancelTicket(u.ID, ticket.ID)
	if err != nil {
		t.Fatalf("CancelTicket sin seat falló: %v", err)
	}
}

func TestTransferTicket_TicketNoActivo(t *testing.T) {
	e := seedEvent(dao.DB)
	u := seedUser(dao.DB, "transfernoactivo@test.com", "11113333")
	s := seedSeat(dao.DB, e.ID)
	tickets, _ := PurchaseTickets(u.ID, e.ID, []uint{s.ID})
	CancelTicket(u.ID, tickets[0].ID)

	target := seedUser(dao.DB, "transfertargetNA@test.com", "22223333")
	err := TransferTicket(u.ID, tickets[0].ID, target.DNI)
	if err == nil {
		t.Error("se esperaba error al transferir ticket no activo")
	}
}

func TestTransferTicket_DestinatarioNoEncontrado(t *testing.T) {
	e := seedEvent(dao.DB)
	u := seedUser(dao.DB, "transfersindest@test.com", "33334444")
	s := seedSeat(dao.DB, e.ID)
	tickets, _ := PurchaseTickets(u.ID, e.ID, []uint{s.ID})

	err := TransferTicket(u.ID, tickets[0].ID, "dninoexiste999")
	if err == nil {
		t.Error("se esperaba error por destinatario no encontrado")
	}
}

func TestGetOccupationReport_VenueSinAsociar(t *testing.T) {
	e := domain.Event{
		Titulo: "Reporte Sin Venue", Categoria: "Test",
		Lugar: "X", Precio: 10, CupoMaximo: 10, CupoDispon: 10,
		VenueID: 77777,
	}
	dao.DB.Create(&e)
	_, err := GetOccupationReport(e.ID)
	if err == nil {
		t.Error("se esperaba error por venue no encontrado")
	}
}

func TestLogin_DBNula(t *testing.T) {
	saved := dao.DB
	dao.DB = nil
	_, _, err := Login(LoginInput{Email: "x@x.com", Password: "pass"})
	dao.DB = saved
	if err == nil {
		t.Error("se esperaba error con DB nula")
	}
}

func TestCancelTicket_DBNula(t *testing.T) {
	saved := dao.DB
	dao.DB = nil
	err := CancelTicket(1, 1)
	dao.DB = saved
	if err == nil {
		t.Error("se esperaba error con DB nula")
	}
}

func TestPurchaseTickets_DBNula(t *testing.T) {
	saved := dao.DB
	dao.DB = nil
	_, err := PurchaseTickets(1, 1, []uint{1})
	dao.DB = saved
	if err == nil {
		t.Error("se esperaba error con DB nula")
	}
}

func TestTransferTicket_DBNula(t *testing.T) {
	saved := dao.DB
	dao.DB = nil
	err := TransferTicket(1, 1, "12345678")
	dao.DB = saved
	if err == nil {
		t.Error("se esperaba error con DB nula")
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
	s := seedSeat(dao.DB, e.ID)
	u := seedUser(dao.DB, "autotransfer@test.com", "12312312")
	tickets, _ := PurchaseTickets(u.ID, e.ID, []uint{s.ID})

	err := TransferTicket(u.ID, tickets[0].ID, u.DNI)
	if err == nil {
		t.Error("no se debería poder transferir un ticket a uno mismo")
	}
}

func TestTransferTicket_NoAutorizado(t *testing.T) {
	e := seedEvent(dao.DB)
	s := seedSeat(dao.DB, e.ID)
	propietario := seedUser(dao.DB, "propietario_ticket@test.com", "PROP0001")
	intruso := seedUser(dao.DB, "intruso_ticket@test.com", "INTR0001")
	tickets, _ := PurchaseTickets(propietario.ID, e.ID, []uint{s.ID})

	err := TransferTicket(intruso.ID, tickets[0].ID, propietario.DNI)
	if err == nil {
		t.Error("un usuario no autorizado no debería poder transferir el ticket")
	}
}

func TestTransferTicket_Exitoso(t *testing.T) {
	e := seedEvent(dao.DB)
	s := seedSeat(dao.DB, e.ID)
	userA := seedUser(dao.DB, "xfer_from@test.com", "XFRA0001")
	userB := seedUser(dao.DB, "xfer_to@test.com", "XFRB0001")
	tickets, err := PurchaseTickets(userA.ID, e.ID, []uint{s.ID})
	if err != nil {
		t.Fatalf("PurchaseTickets falló: %v", err)
	}

	err = TransferTicket(userA.ID, tickets[0].ID, userB.DNI)
	if err != nil {
		t.Fatalf("TransferTicket falló: %v", err)
	}
}

func TestCreateVenue_DBNula(t *testing.T) {
	saved := dao.DB
	dao.DB = nil
	defer func() { dao.DB = saved }()

	_, err := CreateVenue(VenueInput{
		Nombre: "Test", Direccion: "Dir", Filas: 5, ColumnasPorFila: 10,
	})
	if err == nil {
		t.Error("se esperaba error con DB nula")
	}
}
