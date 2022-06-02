package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
)

type Usuario struct {
	ID       int
	Nombre   string
	Email    string
	Password string
	IDRol    int
}

func CreateNewUser(user Usuario, db *sql.DB) (err error) {
	hash := sha1.Sum([]byte(user.Password))

	_, err = db.Exec("INSERT INTO usuarios(nombre, password, correo, id_rol) VALUES (?, ?, ?, ?)",
		user.Nombre, fmt.Sprintf("%X", hash[:]), user.Email, 1,
	)
	if err != nil {
		return
	}

	return
}

func Login(email, password string, db *sql.DB) (user Usuario, err error) {
	hash := sha1.Sum([]byte(password))

	row := db.QueryRow("SELECT id_usuario, nombre FROM usuarios WHERE correo = ? AND password = ?", email, fmt.Sprintf("%X", hash[:]))
	err = row.Scan(&user.ID, &user.Nombre)
	if err != nil {
		return
	}

	return
}

func GetAllUsers(db *sql.DB) (usuarios []Usuario, err error) {
	rows, err := db.Query("SELECT id_usuario, nombre, password, correo, id_rol FROM usuarios")
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		return
	}

	var us Usuario
	for rows.Next() {
		err = rows.Scan(
			&us.ID,
			&us.Nombre,
			&us.Password,
			&us.Email,
			&us.IDRol,
		)
		if err != nil {
			return
		}

		usuarios = append(usuarios, us)
	}

	return
}

func GetUserByID(db *sql.DB, id int) (us Usuario, err error) {
	row := db.QueryRow(`
		SELECT id_usuario, nombre, password, correo, id_rol 
			FROM usuarios 
			WHERE id_usuario = ?`, id)
	if err != nil {
		return
	}

	if row.Err() != nil {
		return
	}

	err = row.Scan(
		&us.ID,
		&us.Nombre,
		&us.Password,
		&us.Email,
		&us.IDRol,
	)
	if err != nil {
		return
	}

	return
}

func DeleteUserByID(db *sql.DB, id string) (err error) {
	row, err := db.Exec("DELETE from usuarios WHERE id_usuario = ?", id)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}

func UpdateUserByID(db *sql.DB, user Usuario, id string) (err error) {
	hash := sha1.Sum([]byte(user.Password))

	row, err := db.Exec("UPDATE usuarios SET password = ?, id_rol = ? WHERE id_usuario = ?",
		fmt.Sprintf("%X", hash[:]), user.IDRol, id,
	)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}
