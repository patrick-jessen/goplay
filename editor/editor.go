package editor

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/patrick-jessen/goplay/engine/scene"
	"github.com/patrick-jessen/goplay/engine/window"
)

var EditorChannel = make(chan func(), 10)

func getScene(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(scene.Current())
}

func setVsync(w http.ResponseWriter, r *http.Request) {
	tmp := struct {
		Vsync bool `json:"vsync"`
	}{}
	json.NewDecoder(r.Body).Decode(&tmp)

	EditorChannel <- func() {
		window.Settings.SetVSync(tmp.Vsync)
		window.Settings.Apply()
	}
	w.WriteHeader(http.StatusOK)
}
func getVsync(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(struct {
		Vsync bool `json:"vsync"`
	}{
		Vsync: window.Settings.VSync(),
	})
}
func setDisplayMode(w http.ResponseWriter, r *http.Request) {
	tmp := struct {
		Fs bool `json:"fs"`
		W  int  `json:"w"`
		H  int  `json:"h"`
	}{}
	json.NewDecoder(r.Body).Decode(&tmp)

	EditorChannel <- func() {
		window.SetVideoMode(tmp.Fs, tmp.W, tmp.H)
	}
	w.WriteHeader(http.StatusOK)
}

func Start() {
	router := mux.NewRouter()
	router.HandleFunc("/scene", getScene).Methods("GET")
	router.HandleFunc("/vsync", getVsync).Methods("GET")
	router.HandleFunc("/vsync", setVsync).Methods("POST")
	router.HandleFunc("/displayMode", setDisplayMode).Methods("POST")

	corsObj := handlers.AllowedOrigins([]string{"*"})

	http.ListenAndServe(":8000", handlers.CORS(corsObj)(router))
}
