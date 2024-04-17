#!/bin/bash

# Get the OS type and machine hardware name
OS="$(uname)"
ARCH="$(uname -m)"
echo "OS: $OS, ARCH: $ARCH"

# Set the download URL based on the OS and ARCH
if [ "$OS" = "Linux" ]; then
    if [ "$ARCH" = "x86_64" ]; then
        URL="https://github.com/yourusername/GAG/releases/download/v1.0/GAG-linux-amd64.tar.gz"
        BINARY="GAG-linux-amd64"
    elif [ "$ARCH" = "arm64" ]; then
        URL="https://github.com/yourusername/GAG/releases/download/v1.0/GAG-linux-arm64.tar.gz"
        BINARY="GAG-linux-arm64"
    else
        echo "Unsupported architecture"
        exit 1
    fi
elif [ "$OS" = "Darwin" ]; then
    URL="https://github.com/yourusername/GAG/releases/download/v1.0/GAG-macos.tar.gz"
    BINARY="GAG-macos"
else
    echo "Unsupported OS"
    exit 1
fi

# Download the tar file
curl -L -o GAG.tar.gz $URL

# Unzip the tar file
tar -xzf GAG.tar.gz

# Make the binary executable
chmod +x $BINARY

# Move the binary to /usr/local/bin
sudo mv $BINARY /usr/local/bin

# Add /usr/local/bin to the PATH if it's not already there
if [[ ":$PATH:" != *":/usr/local/bin:"* ]]; then
    echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc
    source ~/.bashrc
fi

echo "Installation completed successfully"