package request

import (
	"github.com/acya-skulskaya/go-musthave-diploma/internal/validator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserCredits_Valid(t *testing.T) {
	tests := []struct {
		name      string
		userCreds UserCredits
		hasError  bool
		want      validator.Problems
	}{
		{
			name: "valid",
			userCreds: UserCredits{
				Login:    "test",
				Password: "test",
			},
			hasError: false,
			want:     validator.Problems{},
		},
		{
			name: "not valid",
			userCreds: UserCredits{
				Login:    "",
				Password: "",
			},
			hasError: true,
			want: validator.Problems{
				JSONFieldLogin:    ErrorMsgEmptyLogin,
				JSONFieldPassword: ErrorMsgEmptyPassword,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			problems, err := validator.IsValid(tt.userCreds)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, problems)
		})
	}
}
