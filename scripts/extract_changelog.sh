#!/bin/bash
# scripts/extract_changelog.sh
# Extracts changelog for a specific version from CHANGELOG.md

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

# Remove 'v' prefix if present for matching
VERSION_NUMBER=${VERSION#v}

# Extract content between version header and next version header
awk -v ver="$VERSION_NUMBER" '
    /^## / {
        if (found) exit
        if ($2 == ver) found=1
        next
    }
    found { print }
' CHANGELOG.md

