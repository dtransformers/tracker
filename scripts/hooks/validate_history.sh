#!/bin/bash

echo "[HISTORY] Checking HISTORY.md updates..."

CHANGED_FILES=$(git diff --cached --name-only)

if echo "$CHANGED_FILES" | grep -qv "commands/HISTORY.md"; then
  if git diff --cached --name-only | grep -q "internal/"; then
    if ! git diff --cached --name-only | grep -q "commands/HISTORY.md"; then
      echo "[HISTORY][ERROR] Missing update in commands/HISTORY.md"
      exit 1
    fi
  fi
fi

echo "[HISTORY] OK"