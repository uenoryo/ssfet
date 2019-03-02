package ssfet

import (
    "context"
    "strings"

    "github.com/pkg/errors"
    "golang.org/x/oauth2/google"
    "golang.org/x/oauth2/jwt"
    sheets "google.golang.org/api/sheets/v4"
)

var (
    scopes           = []string{sheets.SpreadsheetsScope}
    tokenURL         = google.JWTTokenURL
    defaultExportDir = "./"
)

// Client (､´･ω･)▄︻┻┳═一
type Client interface {
    Get(sheetID, readRange string) (*Sheet, error)
}

// ClientConfig (､´･ω･)▄︻┻┳═一
type ClientConfig struct {
    Email        string
    PrivateKeyID string
    PrivateKey   []byte
}

type sheetAPIClient struct {
    service *sheets.Service
    config  *ClientConfig
}

// NewClient (､´･ω･)▄︻┻┳═一
func NewClient(ctx context.Context, cnf *ClientConfig) (Client, error) {
    clientCnf := &jwt.Config{
        Email:        cnf.Email,
        PrivateKeyID: cnf.PrivateKeyID,
        PrivateKey:   cnf.PrivateKey,
        Scopes:       scopes,
        TokenURL:     tokenURL,
    }

    service, err := sheets.New(clientCnf.Client(ctx))
    if err != nil {
        return nil, errors.Wrap(err, "prepare spreadsheet service failed")
    }
    return &sheetAPIClient{service: service, config: cnf}, nil
}

// Get (､´･ω･)▄︻┻┳═一
func (client *sheetAPIClient) Get(sheetID, readRange string) (*Sheet, error) {
    response, err := client.service.Spreadsheets.Values.Get(sheetID, readRange).Do()
    if err != nil {
        return nil, errors.Wrapf(err, "get spreadsheet failed (sheet ID: %s, range: %s)", sheetID, readRange)
    }

    columns := []interface{}{}
    if len(response.Values) > 0 {
        columns = response.Values[0]
    }

    rows := make([][]string, len(response.Values))
    for i, rowVal := range response.Values {
        rows[i] = make([]string, len(rowVal))
        for j, cell := range rowVal {
            rows[i][j] = cell.(string)
        }

        // カラム数に値が達していなかった場合空白文字で埋める
        for j := len(rows[i]); j < len(columns); j++ {
            rows[i] = append(rows[i], "")
        }
    }
    return NewSheet(sheetID, client.sheetNameByRange(readRange), rows), nil
}

func (client *sheetAPIClient) sheetNameByRange(readRange string) string {
    return strings.Split(readRange, "!")[0]
}
