package converter

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"notify-integrator/internal/types"
	"reflect"
)

func ToObject(b []byte) ([]types.Body, error) {

	csvFile := csv.NewReader(bytes.NewReader(b))
	messages, err := csvFile.ReadAll()
	if err != nil {
		return nil, err
	}

	var bu types.Body
	var bus []types.Body

	header := []string{"nome", "telefone", "idade"}

	buInterface := make(map[string]interface{}, reflect.ValueOf(bu).NumField())
	for _, message := range messages {

		for i, value := range message {
			buInterface[header[i]] = value
		}

		j, err := json.Marshal(buInterface)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(j, &bu)
		if err != nil {
			return nil, err
		}

		bus = append(bus, bu)
	}

	return bus, nil
}
