#!/usr/bin/env bash

set -euo pipefail

# Get version from plugin.yaml
VERSION=$(grep 'version:' plugin.yaml | head -1 | sed 's/.*"\(.*\)"/\1/')
RELEASE_NAME="kubectl-replay_${VERSION}_linux_amd64"
DIST_DIR="dist"
TARBALL="${RELEASE_NAME}.tar.gz"

echo "Building release ${RELEASE_NAME}..."

# Create dist directory
mkdir -p "${DIST_DIR}/${RELEASE_NAME}"

# Build binary
echo "Building binary..."
go build -o "${DIST_DIR}/${RELEASE_NAME}/kubectl-replay" main.go

# Create tarball
echo "Creating tarball..."
cd "${DIST_DIR}"
tar -czf "${TARBALL}" "${RELEASE_NAME}"
cd ..

# Calculate checksum
echo ""
echo "Release created: ${DIST_DIR}/${TARBALL}"
echo ""
echo "SHA256 checksum:"
sha256sum "${DIST_DIR}/${TARBALL}"

echo ""
echo "To create a GitHub release, run:"
echo "  git tag ${VERSION}"
echo "  git push origin ${VERSION}"
echo "  gh release create ${VERSION} ${DIST_DIR}/${TARBALL} --title '${VERSION}' --notes 'Release ${VERSION}'"
