package models

import "database/sql"

func DoDonation(db *sql.DB, cant string) (err error) {
	_, err = db.Exec(
		`INSERT INTO donaciones(cantidad) 
		VALUES (?)
	`, cant)
	if err != nil {
		return
	}

	return
}
