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

func Categorias(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer fmt.Printf("Response from %s\n", r.URL.RequestURI())

		path := strings.TrimPrefix(r.URL.String(), "/categorias?accion=")
		split := strings.Split(path, "&")
		accion := split[0]

		switch accion {
		case "add":
			categoria := r.FormValue("categoria")

			err := models.CreateNewCategoria(db, categoria)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "delete":
			split := strings.Split(split[1], "=")
			id := split[1]

			err := models.DeleteCategoriaByID(db, id)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "update":
			split := strings.Split(split[1], "=")
			id := split[1]

			editCat := r.FormValue("editCat")

			err := models.UpdateCategoriaByID(db, id, editCat)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
		}

		categorias, err := models.GetAllCategorias(db)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		temp, err := template.ParseFiles("./admin/categoria.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = temp.Execute(w, categorias)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	})
}
