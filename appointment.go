package ind

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func ReserveAppointment(venue Venue, slot Slot) (*Slot, error) {
	body, err := json.Marshal(slot)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(reserveURL(venue, slot.Key), `application/json`, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data = bytes.ReplaceAll(data, []byte(")]}',\n"), []byte(""))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w with (status:%d, body:%s)", ErrInvalidHTTPStatus, resp.StatusCode, string(data))
	}

	var reserve = struct {
		Status string `json:"status"`
		Data   Slot   `json:"data"`
	}{}

	err = json.Unmarshal(data, &reserve)
	if err != nil {
		return nil, err
	}

	return &reserve.Data, nil
}

func CreateAppointment(venue Venue, req AppointmentReq) ([]byte, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(appointmentURL(venue), `application/json`, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data = bytes.ReplaceAll(data, []byte(")]}',\n"), []byte(""))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w with (status:%d, body:%s)", ErrInvalidHTTPStatus, resp.StatusCode, string(data))
	}

	return data, nil
}

func appointmentURL(venue Venue) string {
	return fmt.Sprintf(
		"https://oap.ind.nl/oap/api/desks/%s/appointments",
		venue,
	)
}

func reserveURL(venue Venue, key string) string {
	return fmt.Sprintf(
		"https://oap.ind.nl/oap/api/desks/%s/slots/%s",
		venue,
		key,
	)
}
