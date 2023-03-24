package foxyclient

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrieveEmailTemplates(t *testing.T) {
	foxy := newFoxy()
	emailTemplates, _ := foxy.EmailTemplates.List()
	require.Equal(t, "Email Receipt Template", emailTemplates[0].Description)
	require.Equal(t, "<p>New email template</p>", emailTemplates[0].ContentHtml)
	require.Equal(t, "", emailTemplates[0].ContentHtmlUrl)
	require.NotEmpty(t, emailTemplates[0].Id)
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
	require.Nil(t, err, "Error from adding should have been nil")
	require.NotEmpty(t, id, "ID should not be empty")
	createdEmailTemplate, _ := foxy.EmailTemplates.Get(id)
	require.Equal(t, newEmailTemplate.Description, createdEmailTemplate.Description)

	emailTemplates, _ = foxy.EmailTemplates.List()
	newCount := len(emailTemplates)

	require.Equal(t, newCount, initialCount+1)
	err = foxy.EmailTemplates.Delete(id)
	require.Nil(t, err, "Error from deleting should have been nil")
	emailTemplates, _ = foxy.EmailTemplates.List()
	finalCount := len(emailTemplates)
	require.Equal(t, finalCount, initialCount)
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
	require.Nil(t, err, "Error from adding should have been nil")
	newEmailTemplate.Description = "Updated description"
	_, err = foxy.EmailTemplates.Update(id, newEmailTemplate)
	require.Nil(t, err, "Error from updating should have been nil")
	createdEmailTemplate, _ := foxy.EmailTemplates.Get(id)
	require.Equal(t, "Updated description", createdEmailTemplate.Description)

	err = foxy.EmailTemplates.Delete(id)
	require.Nil(t, err, "Error from deleting should have been nil")
}
