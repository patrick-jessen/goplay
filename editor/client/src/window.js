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
      msaa: "4x",
      resolution: [1920, 1080]
    };

    // Get initial state
    api.getVsync(v => {
      if(v) this.setState({vsync:"Enabled"})
      else this.setState({vsync:"Disabled"})
    });

    this.onMode = this.onMode.bind(this);
    this.onRes = this.onRes.bind(this);
    this.onVsync = this.onVsync.bind(this);
    this.onMSAA = this.onMSAA.bind(this);
  }

  onVsync(v) {
    if(v == "Disabled") {
      api.setVsync(false);
      this.setState({vsync:"Disabled"});
    }
    else if (v == "Enabled") {
      api.setVsync(true);
      this.setState({vsync:"Enabled"});
    }
  }

  onMode(m) {
    if(m == "Windowed") {
      api.setDisplayMode(false, -1, -1);
      this.setState({mode:"Windowed"});
      
    }
    else {
      api.setDisplayMode(true, this.state.resolution[0], this.state.resolution[1]);
      this.setState({mode:"Fullscreen"});
    }
  }

  onRes(r) {
    if(r == "1920x1080") {
      this.setState({
        res: "1920x1080",
        resolution:[1920,1080]
      })
    }
    else if (r == "1366x768"){
      this.setState({
        res: "1366x768",
        resolution:[1366,768]
      })
    }
    else if (r == "1280x1024"){
      this.setState({
        res: "1280x1024",
        resolution:[1280,1024]
      })
    }
    else if (r == "1280x800"){
      this.setState({
        res: "1280x800",
        resolution:[1280,800]
      })
    }
    api.setDisplayMode(true, this.state.resolution[0], this.state.resolution[1])
  }

  onMSAA(m) {
    this.setState({msaa:m});
  }

  render({}, {vsync, mode, res, msaa}) {
    return (
      <div>
        <table>
          <Option 
            text="Display Mode"
            selected={mode}
            options={["Windowed","Fullscreen"]} 
            onSelect={this.onMode}
          />

          <Option 
            text="Resolution"
            selected={res}
            options={["1920x1080","1366x768","1280x1024","1280x800"]}
            onSelect={this.onRes}
            disabled={mode != "Fullscreen"}
          />

          <Option 
            text="Multisampling"
            selected={msaa} 
            options={["Off","2x","4x","8x","16x"]} 
            onSelect={this.onMSAA}
          />

          <Option 
            text="Vertical Sync"
            selected={vsync} 
            options={["Disabled","Enabled"]} 
            onSelect={this.onVsync}
          />
        </table>
      </div>
    );
  }
}