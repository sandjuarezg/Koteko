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
}

func Signin(user Usuario, db *sql.DB) (err error) {
	hash := sha1.Sum([]byte(user.Password))

	_, err = db.Exec("INSERT INTO usuarios(nombre, password, correo, id_rol) VALUES (?, ?, ?, ?)", user.Nombre, fmt.Sprintf("%X", hash[:]), user.Email, 1)
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
