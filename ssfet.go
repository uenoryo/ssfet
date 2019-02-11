package ssfet

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/pkg/errors"
    "golang.org/x/sync/errgroup"
)

const (
    columnNameForType   = "type"
    columnNameForName   = "name"
    columnNameForValue1 = "value1"
    columnNameForValue2 = "value2"
    columnNameForValue3 = "value3"

    settingTypeNameData   = "data"
    settingTypeNameOption = "option"
)

var (
    ErrExporterNotSpecified = errors.New("exporter is not specified")
)

// SSfet (､´･ω･)▄︻┻┳═一
type SSfet struct {
    client         Client
    exporter       Exporter
    config         *Config
    sheets         sync.Map
    targets        []string
    DataSettings   []*SettingRow
    OptionSettings []*SettingRow
    Error          error
}

// Config (､´･ω･)▄︻┻┳═一
type Config struct {
    Email            string
    PrivateKeyID     string
    PrivateKey       []byte
    SettingSheetID   string
    SettingSheetName string
    ExportDir        string
}

// SettingRow (､´･ω･)▄︻┻┳═一
type SettingRow struct {
    Name   string
    Value1 string
    Value2 string
    Value3 string
}

// NewSSfet (､´･ω･)▄︻┻┳═一
func NewSSfet(ctx context.Context, cnf *Config) (*SSfet, error) {
    clientCnf := &ClientConfig{
        Email:        cnf.Email,
        PrivateKeyID: cnf.PrivateKeyID,
        PrivateKey:   cnf.PrivateKey,
    }

    client, err := NewClient(ctx, clientCnf)
    if err != nil {
        return nil, errors.Wrap(err, "initialize Google API Client failed")
    }

    defaultExporter := NewCSVExporter()
    return &SSfet{
        config:   cnf,
        client:   client,
        exporter: defaultExporter,
        sheets:   sync.Map{},
    }, nil
}

// LoadSetting (､´･ω･)▄︻┻┳═一
// 設定シートファイルから設定情報を読み込み、SSfetをセットアップする
func (sf *SSfet) LoadSetting() *SSfet {
    if sf.Error != nil {
        return sf
    }

    settingSheetRange := fmt.Sprintf("%s!A1:ZZ", sf.config.SettingSheetName)
    settingSheet, err := sf.client.Get(sf.config.SettingSheetID, settingSheetRange)
    if err != nil {
        return sf.knockingErr(err, "get setting sheet failed")
    }

    var (
        dataSettings   = make([]*SettingRow, 0, len(settingSheet.Rows))
        optionSettings = make([]*SettingRow, 0, len(settingSheet.Rows))
    )
    for _, row := range settingSheet.toMapRows() {
        setting := &SettingRow{
            Name:   row[columnNameForName],
            Value1: row[columnNameForValue1],
            Value2: row[columnNameForValue2],
            Value3: row[columnNameForValue3],
        }
        switch row[columnNameForType] {
        case settingTypeNameData:
            dataSettings = append(dataSettings, setting)
        case settingTypeNameOption:
            optionSettings = append(optionSettings, setting)
        }
    }
    sf.DataSettings = dataSettings
    sf.OptionSettings = optionSettings
    return sf
}

// Fetch (､´･ω･)▄︻┻┳═一
func (sf *SSfet) Fetch() *SSfet {
    if sf.Error != nil {
        return sf
    }

    if sf.exporter == nil {
        return sf.knockingErr(ErrExporterNotSpecified, "export failed")
    }

    eg := errgroup.Group{}
    for _, row := range sf.DataSettings {
        eg.Go(func() error {
            sheet, err := sf.client.Get(row.Name, row.Value1)
            if err != nil {
                return err
            }
            sf.sheets.Store(row.Name, sheet)
            return nil
        })
        time.Sleep(1000)
    }
    if err := eg.Wait(); err != nil {
        return sf.knockingErr(err, "get spreadsheet failed")
    }
    return sf
}

func (sf *SSfet) Sheets() []*Sheet {
    sheets := []*Sheet{}
    sf.sheets.Range(func(key, value interface{}) bool {
        if sheet, ok := value.(*Sheet); ok {
            sheets = append(sheets, sheet)
        }
        return true
    })
    return sheets
}

// Export (､´･ω･)▄︻┻┳═一
func (sf *SSfet) Export() *SSfet {
    return sf
}

func (sf *SSfet) OutputCSV() *SSfet {
    sf.exporter = NewCSVExporter()
    return sf
}

func (sf *SSfet) knockingErr(err error, msg string) *SSfet {
    sf.Error = errors.Wrap(sf.Error, errors.Wrap(err, msg).Error())
    return sf
}
