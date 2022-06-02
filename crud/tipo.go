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

func Tipo(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer fmt.Printf("Response from %s\n", r.URL.RequestURI())

		path := strings.TrimPrefix(r.URL.String(), "/tipo?accion=")
		split := strings.Split(path, "&")
		accion := split[0]

		switch accion {
		case "add":
			tipo := r.FormValue("tipo")

			err := models.CreateNewTipo(db, tipo)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "delete":
			split := strings.Split(split[1], "=")
			id := split[1]

			err := models.DeleteTipoByID(db, id)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "update":
			split := strings.Split(split[1], "=")
			id := split[1]

			editTipo := r.FormValue("editTipo")

			err := models.UpdateTipoByID(db, id, editTipo)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
		}

		tipos, err := models.GetAllTipos(db)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		temp, err := template.ParseFiles("./admin/tipo.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = temp.Execute(w, tipos)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	})
}
