package foxyclient

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newFoxy() Foxy {
	c := readConfig()
	foxy, _ := New(c.BaseUrl, c.ClientID, c.ClientSecret, c.RefreshToken)
	return foxy
}

func TestReadConfig(t *testing.T) {
	c := readConfig()
	assert.Equal(t, "client_1Q6iX3A1UjKNUZxEeV7P", c.ClientID)
	assert.Len(t, c.ClientSecret, 40) // Not asserting what it is, but it must be set
}

func TestRetrieveToken(t *testing.T) {
	conf := readConfig()
	foxy := FoxyHttpClient{baseUrl: conf.BaseUrl}
	result, err := foxy.retrieveToken(conf.ClientID, conf.ClientSecret, conf.RefreshToken)
	assert.Nil(t, err, "Should not have had error")
	assert.NotEmpty(t, result.AccessToken)
}

func TestRetrieveStoreInfo(t *testing.T) {
	foxy := newFoxy()
	storeInfo, _ := foxy.StoreInfo.Get()
	assert.Equal(t, "Terraform Test", storeInfo.StoreName)
}

func TestSetStoreInfo(t *testing.T) {
	foxy := newFoxy()

	_, _ = foxy.StoreInfo.Update(StoreInfo{Language: "english"})
	initialStoreInfo, _ := foxy.StoreInfo.Get()
	assert.Equal(t, "english", initialStoreInfo.Language)

	_, _ = foxy.StoreInfo.Update(StoreInfo{Language: "german"})
	updatedStoreInfo, _ := foxy.StoreInfo.Get()
	assert.Equal(t, "german", updatedStoreInfo.Language)
}

func TestConvertingStoreInfoToJson(t *testing.T) {
	storeInfo := StoreInfo{StoreName: "fish store"}
	bytes, _ := json.Marshal(storeInfo)
	assert.Equal(t, `{"store_name":"fish store"}`, string(bytes))
}

func TestRetrieveWebhooks(t *testing.T) {
	foxy := newFoxy()
	webhooks, _ := foxy.Webhooks.List()
	assert.Equal(t, "Test webhook", webhooks[0].Name)
	assert.Equal(t, "https://example.com/webhook", webhooks[0].Url)
	assert.Equal(t, "json", webhooks[0].Format)
	assert.Equal(t, "334", webhooks[0].Id)
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
	assert.Nil(t, err, "Error from adding should have been nil")
	assert.NotEmpty(t, id, "ID should not be empty")
	createdWebhook, _ := foxy.Webhooks.Get(id)
	assert.Equal(t, newWebhook.Name, createdWebhook.Name)

	webhooks, _ = foxy.Webhooks.List()
	newCount := len(webhooks)

	assert.Equal(t, newCount, initialCount+1)
	err = foxy.Webhooks.Delete(id)
	assert.Nil(t, err, "Error from deleting should have been nil")
	webhooks, _ = foxy.Webhooks.List()
	finalCount := len(webhooks)
	assert.Equal(t, finalCount, initialCount)
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
	assert.Nil(t, err, "Error from adding should have been nil")
	newWebhook.Name = "Updated webhook"
	_, err = foxy.Webhooks.Update(id, newWebhook)
	assert.Nil(t, err, "Error from updating should have been nil")
	createdWebhook, _ := foxy.Webhooks.Get(id)
	assert.Equal(t, "Updated webhook", createdWebhook.Name)

	err = foxy.Webhooks.Delete(id)
	assert.Nil(t, err, "Error from deleting should have been nil")
}
