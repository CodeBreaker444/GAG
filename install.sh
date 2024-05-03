#!/bin/bash

# Get the OS type and machine hardware name
OS="$(uname)"
ARCH="$(uname -m)"
echo "OS: $OS, ARCH: $ARCH"
RELEASE_VERSION="v1.0.0"

# Set the download URL based on the OS and ARCH
if [ "$OS" = "Linux" ]; then
    if [ "$ARCH" = "x86_64" ]; then
        URL="https://github.com/codebreaker444/GAG/releases/download/$RELEASE_VERSION/GAG-linux-amd64.tar.gz"
        BINARY="GAG-linux-amd64"
    elif [ "$ARCH" = "arm64" ]; then
        URL="https://github.com/codebreaker444/GAG/releases/download/$RELEASE_VERSION/GAG-linux-arm64.tar.gz"
        BINARY="GAG-linux-arm64"
    else
        echo "Unsupported architecture"
        exit 1
    fi
elif [ "$OS" = "Darwin" ]; then
    URL="https://github.com/codebreaker444/GAG/releases/download/$RELEASE_VERSION/GAG-macos.tar.gz"
    BINARY="GAG-macos"
else
    echo "Unsupported OS"
    exit 1
fi

# Download the tar file check if the download was successful
curl -L $URL -o GAG.tar.gz
# Check if the downloaded file is a valid tar.gz file
if ! tar tf GAG.tar.gz &> /dev/null; then
    echo "Error: Downloaded file is not a valid tar.gz file"
    rm GAG.tar.gz
    exit 1
fi


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