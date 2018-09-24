package editor

import (
	"encoding/json"
	"net/http"

	"github.com/patrick-jessen/goplay/engine/renderer"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/patrick-jessen/goplay/engine/texture"
	"github.com/patrick-jessen/goplay/engine/window"
)

var Channel = make(chan func(), 10)

func windowSetVsync(w http.ResponseWriter, r *http.Request) {
	tmp := struct {
		Vsync bool `json:"vsync"`
	}{}
	json.NewDecoder(r.Body).Decode(&tmp)

	Channel <- func() {
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

	Channel <- func() {
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

	Channel <- func() {
		window.Settings.SetSize(tmp.Width, tmp.Height)
	}
	w.WriteHeader(http.StatusOK)
}

func windowApply(w http.ResponseWriter, r *http.Request) {
	Channel <- func() {
		window.Settings.Apply()
	}
	w.WriteHeader(http.StatusOK)
}

func textureGetFilter(w http.ResponseWriter, r *http.Request) {
	f, a := texture.Settings.Filter()
	json.NewEncoder(w).Encode(struct {
		Filter int `json:"filter"`
		Aniso  int `json:"aniso"`
	}{
		Filter: int(f),
		Aniso:  a,
	})
}
func textureSetFilter(w http.ResponseWriter, r *http.Request) {
	tmp := struct {
		Filter int `json:"filter"`
		Aniso  int `json:"aniso"`
	}{}
	json.NewDecoder(r.Body).Decode(&tmp)

	Channel <- func() {
		texture.Settings.SetFilter(texture.Filter(tmp.Filter), tmp.Aniso)
	}
	w.WriteHeader(http.StatusOK)
}

func textureGetResolution(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(struct {
		Res int `json:"res"`
	}{
		Res: int(texture.Settings.Resolution()),
	})
}
func textureSetResolution(w http.ResponseWriter, r *http.Request) {
	tmp := struct {
		Res int `json:"res"`
	}{}
	json.NewDecoder(r.Body).Decode(&tmp)

	Channel <- func() {
		texture.Settings.SetResolution(uint(tmp.Res))
	}
	w.WriteHeader(http.StatusOK)
}
func textureApply(w http.ResponseWriter, r *http.Request) {
	Channel <- func() {
		texture.Settings.Apply()
	}
	w.WriteHeader(http.StatusOK)
}

func rendererGetAA(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(struct {
		AA int `json:"aa"`
	}{
		AA: int(renderer.Settings.Antialiasing()),
	})
}
func rendererSetAA(w http.ResponseWriter, r *http.Request) {
	tmp := struct {
		AA int `json:"aa"`
	}{}
	json.NewDecoder(r.Body).Decode(&tmp)

	Channel <- func() {
		renderer.Settings.SetAntialising(renderer.Antialiasing(tmp.AA))
	}
	w.WriteHeader(http.StatusOK)
}
func rendererApply(w http.ResponseWriter, r *http.Request) {
	Channel <- func() {
		renderer.Settings.Apply()
	}
	w.WriteHeader(http.StatusOK)
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

	texture := router.PathPrefix("/texture").Subrouter()
	texture.HandleFunc("/filter", textureGetFilter).Methods("GET")
	texture.HandleFunc("/filter", textureSetFilter).Methods("POST")
	texture.HandleFunc("/resolution", textureGetResolution).Methods("GET")
	texture.HandleFunc("/resolution", textureSetResolution).Methods("POST")
	texture.HandleFunc("/apply", textureApply).Methods("GET")

	renderer := router.PathPrefix("/renderer").Subrouter()
	renderer.HandleFunc("/aa", rendererGetAA).Methods("GET")
	renderer.HandleFunc("/aa", rendererSetAA).Methods("POST")
	renderer.HandleFunc("/apply", rendererApply).Methods("GET")

	corsObj := handlers.AllowedOrigins([]string{"*"})

	http.ListenAndServe(":8000", handlers.CORS(corsObj)(router))
}
