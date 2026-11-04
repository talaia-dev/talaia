#!/bin/sh
# Renders the social cards (1200×630) with headless Chrome: /og-card to
# public/og.png and, for each translation, /<lang>/og-card to
# public/og-<lang>.png. Chrome's window size includes window chrome, so
# the capture is taken taller and cropped to the card.
#   cd web && npm run build && sh scripts/og.sh   (CHROME=… for another binary)
set -e
cd "$(dirname "$0")/.."
[ -f dist/og-card/index.html ] || { echo "og.sh: build first (npm run build)"; exit 1; }
port=8766
(cd dist && python3 -m http.server --bind 127.0.0.1 "$port" >/dev/null 2>&1 &)
sleep 1
for page in og-card $(cd dist && ls -d */og-card 2>/dev/null); do
  [ -f "dist/$page/index.html" ] || continue
  lang=${page%/og-card}; out=og.png; [ "$page" = og-card ] || out="og-$lang.png"
  tmp=$(mktemp --suffix=.png)
  "${CHROME:-google-chrome}" --headless=new --disable-gpu --hide-scrollbars --window-size=1200,900 \
    --virtual-time-budget=4000 --screenshot="$tmp" "http://localhost:$port/$page/" 2>/dev/null
  python3 scripts/pngcrop.py "$tmp" "public/$out" 1200 630
  rm -f "$tmp"
  cp "public/$out" "dist/$out"
  echo "og.sh: wrote public/$out"
done
pid=$(pgrep -f "http\.server $port" || true); if [ -n "$pid" ]; then kill $pid; fi
