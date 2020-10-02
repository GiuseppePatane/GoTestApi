package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"./model"
	"github.com/julienschmidt/httprouter"
)

func GetUser(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		email, err := strconv.Atoi(param.ByName("email"))
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid product ID")
			return
		}

		u := model.User{Email: email}
		if err := u.Get(db); err != nil {
			switch err {
			case sql.ErrNoRows:
				respondWithError(w, http.StatusNotFound, "Product not found")
			default:
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		respondWithJSON(w, http.StatusOK, u)
	}
}

func GetConfig(s []byte) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

		// prettyCfg, _ := json.MarshalIndent(app.Cfg, "", "  ")

		w.Write(s)
	}
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
