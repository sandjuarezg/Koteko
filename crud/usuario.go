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

func Usuario(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer fmt.Printf("Response from %s\n", r.URL.RequestURI())

		path := strings.TrimPrefix(r.URL.String(), "/colores?accion=")
		split := strings.Split(path, "&")
		accion := split[0]

		switch accion {
		case "add":
			color := r.FormValue("color")

			err := models.CreateNewColor(db, color)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "delete":
			split := strings.Split(split[1], "=")
			id := split[1]

			err := models.DeleteColorByID(db, id)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "update":
			split := strings.Split(split[1], "=")
			id := split[1]

			editColor := r.FormValue("editColor")

			err := models.UpdateColorByID(db, id, editColor)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
		}

		colores, err := models.GetAllColores(db)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		temp, err := template.ParseFiles("./admin/colores.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = temp.Execute(w, colores)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	})
}
