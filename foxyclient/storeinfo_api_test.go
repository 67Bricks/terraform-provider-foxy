package foxyclient

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrieveStoreInfo(t *testing.T) {
	foxy := newFoxy()
	storeInfo, _ := foxy.StoreInfo.Get()
	require.Equal(t, "Terraform Test", storeInfo.StoreName)
}

func TestSetStoreInfo(t *testing.T) {
	foxy := newFoxy()

	_, _ = foxy.StoreInfo.Update(StoreInfo{Language: "english"})
	initialStoreInfo, _ := foxy.StoreInfo.Get()
	require.Equal(t, "english", initialStoreInfo.Language)

	_, _ = foxy.StoreInfo.Update(StoreInfo{Language: "german"})
	updatedStoreInfo, _ := foxy.StoreInfo.Get()
	require.Equal(t, "german", updatedStoreInfo.Language)
}

func TestConvertingStoreInfoToJson(t *testing.T) {
	storeInfo := StoreInfo{StoreName: "fish store"}
	bytes, _ := json.Marshal(storeInfo)
	require.Equal(t, `{"store_name":"fish store"}`, string(bytes))
}
