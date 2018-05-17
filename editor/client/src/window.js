import { h, Component } from 'preact';
import api from "./api"

function Option({text, selected, options, onSelect, disabled}) {
  return (
    <tr>
      <td>{text}</td>
      <td>
        <select onChange={(e)=>onSelect(e.target.value)} disabled={disabled}>
          {
            options.map(o => (
              <option selected={o == selected}>{o}</option>
            ))
          }
        </select>
      </td>
    </tr>
  );
}

export default class Window extends Component {
  constructor() {
    super()

    this.state = {
      vsync: "Disabled",
      mode: "Windowed",
      res: "1920x1080",
      resolution: [1920, 1080]
    };

    // Get initial state
    api.window.getVsync(v => {
      if(v) this.setState({vsync:"Enabled"});
      else this.setState({vsync:"Disabled"});
    });
    api.window.getFullScreen(f => {
      if(f) this.setState({mode:"Full screen"});
      else this.setState({mode:"Windowed"});
    });
    api.window.getSize((w,h) => {
      this.setState({
        resolution: [w,h],
        res: `${w}x${h}`
      })
    })

    this.onMode = this.onMode.bind(this);
    this.onRes = this.onRes.bind(this);
    this.onVsync = this.onVsync.bind(this);
    this.onApply = this.onApply.bind(this);
  }

  onVsync(v) {
    if(v == "Disabled")
      api.window.setVsync(false);
    else if (v == "Enabled")
      api.window.setVsync(true);

    this.setState({vsync:v});
  }

  onMode(m) {
    if (m == "Windowed") 
      api.window.setFullScreen(false);
    else
      api.window.setFullScreen(true);

    this.setState({mode:m});
  }

  onRes(r) {
    let spl = r.split("x");

    this.setState({
      res: r,
      resolution:[parseInt(spl[0]), parseInt(spl[1])]
    });

    api.window.setSize(this.state.resolution[0], this.state.resolution[1]);
  }

  onApply() {
    api.window.apply();
  }

  render({}, {vsync, mode, res, msaa}) {
    return (
      <div>
        <table>
          <Option 
            text="Display Mode"
            selected={mode}
            options={["Windowed","Full screen"]} 
            onSelect={this.onMode}
          />

          <Option 
            text="Resolution"
            selected={res}
            options={["1920x1080","1366x768","1280x1024","1280x800", "1024x768"]}
            onSelect={this.onRes}
            disabled={mode != "Full screen"}
          />

          <Option 
            text="Vertical Sync"
            selected={vsync} 
            options={["Disabled","Enabled"]} 
            onSelect={this.onVsync}
          />

          <button onClick={this.onApply}>Apply</button>
        </table>
      </div>
    );
  }
}