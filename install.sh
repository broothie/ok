#!/bin/sh
set -e

# ok installer script
# Usage: curl -sSL https://raw.githubusercontent.com/broothie/ok/main/install.sh | sh

REPO="broothie/ok"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="ok"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    printf "${GREEN}==>${NC} %s\n" "$1"
}

log_warn() {
    printf "${YELLOW}Warning:${NC} %s\n" "$1"
}

log_error() {
    printf "${RED}Error:${NC} %s\n" "$1" >&2
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "Linux";;
        Darwin*)    echo "Darwin";;
        *)
            log_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   echo "x86_64";;
        aarch64|arm64)  echo "arm64";;
        *)
            log_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
}

# Get latest release version from GitHub
get_latest_version() {
    if command -v curl >/dev/null 2>&1; then
        curl -sSf https://api.github.com/repos/${REPO}/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- https://api.github.com/repos/${REPO}/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        log_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
}

# Download file
download_file() {
    url=$1
    output=$2

    if command -v curl >/dev/null 2>&1; then
        curl -sSL "$url" -o "$output"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$url" -O "$output"
    fi
}

main() {
    log_info "Installing ok..."

    # Detect system
    OS=$(detect_os)
    ARCH=$(detect_arch)
    log_info "Detected OS: $OS, Architecture: $ARCH"

    # Get latest version
    log_info "Fetching latest release..."
    VERSION=$(get_latest_version)
    if [ -z "$VERSION" ]; then
        log_error "Failed to fetch latest version"
        exit 1
    fi
    log_info "Latest version: $VERSION"

    # Construct download URL
    ARCHIVE_NAME="ok_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT

    # Download archive
    log_info "Downloading $ARCHIVE_NAME..."
    ARCHIVE_PATH="$TMP_DIR/$ARCHIVE_NAME"
    if ! download_file "$DOWNLOAD_URL" "$ARCHIVE_PATH"; then
        log_error "Failed to download from $DOWNLOAD_URL"
        exit 1
    fi

    # Extract archive
    log_info "Extracting archive..."
    tar -xzf "$ARCHIVE_PATH" -C "$TMP_DIR"

    # Determine install directory
    if [ -w "$INSTALL_DIR" ]; then
        DEST="$INSTALL_DIR/$BINARY_NAME"
    elif [ -w "$HOME/.local/bin" ] || mkdir -p "$HOME/.local/bin" 2>/dev/null; then
        INSTALL_DIR="$HOME/.local/bin"
        DEST="$INSTALL_DIR/$BINARY_NAME"
        log_warn "Installing to $HOME/.local/bin (no write access to /usr/local/bin)"
        log_warn "Make sure $HOME/.local/bin is in your PATH"
    else
        log_error "No writable install directory found. Try running with sudo."
        exit 1
    fi

    # Install binary
    log_info "Installing to $DEST..."
    mv "$TMP_DIR/$BINARY_NAME" "$DEST"
    chmod +x "$DEST"

    # Verify installation
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        VERSION_OUTPUT=$("$BINARY_NAME" --version 2>&1 || echo "")
        log_info "✓ Successfully installed ok!"
        [ -n "$VERSION_OUTPUT" ] && echo "$VERSION_OUTPUT"
    else
        log_warn "Installation complete, but 'ok' is not in your PATH"
        log_warn "Add $INSTALL_DIR to your PATH or run: export PATH=\"$INSTALL_DIR:\$PATH\""
    fi
}

main
