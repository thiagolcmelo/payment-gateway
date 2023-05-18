package memory_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	entity "github.com/thiagolcmelo/payment-gateway/ledger/entities"
	"github.com/thiagolcmelo/payment-gateway/ledger/storage/memory"
)

func TestMemoryLedger_Create(t *testing.T) {
	type testCase struct {
		testName         string
		merchantID       string
		amount           float64
		currency         string
		purchaseTimeUTC  string
		validationMethod string
		cardName         string
		cardNumber       string
		cardExpireMonth  int
		cardExpireYear   int
		cardCvv          int
		metadata         string
		ms               *memory.Storage
		expectedErr      error
	}

	testCases := []testCase{
		{
			testName:         "valid_new_payment",
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
			ms:               memory.NewMemoryStorage(),
			expectedErr:      nil,
		},
		{
			testName:         "invalid_card_cvv_too_big",
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
			cardCvv:          1024,
			ms:               memory.NewMemoryStorage(),
			expectedErr:      entity.ErrInvalidCVV,
		},
		{
			testName:         "invalid_card_cvv_negative",
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
			cardCvv:          -100,
			ms:               memory.NewMemoryStorage(),
			expectedErr:      entity.ErrInvalidCVV,
		},
		{
			testName:         "invalid_card_expired",
			merchantID:       uuid.New().String(),
			amount:           150.00,
			currency:         "USD",
			purchaseTimeUTC:  "2023-05-18T01:00:00.000",
			validationMethod: "push",
			metadata:         "shopper-123",
			cardName:         "name surname",
			cardNumber:       "1111-2222-3333-4444",
			cardExpireMonth:  10,
			cardExpireYear:   2000,
			cardCvv:          123,
			ms:               memory.NewMemoryStorage(),
			expectedErr:      entity.ErrCardExpired,
		},
		{
			testName:         "invalid_card_missing_name",
			merchantID:       uuid.New().String(),
			amount:           150.00,
			currency:         "USD",
			purchaseTimeUTC:  "2023-05-18T01:00:00.000",
			validationMethod: "push",
			metadata:         "shopper-123",
			cardName:         "",
			cardNumber:       "1111-2222-3333-4444",
			cardExpireMonth:  10,
			cardExpireYear:   2099,
			cardCvv:          123,
			ms:               memory.NewMemoryStorage(),
			expectedErr:      entity.ErrInvalidName,
		},
		{
			testName:         "invalid_card_missing_number",
			merchantID:       uuid.New().String(),
			amount:           150.00,
			currency:         "USD",
			purchaseTimeUTC:  "2023-05-18T01:00:00.000",
			validationMethod: "push",
			metadata:         "shopper-123",
			cardName:         "name surname",
			cardNumber:       "",
			cardExpireMonth:  10,
			cardExpireYear:   2099,
			cardCvv:          123,
			ms:               memory.NewMemoryStorage(),
			expectedErr:      entity.ErrInvalidNumber,
		},
		{
			testName:         "invalid_missing_validation_method",
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
			ms:               memory.NewMemoryStorage(),
			expectedErr:      entity.ErrMissingValidationMethod,
		},
		{
			testName:         "invalid_missing_currency",
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
			ms:               memory.NewMemoryStorage(),
			expectedErr:      entity.ErrMissingCurrency,
		},
		{
			testName:         "invalid_negative_amoung",
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
			ms:               memory.NewMemoryStorage(),
			expectedErr:      entity.ErrNegativeAmount,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			payment, _ := entity.NewPayment(
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

			id, err := tc.ms.Create(payment)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}

			if err != nil {
				return
			}

			if id == uuid.Nil {
				t.Error("id was not generated properly")
			}
		})
	}
}

