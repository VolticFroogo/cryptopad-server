package handle

import (
	"log"
	"net/http"

	"github.com/VolticFroogo/config"
	v1 "github.com/VolticFroogo/cryptopad-server/api/v1"
	"github.com/gorilla/mux"
)

const (
	configDirectory = "configs/handle.ini"
)

// Config is the config structure.
type Config struct {
	Port, Certificate, Key string
	SSL                    bool
}

// Start begins listening for all incoming requests.
func Start() {
	// Load the config.
	cfg := Config{}
	err := config.Load(configDirectory, &cfg)
	if err != nil {
		log.Print(err)
		return
	}

	// Create a new Mux Router with strict slash.
	r := mux.NewRouter()
	r.StrictSlash(true)

	// Handle v1 of the API.
	v1.Handle(r)

	// Create a new static file server.
	fileServer := http.FileServer(http.Dir("./static/"))

	// Handle all static files with the file server.
	r.PathPrefix("/").Handler(fileServer)

	if cfg.SSL {
		// If we are using SSL encryption (HTTPS):
		log.Printf("Listening for incoming HTTPS requests on port %v.", cfg.Port)

		// Serve TLS using the certificate and key files from the config.
		http.ListenAndServeTLS(":"+cfg.Port, cfg.Certificate, cfg.Key, r)
	} else {
		// Otherwise:
		log.Printf("Listening for incoming HTTP requests on port %v.", cfg.Port)

		// Serve plain HTTP responses.
		http.ListenAndServe(":"+cfg.Port, r)
	}
}
