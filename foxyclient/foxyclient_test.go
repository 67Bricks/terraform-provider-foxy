package foxyclient

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfig(t *testing.T) {
	c := readConfig()
	assert.Equal(t, "client_1Q6iX3A1UjKNUZxEeV7P", c.ClientID)
}

func TestRetrieveToken(t *testing.T) {
	conf := readConfig()
	foxy := Foxy{baseUrl: conf.BaseUrl}
	result, err := foxy.retrieveToken(conf.ClientID, conf.ClientSecret, conf.RefreshToken)
	assert.Nil(t, err, "Should not have had error")
	assert.NotEmpty(t, result.AccessToken)
}

func TestRetrieveStoreInfo(t *testing.T) {
	storeInfo := RetrieveStoreInfo()
	assert.Equal(t, "Terraform Test", storeInfo.StoreName)
}

func TestSetStoreInfo(t *testing.T) {
	foxy := CreateFoxy()

	_, _ = foxy.UpdateStore(StoreInfo{Language: "english"})
	initialStoreInfo, _ := foxy.GetStore()
	assert.Equal(t, "english", initialStoreInfo.Language)

	_, _ = foxy.UpdateStore(StoreInfo{Language: "german"})
	updatedStoreInfo, _ := foxy.GetStore()
	assert.Equal(t, "german", updatedStoreInfo.Language)
}

func TestConvertingStoreInfoToJson(t *testing.T) {
	storeInfo := StoreInfo{StoreName: "fish store"}
	bytes, _ := json.Marshal(storeInfo)
	assert.Equal(t, `{"store_name":"fish store"}`, string(bytes))
}

func TestRetrieveWebhooks(t *testing.T) {
	foxy := CreateFoxy()
	webhooks, _ := foxy.GetWebhooks()
	assert.Equal(t, "Test webhook", webhooks[0].Name)
	assert.Equal(t, "https://example.com/webhook", webhooks[0].Url)
	assert.Equal(t, "https://api.foxycart.com/webhooks/334", webhooks[0].SelfUrl())
	assert.Equal(t, "334", webhooks[0].GetId())
}

func TestWebhookId(t *testing.T) {
	wh := Webhook{}
	assert.Equal(t, "", wh.GetId())
}

func TestAddAndDeleteWebhook(t *testing.T) {
	foxy := CreateFoxy()
	newWebhook := Webhook{
		Format:        "json",
		Version:       2,
		Name:          "New webhook",
		Url:           "https://example.com/new",
		Query:         "",
		EncryptionKey: "",
		EventResource: "transaction",
		//Events:        []string{"transaction/created", "transaction/captured"},
	}
	webhooks, _ := foxy.GetWebhooks()
	initialCount := len(webhooks)

	id, err := foxy.AddWebhook(newWebhook)
	assert.Nil(t, err, "Error from adding should have been nil")
	assert.NotEmpty(t, id, "ID should not be empty")
	createdWebhook, _ := foxy.GetWebhook(id)
	assert.Equal(t, newWebhook.Name, createdWebhook.Name)

	webhooks, _ = foxy.GetWebhooks()
	newCount := len(webhooks)

	assert.Equal(t, newCount, initialCount+1)
	err = foxy.DeleteWebhook(id)
	assert.Nil(t, err, "Error from deleting should have been nil")
	webhooks, _ = foxy.GetWebhooks()
	finalCount := len(webhooks)
	assert.Equal(t, finalCount, initialCount)
}

func TestAddUpdateAndDeleteWebhook(t *testing.T) {
	foxy := CreateFoxy()
	newWebhook := Webhook{
		Format:        "json",
		Version:       2,
		Name:          "Another new webhook",
		Url:           "https://example.com/newer",
		Query:         "",
		EncryptionKey: "",
		EventResource: "transaction",
	}
	id, err := foxy.AddWebhook(newWebhook)
	assert.Nil(t, err, "Error from adding should have been nil")
	newWebhook.Name = "Updated webhook"
	_, err = foxy.UpdateWebhook(id, newWebhook)
	assert.Nil(t, err, "Error from updating should have been nil")
	createdWebhook, _ := foxy.GetWebhook(id)
	assert.Equal(t, "Updated webhook", createdWebhook.Name)

	err = foxy.DeleteWebhook(id)
	assert.Nil(t, err, "Error from deleting should have been nil")
}
