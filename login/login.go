package login

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/sandjuarezg/koteko/models"
)

func Login(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		defer fmt.Printf("Response from %s\n", r.URL.RequestURI())

		if err := r.ParseForm(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		user, err := models.Login(r.FormValue("email"), r.FormValue("password"), db)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		user, err = models.GetUserByID(db, user.ID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		rol, err := models.GetRolByID(db, user.IDRol)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		var perAdmin []models.Permiso
		var perRol []models.Permiso

		if rol.Rol == "administrador" {
			perRol, err = models.GetPermisosByNameRol(db, "empleado")
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			perAdmin, err = models.GetPermisosByNameRol(db, rol.Rol)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
		} else {
			perRol, err = models.GetPermisosByNameRol(db, rol.Rol)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
		}

		w.Header().Add("session", "true")

		t, err := template.ParseFiles("./admin/index.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		data := struct {
			Name          string
			Rol           string
			PermisosAdmin []models.Permiso
			PermisosRol   []models.Permiso
		}{
			Rol:           rol.Rol,
			PermisosAdmin: perAdmin,
			PermisosRol:   perRol,
			Name:          user.Nombre,
		}

		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	})
}
