"""Crop an 8-bit RGB/RGBA PNG to the top-left W×H. Pure Python, no deps.
    python3 scripts/pngcrop.py in.png out.png 1200 630"""
import struct, sys, zlib

def chunks(data):
    pos = 8
    while pos < len(data):
        n, = struct.unpack('>I', data[pos:pos+4]); t = data[pos+4:pos+8]
        yield t, data[pos+8:pos+8+n]; pos += 12 + n

def unfilter(raw, w, h, bpp):
    stride = w * bpp; out = bytearray(); prev = bytearray(stride); pos = 0
    for _ in range(h):
        f = raw[pos]; line = bytearray(raw[pos+1:pos+1+stride]); pos += 1 + stride
        for i in range(stride):
            a = line[i-bpp] if i >= bpp else 0; b = prev[i]; c = prev[i-bpp] if i >= bpp else 0
            if f == 1: line[i] = (line[i] + a) & 255
            elif f == 2: line[i] = (line[i] + b) & 255
            elif f == 3: line[i] = (line[i] + (a + b) // 2) & 255
            elif f == 4:
                p = a + b - c; pa, pb, pc = abs(p-a), abs(p-b), abs(p-c)
                line[i] = (line[i] + (a if pa <= pb and pa <= pc else b if pb <= pc else c)) & 255
        out += line; prev = line
    return out, stride

def chunk(t, d):
    c = struct.pack('>I', len(d)) + t + d
    return c + struct.pack('>I', zlib.crc32(t + d) & 0xffffffff)

src, dst, W, H = sys.argv[1], sys.argv[2], int(sys.argv[3]), int(sys.argv[4])
data = open(src, 'rb').read()
idat = b''; w = h = 0; ctype = 0
for t, d in chunks(data):
    if t == b'IHDR': w, h, depth, ctype = struct.unpack('>IIBB', d[:10]); assert depth == 8 and ctype in (2, 6)
    elif t == b'IDAT': idat += d
bpp = 4 if ctype == 6 else 3
px, stride = unfilter(zlib.decompress(idat), w, h, bpp)
assert W <= w and H <= h, f"crop {W}x{H} exceeds image {w}x{h}"
rows = b''.join(b'\x00' + bytes(px[y*stride:y*stride+W*bpp]) for y in range(H))
open(dst, 'wb').write(b'\x89PNG\r\n\x1a\n' + chunk(b'IHDR', struct.pack('>IIBBBBB', W, H, 8, ctype, 0, 0, 0))
                      + chunk(b'IDAT', zlib.compress(rows, 9)) + chunk(b'IEND', b''))
print(f"pngcrop: {w}x{h} -> {W}x{H}")
