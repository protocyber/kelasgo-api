#!/bin/bash
echo "🔗 Creating symlinks for Leptonica and Tesseract (if needed)..."

SYMLINKS=(
  "/usr/local/include/leptonica"
  "/usr/local/lib/liblept.a"
  "/usr/local/lib/liblept.dylib"
  "/usr/local/include/tesseract"
  "/usr/local/lib/libtesseract.a"
  "/usr/local/lib/libtesseract.dylib"
)

for LINK in "${SYMLINKS[@]}"; do
  TARGET="/opt/homebrew/opt$(echo $LINK | sed 's|/usr/local||')"
  if [ ! -L "$LINK" ]; then
    echo "➡️  Creating symlink: $LINK → $TARGET"
    sudo ln -s "$TARGET" "$LINK"
  else
    echo "✅ Symlink exists: $LINK"
  fi
done

echo "🎉 macOS symlink setup complete!"
