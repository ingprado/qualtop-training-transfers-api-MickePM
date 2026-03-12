package main

import (
	"database/sql"
	"log"
	"net/http"
	"os" // Para leer variables de entorno
	"transfers-api/internal/logging"
	"transfers-api/internal/repositories"
	"transfers-api/internal/services"
	"transfers-api/internal/transport"

	_ "github.com/go-sql-driver/mysql" // Importante: Driver de MySQL
)

func main() {

	logger := logging.Logger
	logger.Info("inicio del logger ")

	dsn := os.Getenv("DB_SOURCE")
	// Intentar conectar a la DB
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatalf("Error al abrir db: %v", err)
	}
	defer db.Close()

	// 3. Verificar la conexión
	if err := db.Ping(); err != nil {
		logger.Fatalf("db no responde: %v", err)
	}
	logger.Info("conectado a MySQL exitosamente")
	repo := repositories.NewMySQLRepository(db)
	svc := services.NewTransferService(repo)
	h := transport.NewHandler(svc)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/transfers", h.HandleTransfers) //el como vamos a consumir

	addr := ":8080"
	logger.Infof("Servidor listo en %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// package main

// import (
// 	"transfers-api/internal/config"
// 	"transfers-api/internal/handlers"
// 	"transfers-api/internal/logging"
// 	"transfers-api/internal/repositories"
// 	"transfers-api/internal/services"
// 	"transfers-api/internal/transport"
// 	"transfers-api/internal/version"
// )

// func main() {
// 	// init logger
// 	logger := logging.Logger
// 	logger.Info("logger started")

// 	// init config
// 	cfg := config.ParseFromEnv()
// 	logger.Infof("config loaded: %v", cfg.String())

// 	// init repositories
// 	transfersDB := repositories.NewTransfersMongoDBRepository(cfg.MongoDBConfig)
// 	logger.Info("repositories created")

// 	// init services
// 	transfersService := services.NewTransfersService(cfg.Business, transfersDB)
// 	logger.Infof("services created")

// 	// init handlers
// 	transfersHandler := handlers.NewTransfersHandler(transfersService)
// 	logger.Infof("handlers created")

// 	// init server
// 	server := transport.NewHTTPServer(transfersHandler)
// 	server.MapRoutes()
// 	logger.Infof("server created, running %s@%s", version.AppName, version.Version)

// 	// run server
// 	server.Run(":8080")
// }
