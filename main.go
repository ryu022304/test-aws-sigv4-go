package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
)

const region = "ap-northeast-1"
const apiEndpoint = "https://example.execute-api.ap-northeast-1.amazonaws.com/hoge/fuga"
const profileName = "default"

func main() {
    ctx := context.Background()

		// AWS Profile設定読み込み
    cfg, err := config.LoadDefaultConfig(
        ctx,
        config.WithSharedConfigProfile(profileName),
    )
    if err != nil {
        panic(err.Error())
    }

		// Profileからcredential情報読み取り
    credentials, err := cfg.Credentials.Retrieve(ctx)
    if err != nil {
        panic(err.Error())
    }

		// リクエスト内容作成
		body := []byte(`{"body": "test"}`)
    buf  := bytes.NewBuffer(body)

    req, err := http.NewRequest(
        http.MethodPost,
        apiEndpoint,
        buf,
    )
    if err != nil {
        panic(err.Error())
    }
    req.Header.Add("Content-Type", "application/json")

		// リクエストbodyのハッシュ値作成
		b := sha256.Sum256(buf.Bytes())
	  payloadHash := hex.EncodeToString(b[:])

		// SigV4対応
    signer := v4.NewSigner()
    err = signer.SignHTTP(ctx, credentials, req, payloadHash, "execute-api", region, time.Now())

    if err != nil {
        panic(err.Error())
    }

    httpClient := new(http.Client)

		// リクエスト実行
    response, err := httpClient.Do(req)
    if err != nil {
        panic(err.Error())
    }

    defer response.Body.Close()

		// レスポンス取得
    responseBody, err := ioutil.ReadAll(response.Body)
    if err != nil {
        panic(err.Error())
    }

    fmt.Print(string(responseBody))
}
