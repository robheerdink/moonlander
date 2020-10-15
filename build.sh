#!/bin/bash

# build params
# https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63


# cross-compiling wont work, so osx build will fail, unless run on OSX machine
# you can run it in a browser 
# https://ebiten.org/documents/webassembly.html
# GOOS=js GOARCH=wasm go build -o yourgame.wasm github.com/yourname/yourgame

rm -rf ./bin

mkdir -p ./bin/wasm
mkdir -p ./bin/linux
mkdir -p ./bin/osx
mkdir -p ./bin/win

GOOS=js      GOARCH=wasm  go build -o  ./bin/wasm/moonlander.wasm moonlander
GOOS=linux   GOARCH=amd64 go build -o ./bin/linux
GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/osx
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/win

cp -r ./assets ./bin/wasm/
cp -r ./assets ./bin/linux/
cp -r ./assets ./bin/osx/
cp -r ./assets ./bin/win/


# Copying wasm_exec.js to execute the Wasm binary
cp $(go env GOROOT)/misc/wasm/wasm_exec.js ./bin/wasm/

# create main html
echo '<!DOCTYPE html>
<script src="wasm_exec.js"></script>
<script>
// Polyfill
if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

const go = new Go();
WebAssembly.instantiateStreaming(fetch("moonlander.wasm"), go.importObject).then(result => {
  go.run(result.instance);
});
</script>
' > ./bin/wasm/main.html

# create iframe to insert main html (because it will get stretched by default)
echo '<!DOCTYPE html>
<iframe src="main.html" width="1150" height="864"></iframe>
' > ./bin/wasm/index.html