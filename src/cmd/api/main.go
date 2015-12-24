package main

import (
	"api"
	"core"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
)

func main() {
	mux := http.NewServeMux()
	db, err := core.CreateDatabaseConnection("msvcdir", "msvcdir", "microservicesdirtest")

	if err != nil {
		log.Fatalf("Could not connect to the database. Error was %v", err)
	}

	projectRepository := core.ProjectRepository{db}

	mux.Handle("/projects", api.HomeHandler(&projectRepository))

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}
