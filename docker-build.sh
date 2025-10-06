#!/bin/bash
# docker-build.sh - Build Docker image for LazyDnD

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building LazyDnD Docker Image${NC}"
echo ""

# Get version from git tag or default
VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")
echo -e "Version: ${YELLOW}${VERSION}${NC}"
echo ""

# Build the image
echo "Building Docker image..."
docker build -t lazydnd:${VERSION} -t lazydnd:latest .

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ Build successful!${NC}"
    echo ""
    echo "Image tags:"
    echo "  - lazydnd:${VERSION}"
    echo "  - lazydnd:latest"
    echo ""
    echo "Run with:"
    echo "  docker run -it --rm lazydnd:latest"
    echo ""
    echo "Or use docker-compose:"
    echo "  docker-compose up"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

