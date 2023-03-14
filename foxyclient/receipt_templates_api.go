package foxyclient

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type ReceiptTemplatesApi struct {
	apiClient FoxyClient
}

func (foxy *ReceiptTemplatesApi) List() ([]ReceiptTemplate, error) {
	// This is not retrieving all receiptTemplates - only the first 300 - but is it plausible to have more than 300 receiptTemplates?
	path := foxy.storePath() + "/receipt_templates?limit=300"
	body, err := foxy.apiClient.get(path)
	if err != nil {
		return nil, err
	}
	var receiptTemplates []ReceiptTemplate
	embeddedJsonResult := gjson.GetBytes(body, "_embedded.fx:receipt_templates")
	embeddedJson := []byte(embeddedJsonResult.Raw)
	err = json.Unmarshal(embeddedJson, &receiptTemplates)
	if err != nil {
		return nil, err
	}
	for i := range receiptTemplates {
		// Need to modify receiptTemplates[i], rather than accessing wh directly via the loop, because the latter is by value
		receiptTemplates[i].setIdFromSelfUrl()
	}

	return receiptTemplates, err
}

func (foxy *ReceiptTemplatesApi) Get(id string) (ReceiptTemplate, error) {
	path := foxy.receiptTemplatePath(id)
	body, err := foxy.apiClient.get(path)
	if err != nil {
		return ReceiptTemplate{}, err
	}
	var receiptTemplate ReceiptTemplate
	err = json.Unmarshal(body, &receiptTemplate)
	if err != nil {
		return ReceiptTemplate{}, err
	}
	receiptTemplate.setIdFromSelfUrl()
	return receiptTemplate, err
}

func (foxy *ReceiptTemplatesApi) Add(receiptTemplate ReceiptTemplate) (string, error) {
	path := foxy.storePath() + "/receipt_templates"
	updateJson, _ := json.Marshal(receiptTemplate)
	result, err := foxy.apiClient.post(path, string(updateJson))
	if err != nil {
		return "", err
	}
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	id := extractId(selfUrl)
	return id, err
}

func (foxy *ReceiptTemplatesApi) Update(id string, receiptTemplate ReceiptTemplate) (string, error) {
	path := foxy.receiptTemplatePath(id)
	amendedReceiptTemplate := receiptTemplate
	updateJson, _ := json.Marshal(amendedReceiptTemplate)
	result, e := foxy.apiClient.patch(path, string(updateJson))
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	updatedId := extractId(selfUrl)
	return updatedId, e
}

func (foxy *ReceiptTemplatesApi) Delete(id string) error {
	path := foxy.receiptTemplatePath(id)
	_, e := foxy.apiClient.delete(path)
	return e
}

func (foxy *ReceiptTemplatesApi) receiptTemplatePath(id string) string {
	return "/receipt_templates/" + id
}

func (foxy *ReceiptTemplatesApi) storePath() string {
	storeId, _ := foxy.apiClient.retrieveStoreId()
	return "/stores/" + storeId
}

type ReceiptTemplate struct {
	Id          string `json:"-"`
	Description string `json:"description"`
	Content     string `json:"content"`
	ContentUrl  string `json:"content_url"`

	// @todo This and the setIdFromSelfUrl method are a clumsy way of unmarshalling the JSON - could we do this better?
	Links struct {
		Self struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
	} `json:"_links,omitempty"`
}

func (receiptTemplate *ReceiptTemplate) setIdFromSelfUrl() {
	id := extractId(receiptTemplate.Links.Self.Href)
	receiptTemplate.Id = id
}
