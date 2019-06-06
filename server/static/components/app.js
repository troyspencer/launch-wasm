import React from 'react'

export default class App extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      isLoading: true
    }
  }
  
  componentDidMount() {
      const go = new Go()
      const fetchPromise = fetch('/static/main.wasm');
      WebAssembly.instantiateStreaming(fetchPromise, go.importObject).then(async (result) => {
        go.run(result.instance)
        this.setState({isLoading: false})
      });
  }

  render() {
    return this.state.isLoading ? <div>Loading</div> :  <div><button id="myButton">Click to say Hi in console!</button></div>
  }
}
