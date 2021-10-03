package handlers

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/logging/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/iamkey"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net/http"
	"os"
	"time"
)

var client Client

type Client struct {
	sdk        *ycsdk.SDK
	LogGroupId string
}
type Config struct {
	ServiceAccountId  string
	ServiceAccountKey ServiceAccountKey
	LogGroupId        string
}
type ServiceAccountKey struct {
	Id         string
	PublicKey  string
	PrivateKey string
}

func init() {
	cfg, err := newConfig()
	if err != nil {
		log.Fatal(err)
	}
	client, err = newClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func newConfig() (cfg Config, err error) {
	cfg.ServiceAccountId = os.Getenv("SERVICE_ACCOUNT_ID")
	cfg.LogGroupId = os.Getenv("LOG_GROUP_ID")

	publicKey, err := base64Decode(os.Getenv("PUBLIC_KEY"))
	if err != nil {
		return
	}
	privateKey, err := base64Decode(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return
	}
	cfg.ServiceAccountKey = ServiceAccountKey{
		Id:         os.Getenv("KEY_ID"),
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	return
}

func base64Decode(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func newClient(cfg Config) (client Client, err error) {
	var credentials ycsdk.Credentials
	credentials, err = ycsdk.ServiceAccountKey(&iamkey.Key{
		Id:         cfg.ServiceAccountKey.Id,
		Subject:    &iamkey.Key_ServiceAccountId{ServiceAccountId: cfg.ServiceAccountId},
		PublicKey:  cfg.ServiceAccountKey.PublicKey,
		PrivateKey: cfg.ServiceAccountKey.PrivateKey,
	})

	client.LogGroupId = cfg.LogGroupId

	if err != nil {
		return
	}

	client.sdk, err = ycsdk.Build(context.Background(), ycsdk.Config{
		Credentials: credentials,
	})
	return
}

func (h *Handlers) GetLogs(c echo.Context) error {
	type criteriaReqBody struct {
		Levels []logging.LogLevel_Level `json:"levels"`
		Filter string                   `json:"filter"`
		Since  int64                    `json:"since"`
		Until  int64                    `json:"until"`
		Limit  int64                    `json:"limit"`
	}
	criteriaReq := criteriaReqBody{}

	if err := c.Bind(&criteriaReq); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if criteriaReq.Until == 0 {
		criteriaReq.Until = time.Now().Unix()
	}
	if criteriaReq.Since == 0 {
		criteriaReq.Since = time.Unix(criteriaReq.Until, 0).Add(time.Duration(-1) * time.Hour).Unix()
	}
	if criteriaReq.Limit == 0 || criteriaReq.Limit > 100 {
		criteriaReq.Limit = 100
	}

	req := &logging.ReadRequest{
		Selector: &logging.ReadRequest_Criteria{
			Criteria: &logging.Criteria{
				LogGroupId:    client.LogGroupId,
				ResourceTypes: nil,
				ResourceIds:   nil,
				Since:         timestamppb.New(time.Unix(criteriaReq.Since, 0)),
				Until:         timestamppb.New(time.Unix(criteriaReq.Until, 0)),
				Levels:        criteriaReq.Levels,
				Filter:        criteriaReq.Filter,
				PageSize:      criteriaReq.Limit,
			},
		},
	}
	readResponse, err := client.sdk.LogReading().LogReading().Read(context.Background(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, readResponse)
}

func (h *Handlers) GetPageLogs(c echo.Context) error {
	type pageReqBody struct {
		NexPage      string `json:"next_page_token"`
		PreviousPage string `json:"previous_page_token"`
	}
	pageReq := pageReqBody{}

	if err := c.Bind(&pageReq); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	var pageToken string
	if pageReq.NexPage != "" {
		pageToken = pageReq.NexPage
	} else if pageReq.PreviousPage != "" {
		pageToken = pageReq.PreviousPage
	} else {
		return c.JSON(http.StatusBadRequest, errors.New("no page token provided"))
	}
	req := &logging.ReadRequest{
		Selector: &logging.ReadRequest_PageToken{
			PageToken: pageToken,
		},
	}
	readResponse, err := client.sdk.LogReading().LogReading().Read(context.Background(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, readResponse)
}
