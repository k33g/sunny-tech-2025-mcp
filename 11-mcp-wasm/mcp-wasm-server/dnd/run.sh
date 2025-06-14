#!/bin/bash

extism call ../dnd.wasm tools_information \
  --log-level "info" \
  --wasi
echo ""

extism call ../dnd.wasm orc_greetings \
  --input '{"name":"Bob Morane"}' \
  --log-level "info" \
  --wasi
echo ""
