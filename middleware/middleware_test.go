package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ujjwal405/google_login/helper"
)

func TestAuthMiddleware(test *testing.T) {

	mockhandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)

	})
	testcases := []struct {
		name      string
		authsetup func(t *testing.T, r *http.Request)
		response  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			authsetup: func(t *testing.T, r *http.Request) {
				authorization(t, r, time.Minute, "Bearer")
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {

				if recorder.Result().StatusCode != http.StatusOK {
					t.Errorf("expected 200 but got %d", recorder.Result().StatusCode)
				}

			},
		},

		{
			name:      "No header",
			authsetup: func(t *testing.T, r *http.Request) {},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				if recorder.Result().StatusCode != http.StatusUnauthorized {
					t.Errorf("expected 401 but got %d", recorder.Result().StatusCode)
				}
			},
		},

		{
			name: "invalid format",
			authsetup: func(t *testing.T, r *http.Request) {
				authorization(t, r, time.Minute, "")
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				if recorder.Result().StatusCode != http.StatusUnauthorized {
					t.Errorf("expected 401 but got %d", recorder.Result().StatusCode)
				}
			},
		},

		{
			name: "Unsupport",
			authsetup: func(t *testing.T, r *http.Request) {
				authorization(t, r, time.Minute, "unsupport")
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				if recorder.Result().StatusCode != http.StatusUnauthorized {
					t.Errorf("expected 401 but got %d", recorder.Result().StatusCode)
				}
			},
		},

		{
			name: "token expired",
			authsetup: func(t *testing.T, r *http.Request) {
				authorization(t, r, -time.Minute, "Bearer")
			},
			response: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				if recorder.Result().StatusCode != http.StatusUnauthorized {
					t.Errorf("expected 401 but got %d", recorder.Result().StatusCode)
				}
			},
		},
	}
	for i := range testcases {
		tc := testcases[i]
		test.Run(tc.name, func(t *testing.T) {
			path := "/main"
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, path, nil)
			if err != nil {
				t.Error(err)
			}
			tc.authsetup(t, request)
			handler := RecoveryHandler(AuthHandler(mockhandler))
			handler.ServeHTTP(recorder, request)
			tc.response(t, recorder)
		})
	}

}

func authorization(t *testing.T, r *http.Request, duration time.Duration, authtype string) {

	uid := "123"
	token, err := helper.GenerateToken(uid, duration)
	if err != nil {
		t.Error(err)

	}
	authheader := fmt.Sprintf("%s %s", authtype, token)
	r.Header.Set("Authorization", authheader)
}
