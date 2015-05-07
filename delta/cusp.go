package main

import (
	"fmt"
	//	"bytes"
	//"encoding/json"
	_ "github.com/mattn/go-oci8"
	"github.com/ozym/delta"
	//	"io"
	//	"log"
	//	"net/http"
	"encoding/xml"
	//	"time"
)

type CsdStream struct {
	Name        string  `json:"name" xml:"name,attr"`
	Srcname     string  `json:"srcname" xml:"srcname,attr"`
	NetworkId   string  `json:"network_id" xml:"network_id,attr"`
	StationId   string  `json:"station_id" xml:"station_id,attr"`
	LocationId  string  `json:"location_id" xml:"location_id,attr"`
	ChannelId   string  `json:"channel_id" xml:"channel_id,attr"`
	Latitude    float64 `json:"latitude" xml:"latitude,attr"`
	Longitude   float64 `json:"longitude" xml:"longitude,attr"`
	SampleRate  float64 `json:"sample_rate" xml:"sample_rate,attr"`
	Sensitivity float64 `json:"sensitivity" xml:"sensitivity,attr"`
}

type CsdPair struct {
	ChannelNo int64       `json:"channel_no" xml:"channel_no,attr"`
	PinNo     int64       `json:"pin_no" xml:"pin_no,attr"`
	Streams   []CsdStream `json:"stream" xml:"stream"`
}

type CsdInstall struct {
	Start    string    `json:"start" xml:"start,attr"`
	Channels []CsdPair `json:"channel" xml:"channel"`
}

type Csd struct {
	Model    string       `json:"model" xml:"model,attr"`
	Serial   string       `json:"serial" xml:"serial,attr"`
	Installs []CsdInstall `json:"install" xml:"install"`
}

type Csds struct {
	XMLName  xml.Name
	Installs []Csd `json:"serial" xml:"serial"`
}

func CsdConfig(indent bool, models []string) ([]byte, error) {
	csds := Csds{XMLName: xml.Name{"", "csd"}}

	for _, m := range models {

		model, err := delta.FindEquipmentModel(m)
		if err != nil {
			return nil, err
		} else if model == nil {
			continue
		}
		equipment, err := model.FindEquipment()
		if err != nil {
			return nil, err
		}
		for _, e := range equipment {
			if e.SerialNumber == nil {
				continue
			}
			sensor, err := delta.FindSensorByEquipmentId(e.Id)
			if err != nil {
				return nil, err
			}
			components, err := sensor.FindComponents()
			if err != nil {
				return nil, err
			}

			installed, err := sensor.FindInstalledSensors()
			if err != nil {
				return nil, err
			}
			var installs []CsdInstall
			for _, i := range installed {
				start := i.StartTimeStamp.Format("2006-01-02 15:04:05")

				pairs := make(map[int64]CsdPair)
				site, err := i.GetSeismicSite()
				if err != nil {
					return nil, err
				}
				station, err := site.GetSeismicStation()
				if err != nil {
					return nil, err
				}
				if station.LongName == nil {
					continue
				}
				network, err := station.GetNetwork()
				if err != nil {
					return nil, err
				}
				code := network.Code
				if network.ExternalCode != nil {
					code = *network.ExternalCode
				}

				for _, c := range components {
					orient, err := c.GetCompOrient()
					if err != nil {
						return nil, err
					}
					streams, err := c.FindSeismicStreams(i.StartTimeStamp, i.StopTimeStamp)
					if err != nil {
						return nil, err
					}
					for _, s := range streams {
						channel, err := s.GetChannel()
						if err != nil {
							return nil, err
						}
						label, err := s.GetChannelLabel()
						if err != nil {
							return nil, err
						}
						if label == nil {
							continue
						}
						srcname := fmt.Sprintf("%s_%s_%s_%s", code, station.StationId, site.Location, *label)
						k := CsdStream{
							Name:        *station.LongName,
							SampleRate:  s.SampleRate,
							Sensitivity: s.Sensitivity,
							NetworkId:   code,
							StationId:   station.StationId,
							LocationId:  site.Location,
							ChannelId:   *label,
							Latitude:    site.Latitude,
							Longitude:   site.Longitude,
							Srcname:     srcname,
						}
						v := []CsdStream{k}

						_, ok := pairs[channel.PinNo]
						if ok {
							v = append(pairs[channel.PinNo].Streams, k)
						}
						pairs[channel.PinNo] = CsdPair{ChannelNo: channel.PinNo, PinNo: orient.PinNo, Streams: v}
					}

				}
				channels := make([]CsdPair, 0, len(pairs))
				for _, p := range pairs {
					channels = append(channels, p)
				}

				installs = append(installs, CsdInstall{Start: start, Channels: channels})
			}
			if installs == nil {
				continue
			}

			csds.Installs = append(csds.Installs, Csd{Model: m, Serial: *e.SerialNumber, Installs: installs})
		}
	}

	h := []byte(xml.Header)
	var s []byte
	var err error

	if indent {
		s, err = xml.MarshalIndent(csds, "", "  ")
		if err != nil {
			return nil, err
		}
	} else {
		s, err = xml.Marshal(csds)
		if err != nil {
			return nil, err
		}
	}

	return append(append(h[:], s[:]...), '\n'), nil
}
