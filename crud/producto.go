package crud

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/sandjuarezg/koteko/models"
)

func Producto(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer fmt.Printf("Response from %s\n", r.URL.RequestURI())

		path := strings.TrimPrefix(r.URL.String(), "/productoCRUD?accion=")
		split := strings.Split(path, "&")
		accion := split[0]

		switch accion {
		case "add":
			var produc models.Producto

			produc.Producto = r.FormValue("produc")
			produc.Descripcion = r.FormValue("desc")

			precio, err := strconv.ParseFloat(r.FormValue("precio"), 32)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			produc.Precio = precio
			produc.Cantidad = 100
			produc.Tipo = "1"
			produc.Categoria = "1"

			err = models.CreateNewProduct(produc, db)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "delete":
			split := strings.Split(split[1], "=")
			id := split[1]

			err := models.DeleteProductByID(db, id)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "update":
			split := strings.Split(split[1], "=")
			id := split[1]

			var product models.Producto

			product.Producto = r.FormValue("editProduc")
			product.Descripcion = r.FormValue("editDesc")

			precio, err := strconv.ParseFloat(r.FormValue("editPrecio"), 32)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			product.Precio = precio
			product.Cantidad = 100
			product.Tipo = "1"
			product.Categoria = "1"

			err = models.UpdateProductByID(db, product, id)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
		}

		product, err := models.GetAllProducts(db)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		temp, err := template.ParseFiles("./admin/producto.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = temp.Execute(w, product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	})
}
