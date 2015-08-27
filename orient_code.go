package delta

import (
	"database/sql"

	_ "github.com/mattn/go-oci8"
)

type OrientCode struct {
	Id          int64  `json:"network_id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

func GetOrientCode(id int64) (*OrientCode, error) {
	o := OrientCode{}

	p := "SELECT orient_code_id, code, description FROM ORIENT_CODE WHERE orient_code_id = :orient_code_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&o.Id, &o.Code, &o.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &o, nil
}
