package models

import "database/sql"

func AddVenta(db *sql.DB, idUs, cant int, product Producto) (err error) {
	total := int(product.Precio) * cant
	_, err = db.Exec(
		`INSERT INTO ventas(id_usuario, id_producto, cantidad, total) 
		VALUES (?, ?, ?, ?)
	`, idUs, product.IDProducto, cant, total)
	if err != nil {
		return
	}

	return
}
