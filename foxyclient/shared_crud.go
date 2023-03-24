package foxyclient

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type foxyCrud interface {
	GetApiClient() FoxyClient
}

type record interface {
	setIdFromSelfUrl()
}

func DoList[T record](crud foxyCrud, path string) ([]T, error) {
	body, err := crud.GetApiClient().get(path)
	if err != nil {
		return nil, err
	}
	var records []T
	embeddedJsonResult := gjson.GetBytes(body, "_embedded.fx:*")
	embeddedJson := []byte(embeddedJsonResult.Raw)
	err = json.Unmarshal(embeddedJson, &records)
	if err != nil {
		return nil, err
	}
	for i := range records {
		// Need to modify records[i], rather than accessing wh directly via the loop, because the latter is by value
		records[i].setIdFromSelfUrl()
	}

	return records, err
}

func DoGet[T record](crud foxyCrud, path string) (T, error) {
	body, err := crud.GetApiClient().get(path)
	if err != nil {
		empty := new(T)
		return *empty, err
	}
	var record T
	err = json.Unmarshal(body, &record)
	if err != nil {
		empty := new(T)
		return *empty, err
	}
	record.setIdFromSelfUrl()
	return record, err
}

func DoAdd[T record](crud foxyCrud, record T, path string) (string, error) {
	updateJson, _ := json.Marshal(record)
	result, err := crud.GetApiClient().post(path, string(updateJson))
	if err != nil {
		return "", err
	}
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	id := extractId(selfUrl)
	return id, err
}

func DoUpdate[T record](crud foxyCrud, record T, path string) (string, error) {
	updateJson, _ := json.Marshal(record)
	result, e := crud.GetApiClient().patch(path, string(updateJson))
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	updatedId := extractId(selfUrl)
	return updatedId, e
}

func DoDelete[T record](crud foxyCrud, path string) error {
	_, e := crud.GetApiClient().delete(path)
	return e
}

func dereference[T interface{}](original []*T) []T {
	var values []T
	for _, r := range original {
		values = append(values, *r)
	}
	return values
}
