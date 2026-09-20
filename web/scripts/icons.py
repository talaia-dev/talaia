"""Generates web/public/favicon.ico (16/32/48) and apple-touch-icon.png (180)
for the talaia mark. Pure Python, no dependencies, deterministic.

    cd web && python3 scripts/icons.py

Geometry mirrors the definitive brand sheet
(web/brand/logo_talaia_negativo_variantes.svg) and public/favicon.svg:
a cream rounded tile with the navy machicolated tower (door as an
even-odd hole showing the tile) and the amber fire pixel above. The
mark's native units are 60x78; the tile composition is the sheet's
favicon variant: tile 62, mark at translate(15,9) scale(0.52).
Change the mark there, mirror it here, rerun."""
import struct, zlib, sys, os

CREAM=(0xFA,0xF6,0xEF); NAVY=(0x1A,0x22,0x38); FIRE=(0xE8,0xA3,0x3D)

def rect(x,y,w,h):
    return [(x,y),(x+w,y),(x+w,y+h),(x,y+h)]

TOWER=[rect(0,18,12,12),rect(24,18,12,12),rect(48,18,12,12),  # merlons
       rect(0,30,60,12),                                       # machicolated parapet
       rect(6,42,48,36),                                       # body
       rect(14,49,8,13)]                                       # door (even-odd hole)
FIRE_PX=rect(24,0,12,12)                                       # the signal fire

def in_poly(x,y,poly):
    inside=False; n=len(poly)
    for i in range(n):
        x1,y1=poly[i]; x2,y2=poly[(i+1)%n]
        if (y1>y)!=(y2>y):
            xi=x1+(y-y1)*(x2-x1)/(y2-y1)
            if x<xi: inside=not inside
    return inside

def in_evenodd(x,y,polys):
    c=0
    for p in polys:
        if in_poly(x,y,p): c+=1
    return c%2==1

def in_rrect(x,y,size,r):
    cx=cy=size/2; hw=hh=size/2
    dx=abs(x-cx)-(hw-r); dy=abs(y-cy)-(hh-r)
    if dx<=0 or dy<=0: return abs(x-cx)<=hw and abs(y-cy)<=hh
    return dx*dx+dy*dy<=r*r

def render(size, ss=6):
    unit=size/62.0                  # the sheet's favicon tile is 62 units
    scale=0.52*unit
    offx,offy=15*unit,9*unit
    r_tile=12*unit
    px=bytearray()
    for j in range(size):
        row=bytearray([0])
        for i in range(size):
            acc=[0.0,0.0,0.0,0.0]
            for sj in range(ss):
                for si in range(ss):
                    x=i+(si+0.5)/ss; y=j+(sj+0.5)/ss
                    if not in_rrect(x,y,size,r_tile): continue
                    dx=(x-offx)/scale; dy=(y-offy)/scale
                    col=CREAM
                    if in_evenodd(dx,dy,TOWER): col=NAVY
                    if in_poly(dx,dy,FIRE_PX): col=FIRE
                    acc[0]+=col[0]; acc[1]+=col[1]; acc[2]+=col[2]; acc[3]+=255
            n=ss*ss
            a=acc[3]/n
            if a>0:
                cov=acc[3]/255.0
                row+=bytes([int(acc[0]/cov+0.5),int(acc[1]/cov+0.5),int(acc[2]/cov+0.5),int(a+0.5)])
            else:
                row+=bytes([0,0,0,0])
        px+=row
    return png(size,size,bytes(px))

def chunk(t,d):
    c=struct.pack('>I',len(d))+t+d
    return c+struct.pack('>I',zlib.crc32(t+d)&0xffffffff)

def png(w,h,raw):
    return (b'\x89PNG\r\n\x1a\n'+chunk(b'IHDR',struct.pack('>IIBBBBB',w,h,8,6,0,0,0))
            +chunk(b'IDAT',zlib.compress(raw,9))+chunk(b'IEND',b''))

def ico(pngs):
    hdr=struct.pack('<HHH',0,1,len(pngs)); entries=b''; data=b''
    off=6+16*len(pngs)
    for size,p in pngs:
        entries+=struct.pack('<BBBBHHII',size%256,size%256,0,0,1,32,len(p),off+len(data))
        data+=p
    return hdr+entries+data

if __name__=='__main__':
    out=os.path.join(os.path.dirname(os.path.abspath(__file__)),'..','public')
    small=[(16,render(16)),(32,render(32)),(48,render(48))]
    open(os.path.join(out,'favicon.ico'),'wb').write(ico(small))
    open(os.path.join(out,'apple-touch-icon.png'),'wb').write(render(180,ss=4))
    print("wrote favicon.ico (16/32/48) and apple-touch-icon.png (180) to", os.path.normpath(out))
