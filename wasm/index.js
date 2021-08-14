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
  const importObject = go.importObject;
  const wasmModule = await wasmBrowserInstantiate("./wasm/main.wasm", importObject);

  go.run(wasmModule.instance);

  // Call the Add function export from wasm, save the result
  const addResult = wasmModule.instance.exports.add(24, 24);
  document.body.textContent = `Hello World! addResult: ${addResult}`;
};
console.log('wsam Loading...')
wasmLoad();