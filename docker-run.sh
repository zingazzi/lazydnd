#!/bin/bash
# docker-run.sh - Run LazyDnD in Docker

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Running LazyDnD in Docker${NC}"
echo ""

# Check if image exists
if ! docker image inspect lazydnd:latest >/dev/null 2>&1; then
    echo -e "${YELLOW}Image not found. Building...${NC}"
    ./docker-build.sh
fi

echo "Starting LazyDnD..."
echo ""
echo -e "${YELLOW}Press Ctrl+C to exit${NC}"
echo ""

# Run the container with interactive TTY
docker run -it --rm \
    --name lazydnd \
    -e TERM=xterm-256color \
    lazydnd:latest

