package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/go-playground/validator"
	"github.com/ujjwal405/google_login/database"
	"github.com/ujjwal405/google_login/helper"
	model "github.com/ujjwal405/google_login/models"
	"github.com/ujjwal405/google_login/tokenconfig"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Emailexist          = "email already exists"
	Emailnotexist       = "email doesn't exists"
	authtype            = "bearer"
	errEmailNotVerified = errors.New("email is not verified")
)
var validate = validator.New()

var claim = model.ContextKey("claim")

//type Controllers interface {
//	Home(w http.ResponseWriter, r *http.Request)
//MainProfile(w http.ResponseWriter, r *http.Request)
//EditProfile(w http.ResponseWriter, r *http.Request)
//Login(w http.ResponseWriter, r *http.Request)
//Signup(w http.ResponseWriter, r *http.Request)
//LoginEmail(w http.ResponseWriter, r *http.Request)
//Logout(w http.ResponseWriter, r *http.Request)
//Save(w http.ResponseWriter, r *http.Request)
//SignupGmail(w http.ResponseWriter, r *http.Request)
//SingupCallback(w http.ResponseWriter, r *http.Request)
//}

type Controller struct {
	Database database.AllDatabase
}

func NewController(Database database.AllDatabase) *Controller {
	return &Controller{
		Database: Database,
	}
}
func (ctr *Controller) Home(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("template/signuptemplate.tmpl")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	t.Execute(w, nil)

}

func (ctr *Controller) MainProfile(w http.ResponseWriter, r *http.Request) {
	newclaims := r.Context().Value(claim).(helper.Signedetails)
	uid := newclaims.User_id
	log.Println(uid)
	var user model.UserData
	userdata, err := ctr.Database.DbData(uid)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	user.Email = userdata.Email
	user.Phone = userdata.Phone
	user.Username = userdata.Username
	t, err := template.ParseFiles("template/mainprofile.tmpl")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	w.Header().Set("Cache-Control", "no-cache,  max-age=0")

	t.Execute(w, user)

}
func (ctr *Controller) EditProfile(w http.ResponseWriter, r *http.Request) {

	newclaims := r.Context().Value(claim).(helper.Signedetails)
	uid := newclaims.User_id
	var user model.UserData
	userdata, err := ctr.Database.DbData(uid)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	user.Email = userdata.Email
	user.Phone = userdata.Phone
	user.Username = userdata.Username
	t, err := template.ParseFiles("template/editprofile.tmpl")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	w.Header().Set("Cache-Control", "no-cache,max-age=0")
	t.Execute(w, user)

}

func (ctr *Controller) Login(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/login.tmpl")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	t.Execute(w, nil)

}

