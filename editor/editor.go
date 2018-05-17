package editor

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/patrick-jessen/goplay/engine/window"
)

var EditorChannel = make(chan func(), 10)

func windowSetVsync(w http.ResponseWriter, r *http.Request) {
	tmp := struct {
		Vsync bool `json:"vsync"`
	}{}
	json.NewDecoder(r.Body).Decode(&tmp)

	EditorChannel <- func() {
		window.Settings.SetVSync(tmp.Vsync)
	}
	w.WriteHeader(http.StatusOK)
}
func windowGetVsync(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(struct {
		Vsync bool `json:"vsync"`
	}{
		Vsync: window.Settings.VSync(),
	})
}

func windowGetFullScreen(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(struct {
		FullScreen bool `json:"fullScreen"`
	}{
		FullScreen: window.Settings.FullScreen(),
	})
}
func windowSetFullScreen(w http.ResponseWriter, r *http.Request) {
	tmp := struct {
		FullScreen bool `json:"fullScreen"`
	}{}
	json.NewDecoder(r.Body).Decode(&tmp)

	EditorChannel <- func() {
		window.Settings.SetFullScreen(tmp.FullScreen)
	}
	w.WriteHeader(http.StatusOK)
}
func windowGetSize(w http.ResponseWriter, r *http.Request) {
	width, height := window.Settings.Size()
	json.NewEncoder(w).Encode(struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}{
		Width:  width,
		Height: height,
	})
}
func windowSetSize(w http.ResponseWriter, r *http.Request) {
	tmp := struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}{}
	json.NewDecoder(r.Body).Decode(&tmp)

	EditorChannel <- func() {
		window.Settings.SetSize(tmp.Width, tmp.Height)
	}
	w.WriteHeader(http.StatusOK)
}

func windowApply(w http.ResponseWriter, r *http.Request) {
	EditorChannel <- func() {
		window.Settings.Apply()
	}
}

func Start() {
	router := mux.NewRouter()

	window := router.PathPrefix("/window").Subrouter()
	window.HandleFunc("/vsync", windowGetVsync).Methods("GET")
	window.HandleFunc("/vsync", windowSetVsync).Methods("POST")
	window.HandleFunc("/fullScreen", windowGetFullScreen).Methods("GET")
	window.HandleFunc("/fullScreen", windowSetFullScreen).Methods("POST")
	window.HandleFunc("/size", windowGetSize).Methods("GET")
	window.HandleFunc("/size", windowSetSize).Methods("POST")
	window.HandleFunc("/apply", windowApply).Methods("GET")

	corsObj := handlers.AllowedOrigins([]string{"*"})

	http.ListenAndServe(":8000", handlers.CORS(corsObj)(router))
}
