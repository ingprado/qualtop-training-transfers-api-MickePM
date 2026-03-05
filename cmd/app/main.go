package main

import (
	"log"
	"net/http"

	"transfers-api/internal/logging"
	"transfers-api/internal/repositories"
	"transfers-api/internal/services"
	"transfers-api/internal/transport"
	"transfers-api/internal/version"
)

func main() {

	logger := logging.Logger
	logger.Info("inicio del logger ")

	// init repositories podemos tener varios repositories
	repo := repositories.NewMySQLRepository()
	logger.Info("repository created")
	// init services
	svc := services.NewTransferService(repo)
	logger.Infof("service created")

	h := transport.NewHandler(svc)
	logger.Infof("handler created")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/transfers", h.HandleTransfers) //el como vamos a consumir

	addr := ":8080"
	logger.Infof("server starting on %s %s@%s", addr, version.AppName, version.Version)
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
