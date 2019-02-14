package ssfet

import (
    "reflect"
    "testing"
)

func TestSheet_toMapRows(t *testing.T) {
    type Test struct {
        Title  string
        Sheet  Sheet
        Expect []map[string]string
    }

    tests := []Test{
        {
            Title: "normal case",
            Sheet: Sheet{
                Rows: [][]string{
                    {
                        "id",
                        "name",
                        "type",
                    },
                    {
                        "1",
                        "taro",
                        "AAA",
                    },
                    {
                        "2",
                        "jiro",
                        "BBB",
                    },
                    {
                        "3",
                        "saburo",
                        "CCC",
                    },
                },
            },
            Expect: []map[string]string{
                {
                    "id":   "1",
                    "name": "taro",
                    "type": "AAA",
                },
                {
                    "id":   "2",
                    "name": "jiro",
                    "type": "BBB",
                },
                {
                    "id":   "3",
                    "name": "saburo",
                    "type": "CCC",
                },
            },
        },
        {
            Title: "empty case",
            Sheet: Sheet{
                Rows: [][]string{},
            },
            Expect: []map[string]string{},
        },
    }

    for _, test := range tests {
        t.Run(test.Title, func(t *testing.T) {
            if g, w := test.Sheet.toMapRows(), test.Expect; !reflect.DeepEqual(g, w) {
                t.Errorf("unexpected result, got %d, want %d", g, w)
            }
        })
    }
}
