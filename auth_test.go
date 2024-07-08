package twittergraphql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterInput(t *testing.T) {
	testCases := []struct {
		name  string
		input RegisterInput
		err   error
	}{
		{
			name: "valid",
			input: RegisterInput{
				Username:        "newsun",
				Email:           "newsun@gmail.com",
				Password:        "windows",
				ConfirmPassword: "windows",
			},
			err: nil,
		},
		{
			name: "invalid Email",
			input: RegisterInput{
				Username:        "newsun",
				Email:           "newsun@io.io",
				Password:        "windows",
				ConfirmPassword: "windows",
			},
			err: nil,
		},
		{
			name: "too short username",
			input: RegisterInput{
				Username:        "nddddd",
				Email:           "newsun@gmail.com",
				Password:        "windows",
				ConfirmPassword: "windows",
			},
			err: nil,
		},
		{
			name: "short password ",
			input: RegisterInput{
				Username:        "newsun",
				Email:           "newsun@gmail.com",
				Password:        "widdddd",
				ConfirmPassword: "widdddd",
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRegisterInput_Snitize(t *testing.T) {
	input := RegisterInput{
		Username:        " bob",
		Email:           " Me@mailbee.com",
		Password:        "jarvis11",
		ConfirmPassword: "jarvis11",
	}

	want := RegisterInput{
		Username:        "bob",
		Email:           "me@mailbee.com",
		Password:        "jarvis11",
		ConfirmPassword: "jarvis11",
	}

	input.Sanitize()
	require.Equal(t, want, input)

}

func TestLoginInput_Snitize(t *testing.T) {
	input := LoginInput{
		Email:    " Me@mailbee.com",
		Password: "jarvis11",
	}

	want := LoginInput{
		Email:    "me@mailbee.com",
		Password: "jarvis11",
	}

	input.Sanitize()
	require.Equal(t, want, input)

}

func TestLoginInput(t *testing.T) {
	testCases := []struct {
		name  string
		input LoginInput
		err   error
	}{
		{
			name: "valid",
			input: LoginInput{
				Email:    "newsun@gmail.com",
				Password: "windows",
			},
			err: nil,
		},
		{
			name: "invalid Email",
			input: LoginInput{
				Email:    "newsun@io.io",
				Password: "windows",
			},
			err: nil,
		},
		{
			name: "too short username",
			input: LoginInput{
				Email:    "newsun@gmail.com",
				Password: "windows",
			},
			err: nil,
		},
		{
			name: "short password ",
			input: LoginInput{
				Email:    "newsun@gmail.com",
				Password: "widdddd",
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
