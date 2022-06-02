package models

import (
	"database/sql"
)

type Tipo struct {
	IDTipo int
	Tipo   string
}

func GetAllTipos(db *sql.DB) (tipos []Tipo, err error) {
	rows, err := db.Query("SELECT id_tipo, tipo FROM tipos")
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		return
	}

	var tipo Tipo
	for rows.Next() {
		err = rows.Scan(
			&tipo.IDTipo,
			&tipo.Tipo,
		)
		if err != nil {
			return
		}

		tipos = append(tipos, tipo)
	}

	return
}

func CreateNewTipo(db *sql.DB, text string) (err error) {
	_, err = db.Exec("INSERT INTO tipos(tipo) VALUES (?)", text)
	if err != nil {
		return
	}

	return
}

func DeleteTipoByID(db *sql.DB, id string) (err error) {
	row, err := db.Exec("DELETE from tipos WHERE id_tipo = ?", id)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}

func UpdateTipoByID(db *sql.DB, id, newText string) (err error) {
	row, err := db.Exec("UPDATE tipos SET tipo = ? WHERE id_tipo = ?", newText, id)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}
