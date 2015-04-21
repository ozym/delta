package main

import (
	"database/sql"
	"fmt"
	//	"bytes"
	//	"encoding/json"
	_ "github.com/mattn/go-oci8"
	//	"io"
	//	"log"
	//	"net/http"
	"time"
)

type SeismicStream struct {
	Id             int64     `json:"seismic_stream_id"`
	Freq           float64   `json:"freq"`
	MaxDrift       float64   `json:"max_drift"`
	SampleRate     float64   `json:"sample_rate"`
	Sensitivity    float64   `json:"sensitivity"`
	StartTimeStamp time.Time `json:"start_time_stamp"`
	StopTimeStamp  time.Time `json:"stop_time_stamp"`
}

func FindSeismicStream(id int64) (*SeismicStream, error) {
	s := SeismicStream{}

	sql := "SELECT freq, max_drift, sample_rate, sensitivity, start_time_stamp, stop_time_stamp FROM SEISMIC_STREAM WHERE seismic_stream_id = ?"
	_, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func SeismicStreams() ([]SeismicStream, error) {
	var streams []SeismicStream

	sql := "SELECT seismic_stream_id, freq, max_drift, sample_rate, sensitivity, start_time_stamp, stop_time_stamp FROM SEISMIC_STREAM"

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := SeismicStream{}
		if err := rows.Scan(&s.Id, &s.Freq, &s.MaxDrift, &s.SampleRate, &s.Sensitivity, &s.StartTimeStamp, &s.StopTimeStamp); err != nil {
			return nil, err
		}
		streams = append(streams, s)
		break
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return streams, nil
}

func FindCurrentSeismicStreamsByComponentId(id int64) ([]SeismicStream, error) {
	var streams []SeismicStream

	p := "SELECT seismic_stream_id, freq, max_drift, sample_rate, sensitivity, start_time_stamp, stop_time_stamp FROM SEISMIC_STREAM WHERE component_id = :component_id AND stop_time_stamp > :stop_time_stamp"

	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := SeismicStream{}
		if err := rows.Scan(&s.Id, &s.Freq, &s.MaxDrift, &s.SampleRate, &s.Sensitivity, &s.StartTimeStamp, &s.StopTimeStamp); err != nil {
			return nil, err
		}
		streams = append(streams, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return streams, nil

}

func FindSeismicStreamsByComponentId(id int64, start, stop time.Time) ([]SeismicStream, error) {
	var streams []SeismicStream

	if stop.After(time.Now().UTC()) {
		return FindCurrentSeismicStreamsByComponentId(id)
	}

	p := "SELECT seismic_stream_id, freq, max_drift, sample_rate, sensitivity, start_time_stamp, stop_time_stamp FROM SEISMIC_STREAM WHERE component_id = :component_id AND stop_time_stamp > :start_time_stamp AND start_time_stamp < :stop_time_stamp"

	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id, start, stop)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := SeismicStream{}
		if err := rows.Scan(&s.Id, &s.Freq, &s.MaxDrift, &s.SampleRate, &s.Sensitivity, &s.StartTimeStamp, &s.StopTimeStamp); err != nil {
			return nil, err
		}
		streams = append(streams, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return streams, nil
}

func FindCurrentSeismicStreamsByChannelId(id int64) ([]SeismicStream, error) {
	var streams []SeismicStream

	p := "SELECT seismic_stream_id, freq, max_drift, sample_rate, sensitivity, start_time_stamp, stop_time_stamp FROM SEISMIC_STREAM WHERE channel_id = :channel_id AND stop_time_stamp > :stop_time_stamp"

	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := SeismicStream{}
		if err := rows.Scan(&s.Id, &s.Freq, &s.MaxDrift, &s.SampleRate, &s.Sensitivity, &s.StartTimeStamp, &s.StopTimeStamp); err != nil {
			return nil, err
		}
		streams = append(streams, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return streams, nil

}

func FindSeismicStreamsByChannelId(id int64, start, stop time.Time) ([]SeismicStream, error) {
	var streams []SeismicStream

	if stop.After(time.Now().UTC()) {
		return FindCurrentSeismicStreamsByChannelId(id)
	}

	p := "SELECT seismic_stream_id, freq, max_drift, sample_rate, sensitivity, start_time_stamp, stop_time_stamp FROM SEISMIC_STREAM WHERE channel_id = :channel_id AND stop_time_stamp > :start_time_stamp AND start_time_stamp < :stop_time_stamp"

	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id, start, stop)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := SeismicStream{}
		if err := rows.Scan(&s.Id, &s.Freq, &s.MaxDrift, &s.SampleRate, &s.Sensitivity, &s.StartTimeStamp, &s.StopTimeStamp); err != nil {
			return nil, err
		}
		streams = append(streams, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return streams, nil
}

func (s *SeismicStream) GetChannelId() (*int64, error) {
	var id int64

	p := "SELECT channel_id FROM SEISMIC_STREAM WHERE seismic_stream_id = :seismic_stream_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(s.Id).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *SeismicStream) GetChannel() (*Channel, error) {
	id, err := s.GetChannelId()
	if err != nil {
		return nil, err
	} else if id == nil {
		return nil, nil
	}
	return GetChannel(*id)
}

func (s *SeismicStream) GetOrientCodeId() (*int64, error) {
	var id int64

	p := "SELECT orient_code_id FROM SEISMIC_STREAM WHERE seismic_stream_id = :seismic_stream_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(s.Id).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *SeismicStream) GetOrientCode() (*OrientCode, error) {
	id, err := s.GetOrientCodeId()
	if err != nil {
		return nil, err
	} else if id == nil {
		return nil, nil
	}
	return GetOrientCode(*id)
}

func (s *SeismicStream) GetSensorCodeId() (*int64, error) {
	var id int64

	p := "SELECT sensor_code_id FROM SEISMIC_STREAM WHERE seismic_stream_id = :seismic_stream_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(s.Id).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *SeismicStream) GetSensorCode() (*SensorCode, error) {
	id, err := s.GetSensorCodeId()
	if err != nil {
		return nil, err
	} else if id == nil {
		return nil, nil
	}
	return GetSensorCode(*id)
}

func (s *SeismicStream) GetBandCodeId() (*int64, error) {
	var id int64

	p := "SELECT band_code_id FROM SEISMIC_STREAM WHERE seismic_stream_id = :seismic_stream_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(s.Id).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *SeismicStream) GetBandCode() (*BandCode, error) {
	id, err := s.GetBandCodeId()
	if err != nil {
		return nil, err
	} else if id == nil {
		return nil, nil
	}
	return GetBandCode(*id)
}

func (s *SeismicStream) GetOrientFlagId() (*int64, error) {
	var id int64

	p := "SELECT orient_flag_id FROM SEISMIC_STREAM WHERE seismic_stream_id = :seismic_stream_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(s.Id).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *SeismicStream) GetOrientFlag() (*OrientFlag, error) {
	id, err := s.GetOrientFlagId()
	if err != nil {
		return nil, err
	} else if id == nil {
		return nil, nil
	}
	return GetOrientFlag(*id)
}

func (s *SeismicStream) GetChannelLabel() (*string, error) {
	b, err := s.GetBandCode()
	if err != nil {
		return nil, err
	} else if b == nil {
		return nil, nil
	}

	c, err := s.GetSensorCode()
	if err != nil {
		return nil, err
	} else if c == nil {
		return nil, nil
	}

	o, err := s.GetOrientCode()
	if err != nil {
		return nil, err
	} else if o == nil {
		return nil, nil
	}

	l := fmt.Sprintf("%s%s%s", b.Code, c.Code, o.Code)

	return &l, nil
}
