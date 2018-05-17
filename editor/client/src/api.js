const baseURL = "http://localhost:8000/"

export default {

  window: {
    apply() {
      fetch(baseURL + "window/apply");
    },
    getVsync(then) {
      fetch(baseURL + "window/vsync")
        .then(r => r.json()).then(r => {
          then(r.vsync);
        })
    },
    setVsync(val) {
      fetch(baseURL + "window/vsync", {
        method: "POST", 
        body: JSON.stringify({vsync:val})
      });
    },
    getFullScreen(then) {
      fetch(baseURL + "window/fullScreen")
        .then(r => r.json()).then(r => {
          then(r.fullScreen);
        })
    },
    setFullScreen(on) {
      fetch(baseURL + "window/fullScreen", {
        method: "POST", 
        body: JSON.stringify({fullScreen:on})
      });
    },
    getSize(then) {
      fetch(baseURL + "window/size")
        .then(r => r.json()).then(r => {
          then(r.width, r.height);
        })
    },
    setSize(width, height) {
      fetch(baseURL + "window/size", {
        method: "POST", 
        body: JSON.stringify({width,height})
      });
    },
  },

  texture: {
    getFilter(then) {
      fetch(baseURL + "texture/filter")
        .then(r => r.json()).then(r => {
          then(r.filter, r.aniso);
        })
    },
    setFilter(f, a) {
      fetch(baseURL + "texture/filter", {
        method: "POST", 
        body: JSON.stringify({filter:f,aniso:a})
      });
    },
    getResolution(then) {
      fetch(baseURL + "texture/resolution")
        .then(r => r.json()).then(r => {
          then(r.res);
        })
    },
    setResolution(r) {
      fetch(baseURL + "texture/resolution", {
        method: "POST", 
        body: JSON.stringify({res:r})
      });
    },
    apply() {
      fetch(baseURL + "texture/apply");
    }
  }
}