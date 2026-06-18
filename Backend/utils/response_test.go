package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestErrorResponse_EstadoYMensaje(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ErrorResponse(c, http.StatusBadRequest, "algo salió mal")

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado %d, obtenido %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "algo salió mal" {
		t.Errorf("mensaje incorrecto: %s", body["error"])
	}
}

func TestErrorResponse_CodigosHTTP(t *testing.T) {
	codes := []int{
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusConflict,
		http.StatusInternalServerError,
	}
	for _, code := range codes {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		ErrorResponse(c, code, "error")
		if w.Code != code {
			t.Errorf("esperado %d, obtenido %d", code, w.Code)
		}
	}
}

func TestSuccessResponse_Datos(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SuccessResponse(c, http.StatusCreated, gin.H{"id": 1, "estado": "activo"})

	if w.Code != http.StatusCreated {
		t.Errorf("esperado %d, obtenido %d", http.StatusCreated, w.Code)
	}
	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["estado"] != "activo" {
		t.Errorf("campo 'estado' incorrecto en respuesta")
	}
}

func TestSuccessResponse_200(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SuccessResponse(c, http.StatusOK, gin.H{"items": []int{1, 2, 3}})

	if w.Code != http.StatusOK {
		t.Errorf("esperado %d, obtenido %d", http.StatusOK, w.Code)
	}
}
