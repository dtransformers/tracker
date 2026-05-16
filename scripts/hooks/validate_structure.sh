#!/bin/bash

echo "[STRUCTURE] Checking project structure integrity..."

# Example rules
if find . -type f -name "*controller*" | grep -q "domain/"; then
  echo "[STRUCTURE][ERROR] Controllers cannot exist in domain layer"
  exit 1
fi

if find . -type f -name "*repository*" | grep -q "transport/"; then
  echo "[STRUCTURE][ERROR] Repository layer leakage detected"
  exit 1
fi

echo "[STRUCTURE] OK"