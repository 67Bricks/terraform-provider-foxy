package foxyclient

var (
	_ record   = &Webhook{}
	_ foxyCrud = &WebhooksApi{}
)

// ----

type WebhooksApi struct {
	apiClient FoxyClient
}

func (foxy *WebhooksApi) GetApiClient() FoxyClient {
	return foxy.apiClient
}

func (foxy *WebhooksApi) List() ([]Webhook, error) {
	path := foxy.storePath() + "/webhooks?limit=300"
	result, e := DoList[*Webhook](foxy, path)
	return dereference(result), e
}

func (foxy *WebhooksApi) Get(id string) (Webhook, error) {
	path := "/webhooks/" + id
	result, e := DoGet[*Webhook](foxy, path)
	return *result, e
}

func (foxy *WebhooksApi) Add(webhook Webhook) (string, error) {
	path := foxy.storePath() + "/webhooks"
	result, e := DoAdd[*Webhook](foxy, &webhook, path)
	return result, e
}

func (foxy *WebhooksApi) Update(id string, webhook Webhook) (string, error) {
	path := "/webhooks/" + id
	amendedWebhook := webhook
	amendedWebhook.EventResource = "" // This cannot be updated, it can only be set on creation
	result, e := DoUpdate[*Webhook](foxy, &amendedWebhook, path)
	return result, e
}

func (foxy *WebhooksApi) Delete(id string) error {
	path := "/webhooks/" + id
	return DoDelete[*Webhook](foxy, path)
}

func (foxy *WebhooksApi) storePath() string {
	storeId, _ := foxy.apiClient.retrieveStoreId()
	return "/stores/" + storeId
}

// ----

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
