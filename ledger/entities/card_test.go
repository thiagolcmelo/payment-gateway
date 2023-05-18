package entity_test

import (
	"errors"
	"testing"

	entity "github.com/thiagolcmelo/payment-gateway/ledger/entities"
)

func TestCreditCard_NewCreditCard(t *testing.T) {
	type testCase struct {
		testName    string
		name        string
		number      string
		expireMonth int
		expireYear  int
		cvv         int
		expectedErr error
	}

	testCases := []testCase{
		{
			testName:    "valid_credit_card",
			name:        "name surname",
			number:      "1111-2222-3333-4444",
			expireMonth: 10,
			expireYear:  2099,
			cvv:         123,
			expectedErr: nil,
		},
		{
			testName:    "missing_name",
			name:        "",
			number:      "1111-2222-3333-4444",
			expireMonth: 10,
			expireYear:  2099,
			cvv:         123,
			expectedErr: entity.ErrInvalidName,
		},
		{
			testName:    "missing_number",
			name:        "name surname",
			number:      "",
			expireMonth: 10,
			expireYear:  2099,
			cvv:         123,
			expectedErr: entity.ErrInvalidNumber,
		},
		{
			testName:    "expired",
			name:        "name surname",
			number:      "1111-2222-3333-4444",
			expireMonth: 1,
			expireYear:  2020,
			cvv:         123,
			expectedErr: entity.ErrCardExpired,
		},
		{
			testName:    "invalid_cvv_negative",
			name:        "name surname",
			number:      "1111-2222-3333-4444",
			expireMonth: 10,
			expireYear:  2099,
			cvv:         -1,
			expectedErr: entity.ErrInvalidCVV,
		},
		{
			testName:    "invalid_cvv_too_big",
			name:        "name surname",
			number:      "1111-2222-3333-4444",
			expireMonth: 10,
			expireYear:  2099,
			cvv:         1024,
			expectedErr: entity.ErrInvalidCVV,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			_, err := entity.NewCreditCard(tc.number, tc.name, tc.expireMonth, tc.expireYear, tc.cvv)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
