package main

import (
	"exoplanet-service/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/exoplanets", handlers.AddExoplanetHandler).Methods("POST")
	r.HandleFunc("/exoplanets", handlers.ListExoplanetsHandler).Methods("GET")
	r.HandleFunc("/exoplanets/{id}/fuel", handlers.FuelEstimationHandler).Methods("GET")
	r.HandleFunc("/exoplanets/{id}", handlers.GetExoplanetHandler).Methods("GET")
	r.HandleFunc("/exoplanets/{id}", handlers.UpdateExoplanetHandler).Methods("PUT")
	r.HandleFunc("/exoplanets/{id}", handlers.DeleteExoplanetHandler).Methods("DELETE")

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
