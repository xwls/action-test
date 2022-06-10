package main

import (
	"context"
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
		logrus.WithField("err", err.Error()).Error("load config failed")
		panic(err)
	}
	apis = viper.GetStringSlice("apis")
}

func main() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 3; i++ {
		//随机取一个url
		url := apis[rand.Intn(len(apis))]
		//访问
		accessAPI(url)
		time.Sleep(time.Second)
	}
}

// accessAPI 访问API
func accessAPI(url string) {
	logWithUrl := logger.WithField("url", url)
	//读取配置文件中的token
	if err := readToken(token); err != nil {
		logWithUrl.WithField("err", err.Error()).Error("read token failed")
		panic(err)
	}
	//校验token
	ctx := context.Background()
	tokenSource := msOauthConfig.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		logWithUrl.WithField("err", err.Error()).Error("check token failed")
		panic(err)
	}

	//检查token是否更新
	if newToken.AccessToken != token.AccessToken {
		token = newToken
		//更新token
		err := saveToken(token)
		if err != nil {
			logWithUrl.WithField("err", err.Error()).Error("save token failed")
			panic(err)
		}
		logWithUrl.WithField("Expiry", token.Expiry).Info("saved new token")
	}

	//发起请求
	client := oauth2.NewClient(ctx, tokenSource)
	res, err := client.Get(url)
	if err != nil {
		logWithUrl.WithField("err", err.Error()).Error("access api failed")
		panic(err)
	}
	if res.StatusCode != 200 {
		logWithUrl.WithField("status_code", res.StatusCode).Error("response status not ok")
		return
	}
	body := res.Body
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(body)
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		logWithUrl.WithField("err", err.Error()).Error("read body failed")
		panic(err)
	}
	logWithUrl.WithField("body", string(bytes)[0:100]+"...").Info("access success")
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
