import React from 'react'

const styles = {
    mycanvas: {
      position: "fixed",
      backgroundColor: "black",
      opacity: 1.0,
      width: "100%",
      height: "100%",
      top:0,
      right:0,
      bottom:0,
      left:0
    }
  }

export default class CanvasComponent extends React.Component {
    constructor(props) {
        super(props)
    
        this.state = {
          isLoading: true,
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
        return (
        <canvas style={styles.mycanvas} id="mycanvas" />
        );
    }
}

