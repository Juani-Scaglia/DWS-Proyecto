FROM golang:1.21-alpine

WORKDIR /app

# Copiar los archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código del backend
COPY . .

# Compilar el binario de Go
RUN go build -o main .

# Exponer el puerto donde corre tu API (Gin)
EXPOSE 8080

# Comando para arrancar el backend
CMD ["./main"]