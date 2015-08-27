package delta

import (
	"database/sql"

	_ "github.com/mattn/go-oci8"
)

type CompOrient struct {
	Id      int64   `json:"channel_id"`
	Azimuth float64 `json:"azimuth"`
	Dip     float64 `json:"dip"`
	PinNo   int64   `json:"pin_no"`
}

func GetCompOrient(id int64) (*CompOrient, error) {
	c := CompOrient{}

	p := "SELECT comp_orient_id, azimuth, dip, pin_no FROM COMP_ORIENT WHERE comp_orient_id = :comp_orient_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&c.Id, &c.Azimuth, &c.Dip, &c.PinNo)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &c, nil
}
