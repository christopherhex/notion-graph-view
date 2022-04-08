
let resultObj;
const wasmLib = async () => {

    if (resultObj) {
        return resultObj.instance;
    }

    const go = new Go();
    resultObj = await WebAssembly.instantiateStreaming(
        fetch("./assets/main.wasm"),
        go.importObject
    );
    go.run(resultObj.instance);

    return resultObj.instance
}

export default wasmLib;