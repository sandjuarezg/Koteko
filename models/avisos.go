package models

import (
	"database/sql"
	"math/rand"
)

type Aviso struct {
	IDAviso int
	Aviso   string
}

func GetRandomAviso(db *sql.DB) (aviso string, err error) {
	row := db.QueryRow("SELECT count(*) FROM avisos")
	if err != nil {
		return
	}

	var id int

	err = row.Scan(&id)
	if err != nil {
		return
	}

	row = db.QueryRow("SELECT aviso FROM avisos where id_aviso = ?", rand.Intn(id))
	if err != nil {
		return
	}

	err = row.Scan(&aviso)
	if err != nil {
		return
	}

	return
}

func GetAllAvisos(db *sql.DB) (avisos []Aviso, err error) {
	rows, err := db.Query("SELECT id_aviso, aviso FROM avisos")
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		return
	}

	var aviso Aviso
	for rows.Next() {
		err = rows.Scan(
			&aviso.IDAviso,
			&aviso.Aviso,
		)
		if err != nil {
			return
		}

		avisos = append(avisos, aviso)
	}

	return
}

func CreateNewAviso(db *sql.DB, text string) (err error) {
	_, err = db.Exec("INSERT INTO avisos(aviso) VALUES (?)", text)
	if err != nil {
		return
	}

	return
}

func DeleteAvisoByID(db *sql.DB, id string) (err error) {
	row, err := db.Exec("DELETE from avisos WHERE id_aviso = ?", id)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}

func UpdateAvisoByID(db *sql.DB, id, newText string) (err error) {
	row, err := db.Exec("UPDATE avisos SET aviso = ? WHERE id_aviso = ?", newText, id)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}
