package foxyclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"log"
	"strings"
)

type FoxyClient interface {
	get(path string) ([]byte, error)
	put(path string, body string) ([]byte, error)
	post(path string, body string) ([]byte, error)
	patch(path string, body string) ([]byte, error)
	delete(path string) ([]byte, error)

	retrieveStoreId() (string, error)
}

type FoxyHttpClient struct {
	token   oauth2.Token
	baseUrl string
	storeId string
}

var (
	_ FoxyClient = &FoxyHttpClient{}
)

func newFoxyClient(baseUrl string, clientId string, clientSecret string, refreshToken string) (FoxyHttpClient, error) {
	foxy := FoxyHttpClient{baseUrl: baseUrl}
	err := foxy.setToken(clientId, clientSecret, refreshToken)
	_, err = foxy.retrieveStoreId()
	return foxy, err
}

func (foxy *FoxyHttpClient) retrieveStoreId() (string, error) {
	if foxy.storeId == "" {
		rootBody, err := foxy.get("/")
		if err != nil {
			return "", err
		}
		// GJson syntax - see https://github.com/tidwall/gjson
		storeUrl := gjson.GetBytes(rootBody, "_links.fx:store.href").String()
		foxy.storeId = extractId(storeUrl)
	}
	return foxy.storeId, nil
}

func (foxy *FoxyHttpClient) get(path string) ([]byte, error) {
	url := foxy.toUrl(path)
	result, err := foxy.createClient().Get(url)
	return result.Body(), err
}

func (foxy *FoxyHttpClient) patch(path string, body string) ([]byte, error) {
	url := foxy.toUrl(path)
	result, err := foxy.createClient().SetBody(body).Patch(url)
	return result.Body(), err
}

func (foxy *FoxyHttpClient) post(path string, body string) ([]byte, error) {
	url := foxy.toUrl(path)
	result, err := foxy.createClient().SetBody(body).Post(url)
	return result.Body(), err
}

func (foxy *FoxyHttpClient) put(path string, body string) ([]byte, error) {
	url := foxy.toUrl(path)
	result, err := foxy.createClient().SetBody(body).Put(url)
	return result.Body(), err
}

func (foxy *FoxyHttpClient) delete(path string) ([]byte, error) {
	url := foxy.toUrl(path)
	result, err := foxy.createClient().Delete(url)
	return result.Body(), err
}

func (foxy *FoxyHttpClient) toUrl(path string) string {
	var url = path
	if strings.Index(path, foxy.baseUrl) != 0 {
		url = foxy.baseUrl + path
	}
	return url
}

func (foxy *FoxyHttpClient) createClient() *resty.Request {
	// Resty docs - https://github.com/go-resty/resty
	oauthClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&foxy.token))
	client := resty.NewWithClient(oauthClient)
	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		if resp.IsError() {
			return fmt.Errorf("invalid status code %s with response body: %s", resp.Status(), string(resp.Body()))
		}
		return nil
	})
	//client.SetDebug(true)
	result := client.R().
		SetHeader("FOXY-API-VERSION", "1").
		SetAuthToken(foxy.token.AccessToken)
	return result
}

func (foxy *FoxyHttpClient) setToken(clientId string, clientSecret string, refreshToken string) error {
	token, err := foxy.retrieveToken(clientId, clientSecret, refreshToken)
	if err != nil {
		log.Fatalf("Token cannot be retrieved: %s", err.Error())
		return err
	}
	foxy.token = token
	return nil
}

func (foxy *FoxyHttpClient) retrieveToken(clientId string, clientSecret string, refreshToken string) (oauth2.Token, error) {
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
