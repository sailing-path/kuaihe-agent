package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"product_kuaihe/model"
	"strings"
)

// Authorization 获取token
func Authorization() (string, string, error) {
	// 先从redis缓存获取
	accessToken, err := FreeCache.Get([]byte(model.AuthorizationAccessToken))
	if string(accessToken) != "" {
		tokenNew := strings.Split(string(accessToken), ":")
		if len(tokenNew) != 2 {
			return "", "", errors.New("accessToken错误")
		}
		return tokenNew[0], tokenNew[1], nil
	}

	authResponse, err := getAuthorization()
	if err != nil {
		return "", "", err
	}
	if authResponse.AccessToken == "" || authResponse.TokenType == "" || authResponse.ExpiresIn == 0 {
		return "", "", errors.New("accessToken获取失败")
	}
	return authResponse.AccessToken, authResponse.TokenType, nil
}

// CheckAccessToken token校验
func CheckAccessToken(accessToken string) (*model.CheckAccessToken, error) {
	response, err := http.Get(GlobalConfig.AuthAddress + "/oauth/check_token?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bJson, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var accessCheckData model.CheckAccessToken
	if err := json.Unmarshal(bJson, &accessCheckData); err != nil {
		return nil, err
	}
	return &accessCheckData, nil
}

func getAuthorization() (*model.AuthorizationResp, error) {
	postData := url.Values{}
	postData.Add("grant_type", "client_credentials")
	postData.Add("client_id", GlobalConfig.ClientID)
	postData.Add("client_secret", GlobalConfig.ClientSecret)
	response, err := http.Post(GlobalConfig.AuthAddress+"/oauth/token", "application/x-www-form-urlencoded",
		strings.NewReader(postData.Encode()))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bJson, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var authResponse model.AuthorizationResp
	err = json.Unmarshal(bJson, &authResponse)
	if err != nil {
		return nil, err
	}

	// redis缓存 将过期时间扣除秒 提前处理
	if err := FreeCache.Set([]byte(model.AuthorizationAccessToken),
		[]byte(fmt.Sprintf("%v:%v", authResponse.AccessToken, authResponse.TokenType)),
		int(authResponse.ExpiresIn-GlobalConfig.ProcessAuthorizationSeconds)); err != nil {
		return nil, err
	}
	return &authResponse, nil
}
