package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

//token oauth token
var token = &oauth2.Token{}

//apis API
var apis = []string{
	"https://graph.microsoft.com/v1.0/me/",
	"https://graph.microsoft.com/v1.0/me/messages",
}

var endpoint = oauth2.Endpoint{
	AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
	TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
}

var msOauthConfig = &oauth2.Config{
	Endpoint: endpoint,
}

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		logger.Info("Failed to log to file, using default stderr")
	}
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "",
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "",
		FieldMap:          nil,
		CallerPrettyfier:  nil,
		PrettyPrint:       false,
	})
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf")
	err = viper.ReadInConfig()
	if err != nil {
		logger.WithField("err", err.Error()).Error("load config failed")
		panic(err)
	}
	apis = viper.GetStringSlice("apis")
	//从环境变量读取 CLIENT_ID 和 CLIENT_SECRET
	clientId, _ := os.LookupEnv("CLIENT_ID")
	if clientId == "" {
		logger.Error("clientId can not be empty")
		panic(err)
	}
	clientSecret, _ := os.LookupEnv("CLIENT_SECRET")
	if clientSecret == "" {
		logger.Error("clientSecret can not be empty")
		panic(err)
	}
	msOauthConfig.ClientID = clientId
	msOauthConfig.ClientSecret = clientSecret

	msOauthConfig.Scopes = viper.GetStringSlice("scope")
	msOauthConfig.RedirectURL = viper.GetString("redirect_uri")
}

func main() {
	rand.Seed(time.Now().Unix())
	//随机取一个url
	url := apis[rand.Intn(len(apis))]
	logWithUrl := logger.WithField("url", url)
	//访问
	resp, err := accessAPI(url)
	if err != nil {
		logWithUrl.Errorln(err)
		return
	}
	logWithUrl.Infoln(resp)
}

// accessAPI 访问API
func accessAPI(url string) (string, error) {
	logWithUrl := logger.WithField("url", url)
	//读取配置文件中的token
	if err := readToken(token); err != nil {
		return "", err
	}
	//校验token
	ctx := context.Background()
	tokenSource := msOauthConfig.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return "", err
	}

	//检查token是否更新
	if newToken.AccessToken != token.AccessToken {
		token = newToken
		//更新token
		err := saveToken(token)
		if err != nil {
			return "", err
		}
		logWithUrl.WithField("Expiry", token.Expiry).Info("saved new token")
	}

	//发起请求
	client := oauth2.NewClient(ctx, tokenSource)
	res, err := client.Get(url)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		logWithUrl.WithField("status_code", res.StatusCode).Error("response status not ok")
		return "", errors.New("response status not ok")
	}
	body := res.Body
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(body)
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// readToken 从配置文件读取token
func readToken(t *oauth2.Token) error {
	t.AccessToken = viper.GetString("token.access_token")
	t.RefreshToken = viper.GetString("token.refresh_token")
	t.TokenType = viper.GetString("token.token_type")
	t.Expiry = viper.GetTime("token.expiry")

	if token.AccessToken == "" {
		return fmt.Errorf("no access_token loaded")
	}
	return nil
}

//saveToken 将新的token写去配置文件
func saveToken(t *oauth2.Token) error {
	viper.Set("token.access_token", t.AccessToken)
	viper.Set("token.refresh_token", t.RefreshToken)
	viper.Set("token.token_type", t.TokenType)
	viper.Set("token.expiry", t.Expiry)
	return viper.WriteConfig()
}
