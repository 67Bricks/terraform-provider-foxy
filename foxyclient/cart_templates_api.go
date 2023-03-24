package foxyclient

var (
	_ record   = &CartTemplate{}
	_ foxyCrud = &CartTemplatesApi{}
)

type CartTemplatesApi struct {
	apiClient FoxyClient
}

func (foxy *CartTemplatesApi) GetApiClient() FoxyClient {
	return foxy.apiClient
}

func (foxy *CartTemplatesApi) GetListPath() string {
	return foxy.storePath() + "/cart_templates?limit=300"
}

func (foxy *CartTemplatesApi) GetRecordPath(id string) string {
	return "/cart_templates/" + id
}
func (foxy *CartTemplatesApi) GetRecordAddPath() string {
	return foxy.storePath() + "/cart_templates"
}

func (foxy *CartTemplatesApi) List() ([]CartTemplate, error) {
	result, e := DoList[*CartTemplate](foxy)
	return dereference(result), e
}

func (foxy *CartTemplatesApi) Get(id string) (CartTemplate, error) {
	result, e := DoGet[*CartTemplate](foxy, id)
	return *result, e
}

func (foxy *CartTemplatesApi) Add(cartTemplate CartTemplate) (string, error) {
	result, e := DoAdd[*CartTemplate](foxy, &cartTemplate)
	return result, e
}

func (foxy *CartTemplatesApi) Update(id string, cartTemplate CartTemplate) (string, error) {
	result, e := DoUpdate[*CartTemplate](foxy, id, &cartTemplate)
	return result, e
}

func (foxy *CartTemplatesApi) Delete(id string) error {
	return DoDelete[*CartTemplate](foxy, id)
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
