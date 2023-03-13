package foxyclient

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type WebhooksApi struct {
	apiClient FoxyClient
}

func (foxy *WebhooksApi) List() ([]Webhook, error) {
	// This is not retrieving all webhooks - only the first 300 - but is it plausible to have more than 300 webhooks?
	path := foxy.storePath() + "/webhooks?limit=300"
	body, err := foxy.apiClient.get(path)
	if err != nil {
		return nil, err
	}
	var webhooks []Webhook
	embeddedJsonResult := gjson.GetBytes(body, "_embedded.fx:webhooks")
	embeddedJson := []byte(embeddedJsonResult.Raw)
	err = json.Unmarshal(embeddedJson, &webhooks)
	if err != nil {
		return nil, err
	}
	for i := range webhooks {
		// Need to modify webhooks[i], rather than accessing wh directly via the loop, because the latter is by value
		webhooks[i].setIdFromSelfUrl()
	}

	return webhooks, err
}

func (foxy *WebhooksApi) Get(id string) (Webhook, error) {
	path := foxy.webhookPath(id)
	body, err := foxy.apiClient.get(path)
	if err != nil {
		return Webhook{}, err
	}
	var webhook Webhook
	err = json.Unmarshal(body, &webhook)
	if err != nil {
		return Webhook{}, err
	}
	webhook.setIdFromSelfUrl()
	return webhook, err
}

func (foxy *WebhooksApi) Add(webhook Webhook) (string, error) {
	path := foxy.storePath() + "/webhooks"
	updateJson, _ := json.Marshal(webhook)
	result, err := foxy.apiClient.post(path, string(updateJson))
	if err != nil {
		return "", err
	}
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	id := extractId(selfUrl)
	return id, err
}

func (foxy *WebhooksApi) Update(id string, webhook Webhook) (string, error) {
	path := foxy.webhookPath(id)
	amendedWebhook := webhook
	amendedWebhook.EventResource = "" // This cannot be updated, it can only be set on creation
	updateJson, _ := json.Marshal(amendedWebhook)
	result, e := foxy.apiClient.patch(path, string(updateJson))
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	updatedId := extractId(selfUrl)
	return updatedId, e
}

func (foxy *WebhooksApi) Delete(id string) error {
	path := foxy.webhookPath(id)
	_, e := foxy.apiClient.delete(path)
	return e
}

func (foxy *WebhooksApi) webhookPath(id string) string {
	return "/webhooks/" + id
}

func (foxy *WebhooksApi) storePath() string {
	storeId, _ := foxy.apiClient.retrieveStoreId()
	return "/stores/" + storeId
}

type Webhook struct {
	Id            string `json:"-"`
	Format        string `json:"format,omitempty"`
	Version       int    `json:"version,omitempty"`
	Name          string `json:"name,omitempty"`
	Url           string `json:"url,omitempty"`
	Query         string `json:"query,omitempty"`
	EncryptionKey string `json:"encryption_key,omitempty"`
	EventResource string `json:"event_resource,omitempty"`

	// @todo This and the setIdFromSelfUrl method are a clumsy way of unmarshalling the JSON - could we do this better?
	Links struct {
		Self struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
	} `json:"_links,omitempty"`
}

func (webhook *Webhook) setIdFromSelfUrl() {
	id := extractId(webhook.Links.Self.Href)
	webhook.Id = id
}
