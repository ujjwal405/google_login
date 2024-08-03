package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mock_database "github.com/ujjwal405/google_login/database/mock"
	"github.com/ujjwal405/google_login/helper"
	model "github.com/ujjwal405/google_login/models"
)

type eqUserMatcher struct {
	arg      model.UserSignup
	password string
}

func (e eqUserMatcher) Matches(x interface{}) bool {
	arg, ok := x.(model.UserSignup)
	if !ok {
		return false
	}
	istrue, _ := helper.VerifyPassword(e.password, arg.Password)
	if !istrue {
		return false
	}
	e.arg.Password = arg.Password
	e.arg.ID = arg.ID
	e.arg.User_id = arg.User_id
	e.arg.Isvalid = arg.Isvalid
	return reflect.DeepEqual(e.arg, arg)
}
func (e eqUserMatcher) String() string {
	return fmt.Sprintf("match arg %v with password %v", e.arg, e.password)
}
func EqMatcher(arg model.UserSignup, password string) gomock.Matcher {
	return eqUserMatcher{arg, password}
}
func TestSignup(t *testing.T) {
	password := "abcdef"
	errEmail := "silwalujjwal03gmail.com"
	noerrEmail := "silwalujjwal03@gmail.com"

	testcases := []struct {
		name          string
		body          model.UserSignup
		buildstubs    func(alldatabase *mock_database.MockAllDatabase)
		checkresponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "validateerr",
			body: model.UserSignup{
				Email:    errEmail,
				Password: password,
			},
			buildstubs: func(alldatabase *mock_database.MockAllDatabase) {

			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "emailexist",
			body: model.UserSignup{
				Email:    noerrEmail,
				Password: password,
			},
			buildstubs: func(alldatabase *mock_database.MockAllDatabase) {
				alldatabase.EXPECT().
					DBCheckEmail(gomock.Eq(noerrEmail)).
					Times(1).
					Return(errors.New("email already exists"))

			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "DbinsertErr",
			body: model.UserSignup{
				Email:    noerrEmail,
				Password: password,
			},
			buildstubs: func(alldatabase *mock_database.MockAllDatabase) {
				alldatabase.EXPECT().
					DBCheckEmail(gomock.Eq(noerrEmail)).
					Times(1).
					Return(errors.New("insertion error"))
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name: "SeconddbError",
			body: model.UserSignup{
				Email:    noerrEmail,
				Password: password,
			},
			buildstubs: func(alldatabase *mock_database.MockAllDatabase) {
				arg := model.UserSignup{
					Email: noerrEmail,
				}
				alldatabase.EXPECT().
					DBCheckEmail(gomock.Eq(noerrEmail)).
					Times(1).
					Return(errors.New("email doesn't exists"))

				alldatabase.EXPECT().
					DBSignup(EqMatcher(arg, password)).
					Times(1).
					Return(errors.New("insertion error"))

			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "All_ok",
			body: model.UserSignup{
				Email:    noerrEmail,
				Password: password,
			},
			buildstubs: func(alldatabase *mock_database.MockAllDatabase) {
				arg := model.UserSignup{
					Email: noerrEmail,
				}
				alldatabase.EXPECT().
					DBCheckEmail(gomock.Eq(noerrEmail)).
					Times(1).
					Return(errors.New("email doesn't exists"))

				alldatabase.EXPECT().
					DBSignup(EqMatcher(arg, password)).
					Times(1).
					Return(nil)

			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusSeeOther, recorder.Code)
			},
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			alldatabase := mock_database.NewMockAllDatabase(ctrl)
			tc.buildstubs(alldatabase)
			recorder := httptest.NewRecorder()
			url := "/signup"
			body, err := json.Marshal(tc.body)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			control := NewController(alldatabase)
			control.Signup(recorder, req)
			tc.checkresponse(t, recorder)

		})
	}
}

func TestLoginEmail(t *testing.T) {
	password := "abcdef"
	incorrectpass := "klmnop"
	errEmail := "silwalujjwal03gmail.com"
	noerrEmail := "silwalujjwal03@gmail.com"
	hashpass, err := NewHash(password)
	require.NoError(t, err)
	testcases := []struct {
		name          string
		body          model.UserSignup
		buildstubs    func(alldatabase *mock_database.MockAllDatabase)
		checkresponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "validateerr",
			body: model.UserSignup{
				Email:    errEmail,
				Password: password,
			},
			buildstubs: func(alldatabase *mock_database.MockAllDatabase) {

			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},

		{
			name: "Emailnotexist",
			body: model.UserSignup{
				Email:    noerrEmail,
				Password: password,
			},
			buildstubs: func(alldatabase *mock_database.MockAllDatabase) {
				alldatabase.EXPECT().
					DBGetData(gomock.Eq(noerrEmail)).
					Times(1).
					Return(model.UserSignup{}, errors.New("email doesn't exists"))
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "mongoerror",
			body: model.UserSignup{
				Email:    noerrEmail,
				Password: password,
			},
			buildstubs: func(alldatabase *mock_database.MockAllDatabase) {
				alldatabase.EXPECT().
					DBGetData(gomock.Eq(noerrEmail)).
					Times(1).
					Return(model.UserSignup{}, errors.New("mongoerror"))
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "incorrectpass",
			body: model.UserSignup{
				Email:    noerrEmail,
				Password: incorrectpass,
			},
			buildstubs: func(alldatabase *mock_database.MockAllDatabase) {
				dbresponse := model.UserSignup{
					Password: hashpass,
				}
				alldatabase.EXPECT().
					DBGetData(gomock.Eq(noerrEmail)).
					Times(1).
					Return(dbresponse, nil)
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			alldatabase := mock_database.NewMockAllDatabase(ctrl)
			tc.buildstubs(alldatabase)
			recorder := httptest.NewRecorder()
			url := "/login"
			body, err := json.Marshal(tc.body)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			control := NewController(alldatabase)
			control.LoginEmail(recorder, req)
			tc.checkresponse(t, recorder)
		})

	}
}
func NewHash(pass string) (string, error) {
	hashpass, err := helper.HashPassword(pass)
	if err != nil {
		return "", err
	}
	return hashpass, nil
}

func TestSave(t *testing.T) {
	username := "ujjwal"
	errEmail := "silwaujjwal03gmail.com"
	noErrEmail := "silwalujjwal03@gmail.com"
	phone := "9841000000"
	uid := "1234"
	testcases := []struct {
		name          string
		body          model.UserData
		buildstubs    func(alldatabase *mock_database.MockAllDatabase)
		addcontext    func(r *http.Request, val helper.Signedetails) *http.Request
		checkresponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ValidateErr",
			body: model.UserData{
				Email:    errEmail,
				Phone:    phone,
				Username: username,
			},
			buildstubs: func(alldatabase *mock_database.MockAllDatabase) {

			},
			addcontext: func(r *http.Request, val helper.Signedetails) *http.Request {
				req := Addcontext(r, val)
				return req
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "dbErr",
			body: model.UserData{
				Email:    noErrEmail,
				Phone:    phone,
				Username: username,
			},
			buildstubs: func(alldatabase *mock_database.MockAllDatabase) {
				arg := model.UserData{
					Email:    noErrEmail,
					Phone:    phone,
					Username: username,
				}
				alldatabase.EXPECT().
					DBUpdate(gomock.Eq(uid), gomock.Eq(arg)).
					Times(1).
					Return(errors.New("db error"))
			},
			addcontext: func(r *http.Request, val helper.Signedetails) *http.Request {
				req := Addcontext(r, val)
				return req
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "OK",
			body: model.UserData{
				Email:    noErrEmail,
				Phone:    phone,
				Username: username,
			},
			buildstubs: func(alldatabase *mock_database.MockAllDatabase) {
				arg := model.UserData{
					Email:    noErrEmail,
					Phone:    phone,
					Username: username,
				}
				alldatabase.EXPECT().
					DBUpdate(gomock.Eq(uid), gomock.Eq(arg)).
					Times(1).
					Return(nil)
			},
			addcontext: func(r *http.Request, val helper.Signedetails) *http.Request {
				req := Addcontext(r, val)
				return req
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusSeeOther, recorder.Code)
			},
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			userdetails := helper.Signedetails{
				User_id: uid,
			}
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			alldatabase := mock_database.NewMockAllDatabase(ctrl)
			tc.buildstubs(alldatabase)
			recorder := httptest.NewRecorder()
			url := "/save"
			body, err := json.Marshal(tc.body)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(body))
			require.NoError(t, err)
			req = tc.addcontext(req, userdetails)

			req.Header.Set("Content-Type", "application/json")
			control := NewController(alldatabase)
			control.Save(recorder, req)
			tc.checkresponse(t, recorder)
		})
	}
}
func Addcontext(r *http.Request, val helper.Signedetails) *http.Request {
	claim := model.ContextKey("claim")
	ctx := context.WithValue(r.Context(), claim, val)
	req := r.WithContext(ctx)
	log.Println("added")
	return req

}
func TestSignCallback(t *testing.T) {
	uuid := "12345"
	testcases := []struct {
		name          string
		buildstubs    func(alldatabase *mock_database.MockAllDatabase)
		addcontext    func(r *http.Request, uid string) *http.Request
		checkresponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Not_ok",
			buildstubs: func(alldatabase *mock_database.MockAllDatabase) {

				alldatabase.EXPECT().
					DBSignup(gomock.Any()).
					Times(0).
					Return(nil)
			},
			addcontext: func(r *http.Request, uid string) *http.Request {
				return adduid(r, uid)
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, 400, recorder.Code)

			},
		},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			alldatabase := mock_database.NewMockAllDatabase(ctrl)
			tc.buildstubs(alldatabase)
			recorder := httptest.NewRecorder()
			url := "/logincallback"
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			req = tc.addcontext(req, uuid)
			r := req.URL.Query()
			r.Add("code", uuid)
			control := NewController(alldatabase)
			control.SingupCallback(recorder, req)
			tc.checkresponse(t, recorder)
		})
	}
}
func adduid(r *http.Request, uid string) *http.Request {
	status := model.ContextKey("status")
	ctx := context.WithValue(r.Context(), status, uid)
	req := r.WithContext(ctx)

	log.Println("added")
	return req
}
