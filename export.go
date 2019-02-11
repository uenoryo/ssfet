package ssfet

type Exporter interface {
    Export(data, option []*SettingRow, dir string) error
}

type csvExporter struct{}

func NewCSVExporter() Exporter {
    return &csvExporter{}
}

func (exporter *csvExporter) Export(data, option []*SettingRow, dir string) error {
    return nil
}
