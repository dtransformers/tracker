#!/bin/bash

echo "[SCOPE] Validating changes against Root README.md..."

CHANGED_FILES=$(git diff --cached --name-only)

if [ -z "$CHANGED_FILES" ]; then
  echo "[SCOPE] No staged changes."
  exit 0
fi

# Simple heuristic check: ensure README exists and is not ignored
if [ ! -f "Root README.md" ]; then
  echo "[SCOPE][ERROR] Root README.md missing"
  exit 1
fi

# Optional: enforce keyword traceability (basic guard)
for file in $CHANGED_FILES; do
  grep -q "$file" "Root README.md" || echo "[SCOPE][WARN] $file not referenced in scope"
done

echo "[SCOPE] Done"