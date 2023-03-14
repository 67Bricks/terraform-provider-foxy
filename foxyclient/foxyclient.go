package foxyclient

import (
	"fmt"
	"strings"
)

type Foxy struct {
	StoreInfo            StoreInfoApi
	Webhooks             WebhooksApi
	CartTemplates        CartTemplatesApi
	CartIncludeTemplates CartIncludeTemplatesApi
	CheckoutTemplates    CheckoutTemplatesApi
	ReceiptTemplates     ReceiptTemplatesApi
	EmailTemplates       EmailTemplatesApi
}

func New(baseUrl string, clientId string, clientSecret string, refreshToken string) (Foxy, error) {
	if clientSecret == "" {
		return Foxy{}, fmt.Errorf("missing client secret")
	}
	apiClient, err := newFoxyClient(baseUrl, clientId, clientSecret, refreshToken)
	if err != nil {
		return Foxy{}, err
	}
	foxy := Foxy{
		StoreInfo:            StoreInfoApi{apiClient: &apiClient},
		Webhooks:             WebhooksApi{apiClient: &apiClient},
		CartTemplates:        CartTemplatesApi{apiClient: &apiClient},
		CartIncludeTemplates: CartIncludeTemplatesApi{apiClient: &apiClient},
		CheckoutTemplates:    CheckoutTemplatesApi{apiClient: &apiClient},
		ReceiptTemplates:     ReceiptTemplatesApi{apiClient: &apiClient},
		EmailTemplates:       EmailTemplatesApi{apiClient: &apiClient},
	}
	return foxy, nil
}

// -------
// -------

func extractId(selfUrl string) string {
	parts := strings.Split(selfUrl, "/")
	return parts[len(parts)-1]
}
