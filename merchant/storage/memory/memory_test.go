package memory_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	entity "github.com/thiagolcmelo/payment-gateway/merchant/entities"
	"github.com/thiagolcmelo/payment-gateway/merchant/storage/memory"
	"golang.org/x/crypto/bcrypt"
)

func TestMemoryStorage_ReadMerchant(t *testing.T) {
	ms := memory.NewMemoryStorage()

	merchant, err := entity.NewMerchant("username0", "password0", "user 0", true, 100)
	if err != nil {
		t.Fatal(err)
	}

	validId, err := ms.CreateMerchant(merchant)
	if err != nil {
		t.Fatal(err)
	}

	invalidId, err := uuid.Parse("11226ef6-f4f6-11ed-a05b-0242ac120003")
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		testName    string
		id          uuid.UUID
		username    string
		password    string
		name        string
		active      bool
		maxQPS      int
		ms          *memory.Storage
		expectedErr error
	}

	testCases := []testCase{
		{
			testName:    "read_known_merchant",
			id:          validId,
			username:    "username0",
			password:    "password0",
			name:        "user 0",
			active:      true,
			maxQPS:      100,
			ms:          ms,
			expectedErr: nil,
		},
		{
			testName:    "read_unknown_merchant",
			id:          invalidId,
			ms:          ms,
			expectedErr: memory.ErrUnknownMerchantID,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			actualMerchant, err := tc.ms.ReadMerchant(tc.id)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}

			if tc.expectedErr != nil {
				return
			}

			if tc.username != actualMerchant.Username {
				t.Errorf("expected Username=%s, got %s", tc.username, actualMerchant.Username)
			}

			if tc.name != actualMerchant.Name {
				t.Errorf("expected Name=%s, got %s", tc.name, actualMerchant.Name)
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

func TestMemoryStorage_CreateMerchant(t *testing.T) {
	ms := memory.NewMemoryStorage()

	merchant, err := entity.NewMerchant("username0", "password0", "user 0", true, 100)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ms.CreateMerchant(merchant)
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		testName    string
		id          uuid.UUID
		username    string
		password    string
		name        string
		active      bool
		maxQPS      int
		ms          *memory.Storage
		expectedErr error
	}

	testCases := []testCase{
		{
			testName:    "create_merchant_with_valid_values_0",
			id:          uuid.Nil,
			username:    "username",
			password:    "password1",
			name:        "user",
			active:      true,
			maxQPS:      100,
			ms:          memory.NewMemoryStorage(),
			expectedErr: nil,
		},
		{
			testName:    "create_merchant_with_valid_values_1",
			id:          uuid.Nil,
			username:    "username",
			password:    "password1",
			name:        "user",
			active:      false,
			maxQPS:      100,
			ms:          memory.NewMemoryStorage(),
			expectedErr: nil,
		},
		{
			testName:    "create_merchant_with_valid_values_2",
			id:          uuid.Nil,
			username:    "username",
			password:    "password1",
			name:        "user",
			active:      false,
			maxQPS:      0,
			ms:          memory.NewMemoryStorage(),
			expectedErr: nil,
		},
		{
			testName:    "create_merchant_with_invalid_username",
			id:          uuid.Nil,
			username:    "",
			password:    "password1",
			name:        "user",
			active:      true,
			maxQPS:      100,
			ms:          memory.NewMemoryStorage(),
			expectedErr: entity.ErrUsernameOrPasswordEmpty,
		},
		{
			testName:    "create_merchant_with_invalid_password",
			id:          uuid.Nil,
			username:    "username1",
			password:    "",
			name:        "user",
			active:      true,
			maxQPS:      100,
			ms:          memory.NewMemoryStorage(),
			expectedErr: entity.ErrUsernameOrPasswordEmpty,
		},
		{
			testName:    "create_merchant_with_invalid_name",
			id:          uuid.Nil,
			username:    "username1",
			password:    "password1",
			name:        "",
			active:      true,
			maxQPS:      100,
			ms:          memory.NewMemoryStorage(),
			expectedErr: entity.ErrNameEmpty,
		},
		{
			testName:    "create_merchant_with_negative_qps",
			id:          uuid.Nil,
			username:    "username1",
			password:    "password1",
			name:        "user",
			active:      true,
			maxQPS:      -100,
			ms:          memory.NewMemoryStorage(),
			expectedErr: entity.ErrMaxQPSCannotBeNegative,
		},
		{
			testName:    "create_merchant_with_existing_username",
			id:          uuid.Nil,
			username:    "username0",
			password:    "password",
			name:        "user",
			active:      true,
			maxQPS:      100,
			ms:          ms,
			expectedErr: memory.ErrUsernameAlreadyExists,
		},
		{
			testName:    "create_merchant_with_id_already",
			id:          uuid.New(),
			username:    "username",
			password:    "password",
			name:        "user",
			active:      true,
			maxQPS:      100,
			ms:          memory.NewMemoryStorage(),
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			actualId, err := tc.ms.CreateMerchant(entity.Merchant{
				ID:       tc.id,
				Username: tc.username,
				Password: tc.password,
				Name:     tc.name,
				Active:   tc.active,
				MaxQPS:   tc.maxQPS,
			})

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}

			if err != nil {
				return
			}

			if tc.id != uuid.Nil && tc.id != actualId {
				t.Error("storage is changing existing id")
			}
		})
	}
}

func TestMemoryStorage_FindMerchantID(t *testing.T) {
	ms := memory.NewMemoryStorage()

	merchant, err := entity.NewMerchant("username0", "password0", "user", true, 100)
	if err != nil {
		t.Fatal(err)
	}

	validId, err := ms.CreateMerchant(merchant)
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		testName    string
		username    string
		password    string
		ms          *memory.Storage
		expectedId  uuid.UUID
		expectedErr error
	}

	testCases := []testCase{
		{
			testName:    "find_valid_username_and_password",
			username:    "username0",
			password:    "password0",
			ms:          ms,
			expectedId:  validId,
			expectedErr: nil,
		},
		{
			testName:    "find_valid_username_and_wrong_password",
			username:    "username0",
			password:    "password1",
			ms:          ms,
			expectedId:  validId,
			expectedErr: memory.ErrInvalidPassword,
		},
		{
			testName:    "find_invalid_username",
			username:    "username1",
			password:    "password1",
			ms:          ms,
			expectedId:  uuid.Nil,
			expectedErr: memory.ErrUnknownMerchantUsername,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			id, err := tc.ms.FindMerchantID(tc.username, tc.password)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}

			if tc.expectedErr != nil {
				return
			}

			if tc.expectedId != id {
				t.Errorf("expected Id=%s, got %s", tc.expectedId.String(), id.String())
			}
		})
	}
}
