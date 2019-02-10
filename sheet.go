package ssfet

// Sheet (､´･ω･)▄︻┻┳═一
type Sheet struct {
    ID       string
    Title    string
    Rows     [][]string
    _mapRows []map[string]string
}

// NewSheet (､´･ω･)▄︻┻┳═一
func NewSheet(sheetID, title string, rows [][]string) *Sheet {
    return &Sheet{
        ID:    sheetID,
        Title: title,
        Rows:  rows,
    }
}

func (sheet *Sheet) toMapRows() []map[string]string {
    if len(sheet._mapRows) != 0 {
        return sheet._mapRows
    }

    if len(sheet.Rows) == 0 {
        return make([]map[string]string)
    }

    var (
        columns = sheet.Rows[0]
        mapRows = make([]map[string]string)
    )
    for i, row := range sheet.Rows[1:] {
        rowMap := make(map[string]string)
        for i, column := range columns {
            rowMap[column] = row[i]
        }
        mapRows = append(mapRows, rowMap)
    }
    return mapRows
}
