import { h, render, Component } from 'preact';
import Window from "./window"

class App extends Component {
  constructor() {
    super()
  }


  render() {
    return (
      <div>
        <Window />
      </div>
    );
  }
}

render(<App />, document.body);