package foxyclient

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type CartTemplatesApi struct {
	apiClient FoxyClient
}

func (foxy *CartTemplatesApi) List() ([]CartTemplate, error) {
	// This is not retrieving all cartTemplates - only the first 300 - but is it plausible to have more than 300 cartTemplates?
	path := foxy.storePath() + "/cart_templates?limit=300"
	body, err := foxy.apiClient.get(path)
	if err != nil {
		return nil, err
	}
	var cartTemplates []CartTemplate
	embeddedJsonResult := gjson.GetBytes(body, "_embedded.fx:cart_templates")
	embeddedJson := []byte(embeddedJsonResult.Raw)
	err = json.Unmarshal(embeddedJson, &cartTemplates)
	if err != nil {
		return nil, err
	}
	for i := range cartTemplates {
		// Need to modify cartTemplates[i], rather than accessing wh directly via the loop, because the latter is by value
		cartTemplates[i].setIdFromSelfUrl()
	}

	return cartTemplates, err
}

func (foxy *CartTemplatesApi) Get(id string) (CartTemplate, error) {
	path := foxy.cartTemplatePath(id)
	body, err := foxy.apiClient.get(path)
	if err != nil {
		return CartTemplate{}, err
	}
	var cartTemplate CartTemplate
	err = json.Unmarshal(body, &cartTemplate)
	if err != nil {
		return CartTemplate{}, err
	}
	cartTemplate.setIdFromSelfUrl()
	return cartTemplate, err
}

func (foxy *CartTemplatesApi) Add(cartTemplate CartTemplate) (string, error) {
	path := foxy.storePath() + "/cart_templates"
	updateJson, _ := json.Marshal(cartTemplate)
	result, err := foxy.apiClient.post(path, string(updateJson))
	if err != nil {
		return "", err
	}
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	id := extractId(selfUrl)
	return id, err
}

func (foxy *CartTemplatesApi) Update(id string, cartTemplate CartTemplate) (string, error) {
	path := foxy.cartTemplatePath(id)
	amendedCartTemplate := cartTemplate
	updateJson, _ := json.Marshal(amendedCartTemplate)
	result, e := foxy.apiClient.patch(path, string(updateJson))
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	updatedId := extractId(selfUrl)
	return updatedId, e
}

func (foxy *CartTemplatesApi) Delete(id string) error {
	path := foxy.cartTemplatePath(id)
	_, e := foxy.apiClient.delete(path)
	return e
}

func (foxy *CartTemplatesApi) cartTemplatePath(id string) string {
	return "/cart_templates/" + id
}

func (foxy *CartTemplatesApi) storePath() string {
	storeId, _ := foxy.apiClient.retrieveStoreId()
	return "/stores/" + storeId
}

type CartTemplate struct {
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

func (cartTemplate *CartTemplate) setIdFromSelfUrl() {
	id := extractId(cartTemplate.Links.Self.Href)
	cartTemplate.Id = id
}
