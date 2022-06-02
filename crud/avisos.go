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

func Avisos(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer fmt.Printf("Response from %s\n", r.URL.RequestURI())

		path := strings.TrimPrefix(r.URL.String(), "/avisos?accion=")
		split := strings.Split(path, "&")
		accion := split[0]

		switch accion {
		case "add":
			tarea := r.FormValue("aviso")

			err := models.CreateNewAviso(db, tarea)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "delete":
			split := strings.Split(split[1], "=")
			id := split[1]

			err := models.DeleteAvisoByID(db, id)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "update":
			split := strings.Split(split[1], "=")
			id := split[1]

		}

		avisos, err := models.GetAllAvisos(db)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		temp, err := template.ParseFiles("./admin/aviso.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = temp.Execute(w, avisos)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	})
}
