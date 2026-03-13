package http

import (
	"bytes"
	"eduplay-gateway/internal/http/handlers/user/signUp"
	model "eduplay-gateway/internal/lib/models/user"

	errs "eduplay-gateway/internal/storage"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"eduplay-gateway/tests/http/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTest(t *testing.T) (*mocks.UseCase, http.HandlerFunc) {
	t.Helper()

	mockUc := mocks.NewUseCase(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := signUp.New(logger, mockUc)

	return mockUc, handler
}

func TestSignUpHandler_PasswordMismatch(t *testing.T) {
	mockUC, handler := setupTest(t)

	reqBody := model.SignUpRequest{
		Email:          "test@mail.ru",
		Password:       "secret",
		RepeatPassword: "different",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// expectedBody := `"` + storage.ErrPasswordsNotMatch.Error() + `"` // your error string
	// assert.JSONEq(t, expectedBody, w.Body.String())

	mockUC.AssertNotCalled(t, "SignUp")
}

func TestSignUpHandler_UserAlreadyExists(t *testing.T) {
	mockUC, handler := setupTest(t)

	reqBody := model.SignUpRequest{
		Email:          "test@mail.ru",
		Password:       "secret",
		RepeatPassword: "secret",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUC.On("SignUp", req.Context(), &reqBody).Return(nil, errs.ErrUserAlreadyExists)

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	mockUC.AssertCalled(t, "SignUp", req.Context(), &reqBody)
}

func TestSignUpHandler_Success(t *testing.T) {
	mockUC, handler := setupTest(t)

	expectedCredentials := &model.Credentials{
		AccessToken:  "someAccessToken",
		RefreshToken: "someRefreshToken",
	}
	mockUC.On("SignUp", mock.Anything, mock.AnythingOfType("*userModel.SignUpRequest")).
		Return(expectedCredentials, nil).
		Once()

	reqBody := model.SignUpRequest{
		Email:          "new@example.com",
		Password:       "secret",
		RepeatPassword: "secret",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Check that response JSON matches expectedCredentials
	expectedJSON, _ := json.Marshal(expectedCredentials)
	assert.JSONEq(t, string(expectedJSON), w.Body.String())

	mockUC.AssertExpectations(t)
}
