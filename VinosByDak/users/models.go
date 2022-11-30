package users

import (
	"log"
	"net/http"
	"time"

	"github.com/Dakitsune22/go-projects/VinosByDak/config"
	"github.com/Dakitsune22/go-projects/VinosByDak/ec"
	"github.com/Dakitsune22/go-projects/VinosByDak/session"
	"golang.org/x/crypto/bcrypt"
)

const (
	DateYYYYMMDD = "2006-01-02"
	DateDDMMYYYY = "02/01/2006"
)

type User struct {
	// ID     bson.ObjectId // `json:"id" bson:"_id"`
	Id       string    // `json:"id" bson:"id"`
	Password string    // `json:"password" bson:"password"`
	Nickname string    // `json:"nickname" bson:"nickname"`
	Name     string    // `json:"name" bson:"name"`
	Fname    string    // `json:"fname" bson:"fname"`
	Lname    string    // `json:"lname" bson:"lname"`
	Bdate    time.Time // `json:"bdate" bson:"bdate"`
}

func CreateUserAndLogin(rw http.ResponseWriter, req *http.Request) ec.HttpError {
	httpErr := ec.HttpError{}
	user := User{}
	user.Id = req.FormValue("id")
	user.Nickname = req.FormValue("nickname")
	user.Name = req.FormValue("name")
	user.Fname = req.FormValue("fname")
	user.Lname = req.FormValue("lname")
	pwd := req.FormValue("password")
	dbms := req.FormValue("dbmsystem")

	if user.Id == "" || user.Nickname == "" || pwd == "" {
		httpErr.Cod = http.StatusBadRequest
		httpErr.Desc = "Error. Datos obligatorios sin informar."
		return httpErr
	}

	pwdEncrypted, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		httpErr.Cod = http.StatusInternalServerError
		httpErr.Desc = "Error interno. No se ha podido dar de alta el usuario " + user.Id + "."
		return httpErr
	}
	user.Password = string(pwdEncrypted)
	log.Println("*** CreateUserAndLogin *** Password encriptado:", user.Password)

	log.Println("*** CreateUserAndLogin *** Valor de 'dbmsystem' recuperado del formulario:", dbms)
	switch dbms {
	case "postgres", "mysql":
		if err := config.OpenSqlDB(dbms); err != nil {
			httpErr.Cod = http.StatusInternalServerError
			httpErr.Desc = "Error interno. No se ha podido acceder a la base de datos."
			return httpErr
		}
		if err = CreateUserSql(dbms, user); err != nil {
			log.Println("*** CreateUserAndLogin ***", err)
			httpErr.Cod = http.StatusInternalServerError
			httpErr.Desc = "Error interno. No se ha podido dar de alta el usuario " + user.Id + "."
			return httpErr
		}
	case "mongodb":
		if err := config.OpenMongoDB(); err != nil {
			httpErr.Cod = http.StatusInternalServerError
			httpErr.Desc = "Error interno. No se ha podido acceder a la base de datos."
			return httpErr
		}
		if err = CreateUserMongoDB(user); err != nil {
			log.Println("*** CreateUserAndLogin ***", err)
			httpErr.Cod = http.StatusInternalServerError
			httpErr.Desc = "Error interno. No se ha podido dar de alta el usuario " + user.Id + "."
			return httpErr
		}
	}
	// Añadimos el usuario (logged) a la sesión, en la lista de sesiones:
	session.SetSessionUser(req, user.Id, user.Nickname, dbms)

	return ec.HttpError{}
}

func GetUserAndLogIn(rw http.ResponseWriter, req *http.Request) ec.HttpError {
	httpErr := ec.HttpError{}
	idUser := req.FormValue("id")
	pwd := req.FormValue("password")
	dbms := req.FormValue("dbmsystem")

	if idUser == "" || pwd == "" {
		httpErr.Cod = http.StatusBadRequest
		httpErr.Desc = "Error. Datos obligatorios sin informar."
		return httpErr
	}
	user := User{}
	var err error
	log.Println("*** GetUserAndLogin *** Valor de 'dbmsystem' recuperado del formulario:", dbms)
	switch dbms {
	case "postgres", "mysql":
		if err := config.OpenSqlDB(dbms); err != nil {
			httpErr.Cod = http.StatusInternalServerError
			httpErr.Desc = "Error interno. No se ha podido acceder a la base de datos."
			return httpErr
		}
		user, err = GetUserSql(dbms, idUser)
		if err != nil {
			log.Println("*** GetUserAndLogin ***", err)
			httpErr.Cod = http.StatusInternalServerError
			httpErr.Desc = "Error interno. No se ha podido recuperar el usuario " + user.Id + "."
			return httpErr
		}
	case "mongodb":
		if err := config.OpenMongoDB(); err != nil {
			httpErr.Cod = http.StatusInternalServerError
			httpErr.Desc = "Error interno. No se ha podido acceder a la base de datos."
			return httpErr
		}
		user, err = GetUserMongoDB(idUser)
		if err != err {
			log.Println("*** GetUserAndLogin ***", err)
			httpErr.Cod = http.StatusInternalServerError
			httpErr.Desc = "Error interno. No se ha podido recuperar el usuario " + user.Id + "."
			return httpErr
		}
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd)); err != nil {
		httpErr.Cod = http.StatusForbidden
		httpErr.Desc = "Error. El password introducido no es correcto."
		log.Println("*** GetUserAndLogin *** httpErr:", httpErr)
		return httpErr
	}
	// Añadimos el usuario (logged) a la sesión, en la lista de sesiones:
	session.SetSessionUser(req, user.Id, user.Nickname, dbms)

	return ec.HttpError{}
}

