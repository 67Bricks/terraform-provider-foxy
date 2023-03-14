package foxyclient

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type CheckoutTemplatesApi struct {
	apiClient FoxyClient
}

func (foxy *CheckoutTemplatesApi) List() ([]CheckoutTemplate, error) {
	// This is not retrieving all checkoutTemplates - only the first 300 - but is it plausible to have more than 300 checkoutTemplates?
	path := foxy.storePath() + "/checkout_templates?limit=300"
	body, err := foxy.apiClient.get(path)
	if err != nil {
		return nil, err
	}
	var checkoutTemplates []CheckoutTemplate
	embeddedJsonResult := gjson.GetBytes(body, "_embedded.fx:checkout_templates")
	embeddedJson := []byte(embeddedJsonResult.Raw)
	err = json.Unmarshal(embeddedJson, &checkoutTemplates)
	if err != nil {
		return nil, err
	}
	for i := range checkoutTemplates {
		// Need to modify checkoutTemplates[i], rather than accessing wh directly via the loop, because the latter is by value
		checkoutTemplates[i].setIdFromSelfUrl()
	}

	return checkoutTemplates, err
}

func (foxy *CheckoutTemplatesApi) Get(id string) (CheckoutTemplate, error) {
	path := foxy.checkoutTemplatePath(id)
	body, err := foxy.apiClient.get(path)
	if err != nil {
		return CheckoutTemplate{}, err
	}
	var checkoutTemplate CheckoutTemplate
	err = json.Unmarshal(body, &checkoutTemplate)
	if err != nil {
		return CheckoutTemplate{}, err
	}
	checkoutTemplate.setIdFromSelfUrl()
	return checkoutTemplate, err
}

func (foxy *CheckoutTemplatesApi) Add(checkoutTemplate CheckoutTemplate) (string, error) {
	path := foxy.storePath() + "/checkout_templates"
	updateJson, _ := json.Marshal(checkoutTemplate)
	result, err := foxy.apiClient.post(path, string(updateJson))
	if err != nil {
		return "", err
	}
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	id := extractId(selfUrl)
	return id, err
}

func (foxy *CheckoutTemplatesApi) Update(id string, checkoutTemplate CheckoutTemplate) (string, error) {
	path := foxy.checkoutTemplatePath(id)
	amendedCheckoutTemplate := checkoutTemplate
	updateJson, _ := json.Marshal(amendedCheckoutTemplate)
	result, e := foxy.apiClient.patch(path, string(updateJson))
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	updatedId := extractId(selfUrl)
	return updatedId, e
}

func (foxy *CheckoutTemplatesApi) Delete(id string) error {
	path := foxy.checkoutTemplatePath(id)
	_, e := foxy.apiClient.delete(path)
	return e
}

func (foxy *CheckoutTemplatesApi) checkoutTemplatePath(id string) string {
	return "/checkout_templates/" + id
}

func (foxy *CheckoutTemplatesApi) storePath() string {
	storeId, _ := foxy.apiClient.retrieveStoreId()
	return "/stores/" + storeId
}

type CheckoutTemplate struct {
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

func (checkoutTemplate *CheckoutTemplate) setIdFromSelfUrl() {
	id := extractId(checkoutTemplate.Links.Self.Href)
	checkoutTemplate.Id = id
}
