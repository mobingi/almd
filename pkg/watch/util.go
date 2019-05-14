package watch

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const retryInterval = 10 * time.Second

func getAccessToken(c *http.Client, backendURL, id, token string) string {
	dataJSON := map[string]string{
		"id":    id,
		"token": token,
	}

	var accessTokenURL string
	if strings.HasSuffix(backendURL, "/") {
		accessTokenURL = backendURL + "accesstoken"
	} else {
		accessTokenURL = backendURL + "/accesstoken"
	}
	postData, _ := json.Marshal(dataJSON)
	post := func() (string, error) {
		resp, err := c.Post(accessTokenURL, "application/json", bytes.NewReader(postData))
		if err != nil {
			return "", errors.New("get access token,post error:" + err.Error())
		}
		defer resp.Body.Close()

		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		respJSON := make(map[string]string)
		if err := json.Unmarshal(respData, &respJSON); err != nil {
			log.Println("json unmarshal error,resp:", string(respData))
			return "", err
		}

		accessToken, exists := respJSON["access_token"]
		if !exists {
			return "", errors.New("resp body access_token not exists")
		}

		return accessToken, nil
	}

	return util(post)
}

// newHTTPClient new a client don't check x509 cert,because traafic through proxy
func newHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func util(fn func() (string, error)) string {
	t := time.NewTicker(retryInterval)
	for {
		select {
		case <-t.C:
			data, err := fn()
			if err == nil {
				return data
			}
			log.Print("fn error:", err)
		}
	}
}
