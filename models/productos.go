package models

import (
	"database/sql"
)

// Producto structure for Producto
type Producto struct {
	IDProducto  int
	Producto    string
	Descripcion string
	Precio      float32
	Cantidad    int
	Img         string
	Tipo        string
	Categoria   string
	Colores     []string
}

// GetAllProducts Get all products to display the listing on the page
//  @param1 (db): database pointer
//
//  @return1 (products): slice of products
//  @return2 (err): error variable
func GetAllProducts(db *sql.DB) (products []Producto, err error) {
	rows, err := db.Query("SELECT id_producto, producto, descripcion, precio, img FROM productos")
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		return
	}

	var product Producto
	for rows.Next() {
		err = rows.Scan(
			&product.IDProducto,
			&product.Producto,
			&product.Descripcion,
			&product.Precio,
			&product.Img,
		)
		if err != nil {
			return
		}

		products = append(products, product)
	}

	return
}

func GetProductsByCategory(db *sql.DB, category string, limit bool) (products []Producto, err error) {
	sql := `
		SELECT id_producto, producto, descripcion, precio, img 
		FROM productos
		JOIN categorias on productos.id_categoria = categorias.id_categoria
		WHERE categoria =  ? 
	`

	if limit {
		sql += "limit 4"
	}

	rows, err := db.Query(sql, category)
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		return
	}

	var product Producto
	for rows.Next() {
		err = rows.Scan(
			&product.IDProducto,
			&product.Producto,
			&product.Descripcion,
			&product.Precio,
			&product.Img,
		)
		if err != nil {
			return
		}

		products = append(products, product)
	}

	return
}

// GetTopProducts Get top products to display the listing on the page
//  @param1 (db): database pointer
//
//  @return1 (products): slice of products
//  @return2 (err): error variable
func GetTopProducts(db *sql.DB) (products []Producto, err error) {
	rows, err := db.Query("SELECT id_producto, producto, precio, img FROM productos limit 4")
	if err != nil {
		return
	}
	defer rows.Close()

	var product Producto

	for rows.Next() {
		err = rows.Scan(&product.IDProducto, &product.Producto, &product.Precio, &product.Img)
		if err != nil {
			return
		}

		products = append(products, product)
	}

	return
}

// GetProductWithAllDetailsByID Get product rith all details by id
//  @param1 (db): database pointer
//  @param2 (id): id of product
//
//  @return1 (product): product
//  @return2 (err): error variable
func GetProductWithAllDetailsByID(db *sql.DB, id int) (product Producto, err error) {
	row := db.QueryRow(`
	SELECT id_producto, producto, descripcion, precio, cantidad, img , categoria, tipo 
		FROM productos
		INNER JOIN categorias on productos.id_categoria = categorias.id_categoria
		INNER JOIN tipos on productos.id_tipo = tipos.id_tipo
		WHERE id_producto = ?
	`, id)

	err = row.Scan(
		&product.IDProducto,
		&product.Producto,
		&product.Descripcion,
		&product.Precio,
		&product.Cantidad,
		&product.Img,
		&product.Categoria,
		&product.Tipo,
	)
	if err != nil {
		return
	}

	rows, err := db.Query(`
	SELECT color FROM productos_colores
		INNER JOIN productos on productos.id_producto = productos_colores.id_producto
		INNER JOIN colores on colores.id_color = productos_colores.id_color
		WHERE productos.id_producto = ?
	`, id)
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		return
	}

	for rows.Next() {
		var color string

		err = rows.Scan(&color)
		if err != nil {
			return
		}

		product.Colores = append(product.Colores, color)
	}

	return
}

// SearchProductWithAllDetailsByName Search product with all details by name
//  @param1 (db): database pointer
//  @param2 (name): name of product
//
//  @return1 (product): product
//  @return2 (err): error variable
func SearchProductWithAllDetailsByName(db *sql.DB, name string) (products []Producto, err error) {
	rows, err := db.Query("SELECT id_producto, producto, descripcion, precio, img FROM productos WHERE producto like '%" + name + "%'")
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		return
	}

	for rows.Next() {
		var product Producto

		err = rows.Scan(
			&product.IDProducto,
			&product.Producto,
			&product.Descripcion,
			&product.Precio,
			&product.Img,
		)
		if err != nil {
			return
		}

		products = append(products, product)
	}

	return
}

func AddProduct() {

}

func EditProduct() {

}

func DeleteProduct() {

}

func BuyProducts(db *sql.DB, id, cantidad int) (err error) {
	row, err := db.Exec("UPDATE productos SET cantidad = (cantidad - ?) WHERE id_producto = ?", cantidad, id)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}
