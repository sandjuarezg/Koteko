package ws

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sandjuarezg/koteko/models"
)

func RolesWS(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer fmt.Printf("Response from %s:%s\n", r.URL.RequestURI(), r.Method)

		switch r.Method {
		case "GET":
			roles, err := models.GetAllRoles(db)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			err = json.NewEncoder(w).Encode(roles)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "POST":
			var rol models.Rol

			err := json.NewDecoder(r.Body).Decode(&rol)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			if r.Header.Get("id") == "" {
				err = models.CreateNewRol(db, rol.Rol)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)

					return
				}
			} else {
				err = models.UpdateRolByID(db, r.Header.Get("id"), rol.Rol)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)

					return
				}
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
			err := models.DeleteRolByID(db, r.Header.Get("id"))
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
