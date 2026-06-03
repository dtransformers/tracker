#!/bin/bash

echo "[DOCS] Checking documentation consistency..."

CHANGED_FILES=$(git diff --cached --name-only)

if echo "$CHANGED_FILES" | grep -qE "internal/|pkg/|api/"; then
  if ! echo "$CHANGED_FILES" | grep -qE "README|docs/"; then
    echo "[DOCS][ERROR] Code changed without documentation update"
    exit 1
  fi
fi

echo "[DOCS] OK"