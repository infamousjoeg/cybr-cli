#!/bin/bash
set -eou pipefail

read -rp "amd64 or arm64? > " arch
if [[ $arch != "amd64" && $arch != "arm64" ]]; then
  read -rp "amd64 or arm64?" arch
fi

root="$HOME/go/src/github.com/infamousjoeg/cybr-cli/bin/darwin/$arch"

# Sign the executable with Developer ID Application certificate
codesign --deep --force --options=runtime \
  --sign "81176FD6CB590056A09A92B6FA9502A75F0BB3A1" --timestamp "$root/cybr"

mkdir -p "$root/cybr-cli"

pushd "./darwin/$arch" || exit
    tar -czf "cybr-cli_darwin_$arch.tar.gz" cybr
    md5 -qs "cybr-cli_darwin_$arch.tar.gz" > "cybr-cli_darwin_$arch.tar.gz.md5"
    shasum -a 256 "cybr-cli_darwin_$arch.tar.gz" | cut -f1 -d' ' > "cybr-cli_darwin_$arch.tar.gz.sha256"
    open .
popd || exit

# Build pkg installer structure
ditto "$root/cybr" "$root/cybr-cli/usr/local/bin/"

# Build the package
productbuild --identifier "com.github.infamousjoeg.cybr-cli" \
  --sign "42510BD69E802006418644C21E229DC4461D6673" --timestamp \
  --root "$root/cybr-cli" / "$root/cybr-cli_darwin_$arch.pkg"

# Notarize package with Apple
notary_resp=$(xcrun altool --notarize-app --primary-bundle-id com.github.infamousjoeg.cybr-cli \
  --username joe@joe-garcia.com --password "@keychain:altool" \
  --file "$root/cybr-cli_darwin_$arch.pkg")

request_uuid=$(echo "$notary_resp" | sed '1d' | cut -f3 -d' ')

echo ""
echo "Processing notarization with Apple..."
echo "In a new window, check status by running:"
echo -e "$ xcrun altool --notarization-info $request_uuid -u joe@joe-garcia.com -p "@keychain:altool""
echo "When Status is Package Approved, press ENTER."
echo ""
read -n 1 -rs

# Staple notarization to pkg
xcrun stapler staple "$root/cybr-cli_darwin_$arch.pkg"