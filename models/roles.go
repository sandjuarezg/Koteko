package models

import (
	"database/sql"
)

type Rol struct {
	IDRol int
	Rol   string
}

func GetAllRoles(db *sql.DB) (roles []Rol, err error) {
	rows, err := db.Query("SELECT id_rol, rol FROM roles")
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		return
	}

	var rol Rol
	for rows.Next() {
		err = rows.Scan(
			&rol.IDRol,
			&rol.Rol,
		)
		if err != nil {
			return
		}

		roles = append(roles, rol)
	}

	return
}

func GetRolByID(db *sql.DB, id int) (rol Rol, err error) {
	row := db.QueryRow("SELECT id_rol, rol FROM roles where id_rol = ?", id)
	if row.Err() != nil {
		return
	}

	err = row.Scan(
		&rol.IDRol,
		&rol.Rol,
	)
	if err != nil {
		return
	}

	return
}

func CreateNewRol(db *sql.DB, text string) (err error) {
	_, err = db.Exec("INSERT INTO roles(rol) VALUES (?)", text)
	if err != nil {
		return
	}

	return
}

func DeleteRolByID(db *sql.DB, id string) (err error) {
	row, err := db.Exec("DELETE from roles WHERE id_rol = ?", id)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}

func UpdateRolByID(db *sql.DB, id, newText string) (err error) {
	row, err := db.Exec("UPDATE roles SET rol = ? WHERE id_rol = ?", newText, id)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}
