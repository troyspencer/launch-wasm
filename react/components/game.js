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
    if (!WebAssembly.instantiateStreaming) { // polyfill
			WebAssembly.instantiateStreaming = async (resp, importObject) => {
				const source = await (await resp).arrayBuffer();
				return await WebAssembly.instantiate(source, importObject);
			};
		}
    const go = new Go()
    const fetchPromise = fetch('/static/main.wasm');
    WebAssembly.instantiateStreaming(fetchPromise, go.importObject).then(async (result) => {
      props.onLoadingChange(false)
      go.run(result.instance)
    });
  }, [])

  return (
    <canvas style={styles.mycanvas} id="mycanvas" />
  );
}