# Convert Calendar

<script src="wasm_exec.js"></script>
<script>
    if (!WebAssembly.instantiateStreaming) { // polyfill
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }
    const go = new Go()
    WebAssembly.instantiateStreaming(fetch("conv.wasm"), go.importObject).
        then((result) => {
            go.run(result.instance)
            document.getElementById("input").disabled = false;
        })

    function Convert(files) {
        let reader = new FileReader();
        reader.onload = (e) => ConvToICS(new Uint8Array(e.target.result));
        reader.readAsArrayBuffer(files[0]);
    }
</script>
<input type="file" id="input" onchange="Convert(this.files)" disabled>