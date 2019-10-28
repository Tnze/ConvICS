# XI'AN TECHNOLOGICAL UNIVERSITY
西安工业大学课表转换器
> 该页面还未完成，你的.ics文件将不会被导出，敬请期待。

<script src="wasm_exec.js"></script>
<script>
    if (!WebAssembly.instantiateStreaming) { // polyfill
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }
    const go = new Go()
    WebAssembly.instantiateStreaming(fetch("wasm/xatu.wasm"), go.importObject).
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

请稍等，加载转换器可能需要一段时间，完成后上面的按钮就会变得可用。
选择从教务网导出的.xls格式课表，将其转换为.ics文件。转换在浏览器端完成，你的课表将**不会**被上传。
