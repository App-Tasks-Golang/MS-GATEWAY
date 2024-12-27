package main

import (
	"io"
	"log"
	"net/http"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Service struct {
	Name string
	URL  string
}

func main() {
	router := gin.Default()
	
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Origen del frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"}, // Permitir Authorization
		AllowCredentials:  true,
	}))

	// Configuración de microservicios
	services := map[string]Service{
		"tasks":    {Name: "task-service", URL: "http://task-service:8080"},
		"users":    {Name: "user-service", URL: "http://user-service:8082"},
		"auth":     {Name: "auth-service", URL: "http://auth-service:8084"},
	}

	// Registrar rutas dinámicamente
	for prefix, service := range services {
		router.Any("/"+prefix+"/*path", proxyHandler(service))
	}

	// Correr el Gateway
	router.Run(":8083")
}

// proxyHandler devuelve un handler que redirige al servicio correspondiente
func proxyHandler(service Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Omitir el prefijo del servicio
		target := service.URL + c.Request.URL.Path

		req, err := http.NewRequest(c.Request.Method, target, c.Request.Body)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}

		copyHeaders(c.Request.Header, req.Header)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable", "details": err.Error()})
			return
		}
		defer resp.Body.Close()

		// Leer el cuerpo de la respuesta
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error reading response body", "details": err.Error()})
			return
		}

		// Enviar la respuesta tal cual
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	}
}

// copyHeaders copia los headers de una petición a otra
func copyHeaders(src http.Header, dest http.Header) {
	for key, values := range src {
		for _, value := range values {
			dest.Add(key, value)
		}
	}
}
