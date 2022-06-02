package models

import (
	"database/sql"
)

type Color struct {
	IDColor int
	Color   string
}

func GetAllColores(db *sql.DB) (colores []Color, err error) {
	rows, err := db.Query("SELECT id_color, color FROM colores")
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		return
	}

	var color Color
	for rows.Next() {
		err = rows.Scan(
			&color.IDColor,
			&color.Color,
		)
		if err != nil {
			return
		}

		colores = append(colores, color)
	}

	return
}

func CreateNewColor(db *sql.DB, text string) (err error) {
	_, err = db.Exec("INSERT INTO colores(color) VALUES (?)", text)
	if err != nil {
		return
	}

	return
}

func DeleteColorByID(db *sql.DB, id string) (err error) {
	row, err := db.Exec("DELETE from colores WHERE id_color = ?", id)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}

func UpdateColorByID(db *sql.DB, id, newText string) (err error) {
	row, err := db.Exec("UPDATE colores SET color = ? WHERE id_color = ?", newText, id)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}
