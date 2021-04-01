package router

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"scratch-test/api/config"
	"scratch-test/api/middlewares"
	"scratch-test/api/utils"
	"scratch-test/models"
)

func handler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	filter, limit := utils.FilterURLQuery(queries)

	entries, err := config.Database.Read(filter, limit)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := fmt.Fprintf(w, "%s", err.Error()); err != nil {
			log.Println(err)
		}
		return
	}

	clinics, err := models.GetClinics(entries)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(clinics)
	if err != nil {
		log.Println(err)
	}
}

func setupRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/", middlewares.SetLogger(middlewares.SetJSON(handler))).Methods(http.MethodGet)
	return r
}

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	return setupRoutes(r)
}
