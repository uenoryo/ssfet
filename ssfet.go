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
func NewSSfet(client Clinet) *SSfet {
    return &SSfet{client: client}
}

// LoadSetting (､´･ω･)▄︻┻┳═一
// 設定シートファイルから設定情報を読み込み、SSfetをセットアップする
func (ssfet *SSfet) LoadSetting() error {
    settingSheetRange := fmt.Sprintf("%s!A1:ZZ", ssfet.client.SettingSheetName)
    settingSheet, err := ssfet.client.Get(ssfet.client.SettingSheetID, settingSheetRange)
    if err != nil {
        return errors.Wrap(err, "get setting sheet failed")
    }

    ssfet.loadDataSetting(sheet * Sheet)
    ssfet.loadOptionSetting(sheet * Sheet)
}

func (ssfet *SSfet) loadDataSetting(sheet *Sheet) {
    settings := make([]*SettingRow, 0, len(sheet.Rows))
    for _, row := range sheet.toMapRows() {
        if row[columnNameForName] == settingTypeNameData {
            continue
        }
        settings = append(settings, &SettingRow{
            Name:   row[columnNameForName],
            Value1: row[columnNameForValue1],
            Value2: row[columnNameForValue2],
            Value3: row[columnNameForValue3],
        })
    }
    ssfet.DataSettings = settings
}

func (ssfet *SSfet) loadOptionSetting(sheet *Sheet) {
    settings := make([]*SettingRow, 0, len(sheet.Rows))
    for _, row := range sheet.toMapRows() {
        if row[columnNameForName] == settingTypeNameOption {
            continue
        }
        settings = append(settings, &SettingRow{
            Name:   row[columnNameForName],
            Value1: row[columnNameForValue1],
            Value2: row[columnNameForValue2],
            Value3: row[columnNameForValue3],
        })
    }
    ssfet.OptionSettings = settings
}