func TestMemoryLedger_Read(t *testing.T) {
	ms := memory.NewMemoryStorage()

	payment, err := entity.NewPayment(
		uuid.New().String(),
		150.00,
		"USD",
		"2023-05-18T01:00:00.000",
		"push",
		entity.CreditCard{
			Number:      "name surname",
			Name:        "1111-2222-3333-4444",
			ExpireMonth: 10,
			ExpireYear:  2099,
			CVV:         123,
		},
		"shopper-123",
	)
	if err != nil {
		t.Fatal(err)
	}

	paymentID, err := ms.Create(payment)
	if err != nil {
		t.Fatal(err)
	}
	payment, err = ms.Read(paymentID)
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		testName        string
		id              uuid.UUID
		ms              *memory.Storage
		expectedPayment entity.Payment
		expectedErr     error
	}

	testCases := []testCase{
		{
			testName:        "existing_payment_is_found",
			id:              paymentID,
			ms:              ms,
			expectedPayment: payment,
			expectedErr:     nil,
		},
		{
			testName:        "unexisting_payment_returns_error",
			id:              uuid.New(),
			ms:              ms,
			expectedPayment: entity.Payment{},
			expectedErr:     memory.ErrUnknownPayment,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			p, err := tc.ms.Read(tc.id)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}

			if err != nil {
				return
			}

			if !tc.expectedPayment.Equal(p) {
				t.Errorf("expected %v, got %v", tc.expectedPayment, p)
			}
		})
	}
}

func TestMemoryLedger_Update(t *testing.T) {
	ms := memory.NewMemoryStorage()

	payment, err := entity.NewPayment(
		uuid.New().String(),
		150.00,
		"USD",
		"2023-05-18T01:00:00.000",
		"push",
		entity.CreditCard{
			Number:      "name surname",
			Name:        "1111-2222-3333-4444",
			ExpireMonth: 10,
			ExpireYear:  2099,
			CVV:         123,
		},
		"shopper-123",
	)
	if err != nil {
		t.Fatal(err)
	}

	unknonwPayment, err := entity.NewPayment(
		uuid.New().String(),
		150.00,
		"USD",
		"2023-05-18T01:00:00.000",
		"push",
		entity.CreditCard{
			Number:      "name surname",
			Name:        "1111-2222-3333-4444",
			ExpireMonth: 10,
			ExpireYear:  2099,
			CVV:         123,
		},
		"shopper-123",
	)
	if err != nil {
		t.Fatal(err)
	}
	unknonwPayment.ID = uuid.New()

	paymentID, err := ms.Create(payment)
	if err != nil {
		t.Fatal(err)
	}
	payment, err = ms.Read(paymentID)
	if err != nil {
		t.Fatal(err)
	}

	validPayment := payment // by value
	validPayment.Status = entity.Pending

	invalidPayment := payment // by value
	invalidPayment.Amount = -100

	type testCase struct {
		testName        string
		ms              *memory.Storage
		payment         entity.Payment
		expectedPayment entity.Payment
		expectedErr     error
	}

	testCases := []testCase{
		{
			testName:        "update_existing_with_valid_values",
			ms:              ms,
			payment:         validPayment,
			expectedPayment: validPayment,
			expectedErr:     nil,
		},
		{
			testName:        "update_existing_with_invalid_value_negative_amount",
			ms:              ms,
			payment:         invalidPayment,
			expectedPayment: payment,
			expectedErr:     entity.ErrNegativeAmount,
		},
		{
			testName:        "update_unexisting",
			ms:              ms,
			payment:         unknonwPayment,
			expectedPayment: entity.Payment{},
			expectedErr:     memory.ErrUnknownPayment,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			err := tc.ms.Update(tc.payment)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}

			if err != nil {
				return
			}

			p, err := tc.ms.Read(tc.payment.ID)
			if err != nil {
				t.Errorf("update caused error: %v", err)
			}

			if !tc.expectedPayment.Equal(p) {
				t.Errorf("payment was corrupted: expected %v, got %v", tc.expectedPayment, p)
			}
		})
	}
}
