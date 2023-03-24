package foxyclient

var (
	_ record   = &CheckoutTemplate{}
	_ foxyCrud = &CheckoutTemplatesApi{}
)

// ----

type CheckoutTemplatesApi struct {
	apiClient FoxyClient
}

func (foxy *CheckoutTemplatesApi) GetApiClient() FoxyClient {
	return foxy.apiClient
}

func (foxy *CheckoutTemplatesApi) List() ([]CheckoutTemplate, error) {
	path := foxy.storePath() + "/checkout_templates?limit=300"
	result, e := DoList[*CheckoutTemplate](foxy, path)
	return dereference(result), e
}

func (foxy *CheckoutTemplatesApi) Get(id string) (CheckoutTemplate, error) {
	path := "/checkout_templates/" + id
	result, e := DoGet[*CheckoutTemplate](foxy, path)
	return *result, e
}

func (foxy *CheckoutTemplatesApi) Add(checkoutTemplate CheckoutTemplate) (string, error) {
	path := foxy.storePath() + "/checkout_templates"
	result, e := DoAdd[*CheckoutTemplate](foxy, &checkoutTemplate, path)
	return result, e
}

func (foxy *CheckoutTemplatesApi) Update(id string, checkoutTemplate CheckoutTemplate) (string, error) {
	path := "/checkout_templates/" + id
	result, e := DoUpdate[*CheckoutTemplate](foxy, &checkoutTemplate, path)
	return result, e
}

func (foxy *CheckoutTemplatesApi) Delete(id string) error {
	path := "/checkout_templates/" + id
	return DoDelete[*CheckoutTemplate](foxy, path)
}

func (foxy *CheckoutTemplatesApi) storePath() string {
	storeId, _ := foxy.apiClient.retrieveStoreId()
	return "/stores/" + storeId
}

// ----

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
