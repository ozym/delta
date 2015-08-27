package delta

import (
	"database/sql"

	_ "github.com/mattn/go-oci8"
)

type BandCode struct {
	Id           int64   `json:"network_id"`
	Code         string  `json:"code"`
	CornerPeriod *string `json:"corner_period"`
	Description  string  `json:"description"`
	SampleRate   *string `json:"sample_rate"`
}

func GetBandCode(id int64) (*BandCode, error) {
	b := BandCode{}

	p := "SELECT band_code_id, code, corner_period, description, sample_rate FROM BAND_CODE WHERE band_code_id = :band_code_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&b.Id, &b.Code, &b.CornerPeriod, &b.Description, &b.SampleRate)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &b, nil
}
