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
