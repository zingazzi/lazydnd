#!/bin/bash
# scripts/update_version.sh
# Updates version in ui/text_content.go based on git tag

set -e

VERSION=$(git describe --tags --exact-match 2>/dev/null || echo "dev")

echo "Updating version to: $VERSION"

# Update ui/text_content.go
sed -i.bak "s/AppVersion *= *\".*\"/AppVersion   = \"$VERSION\"/" ui/text_content.go
rm -f ui/text_content.go.bak

echo "✅ Version updated to $VERSION in ui/text_content.go"

