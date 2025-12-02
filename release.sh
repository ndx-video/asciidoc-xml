#!/bin/bash
# scripts/release.sh - Create and push a new release

set -e

# Read version from VERSION file
VERSION=$(cat VERSION)
TAG="v${VERSION}"

echo "Creating release for version ${VERSION}"

# 1. Sync version files
echo "Syncing version files..."
make sync-version

# 2. Commit version changes
echo "Committing version changes..."
git add VERSION internal/version/VERSION
git commit -m "Bump version to ${VERSION}" || echo "No changes to commit"

# 3. Create annotated tag
echo "Creating tag ${TAG}..."
git tag -a "${TAG}" -m "Release version ${VERSION}"

# 4. Push changes and tag
echo "Pushing to GitHub..."
git push origin main  # or master, depending on your default branch
git push origin "${TAG}"

echo "✓ Release ${TAG} created and pushed successfully!"
echo ""
echo "Next steps:"
echo "  1. Create distribution packages: make dist-full VERSION=${VERSION}"
echo "  2. Publish Docker images: ./scripts/publish-docker.sh [registry] ${TAG}"
echo "  3. Create GitHub release at: https://github.com/[your-username]/asciidoc-xml/releases/new?tag=${TAG}"