package api

import (
	"io"
	"log"
	"net/http"
	"syscall/js"

	"github.com/snburman/game/config"
)

var ApiClient *API

func init() {
	ApiClient = NewAPI()
}

type API struct {
	userID string
}

func NewAPI() *API {
	return &API{}
}

// disable js during debugging and testing
func (a *API) GetUserID() string {
	// Get user ID from global JS
	fun := js.Global().Get("id")
	id := fun.Invoke().String()
	if id == "" {
		panic("User ID not found")
	}
	// id := "6778d9d1a1a3232f20545d84"
	a.userID = id
	return id
}

func (a *API) UserID() string {
	if a.userID == "" {
		a.GetUserID()
	}
	return a.userID
}

func (a *API) Request(method string, url string) Response {
	client := http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println("error during new request")
		panic(err)
	}
	req.Header.Add("CLIENT_ID", config.Env().CLIENT_ID)
	req.Header.Add("CLIENT_SECRET", config.Env().CLIENT_SECRET)
	res, err := client.Do(req)
	if err != nil {
		return Response{
			Success: false,
			Status:  res.StatusCode,
			Headers: res.Header,
			Error:   err,
		}
	}
	defer res.Body.Close()
	bts, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return Response{
		Success: res.StatusCode >= 200 && res.StatusCode < 300,
		Status:  res.StatusCode,
		Headers: res.Header,
		Body:    bts,
	}
}
