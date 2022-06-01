package models

import (
	"database/sql"
)

type Categoria struct {
	IDCategoria int
	Categoria   string
}

func GetAllCategorias(db *sql.DB) (categorias []Categoria, err error) {
	rows, err := db.Query("SELECT id_categoria, categoria FROM categorias")
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		return
	}

	var categoria Categoria
	for rows.Next() {
		err = rows.Scan(
			&categoria.IDCategoria,
			&categoria.Categoria,
		)
		if err != nil {
			return
		}

		categorias = append(categorias, categoria)
	}

	return
}

func CreateNewCategoria(db *sql.DB, text string) (err error) {
	_, err = db.Exec("INSERT INTO categorias(categoria) VALUES (?)", text)
	if err != nil {
		return
	}

	return
}

func DeleteCategoriaByID(db *sql.DB, id string) (err error) {
	row, err := db.Exec("DELETE from categorias WHERE id_categoria = ?", id)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}

func UpdateCategoriaByID(db *sql.DB, id, newText string) (err error) {
	row, err := db.Exec("UPDATE categorias SET categoria = ? WHERE id_categoria = ?", newText, id)
	if err != nil {
		return
	}

	_, err = row.RowsAffected()
	if err != nil {
		return
	}

	return
}
