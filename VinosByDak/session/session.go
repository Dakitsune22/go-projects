package session

import (
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

type SessionInfo struct {
	LastActivity time.Time
	UserId       string
	Nickname     string
	DBMS         string
}

var Sessions = map[string]SessionInfo{} // [session ID]{last activity, user ID}

const sessionLength = 60
const cookieName = "session"

// NewSession crea nueva cookie de sesión, y añade la sesión a la lista de sesiones con el usuario vacío (not logged).
func newSession(rw http.ResponseWriter) {
	log.Println("*** NEW SESSION ***")
	sID, _ := uuid.NewV4()
	cookie := &http.Cookie{
		Name:   cookieName,
		Value:  sID.String(),
		MaxAge: sessionLength,
		//Path:   "/",
	}
	http.SetCookie(rw, cookie)

	// Se añade sesión a la lista de sesiones:
	Sessions[cookie.Value] = SessionInfo{
		LastActivity: time.Now(),
		UserId:       "",
		Nickname:     "",
		DBMS:         "",
	}
}

// KeepSession comprueba si existe cookie de sesión. Si no existe, crea una nueva.
// Si existe, la refresca y actualiza la hora de última actividad en la lista de sesiones.
func KeepSession(rw http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		newSession(rw)
		return
	}
	// Refresh:
	cookie.MaxAge = sessionLength
	http.SetCookie(rw, cookie)

	// Se actualiza sesión en la lista de sesiones:
	if sInfo, ok := Sessions[cookie.Value]; ok {
		sInfo.LastActivity = time.Now()
		Sessions[cookie.Value] = sInfo
	} else {
		// Si la sesión no existe en la lista (por ejemplo porque se ha reiniciado la app),
		// la añadimos nuevamente:
		Sessions[cookie.Value] = SessionInfo{
			LastActivity: time.Now(),
			UserId:       "",
			Nickname:     "",
			DBMS:         "",
		}
	}
}

// DeleteSession borra la cookie de sesión y elimina la sesión de la lista de sesiones.
func DeleteSession(rw http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		return
	}
	// Eliminamos la sesión de la lista de sesiones:
	delete(Sessions, cookie.Value)

	// Eliminamos la cookie:
	cookie.MaxAge = -1
	http.SetCookie(rw, cookie)
}

// UserIsLogged devuelve si el usuario ha iniciado sesión o no.
func UserIsLogged(req *http.Request) bool {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		return false
	}
	if sInfo, ok := Sessions[cookie.Value]; ok {
		if sInfo.UserId != "" {
			return true
		}
	}
	return false
}

// UserIsAdmin devuelve si el usuario ha iniciado sesión o no.
func UserIsAdmin(req *http.Request) bool {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		return false
	}
	if sInfo, ok := Sessions[cookie.Value]; ok {
		if sInfo.UserId == "admin" {
			return true
		}
	}
	return false
}

// GetSessionUser devuelve el usuario asociado a la sesión (logged user).
// En caso de no haber usuario, se retorna un valor vacío.
func GetSessionUser(req *http.Request) SessionInfo {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		return SessionInfo{}
	}
	if sInfo, ok := Sessions[cookie.Value]; ok {
		return sInfo
	}
	return SessionInfo{}
}

// SetSessionUser asocia un usuario a la sesión.
func SetSessionUser(req *http.Request, userId string, nickname string, dbms string) {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		return
	}
	if sInfo, ok := Sessions[cookie.Value]; ok {
		sInfo.UserId = userId
		sInfo.Nickname = nickname
		sInfo.DBMS = dbms
		Sessions[cookie.Value] = sInfo
	}
}

// UpdateSessionUser actualiza los datos de usuario en la sesión.
// (Se han modificado los datos del usuario en la BBDD, hay que reflejarlo en la sesión)
func UpdateSessionUser(req *http.Request, nickname string) {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		return
	}
	if sInfo, ok := Sessions[cookie.Value]; ok {
		sInfo.Nickname = nickname
		Sessions[cookie.Value] = sInfo
	}
}

// ----- INFO DE SESIONES EN CURSO -----

// ShowSessions muestra en el log todas las sesiones actuales.
func ShowSessions() {
	log.Println()
	log.Printf("*** Sesiones actuales (%d)\n", len(Sessions))
	log.Println("***\t\tClave                                | Id. Usuario")
	log.Println("***\t\t-------------------------------------------------------------------------------")
	i := 0
	for sKey, sInfo := range Sessions {
		i++
		log.Printf("*** %d)\t%s | %s\n", i, sKey, sInfo.UserId)
	}
}

// CleanSessions elimina de la lista de sesiones aquellas cuya última actividad es de hace más de 60 segundos.
func CleanSessions() {
	for sKey, sInfo := range Sessions {
		if time.Since(sInfo.LastActivity) > time.Second*sessionLength {
			log.Println("*** Segundos desde última actividad:", sKey, "-", time.Since(sInfo.LastActivity))
			delete(Sessions, sKey)
		}
	}
}

func GetLoggedUserlist() []SessionInfo {
	var sInfoList []SessionInfo
	for _, sInfo := range Sessions {
		if sInfo.UserId != "" {
			//users = append(users, sInfo.Nickname)
			sInfoList = append(sInfoList, sInfo)
		}
	}
	log.Println("*** Número de usuarios:", len(sInfoList))
	return sInfoList
}
