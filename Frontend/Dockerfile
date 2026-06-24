FROM node:22-alpine

WORKDIR /app

# Instalar dependencias del proyecto
COPY package*.json ./
RUN npm install

# Copiar el código del frontend
COPY . .

# Exponer el puerto por defecto de Vite
EXPOSE 5173

# Comando para correr el servidor de desarrollo en red abierta para Docker
CMD ["npm", "run", "dev", "--", "--host"]