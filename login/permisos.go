package login

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/sandjuarezg/koteko/models"
)

func Permiso(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer fmt.Printf("Response from %s\n", r.URL.RequestURI())

		view := strings.TrimPrefix(r.URL.String(), "/permiso?per=")

		temp, err := template.ParseFiles(fmt.Sprintf("./admin/%s.html", view))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		switch view {
		case "aviso":
			avisos, err := models.GetAllAvisos(db)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			err = temp.Execute(w, avisos)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
		case "categoria":

		case "colores":

		case "producto":

		case "rol":

		case "tipo":

		case "usuario":

		}

		return
	})
}
