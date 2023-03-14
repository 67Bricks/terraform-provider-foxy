package foxyclient

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetrieveCheckoutTemplates(t *testing.T) {
	foxy := newFoxy()
	checkoutTemplates, _ := foxy.CheckoutTemplates.List()
	assert.Equal(t, "Checkout Template", checkoutTemplates[0].Description)
	assert.Equal(t, "", checkoutTemplates[0].Content)
	assert.Equal(t, "", checkoutTemplates[0].ContentUrl)
	assert.NotEmpty(t, checkoutTemplates[0].Id)
}

func TestAddAndDeleteCheckoutTemplate(t *testing.T) {
	foxy := newFoxy()
	newCheckoutTemplate := CheckoutTemplate{
		Description: "Some description",
		Content:     "<p>Some content</p>",
		ContentUrl:  "",
	}
	checkoutTemplates, _ := foxy.CheckoutTemplates.List()
	initialCount := len(checkoutTemplates)

	id, err := foxy.CheckoutTemplates.Add(newCheckoutTemplate)
	assert.Nil(t, err, "Error from adding should have been nil")
	assert.NotEmpty(t, id, "ID should not be empty")
	createdCheckoutTemplate, _ := foxy.CheckoutTemplates.Get(id)
	assert.Equal(t, newCheckoutTemplate.Description, createdCheckoutTemplate.Description)

	checkoutTemplates, _ = foxy.CheckoutTemplates.List()
	newCount := len(checkoutTemplates)

	assert.Equal(t, newCount, initialCount+1)
	err = foxy.CheckoutTemplates.Delete(id)
	assert.Nil(t, err, "Error from deleting should have been nil")
	checkoutTemplates, _ = foxy.CheckoutTemplates.List()
	finalCount := len(checkoutTemplates)
	assert.Equal(t, finalCount, initialCount)
}

func TestAddUpdateAndDeleteCheckoutTemplate(t *testing.T) {
	foxy := newFoxy()
	newCheckoutTemplate := CheckoutTemplate{
		Description: "Some description",
		Content:     "<p>Some content</p>",
		ContentUrl:  "",
	}
	id, err := foxy.CheckoutTemplates.Add(newCheckoutTemplate)
	assert.Nil(t, err, "Error from adding should have been nil")
	newCheckoutTemplate.Description = "Updated description"
	_, err = foxy.CheckoutTemplates.Update(id, newCheckoutTemplate)
	assert.Nil(t, err, "Error from updating should have been nil")
	createdCheckoutTemplate, _ := foxy.CheckoutTemplates.Get(id)
	assert.Equal(t, "Updated description", createdCheckoutTemplate.Description)

	err = foxy.CheckoutTemplates.Delete(id)
	assert.Nil(t, err, "Error from deleting should have been nil")
}
