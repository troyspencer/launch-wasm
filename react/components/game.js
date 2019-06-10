import React from 'react'
import { Spin } from 'antd';
import FlexView from 'react-flexview';

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


export default class Game extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      isLoading: true,
    }
    this.onLoadingChange = this.onLoadingChange.bind(this)
  }

  componentDidMount() {
      const go = new Go()
      const fetchPromise = fetch('/static/main.wasm');
      WebAssembly.instantiateStreaming(fetchPromise, go.importObject).then(async (result) => {
        this.onLoadingChange(false)
        go.run(result.instance)
      });
  }

  onLoadingChange(isLoading) {
    this.setState({isLoading: isLoading})
  }

  render() {
    return (
      <div>
        <canvas style={styles.mycanvas} id="mycanvas" />
        <FlexView column vAlignContent='center' hAlignContent='center' hidden={!this.state.isLoading}>
          <FlexView vAlignContent='center' hAlignContent='center'>
            <Spin tip="Loading..." />
          </FlexView>
        </FlexView>
        
      </div>
    );
  }
}
