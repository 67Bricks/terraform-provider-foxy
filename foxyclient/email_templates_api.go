package foxyclient

var (
	_ record   = &EmailTemplate{}
	_ foxyCrud = &EmailTemplatesApi{}
)

// ----

type EmailTemplatesApi struct {
	apiClient FoxyClient
}

func (foxy *EmailTemplatesApi) GetApiClient() FoxyClient {
	return foxy.apiClient
}

func (foxy *EmailTemplatesApi) List() ([]EmailTemplate, error) {
	path := foxy.storePath() + "/email_templates?limit=300"
	result, e := DoList[*EmailTemplate](foxy, path)
	return dereference(result), e
}

func (foxy *EmailTemplatesApi) Get(id string) (EmailTemplate, error) {
	path := "/email_templates/" + id
	result, e := DoGet[*EmailTemplate](foxy, path)
	return *result, e
}

func (foxy *EmailTemplatesApi) Add(emailTemplate EmailTemplate) (string, error) {
	path := foxy.storePath() + "/email_templates"
	result, e := DoAdd[*EmailTemplate](foxy, &emailTemplate, path)
	return result, e
}

func (foxy *EmailTemplatesApi) Update(id string, emailTemplate EmailTemplate) (string, error) {
	path := "/email_templates/" + id
	result, e := DoUpdate[*EmailTemplate](foxy, &emailTemplate, path)
	return result, e
}

func (foxy *EmailTemplatesApi) Delete(id string) error {
	path := "/email_templates/" + id
	return DoDelete[*EmailTemplate](foxy, path)
}

func (foxy *EmailTemplatesApi) storePath() string {
	storeId, _ := foxy.apiClient.retrieveStoreId()
	return "/stores/" + storeId
}

// ----

type EmailTemplate struct {
	Id             string `json:"-"`
	Description    string `json:"description"`
	Subject        string `json:"subject"`
	ContentHtml    string `json:"content_html"`
	ContentHtmlUrl string `json:"content_html_url"`
	ContentText    string `json:"content_text"`
	ContentTextUrl string `json:"content_text_url"`

	// @todo This and the setIdFromSelfUrl method are a clumsy way of unmarshalling the JSON - could we do this better?
	Links struct {
		Self struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
	} `json:"_links,omitempty"`
}

func (emailTemplate *EmailTemplate) setIdFromSelfUrl() {
	id := extractId(emailTemplate.Links.Self.Href)
	emailTemplate.Id = id
}
