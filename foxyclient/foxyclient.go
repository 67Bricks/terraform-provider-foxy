package foxyclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"log"
	"strings"
)

type FoxyConfig struct {
	ClientID     string `mapstructure:"clientid"`
	ClientSecret string `mapstructure:"clientsecret"`
	RefreshToken string `mapstructure:"refreshtoken"`
	BaseUrl      string `mapstructure:"baseurl"`
}

func readConfig() FoxyConfig {
	viper.SetEnvPrefix("FOXY") // So the env variable "FOXY_CLIENTSECRET" can be used to set the client secret
	viper.AutomaticEnv()
	viper.SetConfigFile("config.toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Unable to read configuration file: %v", err)
	}

	//var config FoxyConfig
	//if err := viper.Unmarshal(&config); err != nil {
	//	log.Fatalf("Unable to parse configuration file: %v", err)
	//}

	// Reading the config manually, rather than via the unmarshalling above, so the environment variable override works
	config := FoxyConfig{
		ClientID:     viper.GetString("clientid"),
		ClientSecret: viper.GetString("clientsecret"),
		RefreshToken: viper.GetString("refreshtoken"),
		BaseUrl:      viper.GetString("baseurl"),
	}

	return config
}

func RetrieveStoreInfo() StoreInfo {
	config := readConfig()
	foxy := Foxy{baseUrl: config.BaseUrl}
	_ = foxy.setToken(config.ClientID, config.ClientSecret, config.RefreshToken)
	var store, _ = foxy.GetStore()
	return store
}

// -------
// -------

type Foxy struct {
	token    oauth2.Token
	baseUrl  string
	storeUrl string
}

func New(baseUrl string, clientId string, clientSecret string, refreshToken string) (Foxy, error) {
	foxy := Foxy{baseUrl: baseUrl}
	err := foxy.setToken(clientId, clientSecret, refreshToken)
	return foxy, err
}

// CreateFoxy should maybe be renamed to New()? @todo initializer?
func CreateFoxy() Foxy {
	config := readConfig()
	foxy := Foxy{baseUrl: config.BaseUrl}
	_ = foxy.setToken(config.ClientID, config.ClientSecret, config.RefreshToken)
	return foxy
}

func (foxy *Foxy) setToken(clientId string, clientSecret string, refreshToken string) error {
	token, err := foxy.retrieveToken(clientId, clientSecret, refreshToken)
	if err != nil {
		fmt.Println("Token cannot be retrieved: " + err.Error())
		return err
	}
	foxy.token = token
	return nil
}

func (foxy *Foxy) retrieveToken(clientId string, clientSecret string, refreshToken string) (oauth2.Token, error) {
	data := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     clientId,
		"client_secret": clientSecret,
	}

	result, err := resty.New().R().SetFormData(data).Post(foxy.baseUrl + "/token")
	if err != nil {
		return oauth2.Token{}, err
	}

	var token oauth2.Token
	err = json.Unmarshal(result.Body(), &token)

	return token, err
}

func (foxy *Foxy) getFromFoxy(path string) ([]byte, error) {
	url := foxy.toUrl(path)
	result, err := foxy.createClient().Get(url)
	fmt.Println("Result is " + string(result.Body()))
	return result.Body(), err
}

func (foxy *Foxy) patchFoxy(path string, body string) ([]byte, error) {
	url := foxy.toUrl(path)
	result, err := foxy.createClient().SetBody(body).Patch(url)
	fmt.Println("Result is " + string(result.Body()))
	return result.Body(), err
}

func (foxy *Foxy) postToFoxy(path string, body string) ([]byte, error) {
	url := foxy.toUrl(path)
	fmt.Println("Request for posting to " + path + " is " + body)
	result, err := foxy.createClient().SetBody(body).Post(url)
	fmt.Println("Results of posting to " + path + " is " + string(result.Body()))
	return result.Body(), err
}

func (foxy *Foxy) putToFoxy(path string, body string) ([]byte, error) {
	url := foxy.toUrl(path)
	result, err := foxy.createClient().SetBody(body).Put(url)
	fmt.Println("Result " + string(result.Body()))
	fmt.Println("---xxx---")
	fmt.Println("---xxx---")

	return result.Body(), err
}

func (foxy *Foxy) deleteFromFoxy(path string) ([]byte, error) {
	url := foxy.toUrl(path)
	fmt.Println("Sending delete to " + path)
	result, err := foxy.createClient().Delete(url)
	fmt.Println("Results of delete are " + string(result.Body()))
	return result.Body(), err
}

func (foxy *Foxy) toUrl(path string) string {
	var url = path
	if strings.Index(path, foxy.baseUrl) != 0 {
		url = foxy.baseUrl + path
	}
	return url
}

