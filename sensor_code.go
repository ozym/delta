package delta

import (
	"database/sql"

	_ "github.com/mattn/go-oci8"
)

type SensorCode struct {
	Id          int64  `json:"network_id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

func GetSensorCode(id int64) (*SensorCode, error) {
	s := SensorCode{}

	p := "SELECT sensor_code_id, code, description FROM SENSOR_CODE WHERE sensor_code_id = :sensor_code_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&s.Id, &s.Code, &s.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &s, nil
}
