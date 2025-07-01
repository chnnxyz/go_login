package handlers

import "net/http"

type SubordinateData struct {
	RegisteredBy string `json:"registeredBy"`
	Username     string `json:"Username"`
	Password     string `json:"Password"`
	Email        string `json:"Email"`
	Role         string `json:"Role"`
}

func RegisterSubordinate(w http.ResponseWriter, r *http.Request) {

}
