const baseURL = "http://localhost:8000/"

export default {
  getVsync(then) {
    return fetch("http://localhost:8000/vsync").then(r => r.json()).then(r => {
      then(r.vsync)
    })
  },
  setVsync(val) {
    fetch(baseURL + "vsync", {
      method: "POST", 
      body: JSON.stringify({vsync:val})
    });
  },
  setDisplayMode(fs, w, h) {
    fetch(baseURL + "displayMode", {
      method: "POST", 
      body: JSON.stringify({fs,w,h})
    });
  }
}