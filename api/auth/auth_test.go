package auth_test

import (
	"github.com/pkg/errors"
	"github.com/steffen25/golang.zone/api/auth"
	"github.com/steffen25/golang.zone/api/auth/jwt"
	"github.com/steffen25/golang.zone/api/auth/models"
	"github.com/steffen25/golang.zone/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuth_Authenticate(t *testing.T) {
	type creds struct {
		email    string
		password string
	}

	cases := []struct {
		name      string
		args      creds
		authMock  *tests.AuthMock
		want      jwt.TokenPair
		mustError bool
	}{
		{
			name: "Fail on finding user",
			args: creds{email: "mail@example.com"},
			authMock: &tests.AuthMock{
				FindByEmailFn: func(email string) (models.User, error) {
					return models.User{}, errors.New("invalid email")
				},
			},
			mustError: true,
		},
		{
			name: "Fail on incorrect password",
			args: creds{email: "mail@example.com", password: "123456"},
			authMock: &tests.AuthMock{
				FindByEmailFn: func(email string) (models.User, error) {
					return models.User{ID: 1, Email: "mail@example.com", Password: "invalid-bcrypt-hash"}, nil
				},
			},
			mustError: true,
		},
		{
			name: "Fail on generate tokens",
			args: creds{email: "mail@example.com", password: "123456"},
			authMock: &tests.AuthMock{
				FindByEmailFn: func(email string) (models.User, error) {
					return models.User{ID: 1, Email: "mail@example.com", Password: "$2a$10$KGrCJ7638juKCdbW4mX1BOuKJpbeHFFYjdH1DQ/SwCNszd/rXIrHW"}, nil
				},
				GenerateTokensFn: func(accessClaims, refreshClaims jwt.APIClaims) (jwt.TokenPair, error) {
					return jwt.TokenPair{}, errors.New("test")
				},
			},
			mustError: true,
		},
		{
			name: "User can login",
			args: creds{email: "mail@example.com", password: "123456"},
			authMock: &tests.AuthMock{
				FindByEmailFn: func(email string) (models.User, error) {
					return models.User{ID: 1, Email: "mail@example.com", Password: "$2a$10$KGrCJ7638juKCdbW4mX1BOuKJpbeHFFYjdH1DQ/SwCNszd/rXIrHW"}, nil
				},
				GenerateTokensFn: func(accessClaims, refreshClaims jwt.APIClaims) (jwt.TokenPair, error) {
					return jwt.TokenPair{}, nil
				},
			},
			mustError: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := auth.New(nil, tt.authMock, tt.authMock, tt.authMock)
			_, err := svc.Authenticate(nil, tt.args.email, tt.args.password)
			assert.Equal(t, tt.mustError, err != nil)
		})
	}

}
