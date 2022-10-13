import {h, Component} from "preact";
import Option from "./option";
import api from "./api";

export default class Renderer extends Component {
  constructor() {
    super();

    this.state = {
      antialiasing: 0
    };

    api.renderer.getAA(a => {
      let val;
      switch(a) {
        case 0: val = "None"; break;
        case 1: val = "FXAA"; break;
        case 2: val = "2x MSAA"; break;
        case 3: val = "4x MSAA"; break;
        case 4: val = "8x MSAA"; break;
        case 5: val = "16x MSAA"; break;
      }
      this.setState({antialiasing:val})
    })

    this.onAntialiasing = this.onAntialiasing.bind(this);
  }

  onAntialiasing(a) {
    let val;
    switch(a) {
      case "None": val=0; break;
      case "FXAA": val=1; break;
      case "2x MSAA": val=2; break;
      case "4x MSAA": val=3; break;
      case "8x MSAA": val=4; break;
      case "16x MSAA": val=5; break;
    }

    api.renderer.setAA(val);
    this.setState({antialiasing: a});
  }

  onApply() {
    api.renderer.apply();
  }

  render({}, {antialiasing}) {
    return (
      <div>
        <Option
          text="Antialiasing"
          options={["None", "FXAA", "2x MSAA", "4x MSAA", "8x MSAA", "16x MSAA"]}
          selected={antialiasing}
          onSelect={this.onAntialiasing}
        />

        <button onClick={this.onApply}>Apply</button>
      </div>
    );
  }
}