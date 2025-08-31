package common

import (

)

// ClientConfig Configuration used by the client
type Apuesta struct {
	Nombre        string
	Apellido      string
	Documento     string
	Nacimiento    string
	Numero	      string
}

func (a Apuesta) toString() string {
	return a.Nombre + "/" +
		a.Apellido + "/" +
		a.Documento + "/" +
		a.Nacimiento + "/" +
		a.Numero
}
