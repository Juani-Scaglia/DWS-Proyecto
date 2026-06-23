package dao

import (
	"fmt"
	"log"
	"os"

	domain "backend/domain/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

const errDBNula = "base de datos no inicializada"

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: No se encontró el archivo .env, se usarán las variables del sistema")
	}

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error crítico: No se pudo conectar a MySQL: %v", err)
	}

	log.Println("¡Conexión segura establecida con MySQL en el puerto " + port + "!")

	err = DB.AutoMigrate(&domain.User{}, &domain.Venue{}, &domain.Event{}, &domain.Ticket{}, &domain.Seat{})
	if err != nil {
		log.Fatalf("Error al mapear las tablas desde models.go: %v", err)
	}

	log.Println("¡Tablas sincronizadas con éxito en MySQL!")
}
