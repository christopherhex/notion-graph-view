import './style';
import App from './components/app';

async function init_wasm() {
    const go = new Go();
    const result = await WebAssembly.instantiateStreaming(
        fetch("/assets/main.wasm"),
        go.importObject
    );
    go.run(result.instance);


}
init_wasm();

export default App;
