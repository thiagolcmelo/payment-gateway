package entity_test

import (
	"errors"
	"testing"

	entity "github.com/thiagolcmelo/payment-gateway/merchant/entities"
	"golang.org/x/crypto/bcrypt"
)

func TestMechant_NewMerchant(t *testing.T) {
	type testCase struct {
		testName    string
		username    string
		password    string
		active      bool
		maxQPS      int
		expectedErr error
	}

	testCases := []testCase{
		{
			testName:    "valid",
			username:    "username0",
			password:    "password0",
			active:      true,
			maxQPS:      100,
			expectedErr: nil,
		},
		{
			testName:    "empty_username",
			username:    "",
			password:    "password0",
			active:      true,
			maxQPS:      100,
			expectedErr: entity.ErrUsernameOrPasswordEmpty,
		},
		{
			testName:    "empty_password",
			username:    "username0",
			password:    "",
			active:      true,
			maxQPS:      100,
			expectedErr: entity.ErrUsernameOrPasswordEmpty,
		},
		{
			testName:    "negative_qps",
			username:    "username0",
			password:    "password0",
			active:      true,
			maxQPS:      -100,
			expectedErr: entity.ErrMaxQPSCannotBeNegative,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			actualMerchant, err := entity.NewMerchant(tc.username, tc.password, tc.active, tc.maxQPS)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}

			if tc.expectedErr != nil {
				return
			}

			if tc.username != actualMerchant.Username {
				t.Errorf("expected Username=%s, got %s", tc.username, actualMerchant.Username)
			}

			if tc.active != actualMerchant.Active {
				t.Errorf("expected Active=%v, got %v", tc.active, actualMerchant.Active)
			}

			if tc.maxQPS != actualMerchant.MaxQPS {
				t.Errorf("expected MaxQPS=%d, got %d", tc.maxQPS, actualMerchant.MaxQPS)
			}

			err = bcrypt.CompareHashAndPassword([]byte(actualMerchant.Password), []byte(tc.password))
			if err != nil {
				t.Error("invalid password")
			}
		})
	}
}
