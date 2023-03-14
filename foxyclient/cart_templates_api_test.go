package foxyclient

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetrieveCartTemplates(t *testing.T) {
	foxy := newFoxy()
	cartTemplates, _ := foxy.CartTemplates.List()
	assert.Equal(t, "Cart Template", cartTemplates[0].Description)
	assert.Equal(t, "", cartTemplates[0].Content)
	assert.Equal(t, "", cartTemplates[0].ContentUrl)
	assert.NotEmpty(t, cartTemplates[0].Id)
}

func TestAddAndDeleteCartTemplate(t *testing.T) {
	foxy := newFoxy()
	newCartTemplate := CartTemplate{
		Description: "Some description",
		Content:     "<p>Some content</p>",
		ContentUrl:  "",
	}
	cartTemplates, _ := foxy.CartTemplates.List()
	initialCount := len(cartTemplates)

	id, err := foxy.CartTemplates.Add(newCartTemplate)
	assert.Nil(t, err, "Error from adding should have been nil")
	assert.NotEmpty(t, id, "ID should not be empty")
	createdCartTemplate, _ := foxy.CartTemplates.Get(id)
	assert.Equal(t, newCartTemplate.Description, createdCartTemplate.Description)

	cartTemplates, _ = foxy.CartTemplates.List()
	newCount := len(cartTemplates)

	assert.Equal(t, newCount, initialCount+1)
	err = foxy.CartTemplates.Delete(id)
	assert.Nil(t, err, "Error from deleting should have been nil")
	cartTemplates, _ = foxy.CartTemplates.List()
	finalCount := len(cartTemplates)
	assert.Equal(t, finalCount, initialCount)
}

func TestAddUpdateAndDeleteCartTemplate(t *testing.T) {
	foxy := newFoxy()
	newCartTemplate := CartTemplate{
		Description: "Some description",
		Content:     "<p>Some content</p>",
		ContentUrl:  "",
	}
	id, err := foxy.CartTemplates.Add(newCartTemplate)
	assert.Nil(t, err, "Error from adding should have been nil")
	newCartTemplate.Description = "Updated description"
	_, err = foxy.CartTemplates.Update(id, newCartTemplate)
	assert.Nil(t, err, "Error from updating should have been nil")
	createdCartTemplate, _ := foxy.CartTemplates.Get(id)
	assert.Equal(t, "Updated description", createdCartTemplate.Description)

	err = foxy.CartTemplates.Delete(id)
	assert.Nil(t, err, "Error from deleting should have been nil")
}
