package request

import (
	"github.com/acya-skulskaya/go-musthave-diploma/internal/validator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrderID_Valid(t *testing.T) {
	tests := []struct {
		name     string
		o        OrderID
		hasError bool
		want     validator.Problems
	}{
		{
			name:     "valid",
			o:        OrderID("12345678903"),
			hasError: false,
			want:     validator.Problems{},
		},
		{
			name:     "not valid luhn",
			o:        OrderID("12345678901"),
			hasError: true,
			want: validator.Problems{
				JSONFieldBody: ErrorMsgInvalidLuhn,
			},
		},
		{
			name:     "not valid, not int",
			o:        OrderID("aaa"),
			hasError: true,
			want: validator.Problems{
				JSONFieldBody: ErrorMsgCouldNotConvertToInt,
			},
		},
		{
			name:     "not valid empty, not int",
			o:        OrderID(""),
			hasError: true,
			want: validator.Problems{
				JSONFieldBody: ErrorMsgCouldNotConvertToInt,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			problems, err := validator.IsValid(tt.o)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, problems)
		})
	}
}
