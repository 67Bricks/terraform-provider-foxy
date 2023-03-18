package foxyclient

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrieveReceiptTemplates(t *testing.T) {
	foxy := newFoxy()
	receiptTemplates, _ := foxy.ReceiptTemplates.List()
	require.Equal(t, "Receipt Template", receiptTemplates[0].Description)
	require.Equal(t, "", receiptTemplates[0].Content)
	require.Equal(t, "", receiptTemplates[0].ContentUrl)
	require.NotEmpty(t, receiptTemplates[0].Id)
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
	require.Nil(t, err, "Error from adding should have been nil")
	require.NotEmpty(t, id, "ID should not be empty")
	createdReceiptTemplate, _ := foxy.ReceiptTemplates.Get(id)
	require.Equal(t, newReceiptTemplate.Description, createdReceiptTemplate.Description)

	receiptTemplates, _ = foxy.ReceiptTemplates.List()
	newCount := len(receiptTemplates)

	require.Equal(t, newCount, initialCount+1)
	err = foxy.ReceiptTemplates.Delete(id)
	require.Nil(t, err, "Error from deleting should have been nil")
	receiptTemplates, _ = foxy.ReceiptTemplates.List()
	finalCount := len(receiptTemplates)
	require.Equal(t, finalCount, initialCount)
}

func TestAddUpdateAndDeleteReceiptTemplate(t *testing.T) {
	foxy := newFoxy()
	newReceiptTemplate := ReceiptTemplate{
		Description: "Some description",
		Content:     "<p>Some content</p>",
		ContentUrl:  "",
	}
	id, err := foxy.ReceiptTemplates.Add(newReceiptTemplate)
	require.Nil(t, err, "Error from adding should have been nil")
	newReceiptTemplate.Description = "Updated description"
	_, err = foxy.ReceiptTemplates.Update(id, newReceiptTemplate)
	require.Nil(t, err, "Error from updating should have been nil")
	createdReceiptTemplate, _ := foxy.ReceiptTemplates.Get(id)
	require.Equal(t, "Updated description", createdReceiptTemplate.Description)

	err = foxy.ReceiptTemplates.Delete(id)
	require.Nil(t, err, "Error from deleting should have been nil")
}
