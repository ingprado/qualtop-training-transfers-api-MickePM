package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	// Añadido para el manejo de esperas
	"transfers-api/internal/logging"
	"transfers-api/internal/repositories"
	"transfers-api/internal/services"
	"transfers-api/internal/transport"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	logger := logging.Logger
	logger.Info("inicio del logger ")

	dsn := os.Getenv("DB_SOURCE")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatalf("Error al abrir db: %v", err)
	}
	defer db.Close()

	logger.Info("conectado a MySQL exitosamente")

	mysqlRepo := repositories.NewMySQLRepository(db)

	repo := repositories.NewCachedRepository(mysqlRepo)

	// --- FIN CONFIGURACIÓN ---

	svc := services.NewTransferService(repo)
	h := transport.NewHandler(svc)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/transfers", h.HandleTransfers)

	addr := ":8080"
	logger.Infof("Servidor listo en %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
