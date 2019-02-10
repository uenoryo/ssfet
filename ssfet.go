package ssfet

import (
    "fmt"

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
    DataSettings   []*SettingRow
    OptionSettings []*SettingRow
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
    return &SSfet{client: client}
}

// LoadSetting (､´･ω･)▄︻┻┳═一
// 設定シートファイルから設定情報を読み込み、SSfetをセットアップする
func (ssfet *SSfet) LoadSetting() error {
    settingSheetRange := fmt.Sprintf("%s!A1:ZZ", ssfet.client.Config().SettingSheetName)
    settingSheet, err := ssfet.client.Get(ssfet.client.Config().SettingSheetID, settingSheetRange)
    if err != nil {
        return errors.Wrap(err, "get setting sheet failed")
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
    ssfet.DataSettings = dataSettings
    ssfet.OptionSettings = optionSettings
    return nil
}
