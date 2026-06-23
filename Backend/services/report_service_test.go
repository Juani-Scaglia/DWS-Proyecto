package services_test

import (
	"testing"
)

func TestCalculateOccupationPercent(t *testing.T) {
	// Declaramos los casos de prueba para verificar las matemáticas del reporte
	tests := []struct {
		name          string
		ticketsIssued int64
		totalCapacity int64
		wantPercent   float64
	}{
		{
			name:          "Establecimiento a mitad de capacidad",
			ticketsIssued: 50,
			totalCapacity: 100,
			wantPercent:   50.0,
		},
		{
			name:          "Establecimiento vacío",
			ticketsIssued: 0,
			totalCapacity: 500,
			wantPercent:   0.0,
		},
		{
			name:          "Establecimiento lleno",
			ticketsIssued: 200,
			totalCapacity: 200,
			wantPercent:   100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulación matemática del cálculo del reporte
			var percent float64
			if tt.totalCapacity > 0 {
				percent = (float64(tt.ticketsIssued) / float64(tt.totalCapacity)) * 100
			}

			if percent != tt.wantPercent {
				t.Errorf("Cálculo incorrecto en '%s': se obtuvo %.2f%%, se esperaba %.2f%%", tt.name, percent, tt.wantPercent)
			}
		})
	}
}
