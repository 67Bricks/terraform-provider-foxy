package foxyclient

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type CartIncludeTemplatesApi struct {
	apiClient FoxyClient
}

func (foxy *CartIncludeTemplatesApi) List() ([]CartIncludeTemplate, error) {
	// This is not retrieving all cartIncludeTemplates - only the first 300 - but is it plausible to have more than 300 cartIncludeTemplates?
	path := foxy.storePath() + "/cart_include_templates?limit=300"
	body, err := foxy.apiClient.get(path)
	if err != nil {
		return nil, err
	}
	var cartIncludeTemplates []CartIncludeTemplate
	embeddedJsonResult := gjson.GetBytes(body, "_embedded.fx:cart_include_templates")
	embeddedJson := []byte(embeddedJsonResult.Raw)
	err = json.Unmarshal(embeddedJson, &cartIncludeTemplates)
	if err != nil {
		return nil, err
	}
	for i := range cartIncludeTemplates {
		// Need to modify cartIncludeTemplates[i], rather than accessing wh directly via the loop, because the latter is by value
		cartIncludeTemplates[i].setIdFromSelfUrl()
	}

	return cartIncludeTemplates, err
}

func (foxy *CartIncludeTemplatesApi) Get(id string) (CartIncludeTemplate, error) {
	path := foxy.cartIncludeTemplatePath(id)
	body, err := foxy.apiClient.get(path)
	if err != nil {
		return CartIncludeTemplate{}, err
	}
	var cartIncludeTemplate CartIncludeTemplate
	err = json.Unmarshal(body, &cartIncludeTemplate)
	if err != nil {
		return CartIncludeTemplate{}, err
	}
	cartIncludeTemplate.setIdFromSelfUrl()
	return cartIncludeTemplate, err
}

func (foxy *CartIncludeTemplatesApi) Add(cartIncludeTemplate CartIncludeTemplate) (string, error) {
	path := foxy.storePath() + "/cart_include_templates"
	updateJson, _ := json.Marshal(cartIncludeTemplate)
	result, err := foxy.apiClient.post(path, string(updateJson))
	if err != nil {
		return "", err
	}
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	id := extractId(selfUrl)
	return id, err
}

func (foxy *CartIncludeTemplatesApi) Update(id string, cartIncludeTemplate CartIncludeTemplate) (string, error) {
	path := foxy.cartIncludeTemplatePath(id)
	amendedCartIncludeTemplate := cartIncludeTemplate
	updateJson, _ := json.Marshal(amendedCartIncludeTemplate)
	result, e := foxy.apiClient.patch(path, string(updateJson))
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	updatedId := extractId(selfUrl)
	return updatedId, e
}

func (foxy *CartIncludeTemplatesApi) Delete(id string) error {
	path := foxy.cartIncludeTemplatePath(id)
	_, e := foxy.apiClient.delete(path)
	return e
}

func (foxy *CartIncludeTemplatesApi) cartIncludeTemplatePath(id string) string {
	return "/cart_include_templates/" + id
}

func (foxy *CartIncludeTemplatesApi) storePath() string {
	storeId, _ := foxy.apiClient.retrieveStoreId()
	return "/stores/" + storeId
}

type CartIncludeTemplate struct {
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

func (cartIncludeTemplate *CartIncludeTemplate) setIdFromSelfUrl() {
	id := extractId(cartIncludeTemplate.Links.Self.Href)
	cartIncludeTemplate.Id = id
}
