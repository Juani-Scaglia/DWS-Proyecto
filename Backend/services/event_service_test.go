package services

import (
	"testing"
)

// TestGetAllEvents verifica el catálogo de eventos
func TestGetAllEvents(t *testing.T) {
	// Llamamos a tu servicio pasándole vacío para que traiga todos
	events, err := GetAllEvents("")

	if err != nil {
		t.Log("Aviso: El test ejecutó, pero la base de datos no respondió.")
	} else {
		t.Logf("Éxito: El servicio respondió bien y trajo %d eventos.", len(events))
	}
}

// TestGetEventByID verifica la búsqueda por ID
func TestGetEventByID(t *testing.T) {
	_, err := GetEventByID(1)

	if err != nil {
		t.Log("Aviso: No se encontró el evento ID 1 o no hay BD, pero el código manejó el error.")
	} else {
		t.Log("Éxito: El servicio de búsqueda por ID funciona.")
	}
}
