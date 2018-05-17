const baseURL = "http://localhost:8000/"

export default {

  window: {
    apply() {
      fetch(baseURL + "window/apply");
    },
    getVsync(then) {
      return fetch(baseURL + "window/vsync")
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
  }
}