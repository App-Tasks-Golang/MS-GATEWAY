# Usa una imagen base oficial de Go
FROM golang:1.23-alpine

# Define el directorio de trabajo en el contenedor
WORKDIR /app

# Copia el archivo `go.mod` y `go.sum` para que se instalen las dependencias correctamente
COPY go.mod go.sum ./

# Instala las dependencias necesarias
RUN go mod tidy

# Copia todo el código fuente al contenedor
COPY . .

# Descargar dockerize
RUN wget https://github.com/jwilder/dockerize/releases/download/v0.6.1/dockerize-linux-amd64-v0.6.1.tar.gz \
    && tar -xzvf dockerize-linux-amd64-v0.6.1.tar.gz \
    && mv dockerize /usr/local/bin/ \
    && rm dockerize-linux-amd64-v0.6.1.tar.gz  # Limpiar el archivo comprimido

# Construye el binario de la aplicación desde la carpeta `cmd`
RUN go build -o main

# Expone el puerto 8080 (ajústalo según tu aplicación)
EXPOSE 8083

# Ejecuta el binario cuando inicie el contenedor
CMD ["dockerize", "-wait", "tcp://task-service:8080", "-wait", "tcp://user-service:8082", "-wait", "tcp://auth-service:8084", "-timeout", "30s", "./main"]
