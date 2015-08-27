package delta

import (
	"database/sql"

	_ "github.com/mattn/go-oci8"
)

type Network struct {
	Id           int64   `json:"network_id"`
	AlbumPostfix *string `json:"album_postfix"`
	Code         string  `json:"code"`
	Description  string  `json:"description"`
	ExternalCode *string `json:"external_code"`
	PublicFlag   *bool   `json:"public_flag"`
}

func GetNetwork(id int64) (*Network, error) {
	n := Network{}

	p := "SELECT network_id, album_postfix, code, description, external_code, public_flag FROM NETWORK WHERE network_id = :network_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&n.Id, &n.AlbumPostfix, &n.Code, &n.Description, &n.ExternalCode, &n.PublicFlag)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &n, nil
}
