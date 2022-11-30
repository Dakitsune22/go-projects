package ec

import "net/http"

type HttpError struct {
	Cod  int
	Desc string
}

func (e HttpError) OnErrorShowHttpError(rw http.ResponseWriter) {
	if e.Cod != 0 {
		http.Error(rw, e.Desc, e.Cod)
	}
}

func (e HttpError) ExistHttpError() bool {
	return e.Cod != 0
}

/*func (e HttpError) Redirect(pattern string) {
if e.Cod == 0{

}
*/
