import React, {useEffect} from 'react'

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

export default function Game(props) {
  useEffect(() => {
    const go = new Go()
    fetch('/static/main.wasm').then(response =>
      response.arrayBuffer()
    ).then(bytes =>
      WebAssembly.instantiate(bytes, go.importObject)
    ).then(results => {
      props.onLoadingChange(false)
      go.run(results.instance)
    });
  }, [])

  return (
    <canvas style={styles.mycanvas} id="mycanvas" />
  );
}