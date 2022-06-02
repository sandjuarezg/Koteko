package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sandjuarezg/koteko/crud"
	"github.com/sandjuarezg/koteko/login"
	"github.com/sandjuarezg/koteko/models"
	"github.com/sandjuarezg/koteko/ws"
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
	http.Handle("/buy", buy(db))
	http.Handle("/donation", donation(db))

	// crud
	http.Handle("/avisos", crud.Avisos(db))

	// ws
	http.Handle("/avisosWS", ws.AvisosWS(db))
	http.Handle("/categoriasWS", ws.CategoriasWS(db))
	http.Handle("/coloresWS", ws.ColoresWS(db))
	http.Handle("/permisosWS", ws.PermisosWS(db))
	http.Handle("/productosWS", ws.ProductosWS(db))
	http.Handle("/rolesWS", ws.RolesWS(db))
	http.Handle("/tiposWS", ws.TiposWS(db))
	http.Handle("/usuariosWS", ws.UsuariosWS(db))
	http.Handle("/ventasWS", ws.VentasWS(db))

	//login
	http.Handle("/login", login.Login(db))
	http.Handle("/permiso", login.Permiso(db))

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

		err := models.CreateNewUser(us, db)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		t, err := template.ParseFiles("./public/html/account.html")
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

		err = models.CreateNewVenta(db, cantidad, product)
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

		err := models.DoDonation(db, cantidad)
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
