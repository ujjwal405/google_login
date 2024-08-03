package main

import (
	"fmt"
	"log"
	"net/http"

	controller "github.com/ujjwal405/google_login/controllers"
	"github.com/ujjwal405/google_login/database"
	"github.com/ujjwal405/google_login/middleware"
)

const CollectionName = "usertest"

func main() {
	client, err := database.DbInstance()
	if err != nil {
		fmt.Println(err.Error())
	}
	collection := database.OpenCollection(client, CollectionName)
	newdatabase := database.NewDatabase(client, collection)
	Controller := controller.NewController(newdatabase)
	http.HandleFunc("/", Controller.Home)
	http.HandleFunc("/signup", Controller.Signup)
	http.HandleFunc("/logingoogle", Controller.SignupGmail)
	http.HandleFunc("/logincallback", Controller.SingupCallback)
	http.HandleFunc("/login", Controller.Login)
	http.HandleFunc("/loginemail", Controller.LoginEmail)
	http.HandleFunc("/main", middleware.RecoveryHandler(middleware.AuthHandler(http.HandlerFunc(Controller.MainProfile))))
	http.HandleFunc("/edit", middleware.RecoveryHandler(middleware.AuthHandler(http.HandlerFunc(Controller.EditProfile))))
	http.HandleFunc("/logout", Controller.Logout)
	http.HandleFunc("/save", middleware.RecoveryHandler(middleware.AuthHandler(http.HandlerFunc(Controller.Save))))
	http.HandleFunc("/cancel", Controller.Cancel)
	fmt.Println("hello")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}
