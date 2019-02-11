package ssfet

import (
    "fmt"
    "time"

    "github.com/pkg/errors"
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

// SSfet (､´･ω･)▄︻┻┳═一
type SSfet struct {
    client         Client
    exporter       Exporter
    sheets         map[string]*Sheet
    targets        []string
    DataSettings   []*SettingRow
    OptionSettings []*SettingRow
    Error          error
}

// SettingRow (､´･ω･)▄︻┻┳═一
type SettingRow struct {
    Name   string
    Value1 string
    Value2 string
    Value3 string
}

// NewSSfet (､´･ω･)▄︻┻┳═一
func NewSSfet(client Client) *SSfet {
    defaultExporter := NewCSVExporter()
    return &SSfet{client: client, exporter: defaultExporter}
}

// LoadSetting (､´･ω･)▄︻┻┳═一
// 設定シートファイルから設定情報を読み込み、SSfetをセットアップする
func (sf *SSfet) LoadSetting() *SSfet {
    if sf.Error != nil {
        return sf
    }

    settingSheetRange := fmt.Sprintf("%s!A1:ZZ", sf.client.Config().SettingSheetName)
    settingSheet, err := sf.client.Get(sf.client.Config().SettingSheetID, settingSheetRange)
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

// Target  (､´･ω･)▄︻┻┳═一
func (sf *SSfet) Target() *SSfet {
    sf.targets = names
    return sf
}

// Export (､´･ω･)▄︻┻┳═一
func (sf *SSfet) Export() *SSfet {
    if sf.Error != nil {
        return sf
    }

    if sf.exporter == nil {
        err := errors.New("exporter is not specified")
        return sf.knockingErr(err, "export failed")
    }

    eg := errgroup.Group{}
    for _, name := range sf.targets {
        eg.Go(func() error {
            return nil
        })
        time.Sleep(WaitTime)
    }

    if err := sf.exporter.Export(sf.OptionSettings, sf.client.Config().ExportDir); err != nil {
        return sf.knockingErr(err, "export failed")
    }
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
