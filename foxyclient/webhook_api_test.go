package foxyclient

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrieveWebhooks(t *testing.T) {
	foxy := newFoxy()
	webhooks, _ := foxy.Webhooks.List()
	require.Equal(t, "Test webhook", webhooks[0].Name)
	require.Equal(t, "https://example.com/webhook", webhooks[0].Url)
	require.Equal(t, "json", webhooks[0].Format)
	require.Equal(t, "334", webhooks[0].Id)
}

func TestAddAndDeleteWebhook(t *testing.T) {
	foxy := newFoxy()
	newWebhook := Webhook{
		Format:        "json",
		Version:       2,
		Name:          "New webhook",
		Url:           "https://example.com/new",
		Query:         "",
		EncryptionKey: "",
		EventResource: "transaction",
	}
	webhooks, _ := foxy.Webhooks.List()
	initialCount := len(webhooks)

	id, err := foxy.Webhooks.Add(newWebhook)
	require.Nil(t, err, "Error from adding should have been nil")
	require.NotEmpty(t, id, "ID should not be empty")
	createdWebhook, _ := foxy.Webhooks.Get(id)
	require.Equal(t, newWebhook.Name, createdWebhook.Name)

	webhooks, _ = foxy.Webhooks.List()
	newCount := len(webhooks)

	require.Equal(t, newCount, initialCount+1)
	err = foxy.Webhooks.Delete(id)
	require.Nil(t, err, "Error from deleting should have been nil")
	webhooks, _ = foxy.Webhooks.List()
	finalCount := len(webhooks)
	require.Equal(t, finalCount, initialCount)
}

func TestAddUpdateAndDeleteWebhook(t *testing.T) {
	foxy := newFoxy()
	newWebhook := Webhook{
		Format:        "json",
		Version:       2,
		Name:          "Another new webhook",
		Url:           "https://example.com/newer",
		Query:         "",
		EncryptionKey: "",
		EventResource: "transaction",
	}
	id, err := foxy.Webhooks.Add(newWebhook)
	require.Nil(t, err, "Error from adding should have been nil")
	newWebhook.Name = "Updated webhook"
	_, err = foxy.Webhooks.Update(id, newWebhook)
	require.Nil(t, err, "Error from updating should have been nil")
	createdWebhook, _ := foxy.Webhooks.Get(id)
	require.Equal(t, "Updated webhook", createdWebhook.Name)

	err = foxy.Webhooks.Delete(id)
	require.Nil(t, err, "Error from deleting should have been nil")
}
