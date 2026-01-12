package request

import (
	"github.com/acya-skulskaya/go-musthave-diploma/internal/validator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithdrawal_Valid(t *testing.T) {
	tests := []struct {
		name       string
		withdrawal Withdrawal
		hasError   bool
		want       validator.Problems
	}{
		{
			name: "valid",
			withdrawal: Withdrawal{
				Order: "12345678903",
				Sum:   123.3,
			},
			hasError: false,
			want:     validator.Problems{},
		},
		{
			name: "not valid luhn, should be more than zero",
			withdrawal: Withdrawal{
				Order: "12345678901",
				Sum:   -123.3,
			},
			hasError: true,
			want: validator.Problems{
				JSONFieldOrder: ErrorMsgInvalidLuhn,
				JSONFieldSum:   ErrorMsgShouldBeMoreThanZero,
			},
		},
		{
			name: "not valid could not convert to int",
			withdrawal: Withdrawal{
				Order: "1234567890ddd",
				Sum:   123.3,
			},
			hasError: true,
			want: validator.Problems{
				JSONFieldOrder: ErrorMsgCouldNotConvertToInt,
			},
		},
		{
			name:       "not valid empty",
			withdrawal: Withdrawal{},
			hasError:   true,
			want: validator.Problems{
				JSONFieldOrder: ErrorMsgCouldNotConvertToInt,
				JSONFieldSum:   ErrorMsgShouldBeMoreThanZero,
			},
		},
		{
			name: "not valid empty 2",
			withdrawal: Withdrawal{
				Order: "",
				Sum:   0,
			},
			hasError: true,
			want: validator.Problems{
				JSONFieldOrder: ErrorMsgCouldNotConvertToInt,
				JSONFieldSum:   ErrorMsgShouldBeMoreThanZero,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			problems, err := validator.IsValid(tt.withdrawal)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, problems)
		})
	}
}
