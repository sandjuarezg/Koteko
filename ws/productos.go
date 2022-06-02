package ws

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sandjuarezg/koteko/models"
)

func ProductosWS(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer fmt.Printf("Response from %s:%s\n", r.URL.RequestURI(), r.Method)

		switch r.Method {
		case "GET":
			if r.Header.Get("id") == "" {
				products, err := models.GetAllProducts(db)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)

					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				err = json.NewEncoder(w).Encode(products)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)

					return
				}
			} else {
				id, err := strconv.Atoi(r.Header.Get("id"))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)

					return
				}

				product, err := models.GetProductWithAllDetailsByID(db, id)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)

					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				err = json.NewEncoder(w).Encode(product)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)

					return
				}
			}

		case "POST":
			var product models.Producto

			err := json.NewDecoder(r.Body).Decode(&product)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			err = models.CreateNewProduct(product, db)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			err = json.NewEncoder(w).Encode(Message{Body: "successful action"})
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "PUT":
			var product models.Producto

			err := json.NewDecoder(r.Body).Decode(&product)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			err = models.UpdateProductByID(db, product, r.Header.Get("id"))
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
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
			err := models.DeleteProductByID(db, r.Header.Get("id"))
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
