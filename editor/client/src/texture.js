import {h, Component} from "preact";
import Option from "./option";
import api from "./api";

export default class Texture extends Component {
  constructor() {
    super();

    this.state = {
      filter: "Bilinear",
      res: "Normal",
    };

    this.onFilter = this.onFilter.bind(this);
    this.onRes = this.onRes.bind(this);
  }

  onFilter(f) {
    let filter, aniso;

    switch(f) {
      case "Bilinear": 
        filter = 0x2701;
        aniso = 0;
        break;
      case "Trilinear":
        filter = 0x2703;
        aniso = 0;
      case "2x Anisotropic":
        aniso = 2;
        break;
      case "4x Anisotropic":
        aniso = 4;
        break;
      case "8x Anisotropic":
        aniso = 8;
        break;
      case "16x Anisotropic":
        aniso = 16;
        break;
    }
    api.texture.setFilter(filter, aniso);

    this.setState({filter:f});
  }

  onRes(r) {
    this.setState({res:r});

    if(r == "Low")
      api.texture.setResolution(0);
    else if(r == "Normal")
      api.texture.setResolution(1);
    else if (r == "High")
      api.texture.setResolution(2);
  }

  onApply() {
    api.texture.apply();
  }

  render({}, {filter, res}) {
    return (
      <div>
        <table>
          <Option
            text="Filtering"
            selected={filter}
            options={[
              "Bilinear", 
              "Trilinear", 
              "2x Anisotropic", 
              "4x Anisotropic", 
              "8x Anisotropic", 
              "16x Anisotropic"]}
            onSelect={this.onFilter}
          />

          <Option
            text="Resolution"
            selected={res}
            options={[
              "Low", 
              "Normal", 
              "High"]}
            onSelect={this.onRes}
          />
        </table>
        <button onClick={this.onApply}>Apply</button>
      </div>
    );
  }
}