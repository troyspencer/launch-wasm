import React, {useEffect} from 'react'
import Title from 'antd/lib/typography/Title';

const styles = {
  title: {
    marginTop: "5em",
    color: "rgb(180,180,180)",
    textAlign: "center"
  },
  noWasm: {
    position: "fixed",
    backgroundColor: "black",
    width: "100%",
    height: "100%",
  },
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
  const wasmSupported = (typeof WebAssembly === "object");
  if (!wasmSupported) {
    props.onLoadingChange(false)
    return (
      <div style={styles.noWasm}>
        <Title level={2} style={styles.title}>WebAssembly is not supported on this device or browser.</Title>
      </div>
    )
  }

  useEffect(() => {
    const go = new Go()
    const fetchPromise = fetch('/static/main.wasm');
    WebAssembly.instantiateStreaming(fetchPromise, go.importObject).then(async (result) => {
      props.onLoadingChange(false)
      props.onLoadedChange(true)
      go.run(result.instance)
    });
  }, [])

  return (
    <canvas style={styles.mycanvas} id="mycanvas" />
  );
}