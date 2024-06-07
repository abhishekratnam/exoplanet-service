package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ExoplanetType string

const (
	GasGiant    ExoplanetType = "GasGiant"
	Terrestrial ExoplanetType = "Terrestrial"
)

type Exoplanet struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Distance    int           `json:"distance"`
	Radius      float64       `json:"radius"`
	Mass        float64       `json:"mass,omitempty"`
	Type        ExoplanetType `json:"type"`
}

var exoplanets = make(map[string]Exoplanet)

func (e Exoplanet) Validate() error {
	if e.Name == "" || e.Description == "" {
		return errors.New("name and description are required")
	}
	if e.Distance <= 10 || e.Distance >= 1000 {
		return errors.New("distance must be between 10 and 1000 light years")
	}
	if e.Radius <= 0.1 || e.Radius >= 10 {
		return errors.New("radius must be between 0.1 and 10 Earth-radius units")
	}
	if e.Type == Terrestrial && (e.Mass <= 0.1 || e.Mass >= 10) {
		return errors.New("mass must be between 0.1 and 10 Earth-mass units for terrestrial planets")
	}
	return nil
}

func (e Exoplanet) Gravity() float64 {
	if e.Type == GasGiant {
		return 0.5 / (e.Radius * e.Radius)
	}
	return e.Mass / (e.Radius * e.Radius)
}

func FuelEstimation(distance int, gravity float64, crewCapacity int) float64 {
	return float64(distance) / (gravity * gravity) * float64(crewCapacity)
}

func AddExoplanetHandler(w http.ResponseWriter, r *http.Request) {
	var exoplanet Exoplanet
	if err := json.NewDecoder(r.Body).Decode(&exoplanet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exoplanet.ID = uuid.New().String()
	if err := exoplanet.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exoplanets[exoplanet.ID] = exoplanet
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exoplanet)
}

func ListExoplanetsHandler(w http.ResponseWriter, r *http.Request) {
	var list []Exoplanet
	for _, v := range exoplanets {
		list = append(list, v)
	}
	json.NewEncoder(w).Encode(list)
}

func GetExoplanetHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	exoplanet, exists := exoplanets[id]
	if !exists {
		http.Error(w, "exoplanet not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(exoplanet)
}

func UpdateExoplanetHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var exoplanet Exoplanet
	if err := json.NewDecoder(r.Body).Decode(&exoplanet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, exists := exoplanets[id]; !exists {
		http.Error(w, "exoplanet not found", http.StatusNotFound)
		return
	}

	exoplanet.ID = id
	if err := exoplanet.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exoplanets[id] = exoplanet
	json.NewEncoder(w).Encode(exoplanet)
}

func DeleteExoplanetHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if _, exists := exoplanets[id]; !exists {
		http.Error(w, "exoplanet not found", http.StatusNotFound)
		return
	}

	delete(exoplanets, id)
	w.WriteHeader(http.StatusNoContent)
}

func FuelEstimationHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	crewCapacity := r.URL.Query().Get("crew")
	if crewCapacity == "" {
		http.Error(w, "crew capacity is required", http.StatusBadRequest)
		return
	}
	crew, err := strconv.Atoi(crewCapacity)
	if err != nil {
		http.Error(w, "invalid crew capacity", http.StatusBadRequest)
		return
	}

	exoplanet, exists := exoplanets[id]
	if !exists {
		http.Error(w, "exoplanet not found", http.StatusNotFound)
		return
	}

	gravity := exoplanet.Gravity()
	fuel := FuelEstimation(exoplanet.Distance, gravity, crew)
	json.NewEncoder(w).Encode(map[string]float64{"fuel_estimation": fuel})
}
