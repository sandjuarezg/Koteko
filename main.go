package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sandjuarezg/koteko/models"
)

const (
	user     = "koteko"
	password = "1234"
	host     = "localhost"
	port     = "3306"
	db       = "koteko"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, db))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/index", index(db))
	http.Handle("/producto", producto(db))
	http.Handle("/search", search(db))
	http.Handle("/categoria", categoria(db))
	http.Handle("/signin", signin(db))
	http.Handle("/login", login(db))
	http.Handle("/buy", buy(db))
	http.Handle("/donation", donation(db))
	http.Handle("/avisos", avisos(db))

	http.Handle("/avisosWS", avisosWS(db))

	fmt.Println("Listening on localhost:8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Print(err)

		return
	}
}

func index(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer fmt.Printf("Response from %s\n", r.URL.RequestURI())

		temp, err := template.ParseFiles("./public/index.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		aviso, err := models.GetRandomAviso(db)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		products, err := models.GetTopProducts(db)
		if err != nil {
			log.Fatal(err)
		}

		data := struct {
			Products []models.Producto
			Aviso    string
		}{
			Products: products,
			Aviso:    aviso,
		}

		err = temp.Execute(w, data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	})
}

func search(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		defer fmt.Printf("Response from %s\n", r.URL.RequestURI())

		if err := r.ParseForm(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		products, err := models.SearchProductWithAllDetailsByName(db, r.FormValue("busqueda"))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		aviso, err := models.GetRandomAviso(db)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		data := struct {
			Products []models.Producto
			Aviso    string
		}{
			Products: products,
			Aviso:    aviso,
		}

		temp, err := template.ParseFiles("./public/html/list_products.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = temp.Execute(w, data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	})
}

func categoria(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		defer fmt.Printf("Response from %s\n", r.URL.RequestURI())

		aviso, err := models.GetRandomAviso(db)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		category := strings.TrimPrefix(r.URL.String(), "/categoria?categoria=")
		var products []models.Producto

		if category == "todos" {
			products, err = models.GetAllProducts(db)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
		} else {
			products, err = models.GetProductsByCategory(db, category, false)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
		}

		data := struct {
			Products []models.Producto
			Aviso    string
		}{
			Products: products,
			Aviso:    aviso,
		}

		temp, err := template.ParseFiles("./public/html/list_products.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = temp.Execute(w, data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	})
}

func producto(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer fmt.Printf("Response from %s\n", r.URL.RequestURI())

		temp, err := template.ParseFiles("./public/html/product.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		id, err := strconv.Atoi(strings.Trim(r.URL.String(), "/producto?id="))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		product, err := models.GetProductWithAllDetailsByID(db, id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		products, err := models.GetProductsByCategory(db, product.Categoria, true)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		data := struct {
			Products []models.Producto
			Product  models.Producto
		}{
			Products: products,
			Product:  product,
		}

		err = temp.Execute(w, data)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	})
}

func signin(db *sql.DB) http.Handler {
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

		us := models.Usuario{
			Nombre:   r.FormValue("name"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		err := models.Signin(us, db)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	})
}

func login(db *sql.DB) http.Handler {
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

		fmt.Println(user)

		return
	})
}

func buy(db *sql.DB) http.Handler {
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

		cantidad, err := strconv.Atoi(r.FormValue("cantidad"))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		id, err := strconv.Atoi(strings.Trim(r.URL.String(), "/buy?id="))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = models.BuyProducts(db, id, cantidad)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		product, err := models.GetProductWithAllDetailsByID(db, id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = models.AddVenta(db, 1, cantidad, product)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = models.GenerateBuyPDF(product, cantidad)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		products, err := models.GetProductsByCategory(db, product.Categoria, true)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		data := struct {
			Products []models.Producto
			Product  models.Producto
		}{
			Products: products,
			Product:  product,
		}

		t, err := template.ParseFiles("./public/html/product.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
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

func donation(db *sql.DB) http.Handler {
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

		cantidad := r.FormValue("cantidad")

		err := models.DoDonation(db, "1", cantidad)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		t, err := template.ParseFiles("./public/html/donation.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = t.Execute(w, nil)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	})
}

func avisos(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		defer fmt.Printf("Response from %s\n", r.URL.RequestURI())

		avisos, err := models.GetAllAvisos(db)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		temp, err := template.ParseFiles("./public/admin/avisos.html")
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

func avisosWS(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			avisos, err := models.GetAllAvisos(db)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			data, err := json.Marshal(avisos)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err = io.WriteString(w, string(data))
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case "POST":
			var aviso models.Aviso

			err := json.NewDecoder(r.Body).Decode(&aviso)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			fmt.Println(aviso)

		case "DELETE":
		}
	})
}

// GET    		- http://localhost:8080/avisosWS
// POST   		- http://localhost:8080/avisosWS
// POST (edit) 	- http://localhost:8080/avisosWS?id=X
// DELETE 		- http://localhost:8080/avisosWS?id=X
