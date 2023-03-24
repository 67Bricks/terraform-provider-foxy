package foxyclient

var (
	_ record   = &ReceiptTemplate{}
	_ foxyCrud = &ReceiptTemplatesApi{}
)

// ----

type ReceiptTemplatesApi struct {
	apiClient FoxyClient
}

func (foxy *ReceiptTemplatesApi) GetApiClient() FoxyClient {
	return foxy.apiClient
}

func (foxy *ReceiptTemplatesApi) List() ([]ReceiptTemplate, error) {
	path := foxy.storePath() + "/receipt_templates?limit=300"
	result, e := DoList[*ReceiptTemplate](foxy, path)
	return dereference(result), e
}

func (foxy *ReceiptTemplatesApi) Get(id string) (ReceiptTemplate, error) {
	path := "/receipt_templates/" + id
	result, e := DoGet[*ReceiptTemplate](foxy, path)
	return *result, e
}

func (foxy *ReceiptTemplatesApi) Add(receiptTemplate ReceiptTemplate) (string, error) {
	path := foxy.storePath() + "/receipt_templates"
	result, e := DoAdd[*ReceiptTemplate](foxy, &receiptTemplate, path)
	return result, e
}

func (foxy *ReceiptTemplatesApi) Update(id string, receiptTemplate ReceiptTemplate) (string, error) {
	path := "/receipt_templates/" + id
	result, e := DoUpdate[*ReceiptTemplate](foxy, &receiptTemplate, path)
	return result, e
}

func (foxy *ReceiptTemplatesApi) Delete(id string) error {
	path := "/receipt_templates/" + id
	return DoDelete[*ReceiptTemplate](foxy, path)
}

func (foxy *ReceiptTemplatesApi) storePath() string {
	storeId, _ := foxy.apiClient.retrieveStoreId()
	return "/stores/" + storeId
}

// ----

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
