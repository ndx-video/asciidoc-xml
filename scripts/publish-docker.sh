#!/bin/bash
# Script to build and publish Docker images for AsciiDoc XML Converter
# Usage: ./scripts/publish-docker.sh [registry] [version]
# Example: ./scripts/publish-docker.sh docker.io/username v1.0.0

set -e

# Default values
REGISTRY="${1:-asciidoc-xml}"
VERSION="${2:-$(git describe --tags --always --dirty 2>/dev/null || echo 'latest')}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building and publishing Docker images${NC}"
echo -e "Registry: ${YELLOW}${REGISTRY}${NC}"
echo -e "Version: ${YELLOW}${VERSION}${NC}"
echo ""

# Build web image
echo -e "${GREEN}Building web image...${NC}"
docker build -f Dockerfile.web -t ${REGISTRY}/web:${VERSION} -t ${REGISTRY}/web:latest .

# Build watcher image
echo -e "${GREEN}Building watcher image...${NC}"
docker build -f Dockerfile.watcher -t ${REGISTRY}/watcher:${VERSION} -t ${REGISTRY}/watcher:latest .

# Ask for confirmation before pushing
read -p "Push images to registry? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${GREEN}Pushing web image...${NC}"
    docker push ${REGISTRY}/web:${VERSION}
    docker push ${REGISTRY}/web:latest
    
    echo -e "${GREEN}Pushing watcher image...${NC}"
    docker push ${REGISTRY}/watcher:${VERSION}
    docker push ${REGISTRY}/watcher:latest
    
    echo -e "${GREEN}Images pushed successfully!${NC}"
    
    # Create docker-compose.prod.yml with published images
    echo -e "${GREEN}Creating docker-compose.prod.yml...${NC}"
    cat > docker-compose.prod.yml <<EOF
version: '3.8'

services:
  web:
    image: ${REGISTRY}/web:${VERSION}
    container_name: asciidoc-xml-web
    ports:
      - "8005:8005"
    environment:
      - PORT=8005
    volumes:
      - ./examples:/app/examples:ro
      - ./docs:/app/docs:ro
      - ./xslt:/app/xslt:ro
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8005/"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 10s
    restart: unless-stopped
    networks:
      - asciidoc-xml-network

  watcher:
    image: ${REGISTRY}/watcher:${VERSION}
    container_name: asciidoc-xml-watcher
    ports:
      - "8006:8006"
    environment:
      - WATCH_DIR=/watch
    volumes:
      - ./watch:/watch
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8006/status"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 10s
    restart: unless-stopped
    networks:
      - asciidoc-xml-network

networks:
  asciidoc-xml-network:
    driver: bridge
EOF
    echo -e "${GREEN}Created docker-compose.prod.yml${NC}"
    echo -e "${YELLOW}To use the published images, run: docker-compose -f docker-compose.prod.yml up${NC}"
else
    echo -e "${YELLOW}Images built but not pushed.${NC}"
    echo -e "${YELLOW}To push manually, run:${NC}"
    echo -e "  docker push ${REGISTRY}/web:${VERSION}"
    echo -e "  docker push ${REGISTRY}/web:latest"
    echo -e "  docker push ${REGISTRY}/watcher:${VERSION}"
    echo -e "  docker push ${REGISTRY}/watcher:latest"
fi

echo ""
echo -e "${GREEN}Done!${NC}"

