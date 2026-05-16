#!/bin/bash

echo "[TESTS] Running test suite..."

go test ./...   # replace with npm test / pytest if needed

if [ $? -ne 0 ]; then
  echo "[TESTS][ERROR] Tests failed"
  exit 1
fi

echo "[TESTS] OK"