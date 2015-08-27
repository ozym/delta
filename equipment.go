package delta

import (
	"database/sql"

	_ "github.com/mattn/go-oci8"
)

type Equipment struct {
	Id           int64   `json:"equipment_id"`
	AssetNumber  *string `json:"asset_number"`
	SerialNumber *string `json:"serial_number"`
	Notes        *string `json:"notes"`
}

func GetEquipment(id int64) (*Equipment, error) {
	e := Equipment{Id: id}

	p := "SELECT asset_number, serial_number, notes FROM equipment WHERE equipment_id = :equipment_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&e.AssetNumber, &e.SerialNumber, &e.Notes)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &e, nil
}

func GetEquipments() ([]Equipment, error) {
	var equipment []Equipment

	q := "SELECT equipment_id, asset_number, serial_number, notes FROM EQUIPMENT"

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := Equipment{}
		if err := rows.Scan(&e.Id, &e.AssetNumber, &e.SerialNumber, &e.Notes); err != nil {
			return nil, err
		}
		equipment = append(equipment, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return equipment, nil
}

func FindEquipmentByModelId(id int64) ([]Equipment, error) {
	var equipment []Equipment

	p := "SELECT equipment_id, asset_number, serial_number, notes FROM EQUIPMENT e WHERE e.equipment_model_id = :equipment_model_id"

	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := Equipment{}
		if err := rows.Scan(&e.Id, &e.AssetNumber, &e.SerialNumber, &e.Notes); err != nil {
			return nil, err
		}
		equipment = append(equipment, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return equipment, nil
}
