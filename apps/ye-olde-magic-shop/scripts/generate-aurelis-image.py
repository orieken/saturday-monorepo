#!/usr/bin/env python3
"""Regenerates public/aurelis.jpg, the equirectangular surface of Aurelis.

Land is drawn under every reagent site in apps/mock-api/data/reagent-sites.json, plus a few
extra landmasses so the planet does not look built around its sites. The output is committed;
rerun this only when the sites change. Needs Pillow. Run from apps/ye-olde-magic-shop.
"""
import json
from pathlib import Path

from PIL import Image, ImageDraw

WIDTH, HEIGHT = 1024, 512  # larger images made software rendering slow enough to time out
SEA, LAND, HIGHLAND, GRID = (27, 73, 130), (92, 138, 62), (74, 110, 50), (24, 58, 100)
SITE_RADIUS_DEGREES = (16, 11)  # half-width in longitude, half-height in latitude
EXTRA_LAND = [(-140, 35, 20, 14), (-5, -55, 24, 10), (95, 45, 22, 15), (-75, -5, 16, 12)]


def to_pixel(lng, lat):
    return (lng + 180) / 360 * WIDTH, (90 - lat) / 180 * HEIGHT


def ellipse(draw, lng, lat, half_lng, half_lat, color):
    # Drawn three times so land crossing the 180th meridian wraps around the image.
    for shift in (-360, 0, 360):
        left, top = to_pixel(lng + shift - half_lng, lat + half_lat)
        right, bottom = to_pixel(lng + shift + half_lng, lat - half_lat)
        draw.ellipse([left, top, right, bottom], fill=color)


def land_at(draw, lng, lat, half_lng, half_lat):
    # A pole's row spans the whole image, so land near one needs far more longitude to read as land.
    widen = 1 + (abs(lat) / 90) ** 8 * 8
    ellipse(draw, lng, lat, half_lng * widen, half_lat, LAND)
    ellipse(draw, lng, lat, half_lng * widen * 0.45, half_lat * 0.45, HIGHLAND)


def graticule(draw):
    for lng in range(-180, 180, 30):
        draw.line([to_pixel(lng, 90), to_pixel(lng, -90)], fill=GRID, width=1)
    for lat in range(-60, 90, 30):
        draw.line([to_pixel(-180, lat), to_pixel(180, lat)], fill=GRID, width=1)


def main():
    sites = json.loads(Path('../mock-api/data/reagent-sites.json').read_text())
    image = Image.new('RGB', (WIDTH, HEIGHT), SEA)
    draw = ImageDraw.Draw(image)
    for lng, lat, half_lng, half_lat in EXTRA_LAND:
        land_at(draw, lng, lat, half_lng, half_lat)
    for site in sites:
        land_at(draw, site['lng'], site['lat'], *SITE_RADIUS_DEGREES)
    graticule(draw)
    image.save('public/aurelis.jpg', quality=85)


if __name__ == '__main__':
    main()
