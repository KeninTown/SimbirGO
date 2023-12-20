package authHandler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"simbirGo/internal/entities"
	mock_authHandler "simbirGo/internal/server/handlers/authHandler/mock"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_SignUp(t *testing.T) {
	type userData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		IsAdmin  bool   `json:"isAdmin"`
	}

	type mockBehavior func(s *mock_authHandler.MockAuthUsecase, user entities.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           entities.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"foo","password":"bar","isAdmin":true}`,
			inputUser: entities.User{
				Id:       0,
				Username: "foo",
				Password: "bar",
				IsAdmin:  true,
				Balance:  0,
			},
			mockBehavior: func(s *mock_authHandler.MockAuthUsecase, user entities.User) {
				s.EXPECT().SignUp(user).Return(entities.User{
					Id:       1,
					Username: "foo",
					Password: "bar",
					IsAdmin:  true,
					Balance:  0,
				}, "token", nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: `{"id":1,"username":"foo","password":"bar","isAdmin":true,"balance":0}`,
		},
		{
			name:                "Empty fields",
			inputBody:           `{"password":"bar","isAdmin":true}`,
			mockBehavior:        func(s *mock_authHandler.MockAuthUsecase, user entities.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"err":"Key: 'userData.Username' Error:Field validation for 'Username' failed on the 'required' tag"}`,
		},
		{
			name:      "Usecase error",
			inputBody: `{"username":"foo","password":"bar","isAdmin":true}`,
			inputUser: entities.User{
				Id:       0,
				Username: "foo",
				Password: "bar",
				IsAdmin:  true,
				Balance:  0,
			},
			mockBehavior: func(s *mock_authHandler.MockAuthUsecase, user entities.User) {
				s.EXPECT().SignUp(user).Return(entities.User{}, "token", errors.New("something went wrong"))
			}, expectedStatusCode: 400,
			expectedRequestBody: `{"err":"something went wrong"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_authHandler.NewMockAuthUsecase(c)

			testCase.mockBehavior(auth, testCase.inputUser)
			handler := New(auth)

			// Test server
			r := gin.New()
			r.POST("/signUp", handler.UserSignUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/signUp",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode, "invalid status code, expected %d have %d", testCase.expectedStatusCode, w.Code)

			assert.Equal(t, testCase.expectedRequestBody, w.Body.String(), "invalid requset body, expected %s have %s", testCase.expectedRequestBody, w.Body.String())
		})

	}
}