func GetUser(dbms string, userId string) (User, ec.HttpError) {
	var err error
	user := User{}
	switch dbms {
	case "postgres", "mysql":
		user, err = GetUserSql(dbms, userId)
	case "mongodb":
		user, err = GetUserMongoDB(userId)
	}
	if err != nil {
		log.Println("*** GetUser ***", err)
		httpErr := ec.HttpError{}
		httpErr.Cod = http.StatusInternalServerError
		httpErr.Desc = "Error interno. No se ha podido recuperar el usuario " + user.Id + "."
		return User{}, httpErr
	}
	return user, ec.HttpError{}
}

func GetAllUsers(dbms string) ([]User, ec.HttpError) {
	var err error
	httpErr := ec.HttpError{}
	users := make([]User, 0)
	switch dbms {
	case "postgres", "mysql":
		users, err = GetAllUsersSql(dbms)
	case "mongodb":
		users, err = GetAllUsersMongoDB()
	}
	if err != nil {
		log.Println("*** GetAllUsers ***", err)
		httpErr.Cod = http.StatusInternalServerError
		httpErr.Desc = "Error interno. No se han podido recuperar los usuarios."
		return nil, httpErr
	}
	return users, httpErr
}

func DeleteUser(dbms string, userId string) ec.HttpError {
	var err error
	httpErr := ec.HttpError{}
	log.Println("*** DeleteUser *** Usuario a borrar:", userId)
	switch dbms {
	case "postgres", "mysql":
		err = DeleteUserSql(dbms, userId)
	case "mongodb":
		err = DeleteUserMongoDB(userId)
	}
	if err != nil {
		log.Println("*** DeleteUser ***", err)
		httpErr.Cod = http.StatusInternalServerError
		httpErr.Desc = "Error interno. No se ha podido eliminar el usuario."
	}
	return httpErr
}

func UpdateUser(dbms string, user User) ec.HttpError {
	var err error
	httpErr := ec.HttpError{}
	if user.Nickname == "" {
		httpErr.Cod = http.StatusBadRequest
		httpErr.Desc = "Error. Datos obligatorios sin informar."
		return httpErr
	}
	log.Println("*** UpdateUser *** Usuario a actualizar:", user.Id)
	switch dbms {
	case "postgres", "mysql":
		err = UpdateUserSql(dbms, user)
	case "mongodb":
		err = UpdateUserMongoDB(user)
	}
	if err != nil {
		log.Println("*** UpdateUser ***", err)
		httpErr.Cod = http.StatusInternalServerError
		httpErr.Desc = "Error interno. No se ha podido actualizar el usuario."
	}
	return httpErr
}

func ChangeUserPassword(rw http.ResponseWriter, req *http.Request) ec.HttpError {
	httpErr := ec.HttpError{}
	pass := req.FormValue("pass")
	npass := req.FormValue("npass")
	rpass := req.FormValue("rpass")
	if npass != rpass {
		httpErr.Cod = http.StatusBadRequest
		httpErr.Desc = "Error. El nuevo password y su repetición no coinciden."
		log.Println("*** ChangeUserPassword *** httpErr:", httpErr)
		return httpErr
	}
	sInfo := session.GetSessionUser(req)
	user, ret := GetUser(sInfo.DBMS, req.FormValue("id"))
	if ret.ExistHttpError() {
		return ret
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
		httpErr.Cod = http.StatusForbidden
		httpErr.Desc = "Error. El password introducido no es correcto."
		return httpErr
	}
	npassEnc, err := bcrypt.GenerateFromPassword([]byte(npass), bcrypt.MinCost)
	if err != nil {
		httpErr.Cod = http.StatusInternalServerError
		httpErr.Desc = "Error interno. No se ha cambiar el password."
		return httpErr
	}
	log.Println("*** ChangeUserPassword *** Password (encriptado) previo al cambio:", user.Password)
	user.Password = string(npassEnc)
	log.Println("*** ChangeUserPassword *** Password (encriptado) posterior al cambio:", user.Password)
	ret = UpdateUser(sInfo.DBMS, user)
	if ret.ExistHttpError() {
		return ret
	}
	return httpErr
}
