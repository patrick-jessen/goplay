import { h, render, Component } from "preact";
import Window from "./window";
import Texture from "./texture";
import Renderer from "./renderer";

class App extends Component {
  constructor() {
    super();
  }

  render() {
    return (
      <div>
        <h4>Window</h4>
        <Window />

        <h4>Texture</h4>
        <Texture />

        <h4>Renderer</h4>
        <Renderer />
      </div>
    );
  }
}

render(<App />, document.body);