func (foxy *Foxy) createClient() *resty.Request {
	// Resty docs - https://github.com/go-resty/resty
	oauthClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&foxy.token))
	client := resty.NewWithClient(oauthClient)
	//client.SetDebug(true)
	result := client.R().
		SetHeader("FOXY-API-VERSION", "1").
		SetAuthToken(foxy.token.AccessToken)
	return result
}

func (foxy *Foxy) GetStore() (StoreInfo, error) {
	storeUrl := foxy.ensureStoreUrl()
	body, e := foxy.getFromFoxy(storeUrl)
	var storeInfo StoreInfo
	e = json.Unmarshal(body, &storeInfo)
	return storeInfo, e
}

func (foxy *Foxy) UpdateStore(storeInfo StoreInfo) (string, error) {
	storeUrl := foxy.ensureStoreUrl()

	updateJson, _ := json.Marshal(storeInfo)
	body, e := foxy.patchFoxy(storeUrl, string(updateJson))

	return string(body), e
}

func (foxy *Foxy) GetWebhooks() ([]Webhook, error) {
	storeUrl := foxy.ensureStoreUrl()
	// @todo This is not retrieving all webhooks - only the first 300 - needs to iterate through the list if more
	fmt.Println("Retrieving webhooks from " + storeUrl + "/webhooks")
	body, e := foxy.getFromFoxy(storeUrl + "/webhooks?limit=300")
	var webhooks []Webhook
	embeddedJsonResult := gjson.GetBytes(body, "_embedded.fx:webhooks")
	embeddedJson := []byte(embeddedJsonResult.Raw)
	e = json.Unmarshal(embeddedJson, &webhooks)
	for _, wh := range webhooks {
		wh.SetIdFromSelfUrl()
	}
	return webhooks, e
}

func (foxy *Foxy) GetWebhook(id string) (Webhook, error) {
	url := foxy.baseUrl + "/webhooks/" + id
	fmt.Println("Retrieving webhook from " + url)
	body, e := foxy.getFromFoxy(url)
	var webhook Webhook
	e = json.Unmarshal(body, &webhook)
	webhook.SetIdFromSelfUrl()
	return webhook, e
}

func (foxy *Foxy) AddWebhook(webhook Webhook) (string, error) {
	storeUrl := foxy.ensureStoreUrl()
	updateJson, _ := json.Marshal(webhook)
	result, e := foxy.postToFoxy(storeUrl+"/webhooks", string(updateJson))
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	id := extractId(selfUrl)
	return id, e
}

func (foxy *Foxy) UpdateWebhook(id string, webhook Webhook) (string, error) {
	url := foxy.baseUrl + "/webhooks/" + id
	amendedWebhook := webhook
	amendedWebhook.EventResource = "" // This cannot be updated, it can only be set on creation
	updateJson, _ := json.Marshal(amendedWebhook)
	result, e := foxy.patchFoxy(url, string(updateJson))
	selfUrl := gjson.GetBytes(result, "_links.self.href").String()
	updatedId := extractId(selfUrl)
	return updatedId, e
}

func (foxy *Foxy) DeleteWebhook(id string) error {
	url := foxy.baseUrl + "/webhooks/" + id
	_, e := foxy.deleteFromFoxy(url)
	return e
}

func (foxy *Foxy) ensureStoreUrl() string {
	if foxy.storeUrl == "" {
		rootBody, _ := foxy.getFromFoxy("/")
		// GJson syntax - see https://github.com/tidwall/gjson
		storeUrl := gjson.GetBytes(rootBody, "_links.fx:store.href").String()
		foxy.storeUrl = storeUrl
	}
	return foxy.storeUrl
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

type Webhook struct {
	Id            string `json:"-"`
	Format        string `json:"format,omitempty"`
	Version       int    `json:"version,omitempty"`
	Name          string `json:"name,omitempty"`
	Url           string `json:"url,omitempty"`
	Query         string `json:"query,omitempty"`
	EncryptionKey string `json:"encryption_key,omitempty"`
	//Events        []string `json:"events,omitempty"`
	EventResource string `json:"event_resource,omitempty"`

	// @todo This and the SelfUrl method are a clumsy way of unmarshalling the JSON - should we do this better?
	Links struct {
		Self struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
	} `json:"_links,omitempty"`
}

func (webhook *Webhook) SetIdFromSelfUrl() {
	id := extractId(webhook.SelfUrl())
	webhook.Id = id
}

func (webhook *Webhook) SelfUrl() string {
	return webhook.Links.Self.Href
}

func (webhook *Webhook) GetId() string {
	return extractId(webhook.SelfUrl())
}

func extractId(selfUrl string) string {
	parts := strings.Split(selfUrl, "/")
	return parts[len(parts)-1]
}
