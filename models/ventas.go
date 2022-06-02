package models

import "database/sql"

type Venta struct {
	IDVenta    int
	IDProducto int
	Cantidad   int
	Total      float32
	Fecha      string
}

func CreateNewVenta(db *sql.DB, cant int, product Producto) (err error) {
	total := int(product.Precio) * cant

	_, err = db.Exec(
		`INSERT INTO ventas(id_producto, cantidad, total) 
		VALUES (?, ?, ?)
	`, product.IDProducto, cant, total)
	if err != nil {
		return
	}

	return
}

func GetAllVentas(db *sql.DB) (ventas []Venta, err error) {
	rows, err := db.Query(`
		SELECT id_venta, id_producto, cantidad, total, fecha 
			FROM ventas`,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		return
	}

	var venta Venta
	for rows.Next() {
		err = rows.Scan(
			&venta.IDVenta,
			&venta.IDProducto,
			&venta.Cantidad,
			&venta.Total,
			&venta.Fecha,
		)
		if err != nil {
			return
		}

		ventas = append(ventas, venta)
	}

	return
}

func DeleteVentaByID(db *sql.DB, id string) (err error) {
	row, err := db.Exec("DELETE from ventas WHERE id_venta = ?", id)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}

func UpdateVentaByID(db *sql.DB, id string, cantidad int, product Producto) (err error) {
	total := int(product.Precio) * cantidad

	row, err := db.Exec(`
		UPDATE ventas 
			SET id_producto = ?, cantidad = ?, total = ? 
			WHERE id_venta = ?
		`, product.IDProducto, cantidad, total, id)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}
