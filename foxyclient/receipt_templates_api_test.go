package foxyclient

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetrieveReceiptTemplates(t *testing.T) {
	foxy := newFoxy()
	receiptTemplates, _ := foxy.ReceiptTemplates.List()
	assert.Equal(t, "Receipt Template", receiptTemplates[0].Description)
	assert.Equal(t, "", receiptTemplates[0].Content)
	assert.Equal(t, "", receiptTemplates[0].ContentUrl)
	assert.NotEmpty(t, receiptTemplates[0].Id)
}

func TestAddAndDeleteReceiptTemplate(t *testing.T) {
	foxy := newFoxy()
	newReceiptTemplate := ReceiptTemplate{
		Description: "Some description",
		Content:     "<p>Some content</p>",
		ContentUrl:  "",
	}
	receiptTemplates, _ := foxy.ReceiptTemplates.List()
	initialCount := len(receiptTemplates)

	id, err := foxy.ReceiptTemplates.Add(newReceiptTemplate)
	assert.Nil(t, err, "Error from adding should have been nil")
	assert.NotEmpty(t, id, "ID should not be empty")
	createdReceiptTemplate, _ := foxy.ReceiptTemplates.Get(id)
	assert.Equal(t, newReceiptTemplate.Description, createdReceiptTemplate.Description)

	receiptTemplates, _ = foxy.ReceiptTemplates.List()
	newCount := len(receiptTemplates)

	assert.Equal(t, newCount, initialCount+1)
	err = foxy.ReceiptTemplates.Delete(id)
	assert.Nil(t, err, "Error from deleting should have been nil")
	receiptTemplates, _ = foxy.ReceiptTemplates.List()
	finalCount := len(receiptTemplates)
	assert.Equal(t, finalCount, initialCount)
}

func TestAddUpdateAndDeleteReceiptTemplate(t *testing.T) {
	foxy := newFoxy()
	newReceiptTemplate := ReceiptTemplate{
		Description: "Some description",
		Content:     "<p>Some content</p>",
		ContentUrl:  "",
	}
	id, err := foxy.ReceiptTemplates.Add(newReceiptTemplate)
	assert.Nil(t, err, "Error from adding should have been nil")
	newReceiptTemplate.Description = "Updated description"
	_, err = foxy.ReceiptTemplates.Update(id, newReceiptTemplate)
	assert.Nil(t, err, "Error from updating should have been nil")
	createdReceiptTemplate, _ := foxy.ReceiptTemplates.Get(id)
	assert.Equal(t, "Updated description", createdReceiptTemplate.Description)

	err = foxy.ReceiptTemplates.Delete(id)
	assert.Nil(t, err, "Error from deleting should have been nil")
}
