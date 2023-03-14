package foxyclient

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type EmailTemplatesApi struct {
	apiClient FoxyClient
}

func (foxy *EmailTemplatesApi) List() ([]EmailTemplate, error) {
	// This is not retrieving all emailTemplates - only the first 300 - but is it plausible to have more than 300 emailTemplates?
	path := foxy.storePath() + "/email_templates?limit=300"
	body, err := foxy.apiClient.get(path)
	if err != nil {
		return nil, err
	}
	var emailTemplates []EmailTemplate
	embeddedJsonResult := gjson.GetBytes(body, "_embedded.fx:email_templates")
	embeddedJson := []byte(embeddedJsonResult.Raw)
	err = json.Unmarshal(embeddedJson, &emailTemplates)
	if err != nil {
		return nil, err
	}
	for i := range emailTemplates {
		// Need to modify emailTemplates[i], rather than accessing wh directly via the loop, because the latter is by value
		emailTemplates[i].setIdFromSelfUrl()
	}

	return emailTemplates, err
}

func (foxy *EmailTemplatesApi) Get(id string) (EmailTemplate, error) {
	path := foxy.emailTemplatePath(id)
	body, err := foxy.apiClient.get(path)
	if err != nil {
		return EmailTemplate{}, err
	}
	var emailTemplate EmailTemplate
	err = json.Unmarshal(body, &emailTemplate)
	if err != nil {
		return EmailTemplate{}, err
	}
	emailTemplate.setIdFromSelfUrl()
	return emailTemplate, err
}

func (foxy *EmailTemplatesApi) Add(emailTemplate EmailTemplate) (string, error) {
	path := foxy.storePath() + "/email_templates"
	updateJson, _ := json.Marshal(emailTemplate)
	result, err := foxy.apiClient.post(path, string(updateJson))
	if err != nil {
		return "", err
	}
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	id := extractId(selfUrl)
	return id, err
}

func (foxy *EmailTemplatesApi) Update(id string, emailTemplate EmailTemplate) (string, error) {
	path := foxy.emailTemplatePath(id)
	amendedEmailTemplate := emailTemplate
	updateJson, _ := json.Marshal(amendedEmailTemplate)
	result, e := foxy.apiClient.patch(path, string(updateJson))
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	updatedId := extractId(selfUrl)
	return updatedId, e
}

func (foxy *EmailTemplatesApi) Delete(id string) error {
	path := foxy.emailTemplatePath(id)
	_, e := foxy.apiClient.delete(path)
	return e
}

func (foxy *EmailTemplatesApi) emailTemplatePath(id string) string {
	return "/email_templates/" + id
}

func (foxy *EmailTemplatesApi) storePath() string {
	storeId, _ := foxy.apiClient.retrieveStoreId()
	return "/stores/" + storeId
}

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
