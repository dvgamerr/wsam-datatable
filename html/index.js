const wasmBrowserInstantiate = async (wasmModuleUrl, importObject) => {
  let response = undefined;

  if (WebAssembly.instantiateStreaming) {
    // Fetch the module, and instantiate it as it is downloading
    response = await WebAssembly.instantiateStreaming(fetch(wasmModuleUrl), importObject);
  } else {
    const fetchAndInstantiateTask = async () => {
      const wasmArrayBuffer = await fetch(wasmModuleUrl).then(response =>
        response.arrayBuffer()
      );
      return WebAssembly.instantiate(wasmArrayBuffer, importObject);
    };
    response = await fetchAndInstantiateTask();
  }
  return response;
}

const go = new Go();
const wasmLoad = async () => {
  const wasmModule = await wasmBrowserInstantiate("./main.wasm", go.importObject);
  var wasm = wasmModule.instance;
  go.run(wasm);
  wasm.exports.update();

  document.querySelector('#a').oninput = wasm.exports.update;
  document.querySelector('#b').oninput = wasm.exports.update;

  console.log('WASM Complated.')
};

console.log('WASM Loading...')
wasmLoad();
