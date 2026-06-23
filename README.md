# DWS-Proyecto — Sistema de Venta de Tickets

Aplicación web full-stack para la compra, cancelación y transferencia de entradas a eventos. Desarrollada con Go (Gin + GORM) en el backend y React (Vite) en el frontend.

---

## Requisitos previos

| Herramienta | Versión mínima | Descarga |
|-------------|---------------|----------|
| Go          | 1.21+         | https://go.dev/dl |
| Node.js     | 18+           | https://nodejs.org |
| MySQL       | 8.0+          | https://dev.mysql.com/downloads |
| Git         | cualquiera    | https://git-scm.com |

---

## Estructura del proyecto

```
DWS-Proyecto/
├── Backend/
│   ├── controllers/      # Handlers HTTP (auth, events, tickets)
│   ├── dao/              # Acceso a la base de datos (GORM)
│   ├── database/         # seed.sql con eventos de prueba
│   ├── domain/models/    # Structs: User, Event, Ticket
│   ├── middlewares/      # JWT AuthMiddleware
│   ├── services/         # Lógica de negocio y transacciones
│   ├── .env.example      # Plantilla de variables de entorno
│   ├── go.mod
│   └── main.go
└── Frontend/
    ├── src/
    ├── package.json
    └── vite.config.js
```

---

## Instalación y ejecución — Backend

### 1. Crear la base de datos en MySQL

```sql
CREATE DATABASE ticket_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

> Las tablas (`users`, `events`, `tickets`) se crean automáticamente cuando el backend arranca por primera vez gracias a **GORM AutoMigrate**.

### 2. Configurar variables de entorno

```bash
cd Backend
cp .env.example .env
```

Editar `.env` con los datos reales:

```env
DB_USER=root
DB_PASSWORD=tu_password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=ticket_system
JWT_SECRET=una_clave_secreta_larga_y_segura
```

> `.env` está en `.gitignore` y **nunca debe commitearse**.

### 3. Instalar dependencias Go

```bash
go mod tidy
```

### 4. (Opcional) Cargar eventos de prueba

```bash
mysql -u root -p ticket_system < database/seed.sql
```

### 5. Ejecutar el servidor

```bash
go run main.go
```

El servidor queda escuchando en `http://localhost:8080`.

---

## Instalación y ejecución — Frontend

### 1. Instalar dependencias

```bash
cd Frontend
npm install
```

### 2. Ejecutar en modo desarrollo

```bash
npm run dev
```

La app queda disponible en `http://localhost:5173`.

### Otros comandos frontend

```bash
npm run build    # Build de producción (genera dist/)
npm run preview  # Preview del build de producción
npm run lint     # Análisis estático con ESLint
```

---

## Endpoints de la API

**Base URL:** `http://localhost:8080/api`

### Públicos

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/auth/register` | Registrar nuevo usuario |
| POST | `/auth/login` | Login, devuelve JWT |
| GET | `/events` | Listar eventos (`?category=Recitales`) |
| GET | `/events/:id` | Detalle de un evento |

### Protegidos (requieren `Authorization: Bearer <token>`)

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/tickets/purchase` | Comprar ticket `{"event_id": 1}` |
| GET | `/tickets/my-tickets` | Ver mis tickets |
| POST | `/tickets/:id/cancel` | Cancelar un ticket |
| POST | `/tickets/:id/transfer` | Transferir por DNI `{"dni": "12345678"}` |

---

## Tests

### Backend

Correr desde la carpeta `Backend/`:

```bash
cd Backend
go test ./middlewares/... ./services/... -v -cover
```

Para un reporte detallado por función:

```bash
go test ./middlewares/... ./services/... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### Frontend

```bash
cd Frontend
npm test
```

---

## Variables de entorno — referencia completa

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| `DB_USER` | Usuario de MySQL | `root` |
| `DB_PASSWORD` | Contraseña de MySQL | `mi_password` |
| `DB_HOST` | Host de MySQL | `127.0.0.1` |
| `DB_PORT` | Puerto de MySQL | `3306` |
| `DB_NAME` | Nombre de la base de datos | `ticket_system` |
| `JWT_SECRET` | Clave para firmar tokens JWT | `clave_larga_y_aleatoria` |

---

## Integrantes

| Nombre | Rol |
|--------|-----|
| Juan   | Backend — Auth, Tickets, Transacciones |
| Simón  | Backend — Endpoints protegidos |
| Facu   | Frontend — Vistas y consumo de API |



## 🐳 Instrucciones de Despliegue con Docker

Para levantar la infraestructura completa de la aplicación (Backend en Go, Frontend en React y Base de Datos MySQL) en un entorno aislado y contenedorizado, asegúrese de tener Docker Desktop ejecutándose y ejecute el siguiente comando en la raíz del proyecto:

```bash
docker-compose up --build


## 🧠 Decisiones de Diseño e Ingeniería (Ingeniería Inversa)

Para garantizar la estabilidad y seguridad del sistema en un entorno de producción, se analizaron los módulos ya desarrollados y se documentaron las siguientes soluciones arquitectónicas:

### 1. Gestión de Aforo Dinámico y Consistencia Transaccional (Evitar Sobreventa)
Para mitigar **condiciones de carrera (Race Conditions)** en escenarios de alta concurrencia (por ejemplo, muchos usuarios comprando la última entrada al mismo tiempo), se aplicó **Bloqueo Pesimista (Pessimistic Locking)** a nivel de base de datos utilizando la cláusula `FOR UPDATE` mediante GORM en el módulo de tickets:

```go
// Ejemplo conceptual del bloqueo aplicado en services/tickets.go
tx := db.Begin()
var event models.Event
if err := tx.Raw("SELECT * FROM events WHERE id = ? FOR UPDATE", eventID).Scan(&event).Error; err != nil {
    tx.Rollback()
    return err
}