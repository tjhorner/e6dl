#!/bin/bash

latest_release="v0.0.9"

case $(uname -s) in
  Darwin*) os="darwin";;
  Linux*) os="linux";;
  *)
    echo "Unsupported system type! Please build directly from source."
    exit 1
    ;;
esac

case $(uname -m) in
  x86_64*) arch="amd64";;
  *)
    echo "Unsupported processor architecture! Please build directly from source."
    exit 1
    ;;
esac

echo "Detected system: $os, $arch"
echo "Downloading binary for release $latest_release from GitHub..."

url="https://github.com/tjhorner/e6dl/releases/download/$latest_release/e6dl-$os-$arch"

curl $url -Lo /usr/local/bin/e6dl
chmod +x /usr/local/bin/e6dl

echo "Done! Installed to /usr/local/bin/e6dl. Simply remove this binary to uninstall."