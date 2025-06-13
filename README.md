# shapes
[![Go Reference](https://pkg.go.dev/badge/github.com/erparts/go-shapes.svg)](https://pkg.go.dev/github.com/erparts/go-shapes)

`shapes` is a package for [**Ebitengine**](https://github.com/hajimehoshi/ebiten) (the 2D game library written by Hajime Hoshi) that allows rendering some common shapes and effects in complementary ways to the official [`vector`](github.com/hajimehoshi/ebiten/v2/vector) package. Some examples:

- Smooth triangles, hexagons, ellipses and rects (all with the option of rounding).
- High quality image outlines, expansions and blurs.
- Simple patterns like dots, rhombuses and triangles.
- Some color functions and gradients.

Unlike `vector`, `shapes` relies more on simple [Kage shaders](https://github.com/tinne26/kage-desk) instead of raw triangles rasterization. This means rendering tends to be smoother, but extra care has to be taken as changing shaders or some of its parameters will break [batching](https://github.com/tinne26/efficient-ebitengine).

[TODO: simple UI surface example image]

# Credit

Many of the SDFs are based on https://iquilezles.org/articles/distfunctions2d.
