#!/bin/sh
# Renders public/og.png (1200×630) from the /og-card page with headless
# Chrome. Chrome's window size includes window chrome, so the capture is
# taken taller and cropped to the card.  cd web && npm run build && sh scripts/og.sh
set -e
cd "$(dirname "$0")/.."
[ -f dist/og-card/index.html ] || { echo "og.sh: build first (npm run build)"; exit 1; }
port=8766
(cd dist && python3 -m http.server "$port" >/dev/null 2>&1 &)
sleep 1
tmp=$(mktemp --suffix=.png)
google-chrome --headless=new --disable-gpu --hide-scrollbars --window-size=1200,900 \
  --virtual-time-budget=4000 --screenshot="$tmp" "http://localhost:$port/og-card/" 2>/dev/null
pid=$(pgrep -f "http\.server $port" || true); [ -n "$pid" ] && kill $pid
python3 scripts/pngcrop.py "$tmp" public/og.png 1200 630
rm -f "$tmp"
cp public/og.png dist/og.png
echo "og.sh: wrote public/og.png"
