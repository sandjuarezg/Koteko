package ws

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sandjuarezg/koteko/models"
)

func TiposWS(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer fmt.Printf("Response from %s:%s\n", r.URL.RequestURI(), r.Method)

		switch r.Method {
		case "GET":
			tipos, err := models.GetAllTipos(db)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			err = json.NewEncoder(w).Encode(tipos)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "POST":
			var tipo models.Tipo

			err := json.NewDecoder(r.Body).Decode(&tipo)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			err = models.CreateNewTipo(db, tipo.Tipo)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			err = json.NewEncoder(w).Encode(Message{Body: "successful action"})
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "PUT":
			var tipo models.Tipo

			err := json.NewDecoder(r.Body).Decode(&tipo)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			err = models.UpdateTipoByID(db, r.Header.Get("id"), tipo.Tipo)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			err = json.NewEncoder(w).Encode(Message{Body: "successful action"})
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "DELETE":
			err := models.DeleteTipoByID(db, r.Header.Get("id"))
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			err = json.NewEncoder(w).Encode(Message{Body: "data delete successfully"})
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
		}
	})
}
