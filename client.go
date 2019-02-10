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
    scopes   = []string{sheets.SpreadsheetsScope}
    tokenURL = google.JWTTokenURL
)

// Client (､´･ω･)▄︻┻┳═一
type Client interface {
    Get(sheetID, readRange string) (*Sheet, error)
}

// Config (､´･ω･)▄︻┻┳═一
type Config struct {
    Email        string
    PrivateKeyID string
    PrivateKey   []byte
}

type sheetAPIClient struct {
    service *sheets.Service
}

// NewClient (､´･ω･)▄︻┻┳═一
func NewClient(ctx context.Context, cnf *Config) (Client, error) {
    cliCnf := &jwt.Config{
        Email:        cnf.Email,
        PrivateKeyID: cnf.PrivateKeyID,
        PrivateKey:   cnf.PrivateKey,
        Scopes:       scopes,
        TokenURL:     tokenURL,
    }

    service, err := sheets.New(cliCnf.Client(ctx))
    if err != nil {
        return nil, errors.Wrap(err, "prepare spreadsheet service failed")
    }
    return &sheetAPIClient{service: service}, nil
}

// Get (､´･ω･)▄︻┻┳═一
func (client *sheetAPIClient) Get(sheetID, readRange string) (*Sheet, error) {
    response, err := client.service.Spreadsheets.Values.Get(sheetID, readRange).Do()
    if err != nil {
        return nil, errors.Wrapf(err, "get spreadsheet failed (sheet ID: %s, range: %s)", sheetID, readRange)
    }

    rows := make([][]string, len(response.Values))
    for i, rowVal := range response.Values {
        rows[i] = make([]string, len(rowVal))
        for j, cell := range rowVal {
            rows[i][j] = cell.(string)
        }
    }
    return NewSheet(sheetID, client.sheetNameByRange(readRange), rows), nil
}

func (client *sheetAPIClient) sheetNameByRange(readRange string) string {
    return strings.Split(readRange, "!")[0]
}
