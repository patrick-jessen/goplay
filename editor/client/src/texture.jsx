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

    api.texture.getFilter((f, a) => {
      let filter;
      if(a == 16) filter = "16x Anisotropic";
      else if(a == 8) filter = "8x Anisotropic";
      else if(a == 4) filter = "4x Anisotropic";
      else if(a == 2) filter = "2x Anisotropic";
      else if(f == 0x2703) filter = "Trilinear";
      else filter = "Bilinear";
      this.setState({filter});
    });
    api.texture.getResolution(r => {
      let res;
      if (r == 16) res = "Low";
      else if (r == 8) res = "Medium";
      else if (r == 1) res = "High";
      this.setState({res});
    });
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
      api.texture.setResolution(16);
    else if(r == "Medium")
      api.texture.setResolution(8);
    else if (r == "High")
      api.texture.setResolution(1);
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
              "Medium", 
              "High"]}
            onSelect={this.onRes}
          />
        </table>
        <button onClick={this.onApply}>Apply</button>
      </div>
    );
  }
}