package delta

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-oci8"
)

type EquipmentModel struct {
	Id          int64   `json:"equipment_model_id"`
	Model       string  `json:"model"`
	ModelAlias  *string `json:"model_alias"`
	Description string  `json:"description"`
	Notes       *string `json:"notes"`
	ModelNmbr   *string `json:"model_nmbr"`
}

func GetEquipmentModel(id int64) (*EquipmentModel, error) {
	m := EquipmentModel{Id: id}

	p := "SELECT model, model_alias, description, notes, model_nmbr FROM EQUIPMENT_MODEL WHERE equipment_model_id = :equipment_model_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&m.Model, &m.ModelAlias, &m.Description, &m.Notes, &m.ModelNmbr)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &m, nil
}

func GetEquipmentModels() ([]EquipmentModel, error) {
	var models []EquipmentModel

	q := "SELECT model, equipment_model_id, model_alias, description, notes, model_nmbr FROM EQUIPMENT_MODEL"

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := EquipmentModel{}
		if err := rows.Scan(&m.Id, &m.Model, &m.ModelAlias, &m.Description, &m.Notes, &m.ModelNmbr); err != nil {
			return nil, err
		}
		models = append(models, m)
		break
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return models, nil
}

func FindEquipmentModel(model string) (*EquipmentModel, error) {
	m := EquipmentModel{Model: model}

	p := "SELECT equipment_model_id, model_alias, description, notes, model_nmbr FROM EQUIPMENT_MODEL WHERE model = :model"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	fmt.Println(model)
	err = stmt.QueryRow(model).Scan(&m.Id, &m.ModelAlias, &m.Description, &m.Notes, &m.ModelNmbr)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &m, nil
}

func (m *EquipmentModel) FindEquipment() ([]Equipment, error) {
	return FindEquipmentByModelId(m.Id)
}
