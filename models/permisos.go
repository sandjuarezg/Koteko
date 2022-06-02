package models

import (
	"database/sql"
)

type Permiso struct {
	IDPermiso int
	Permiso   string
}

func GetAllPermisos(db *sql.DB) (permisos []Permiso, err error) {
	rows, err := db.Query("SELECT id_permiso, permiso FROM permisos")
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		return
	}

	var permiso Permiso
	for rows.Next() {
		err = rows.Scan(
			&permiso.IDPermiso,
			&permiso.Permiso,
		)
		if err != nil {
			return
		}

		permisos = append(permisos, permiso)
	}

	return
}

func GetPermisosByIDRol(db *sql.DB, id int) (permisos []Permiso, err error) {
	rows, err := db.Query(`
		SELECT permiso 
			FROM permisos 
			join roles_permiso on roles_permiso.id_permiso = permisos.id_permiso
			WHERE id_rol = ?
	`, id)
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		return
	}

	var permiso Permiso
	for rows.Next() {
		err = rows.Scan(
			&permiso.Permiso,
		)
		if err != nil {
			return
		}

		permisos = append(permisos, permiso)
	}

	return
}

func GetPermisosByNameRol(db *sql.DB, rol string) (permisos []Permiso, err error) {
	rows, err := db.Query(`
		SELECT permiso 
			FROM roles_permiso 
			join permisos on roles_permiso.id_permiso = permisos.id_permiso 
			join roles on roles_permiso.id_rol = roles.id_rol 
			WHERE rol = ?
	`, rol)
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		return
	}

	var permiso Permiso
	for rows.Next() {
		err = rows.Scan(
			&permiso.Permiso,
		)
		if err != nil {
			return
		}

		permisos = append(permisos, permiso)
	}

	return
}

func CreateNewPermiso(db *sql.DB, text string) (err error) {
	_, err = db.Exec("INSERT INTO permisos(permiso) VALUES (?)", text)
	if err != nil {
		return
	}

	return
}

func DeletePermisoByID(db *sql.DB, id string) (err error) {
	row, err := db.Exec("DELETE from permisos WHERE id_permiso = ?", id)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}

func UpdatePermisoByID(db *sql.DB, id, newText string) (err error) {
	row, err := db.Exec("UPDATE permisos SET permiso = ? WHERE id_permiso = ?", newText, id)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}