func (ctr *Controller) Signup(w http.ResponseWriter, r *http.Request) {
	var user model.UserSignup
	if r.Header.Get("Content-Type") == "application/json" {

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, http.StatusText(500), 500)
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		email := r.Form.Get("email")
		password := r.Form.Get("password")
		user.Email = email
		user.Password = password
	}

	if validateerr := validate.Struct(user); validateerr != nil {
		http.Error(w, http.StatusText(400), 400)
		w.Write([]byte(validateerr.Error()))
		return
	}

	if err := ctr.Database.DBCheckEmail(user.Email); err != nil {
		if err.Error() == Emailexist {
			http.Error(w, http.StatusText(400), 400)
			w.Write([]byte(err.Error()))
			return
		} else if err.Error() != Emailexist && err.Error() != Emailnotexist {
			http.Error(w, http.StatusText(500), 500)
			w.Write([]byte(err.Error()))
			return
		}
	}

	hashedpassword, err := helper.HashPassword(user.Password)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		w.Write([]byte(err.Error()))
		return
	}
	user.Password = hashedpassword
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	user.Isvalid = false
	if err := ctr.Database.DBSignup(user); err != nil {
		http.Error(w, http.StatusText(500), 500)
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
func (ctr *Controller) LoginEmail(w http.ResponseWriter, r *http.Request) {
	var user model.UserSignup
	if r.Header.Get("Content-Type") == "application/json" {

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, http.StatusText(500), 500)
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		email := r.Form.Get("email")
		password := r.Form.Get("password")
		user.Email = email
		user.Password = password
	}

	if validateerr := validate.Struct(user); validateerr != nil {
		http.Error(w, http.StatusText(400), 400)
		w.Write([]byte(validateerr.Error()))
		return
	}
	founduser, err := ctr.Database.DBGetData(user.Email)
	if err != nil {
		if err.Error() == Emailnotexist {
			http.Error(w, http.StatusText(400), 400)
			w.Write([]byte(err.Error()))
			return
		} else if err.Error() != Emailexist && err.Error() != Emailnotexist {
			http.Error(w, http.StatusText(500), 500)
			w.Write([]byte(err.Error()))
			return
		}

	}
	isok, _ := helper.VerifyPassword(user.Password, founduser.Password)
	if !isok {
		http.Error(w, http.StatusText(http.StatusBadRequest), 400)
		w.Write([]byte(" provided incorrect password"))
		return
	}
	token, err := helper.GenerateToken(founduser.User_id, 24*time.Hour)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
		w.Write([]byte(err.Error()))
		return
	}
	jwt_token := fmt.Sprintf("%s %s", authtype, token)
	helper.SetAuthCookie(w, jwt_token)
	http.Redirect(w, r, "/main", http.StatusSeeOther)

}
func (ctr *Controller) Logout(w http.ResponseWriter, r *http.Request) {
	emptytoken := fmt.Sprintf(" %s %s", authtype, "")
	helper.DeleteAuthCookie(w, emptytoken)
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}
func (ctr *Controller) Save(w http.ResponseWriter, r *http.Request) {
	userinfo := r.Context().Value(claim).(helper.Signedetails)
	uid := userinfo.User_id

	var user model.UserData
	if r.Header.Get("Content-Type") == "application/json" {

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, http.StatusText(500), 500)
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		email := r.Form.Get("email")
		username := r.Form.Get("username")
		phone := r.Form.Get("phone")
		user.Username = username
		user.Email = email
		user.Phone = phone
	}
	if validateerr := validate.Struct(user); validateerr != nil {
		http.Error(w, validateerr.Error(), http.StatusBadRequest)
		return
	}
	if err := ctr.Database.DBUpdate(uid, user); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	http.Redirect(w, r, "/main", http.StatusSeeOther)
}

func (ctr *Controller) SignupGmail(w http.ResponseWriter, r *http.Request) {

	uid, err := helper.GenerateRandom()
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	config, err := tokenconfig.LoginConfig()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	url := config.AuthCodeURL(uid)
	helper.SetCookie(w, uid)
	http.Redirect(w, r, url, http.StatusSeeOther)
}
func (ctr *Controller) SingupCallback(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("uid")
	if err != nil {
		http.Error(w, "status doesnt match", http.StatusBadRequest)
		return
	}
	uid := cookie.Value
	state := r.URL.Query().Get("state")
	if state != uid {
		http.Error(w, "status doesnt match", http.StatusInternalServerError)
		return
	}
	code := r.URL.Query().Get("code")
	config, err := tokenconfig.LoginConfig()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	userdata, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var data map[string]interface{}
	json.Unmarshal(userdata, &data)
	email := data["email"].(string)
	name := data["name"].(string)
	isverified := data["verified_email"].(bool)
	if !isverified {
		http.Error(w, errEmailNotVerified.Error(), http.StatusBadRequest)
		return
	}
	var user model.UserSignup
	user.Email = email
	user.Username = name
	user.Isvalid = true
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	if err := ctr.Database.DBSignup(user); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	jwttoken, err := helper.GenerateToken(user.User_id, 24*time.Hour)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
		w.Write([]byte(err.Error()))
		return
	}
	jwt_token := fmt.Sprintf("%s %s", authtype, jwttoken)

	helper.SetAuthCookie(w, jwt_token)
	http.Redirect(w, r, "/main", http.StatusSeeOther)
}
func (ctr *Controller) Cancel(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/main", http.StatusSeeOther)
}
