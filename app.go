package main

import (
	"database/sql"
	"fmt"
	"log"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

func (a *App) getRecipe(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	r := recipe{ID: id}
	if err := r.getRecipe(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(writer, http.StatusNotFound, "Recipe not found")
		default:
			respondWithError(writer, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(writer, http.StatusOK, r)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getRecipes(writer http.ResponseWriter, request *http.Request) {
	count, _ := strconv.Atoi(request.FormValue("count"))
	start, _ := strconv.Atoi(request.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getRecipes(a.DB, start, count)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusOK, products)
}

func (a *App) createRecipe(writer http.ResponseWriter, request *http.Request) {
	var r recipe
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&r); err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer request.Body.Close()

	if err := r.createRecipe(a.DB); err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusCreated, r)
}

func (a *App) updateRecipe(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var r recipe
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&r); err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer request.Body.Close()
	r.ID = id

	if err := r.updateRecipe(a.DB); err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusOK, r)
}

func (a *App) deleteRecipe(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	r := recipe{ID: id}
	if err := r.deleteRecipe(a.DB); err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/recipes", a.getRecipes).Methods("GET")
	a.Router.HandleFunc("/recipe", a.createRecipe).Methods("POST")
	a.Router.HandleFunc("/recipe/{id:[0-9]+}", a.getRecipe).Methods("GET")
	a.Router.HandleFunc("/recipe/{id:[0-9]+}", a.updateRecipe).Methods("PUT")
	a.Router.HandleFunc("/recipe/{id:[0-9]+}", a.deleteRecipe).Methods("DELETE")
}
