package foxyclient

import "encoding/json"

type StoreInfoApi struct {
	apiClient FoxyClient
}

func (foxy *StoreInfoApi) Get() (StoreInfo, error) {
	path := foxy.storePath()
	body, e := foxy.apiClient.get(path)
	var storeInfo StoreInfo
	e = json.Unmarshal(body, &storeInfo)
	return storeInfo, e
}

func (foxy *StoreInfoApi) Update(storeInfo StoreInfo) (string, error) {
	updateJson, _ := json.Marshal(storeInfo)
	path := foxy.storePath()
	body, e := foxy.apiClient.patch(path, string(updateJson))
	return string(body), e
}

func (foxy *StoreInfoApi) storePath() string {
	storeId, _ := foxy.apiClient.retrieveStoreId()
	return "/stores/" + storeId
}

type StoreInfo struct {
	StoreVersionUri                string      `json:"store_version_uri,omitempty"`
	StoreName                      string      `json:"store_name,omitempty"`
	StoreDomain                    string      `json:"store_domain,omitempty"`
	UseRemoteDomain                bool        `json:"use_remote_domain,omitempty"`
	StoreUrl                       string      `json:"store_url,omitempty"`
	ReceiptContinueUrl             string      `json:"receipt_continue_url,omitempty"`
	StoreEmail                     string      `json:"store_email,omitempty"`
	FromEmail                      string      `json:"from_email,omitempty"`
	UseEmailDns                    bool        `json:"use_email_dns,omitempty"`
	BccOnReceiptEmail              bool        `json:"bcc_on_receipt_email,omitempty"`
	SmtpConfig                     string      `json:"smtp_config,omitempty"`
	PostalCode                     string      `json:"postal_code,omitempty"`
	Region                         string      `json:"region,omitempty"`
	Country                        string      `json:"country,omitempty"`
	LocaleCode                     string      `json:"locale_code,omitempty"`
	Timezone                       string      `json:"timezone,omitempty"`
	HideCurrencySymbol             bool        `json:"hide_currency_symbol,omitempty"`
	HideDecimalCharacters          bool        `json:"hide_decimal_characters,omitempty"`
	UseInternationalCurrencySymbol bool        `json:"use_international_currency_symbol,omitempty"`
	Language                       string      `json:"language,omitempty"`
	LogoUrl                        string      `json:"logo_url,omitempty"`
	CheckoutType                   string      `json:"checkout_type,omitempty"`
	UseWebhook                     bool        `json:"use_webhook,omitempty"`
	WebhookUrl                     string      `json:"webhook_url,omitempty"`
	WebhookKey                     string      `json:"webhook_key,omitempty"`
	UseCartValidation              bool        `json:"use_cart_validation,omitempty"`
	UseSingleSignOn                bool        `json:"use_single_sign_on,omitempty"`
	SingleSignOnUrl                string      `json:"single_sign_on_url,omitempty"`
	CustomerPasswordHashType       string      `json:"customer_password_hash_type,omitempty"`
	CustomerPasswordHashConfig     string      `json:"customer_password_hash_config,omitempty"`
	FeaturesMultiship              bool        `json:"features_multiship,omitempty"`
	ProductsRequireExpiresProperty bool        `json:"products_require_expires_property,omitempty"`
	AppSessionTime                 int         `json:"app_session_time,omitempty"`
	ShippingAddressType            string      `json:"shipping_address_type,omitempty"`
	RequireSignedShippingRates     bool        `json:"require_signed_shipping_rates,omitempty"`
	UnifiedOrderEntryPassword      string      `json:"unified_order_entry_password,omitempty"`
	CustomDisplayIdConfig          interface{} `json:"custom_display_id_config,omitempty"`
	AffiliateId                    int         `json:"affiliate_id,omitempty"`
	IsMaintenanceMode              bool        `json:"is_maintenance_mode,omitempty"`
	IsActive                       bool        `json:"is_active,omitempty"`
	FirstPaymentDate               interface{} `json:"first_payment_date,omitempty"`
	Features                       interface{} `json:"features,omitempty"`
}
