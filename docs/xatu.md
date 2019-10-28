# XI'AN TECHNOLOGICAL UNIVERSITY
西安工业大学课表转换器;

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
            document.getElementById("output").disabled = false;
        })
    function Convert() {
        let reader = new FileReader();
        reader.onload = (e) => ConvToICS(new Uint8Array(e.target.result), Download, Show);
        reader.readAsArrayBuffer(document.getElementById("input").files[0]);
    }
    function Download(output) {
        let blob = new Blob([output], { type: "text/calendar" });
        let a = document.createElement('a');
        a.download = document.getElementById("input").files[0].name + ".ics";
        a.href = URL.createObjectURL(blob);
        document.body.appendChild(a);
        // a.click();
        document.body.removeChild(a);
    }
    function Show(info) {
        let row = (name, value) => {
            let r = document.createElement('tr');
            let Name = document.createElement('td');
            let Value = document.createElement('td');
            Name.appendChild(document.createTextNode(name));
            Value.appendChild(document.createTextNode(value));
            r.appendChild(Name);
            r.appendChild(Value);
            return r;
        }
        let it = document.getElementById("infotable");
        it.appendChild(row("学年", info.year));
        it.appendChild(row("姓名", info.name));
        it.appendChild(row("学号", info.id));
        it.appendChild(row("学分", info.score));
    }
</script>

<input type="file" id="input">
<br><input type="button" id="Run" value="output" onclick="Convert()" disabled>
<table id="infotable"></table>

请稍等，加载转换器可能需要一段时间，完成后上面的按钮就会变得可用。
选择从教务网导出的.xls格式课表，点击`Run`按钮，将其转换为.ics文件。转换在浏览器端完成，你的课表将**不会**被上传。
