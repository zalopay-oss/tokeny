package password

import (
	"database/sql"
	"encoding/hex"
	"github.com/golang/mock/gomock"
	"github.com/zalopay-oss/tokeny/pkg/keyvalue"
	mock_keyvalue "github.com/zalopay-oss/tokeny/pkg/keyvalue/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestKeyPasswordValue(t *testing.T) {
	assert.Equal(t, "password", keyPassword)
}

func TestManager_IsRegistered(t *testing.T) {
	tests := []struct {
		name             string
		kvExpectedResult string
		kvExpectedError  error
		expectedResult   bool
		expectedError    error
	}{
		{
			name: "Success",
			kvExpectedResult: "dump result",
			kvExpectedError: nil,
			expectedResult: true,
			expectedError: nil,
		},
		{
			name: "Not found",
			kvExpectedResult: "dump result",
			kvExpectedError: keyvalue.ErrNoRecord,
			expectedResult: false,
			expectedError: nil,
		},
		{
			name: "Unknown error",
			kvExpectedResult: "dump result",
			kvExpectedError: sql.ErrNoRows,
			expectedResult: false,
			expectedError: sql.ErrNoRows,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			kvMock := mock_keyvalue.NewMockStore(ctrl)
			defer ctrl.Finish()

			kvMock.
				EXPECT().
				Get(keyPassword).
				Return(test.kvExpectedResult, test.kvExpectedError).
				Times(1)

			manager := NewManager(kvMock)
			result, err := manager.IsRegistered()
			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test. expectedResult, result)
			}
		})
	}
}

func TestManager_Register(t *testing.T) {
	tests := []struct {
		name string
		password string
		rePassword string
		kvStoreExpectedCall int
		kvStoreExpectedError error
		expectedError error
	}{
		{
			name: "Passwords mismatch",
			password: "123",
			rePassword: "234",
			kvStoreExpectedCall: 0,
			kvStoreExpectedError: nil,
			expectedError: ErrPasswordsMismatch,
		},
		{
			name: "KVStore error",
			password: "123",
			rePassword: "123",
			kvStoreExpectedCall: 1,
			kvStoreExpectedError: sql.ErrNoRows,
			expectedError: sql.ErrNoRows,
		},
		{
			name: "Success",
			password: "123",
			rePassword: "123",
			kvStoreExpectedCall: 1,
			kvStoreExpectedError: nil,
			expectedError: nil,
		},
	}
		for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			kvMock := mock_keyvalue.NewMockStore(ctrl)

			kvMock.EXPECT().
				Set(keyPassword, gomock.Any()).
				Return(test.kvStoreExpectedError).
				Times(test.kvStoreExpectedCall)

			err := NewManager(kvMock).Register(test.password, test.rePassword)
			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestManager_Login(t *testing.T) {
	tests := []struct {
		name string
		inputPassword string
		correctPassword string
		kvError error
		expectedError error
	}{
		{
			name: "Not registered",
			inputPassword: "123",
			correctPassword: "",
			kvError: keyvalue.ErrNoRecord,
			expectedError: ErrNotRegistered,
		},
		{
			name: "KVStore error",
			inputPassword: "123",
			correctPassword: "",
			kvError: sql.ErrConnDone,
			expectedError: sql.ErrConnDone,
		},
		{
			name: "Password mismatch",
			inputPassword: "123",
			correctPassword: "456",
			kvError: nil,
			expectedError: ErrWrongPassword,
		},
		{
			name: "Success",
			inputPassword: "123",
			correctPassword: "123",
			kvError: nil,
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hashedPass, err := bcrypt.GenerateFromPassword([]byte(test.correctPassword), bcrypt.DefaultCost)
			assert.NoError(t, err)
			hashedPassString := hex.EncodeToString(hashedPass)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			kvMock := mock_keyvalue.NewMockStore(ctrl)

			kvMock.EXPECT().
				Get(keyPassword).
				Return(hashedPassString, test.kvError).
				Times(1)

			err = NewManager(kvMock).Login(test.inputPassword)
			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
