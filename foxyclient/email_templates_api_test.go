package foxyclient

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetrieveEmailTemplates(t *testing.T) {
	foxy := newFoxy()
	emailTemplates, _ := foxy.EmailTemplates.List()
	assert.Equal(t, "Email Receipt Template", emailTemplates[0].Description)
	assert.Equal(t, "", emailTemplates[0].ContentHtml)
	assert.Equal(t, "", emailTemplates[0].ContentHtmlUrl)
	assert.NotEmpty(t, emailTemplates[0].Id)
}

func TestAddAndDeleteEmailTemplate(t *testing.T) {
	foxy := newFoxy()
	newEmailTemplate := EmailTemplate{
		Description:    "Some description",
		Subject:        "Some subject",
		ContentHtml:    "<p>Some content</p>",
		ContentHtmlUrl: "",
		ContentText:    "Some content",
		ContentTextUrl: "",
	}
	emailTemplates, _ := foxy.EmailTemplates.List()
	initialCount := len(emailTemplates)

	id, err := foxy.EmailTemplates.Add(newEmailTemplate)
	assert.Nil(t, err, "Error from adding should have been nil")
	assert.NotEmpty(t, id, "ID should not be empty")
	createdEmailTemplate, _ := foxy.EmailTemplates.Get(id)
	assert.Equal(t, newEmailTemplate.Description, createdEmailTemplate.Description)

	emailTemplates, _ = foxy.EmailTemplates.List()
	newCount := len(emailTemplates)

	assert.Equal(t, newCount, initialCount+1)
	err = foxy.EmailTemplates.Delete(id)
	assert.Nil(t, err, "Error from deleting should have been nil")
	emailTemplates, _ = foxy.EmailTemplates.List()
	finalCount := len(emailTemplates)
	assert.Equal(t, finalCount, initialCount)
}

func TestAddUpdateAndDeleteEmailTemplate(t *testing.T) {
	foxy := newFoxy()
	newEmailTemplate := EmailTemplate{
		Description:    "Some description",
		ContentHtml:    "<p>Some content</p>",
		ContentHtmlUrl: "",
		ContentText:    "Some content",
		ContentTextUrl: "",
	}
	id, err := foxy.EmailTemplates.Add(newEmailTemplate)
	assert.Nil(t, err, "Error from adding should have been nil")
	newEmailTemplate.Description = "Updated description"
	_, err = foxy.EmailTemplates.Update(id, newEmailTemplate)
	assert.Nil(t, err, "Error from updating should have been nil")
	createdEmailTemplate, _ := foxy.EmailTemplates.Get(id)
	assert.Equal(t, "Updated description", createdEmailTemplate.Description)

	err = foxy.EmailTemplates.Delete(id)
	assert.Nil(t, err, "Error from deleting should have been nil")
}
