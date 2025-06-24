#!/bin/bash
tinygo build -scheduler=none --no-debug \
  -o gen-name.wasm \
  -target wasi main.go

