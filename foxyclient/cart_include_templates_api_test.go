package foxyclient

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetrieveCartIncludeTemplates(t *testing.T) {
	foxy := newFoxy()
	cartIncludeTemplates, _ := foxy.CartIncludeTemplates.List()
	assert.Equal(t, "Cart Include Template", cartIncludeTemplates[0].Description)
	assert.Equal(t, "", cartIncludeTemplates[0].Content)
	assert.Equal(t, "", cartIncludeTemplates[0].ContentUrl)
	assert.NotEmpty(t, cartIncludeTemplates[0].Id)
}

func TestAddAndDeleteCartIncludeTemplate(t *testing.T) {
	foxy := newFoxy()
	newCartIncludeTemplate := CartIncludeTemplate{
		Description: "Some description",
		Content:     "<p>Some content</p>",
		ContentUrl:  "",
	}
	cartIncludeTemplates, _ := foxy.CartIncludeTemplates.List()
	initialCount := len(cartIncludeTemplates)

	id, err := foxy.CartIncludeTemplates.Add(newCartIncludeTemplate)
	assert.Nil(t, err, "Error from adding should have been nil")
	assert.NotEmpty(t, id, "ID should not be empty")
	createdCartIncludeTemplate, _ := foxy.CartIncludeTemplates.Get(id)
	assert.Equal(t, newCartIncludeTemplate.Description, createdCartIncludeTemplate.Description)

	cartIncludeTemplates, _ = foxy.CartIncludeTemplates.List()
	newCount := len(cartIncludeTemplates)

	assert.Equal(t, newCount, initialCount+1)
	err = foxy.CartIncludeTemplates.Delete(id)
	assert.Nil(t, err, "Error from deleting should have been nil")
	cartIncludeTemplates, _ = foxy.CartIncludeTemplates.List()
	finalCount := len(cartIncludeTemplates)
	assert.Equal(t, finalCount, initialCount)
}

func TestAddUpdateAndDeleteCartIncludeTemplate(t *testing.T) {
	foxy := newFoxy()
	newCartIncludeTemplate := CartIncludeTemplate{
		Description: "Some description",
		Content:     "<p>Some content</p>",
		ContentUrl:  "",
	}
	id, err := foxy.CartIncludeTemplates.Add(newCartIncludeTemplate)
	assert.Nil(t, err, "Error from adding should have been nil")
	newCartIncludeTemplate.Description = "Updated description"
	_, err = foxy.CartIncludeTemplates.Update(id, newCartIncludeTemplate)
	assert.Nil(t, err, "Error from updating should have been nil")
	createdCartIncludeTemplate, _ := foxy.CartIncludeTemplates.Get(id)
	assert.Equal(t, "Updated description", createdCartIncludeTemplate.Description)

	err = foxy.CartIncludeTemplates.Delete(id)
	assert.Nil(t, err, "Error from deleting should have been nil")
}
