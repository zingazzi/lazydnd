#!/bin/bash
# scripts/extract_version.sh
# Extracts version from git tag or returns 'dev'

VERSION=$(git describe --tags --exact-match 2>/dev/null || echo "dev")
echo "$VERSION"

