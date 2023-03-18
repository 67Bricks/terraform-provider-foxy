package foxyclient

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrieveCartTemplates(t *testing.T) {
	foxy := newFoxy()
	cartTemplates, _ := foxy.CartTemplates.List()
	require.Equal(t, "Cart Template", cartTemplates[0].Description)
	require.Equal(t, "", cartTemplates[0].Content)
	require.Equal(t, "", cartTemplates[0].ContentUrl)
	require.NotEmpty(t, cartTemplates[0].Id)
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
	require.Nil(t, err, "Error from adding should have been nil")
	require.NotEmpty(t, id, "ID should not be empty")
	createdCartTemplate, _ := foxy.CartTemplates.Get(id)
	require.Equal(t, newCartTemplate.Description, createdCartTemplate.Description)

	cartTemplates, _ = foxy.CartTemplates.List()
	newCount := len(cartTemplates)

	require.Equal(t, newCount, initialCount+1)
	err = foxy.CartTemplates.Delete(id)
	require.Nil(t, err, "Error from deleting should have been nil")
	cartTemplates, _ = foxy.CartTemplates.List()
	finalCount := len(cartTemplates)
	require.Equal(t, finalCount, initialCount)
}

func TestAddUpdateAndDeleteCartTemplate(t *testing.T) {
	foxy := newFoxy()
	newCartTemplate := CartTemplate{
		Description: "Some description",
		Content:     "<p>Some content</p>",
		ContentUrl:  "",
	}
	id, err := foxy.CartTemplates.Add(newCartTemplate)
	require.Nil(t, err, "Error from adding should have been nil")
	newCartTemplate.Description = "Updated description"
	_, err = foxy.CartTemplates.Update(id, newCartTemplate)
	require.Nil(t, err, "Error from updating should have been nil")
	createdCartTemplate, _ := foxy.CartTemplates.Get(id)
	require.Equal(t, "Updated description", createdCartTemplate.Description)

	err = foxy.CartTemplates.Delete(id)
	require.Nil(t, err, "Error from deleting should have been nil")
}
