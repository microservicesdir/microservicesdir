package main

import (
	"api"
	"core"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
)

func main() {
	mux := http.NewServeMux()
	p := core.Project{}
	mux.Handle("/projects", api.HomeHandler(&p))

	log.Info(p.GetAllProjects())
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}
