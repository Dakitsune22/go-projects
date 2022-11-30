package main

import (
	"log"
	"net/http"

	"github.com/Dakitsune22/go-projects/VinosByDak/users"
)

func main() {
	log.Println("*** Main lanzado.")
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", users.Index)
	http.HandleFunc("/login", users.LogIn)
	http.HandleFunc("/login/process", users.LogInProcess)
	http.HandleFunc("/signup", users.SignUp)
	http.HandleFunc("/signup/process", users.SignUpProcess)
	http.HandleFunc("/logout/process", users.LogOutProcess)
	http.HandleFunc("/userdata/logged", users.ShowLoggedUsers)
	http.HandleFunc("/userdata/registered", users.ShowRegisteredUsers)
	http.HandleFunc("/delete/process", users.DeleteProcess)
	http.HandleFunc("/update", users.Update)
	http.HandleFunc("/update/process", users.UpdateProcess)
	http.HandleFunc("/changepwd", users.ChangePassword)
	http.HandleFunc("/changepwd/process", users.ChangePasswordProcess)

	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("assets"))))

	http.ListenAndServe(":8080", nil)
}
