package foxyclient

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type foxyCrud interface {
	GetListPath() string
	GetRecordPath(id string) string
	GetRecordAddPath() string
	GetApiClient() FoxyClient
}

type record interface {
	setIdFromSelfUrl()
}

func DoList[T record](crud foxyCrud) ([]T, error) {
	// This is not retrieving all records - only the first 300 - but is it plausible to have more than 300 records?
	path := crud.GetListPath()
	body, err := crud.GetApiClient().get(path)
	if err != nil {
		return nil, err
	}
	var records []T
	embeddedJsonResult := gjson.GetBytes(body, "_embedded.fx:cart_templates")
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

func DoGet[T record](crud foxyCrud, id string) (T, error) {
	path := crud.GetRecordPath(id)
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

func DoAdd[T record](crud foxyCrud, record T) (string, error) {
	path := crud.GetRecordAddPath()
	updateJson, _ := json.Marshal(record)
	result, err := crud.GetApiClient().post(path, string(updateJson))
	if err != nil {
		return "", err
	}
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	id := extractId(selfUrl)
	return id, err
}

func DoUpdate[T record](crud foxyCrud, id string, record T) (string, error) {
	path := crud.GetRecordPath(id)
	amendedCartTemplate := record
	updateJson, _ := json.Marshal(amendedCartTemplate)
	result, e := crud.GetApiClient().patch(path, string(updateJson))
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	updatedId := extractId(selfUrl)
	return updatedId, e
}

func DoDelete[T record](crud foxyCrud, id string) error {
	path := crud.GetRecordPath(id)
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
