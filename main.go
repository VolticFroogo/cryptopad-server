package main

import (
	"log"

	"github.com/VolticFroogo/cryptopad-server/db"
	"github.com/VolticFroogo/cryptopad-server/handle"
)

const (
	dbCfgDir = "configs/db.ini"
)

func main() {
	// Initialise the DB.
	err := db.Init(dbCfgDir)
	if err != nil {
		log.Print(err)
		return
	}

	// Start handling incoming requests.
	go handle.Start()
}
