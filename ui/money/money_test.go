package money

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestMoney(t *testing.T) {
	tests := []struct {
		name      string
		props     Props
		wantValue string
		wantCurr  string
	}{
		{
			name:      "simple default",
			props:     Props{Amount: 1234.56},
			wantValue: "1,234.56",
			wantCurr:  "",
		},
		{
			name:      "with currency",
			props:     Props{Amount: 1234.56, Currency: "$"},
			wantValue: "1,234.56",
			wantCurr:  "$",
		},
		{
			name:      "custom precision",
			props:     Props{Amount: 123.4567, Precision: P(3)},
			wantValue: "123.457", // rounded
			wantCurr:  "",
		},
		{
			name:      "zero precision",
			props:     Props{Amount: 123.4567, Precision: P(0)},
			wantValue: "123", // now strictly 0 decimals
			wantCurr:  "",
		},
		{
			name:      "string input",
			props:     Props{Amount: "123.45"},
			wantValue: "123.45",
			wantCurr:  "",
		},
		{
			name:      "string input precise",
			props:     Props{Amount: "123.456789", Precision: P(4)},
			wantValue: "123.4568",
			wantCurr:  "",
		},
		{
			name:      "empty string",
			props:     Props{Amount: ""},
			wantValue: "0.00",
			wantCurr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w := io.Pipe()
			go func() {
				Money(tt.props).Render(context.Background(), w)
				w.Close()
			}()
			buf := new(strings.Builder)
			io.Copy(buf, r)
			got := buf.String()

			if !strings.Contains(got, tt.wantValue) {
				t.Errorf("Money() got = %v, want value %v", got, tt.wantValue)
			}
			if tt.wantCurr != "" && !strings.Contains(got, tt.wantCurr) {
				t.Errorf("Money() got = %v, want currency %v", got, tt.wantCurr)
			}
		})
	}
}
