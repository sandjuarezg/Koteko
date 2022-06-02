package ws

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sandjuarezg/koteko/models"
)

func ColoresWS(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer fmt.Printf("Response from %s:%s\n", r.URL.RequestURI(), r.Method)

		switch r.Method {
		case "GET":
			colores, err := models.GetAllColores(db)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			err = json.NewEncoder(w).Encode(colores)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "POST":
			var color models.Color

			err := json.NewDecoder(r.Body).Decode(&color)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			err = models.CreateNewColor(db, color.Color)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			err = json.NewEncoder(w).Encode(Message{Body: "successful insert"})
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "PUT":
			var color models.Color

			err := json.NewDecoder(r.Body).Decode(&color)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			err = models.UpdateColorByID(db, r.Header.Get("id"), color.Color)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			err = json.NewEncoder(w).Encode(Message{Body: "successful update"})
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "DELETE":
			err := models.DeleteColorByID(db, r.Header.Get("id"))
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
