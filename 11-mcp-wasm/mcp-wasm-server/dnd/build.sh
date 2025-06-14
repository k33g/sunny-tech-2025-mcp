#!/bin/bash
tinygo build -scheduler=none --no-debug \
  -o ../dnd.wasm \
  -target wasi main.go

