package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time" // <--- 1. IMPORTANTE AGREGAR ESTO

	"transfers-api/internal/logging"
	"transfers-api/internal/queue"
	"transfers-api/internal/repositories"
	"transfers-api/internal/services"
	"transfers-api/internal/transport"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	logger := logging.Logger
	logger.Info("Inicio del servicio...")

	// --- MySQL ---
	dsn := os.Getenv("DB_SOURCE")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatalf("Error al abrir db: %v", err)
	}
	defer db.Close()

	logger.Info("Conectado a MySQL exitosamente")

	// --- RabbitMQ Config ---
	rabbitURL := os.Getenv("RABBITMQ_URL")
	var producer *queue.RabbitMQProducer

	// --- Repositorios ---
	mysqlRepo := repositories.NewMySQLRepository(db)
	repo := repositories.NewCachedRepository(mysqlRepo, logger)

	for i := 0; i < 15; i++ {
		producer, err = queue.NewRabbitMQProducer(rabbitURL, "transfer_events")

		if err == nil {
			logger.Info("¡Conexión a RabbitMQ establecida!")
			break
		}

		logger.Warnf("RabbitMQ no está listo, reintentando en 5s... Error: %v", err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		logger.Errorf("No se pudo conectar a RabbitMQ tras 15 intentos: %v", err)
	}

	// --- Inyección de Dependencias ---
	svc := services.NewTransferService(repo, producer)
	h := transport.NewHandler(svc)

	// --- Rutas ---
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/transfers", h.HandleTransfers)

	addr := ":8080"
	logger.Infof("Servidor listo en %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
