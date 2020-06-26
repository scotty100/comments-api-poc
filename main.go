package main

import (
	"comments-api/api"
	"context"
	"github.com/looplab/eventhorizon/repo/mongodb"
	"log"
	"net/http"
)

func main() {
	log.Println("starting comments backend")

	h, err := api.NewHandler()

	// NOTE: Temp clear of DB on startup.
	repo, ok := h.Repo.Parent().(*mongodb.Repo)
	if !ok {
		log.Fatal("incorrect repo type")
	}
	if err := repo.Clear(context.Background()); err != nil {
		log.Println("could not clear DB:", err)
	}

	if err != nil {
		log.Fatal(err)
	}
	log.Println(http.ListenAndServe(":8080", h))
}