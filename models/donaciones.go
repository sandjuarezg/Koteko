package models

import "database/sql"

func DoDonation(db *sql.DB, idUs, cant string) (err error) {
	_, err = db.Exec(
		`INSERT INTO donaciones(id_usuario, cantidad) 
		VALUES (?, ?)
	`, idUs, cant)
	if err != nil {
		return
	}

	return
}
