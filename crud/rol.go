package crud

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/sandjuarezg/koteko/models"
)

func Rol(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer fmt.Printf("Response from %s\n", r.URL.RequestURI())

		path := strings.TrimPrefix(r.URL.String(), "/rol?accion=")
		split := strings.Split(path, "&")
		accion := split[0]

		switch accion {
		case "add":
			rol := r.FormValue("rol")

			err := models.CreateNewRol(db, rol)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "delete":
			split := strings.Split(split[1], "=")
			id := split[1]

			err := models.DeleteRolByID(db, id)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "update":
			split := strings.Split(split[1], "=")
			id := split[1]

			editRol := r.FormValue("editRol")

			err := models.UpdateRolByID(db, id, editRol)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
		}

		roles, err := models.GetAllRoles(db)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		temp, err := template.ParseFiles("./admin/rol.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = temp.Execute(w, roles)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	})
}
