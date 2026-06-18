package dao

import (
	"fmt"
	"log"
	"os"

<<<<<<< HEAD
	"backend/domain/models"
=======
	"backend/domain/models" // Importa tu models.go desde la carpeta domain [cite: 30]
>>>>>>> feature/eventos

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB es la variable global que usarán para las consultas en los servicios [cite: 138, 141]
var DB *gorm.DB

func InitDB() {
	// 1. Intentar cargar el archivo .env local [cite: 161]
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: No se encontró el archivo .env, se usarán las variables del sistema")
	}

	// 2. Extraer las variables del entorno
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	// 3. Armar la cadena de conexión (DSN) [cite: 26]
<<<<<<< HEAD
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", 
=======
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
>>>>>>> feature/eventos
		username, password, host, port, dbname,
	)

	// 4. Abrir la conexión con GORM [cite: 30]
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error crítico: No se pudo conectar a MySQL: %v", err)
	}

	log.Println("¡Conexión segura establecida con MySQL en el puerto " + port + "!")

	// 5. AUTOMIGRATE: Lee tu models.go y crea/actualiza las tablas automáticamente [cite: 30]
	err = DB.AutoMigrate(&domain.User{}, &domain.Event{}, &domain.Ticket{})
	if err != nil {
		log.Fatalf("Error al mapear las tablas desde models.go: %v", err)
	}

	log.Println("¡Tablas sincronizadas con éxito en MySQL!")
<<<<<<< HEAD
}
=======
}
>>>>>>> feature/eventos
