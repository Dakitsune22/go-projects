package users

import (
	"log"
	"net/http"

	"github.com/Dakitsune22/go-projects/VinosByDak/config"
	"github.com/Dakitsune22/go-projects/VinosByDak/session"
)

func Index(rw http.ResponseWriter, req *http.Request) {
	session.KeepSession(rw, req)
	sInfo := session.GetSessionUser(req)
	config.Tpl.ExecuteTemplate(rw, "index.gohtml", sInfo)
	session.ShowSessions()
}

func LogIn(rw http.ResponseWriter, req *http.Request) {
	if session.UserIsLogged(req) {
		http.Redirect(rw, req, "/", http.StatusSeeOther)
	}
	session.KeepSession(rw, req)
	config.Tpl.ExecuteTemplate(rw, "login.gohtml", nil)
}

func LogInProcess(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	GetUserAndLogIn(rw, req).OnErrorShowHttpError(rw)
	http.Redirect(rw, req, "/", http.StatusSeeOther)
}

func SignUp(rw http.ResponseWriter, req *http.Request) {
	if session.UserIsLogged(req) {
		http.Redirect(rw, req, "/", http.StatusSeeOther)
	}
	session.KeepSession(rw, req)
	config.Tpl.ExecuteTemplate(rw, "signup.gohtml", nil)
}

func SignUpProcess(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	CreateUserAndLogin(rw, req).OnErrorShowHttpError(rw)
	http.Redirect(rw, req, "/", http.StatusSeeOther)
}

func LogOutProcess(rw http.ResponseWriter, req *http.Request) {
	if session.UserIsLogged(req) {
		log.Println("*** ANTES DE BORRAR SESION ***")
		session.ShowSessions()

		session.DeleteSession(rw, req)

		log.Println("*** DESPUES DE BORRAR SESION ***")
		session.ShowSessions()
	}
	http.Redirect(rw, req, "/", http.StatusSeeOther)
}

func ShowLoggedUsers(rw http.ResponseWriter, req *http.Request) {
	session.KeepSession(rw, req)
	session.CleanSessions()
	config.Tpl.ExecuteTemplate(rw, "loggedusers.gohtml", session.GetLoggedUserlist())
}

func ShowRegisteredUsers(rw http.ResponseWriter, req *http.Request) {
	session.KeepSession(rw, req)
	if !session.UserIsAdmin(req) {
		http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	sInfo := session.GetSessionUser(req)
	users, ret := GetAllUsers(sInfo.DBMS)
	ret.OnErrorShowHttpError(rw)
	type UDb struct {
		DBMS  string
		Users []User
	}
	usersDb := UDb{sInfo.DBMS, users}
	log.Println("******************************\n", usersDb)
	config.Tpl.ExecuteTemplate(rw, "registeredusers.gohtml", usersDb)
}

func DeleteProcess(rw http.ResponseWriter, req *http.Request) {
	if !session.UserIsAdmin(req) {
		http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	sInfo := session.GetSessionUser(req)
	userId := req.FormValue("id")
	DeleteUser(sInfo.DBMS, userId).OnErrorShowHttpError(rw)
	http.Redirect(rw, req, "/userdata/registered", http.StatusSeeOther)
}

func Update(rw http.ResponseWriter, req *http.Request) {
	if !session.UserIsLogged(req) {
		http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	session.KeepSession(rw, req)
	sInfo := session.GetSessionUser(req)
	userId := req.FormValue("id")
	log.Println("*** Update *** UserId:", userId)
	user, ret := GetUser(sInfo.DBMS, userId)
	ret.OnErrorShowHttpError(rw)
	config.Tpl.ExecuteTemplate(rw, "update.gohtml", user)
}

func UpdateProcess(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if !session.UserIsLogged(req) {
		http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	sInfo := session.GetSessionUser(req)
	user, ret := GetUser(sInfo.DBMS, req.FormValue("id"))
	ret.OnErrorShowHttpError(rw)
	user.Nickname = req.FormValue("nickname")
	user.Name = req.FormValue("name")
	user.Fname = req.FormValue("fname")
	user.Lname = req.FormValue("lname")
	log.Println("*** UpdateProcess *** Usuario a actualizar:", user.Id)
	UpdateUser(sInfo.DBMS, user).OnErrorShowHttpError(rw)
	if sInfo.UserId == user.Id {
		session.UpdateSessionUser(req, user.Nickname)
		if session.UserIsAdmin(req) {
			http.Redirect(rw, req, "/userdata/registered", http.StatusSeeOther)
		} else {
			http.Redirect(rw, req, "/", http.StatusSeeOther)
		}
	} else {
		http.Redirect(rw, req, "/userdata/registered", http.StatusSeeOther)
	}
}

func ChangePassword(rw http.ResponseWriter, req *http.Request) {
	if !session.UserIsLogged(req) {
		http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	session.KeepSession(rw, req)
	config.Tpl.ExecuteTemplate(rw, "changepass.gohtml", req.FormValue("id"))
}

func ChangePasswordProcess(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if !session.UserIsLogged(req) {
		http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	ChangeUserPassword(rw, req).OnErrorShowHttpError(rw)
	http.Redirect(rw, req, "/", http.StatusSeeOther)
}
