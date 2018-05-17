import { h, render, Component } from 'preact';
import Window from "./window";
import Texture from "./texture";

class App extends Component {
  constructor() {
    super()
  }


  render() {
    return (
      <div>
        <h4>Window</h4>
        <Window />

        <h4>Texture</h4>
        <Texture />
      </div>
    );
  }
}

render(<App />, document.body);