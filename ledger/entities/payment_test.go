package entity_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	entity "github.com/thiagolcmelo/payment-gateway/ledger/entities"
)

func TestPayment_NewPayment(t *testing.T) {
	type testCase struct {
		testName         string
		merchantID       string
		amount           float64
		currency         string
		purchaseTimeUTC  string
		validationMethod string
		metadata         string
		cardName         string
		cardNumber       string
		cardExpireMonth  int
		cardExpireYear   int
		cardCvv          int
		expectedErr      error
	}

	testCases := []testCase{
		{
			testName:         "valid_payment",
			merchantID:       uuid.New().String(),
			amount:           150.00,
			currency:         "USD",
			purchaseTimeUTC:  "2023-05-18T01:00:00.000",
			validationMethod: "push",
			metadata:         "shopper-123",
			cardName:         "name surname",
			cardNumber:       "1111-2222-3333-4444",
			cardExpireMonth:  10,
			cardExpireYear:   2099,
			cardCvv:          123,
			expectedErr:      nil,
		},
		{
			testName:         "negative_amount",
			merchantID:       uuid.New().String(),
			amount:           -150.00,
			currency:         "USD",
			purchaseTimeUTC:  "2023-05-18T01:00:00.000",
			validationMethod: "push",
			metadata:         "shopper-123",
			cardName:         "name surname",
			cardNumber:       "1111-2222-3333-4444",
			cardExpireMonth:  10,
			cardExpireYear:   2099,
			cardCvv:          123,
			expectedErr:      entity.ErrNegativeAmount,
		},
		{
			testName:         "missing_currency",
			merchantID:       uuid.New().String(),
			amount:           150.00,
			currency:         "",
			purchaseTimeUTC:  "2023-05-18T01:00:00.000",
			validationMethod: "push",
			metadata:         "shopper-123",
			cardName:         "name surname",
			cardNumber:       "1111-2222-3333-4444",
			cardExpireMonth:  10,
			cardExpireYear:   2099,
			cardCvv:          123,
			expectedErr:      entity.ErrMissingCurrency,
		},
		{
			testName:         "valid_payment",
			merchantID:       uuid.New().String(),
			amount:           150.00,
			currency:         "USD",
			purchaseTimeUTC:  "2023-05-18T01:00:00.000",
			validationMethod: "",
			metadata:         "shopper-123",
			cardName:         "name surname",
			cardNumber:       "1111-2222-3333-4444",
			cardExpireMonth:  10,
			cardExpireYear:   2099,
			cardCvv:          123,
			expectedErr:      entity.ErrMissingValidationMethod,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			_, err := entity.NewPayment(
				tc.merchantID,
				tc.amount,
				tc.currency,
				tc.purchaseTimeUTC,
				tc.validationMethod,
				// this is created explicitly to permit validation based on invalid card
				entity.CreditCard{
					Number:      tc.cardNumber,
					Name:        tc.cardName,
					ExpireMonth: tc.cardExpireMonth,
					ExpireYear:  tc.cardExpireYear,
					CVV:         tc.cardCvv,
				},
				tc.metadata,
			)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
