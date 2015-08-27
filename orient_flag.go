package delta

import (
	"database/sql"

	_ "github.com/mattn/go-oci8"
)

type OrientFlag struct {
	Id          int64  `json:"network_id"`
	Description string `json:"description"`
	Flag        string `json:"flag"`
}

func GetOrientFlag(id int64) (*OrientFlag, error) {
	o := OrientFlag{}

	p := "SELECT orient_flag_id, description, flag FROM ORIENT_FLAG WHERE orient_flag_id = :orient_flag_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&o.Id, &o.Description, &o.Flag)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &o, nil
}
