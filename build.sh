#!/bin/bash

# build params
# https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63


# cross-compiling wont work, so osx build will fail, unless run on OSX machine
# you can run it in a browser 
# https://ebiten.org/documents/webassembly.html

rm -rf ./bin

mkdir -p ./bin/linux
mkdir -p ./bin/osx
mkdir -p ./bin/win

GOOS=linux   GOARCH=amd64 go build -o ./bin/linux
GOOS=darwin  GOARCH=amd64 go build -o ./bin/osx
GOOS=windows GOARCH=amd64 go build -o ./bin/win

cp -r ./assets ./bin/linux/
cp -r ./assets ./bin/osx/
cp -r ./assets ./bin/win/


