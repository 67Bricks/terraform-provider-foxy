package foxyclient

var (
	_ record   = &CartTemplate{}
	_ foxyCrud = &CartTemplatesApi{}
)

// ----

type CartTemplatesApi struct {
	apiClient FoxyClient
}

func (foxy *CartTemplatesApi) GetApiClient() FoxyClient {
	return foxy.apiClient
}

func (foxy *CartTemplatesApi) List() ([]CartTemplate, error) {
	path := foxy.storePath() + "/cart_templates?limit=300"
	result, e := DoList[*CartTemplate](foxy, path)
	return dereference(result), e
}

func (foxy *CartTemplatesApi) Get(id string) (CartTemplate, error) {
	path := "/cart_templates/" + id
	result, e := DoGet[*CartTemplate](foxy, path)
	return *result, e
}

func (foxy *CartTemplatesApi) Add(cartTemplate CartTemplate) (string, error) {
	path := foxy.storePath() + "/cart_templates"
	result, e := DoAdd[*CartTemplate](foxy, &cartTemplate, path)
	return result, e
}

func (foxy *CartTemplatesApi) Update(id string, cartTemplate CartTemplate) (string, error) {
	path := "/cart_templates/" + id
	result, e := DoUpdate[*CartTemplate](foxy, &cartTemplate, path)
	return result, e
}

func (foxy *CartTemplatesApi) Delete(id string) error {
	path := "/cart_templates/" + id
	return DoDelete[*CartTemplate](foxy, path)
}

func (foxy *CartTemplatesApi) storePath() string {
	storeId, _ := foxy.apiClient.retrieveStoreId()
	return "/stores/" + storeId
}

// ----

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
