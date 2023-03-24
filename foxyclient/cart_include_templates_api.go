package foxyclient

var (
	_ record   = &CartIncludeTemplate{}
	_ foxyCrud = &CartIncludeTemplatesApi{}
)

type CartIncludeTemplatesApi struct {
	apiClient FoxyClient
}

func (foxy *CartIncludeTemplatesApi) GetApiClient() FoxyClient {
	return foxy.apiClient
}

func (foxy *CartIncludeTemplatesApi) List() ([]CartIncludeTemplate, error) {
	path := foxy.storePath() + "/cart_include_templates?limit=300"
	result, e := DoList[*CartIncludeTemplate](foxy, path)
	return dereference(result), e
}

func (foxy *CartIncludeTemplatesApi) Get(id string) (CartIncludeTemplate, error) {
	path := "/cart_include_templates/" + id
	result, e := DoGet[*CartIncludeTemplate](foxy, path)
	return *result, e
}

func (foxy *CartIncludeTemplatesApi) Add(cartIncludeTemplate CartIncludeTemplate) (string, error) {
	path := foxy.storePath() + "/cart_include_templates"
	result, e := DoAdd[*CartIncludeTemplate](foxy, &cartIncludeTemplate, path)
	return result, e
}

func (foxy *CartIncludeTemplatesApi) Update(id string, cartIncludeTemplate CartIncludeTemplate) (string, error) {
	path := "/cart_include_templates/" + id
	result, e := DoUpdate[*CartIncludeTemplate](foxy, &cartIncludeTemplate, path)
	return result, e
}

func (foxy *CartIncludeTemplatesApi) Delete(id string) error {
	path := "/cart_include_templates/" + id
	return DoDelete[*CartIncludeTemplate](foxy, path)
}

func (foxy *CartIncludeTemplatesApi) storePath() string {
	storeId, _ := foxy.apiClient.retrieveStoreId()
	return "/stores/" + storeId
}

// ---

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
