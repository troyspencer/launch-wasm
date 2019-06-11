import React, {useState, useEffect} from 'react'
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

export default function Game() {
  const [loading,setLoading] = useState(true)
  useEffect(() => {
    const go = new Go()
    const fetchPromise = fetch('/static/main.wasm');
    WebAssembly.instantiateStreaming(fetchPromise, go.importObject).then(async (result) => {
      setLoading(false)
      go.run(result.instance)
    });
  }, [])

  return (
    <div>
      <canvas style={styles.mycanvas} id="mycanvas" />
      <FlexView column vAlignContent='center' hAlignContent='center' hidden={!loading}>
        <FlexView vAlignContent='center' hAlignContent='center'>
          <Spin tip="Loading..." />
        </FlexView>
      </FlexView>
      
    </div>
  );
}