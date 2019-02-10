package ssfet

// Sheet (､´･ω･)▄︻┻┳═一
type Sheet struct {
    ID    string
    Title string
    Rows  [][]string
}

// NewSheet (､´･ω･)▄︻┻┳═一
func NewSheet(sheetID, title string, rows [][]string) *Sheet {
    return &Sheet{
        ID:    sheetID,
        Title: title,
        Rows:  rows,
    }
}
