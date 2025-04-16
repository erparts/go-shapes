// shapes is a package for Ebitengine (the 2D game library written by Hajime Hoshi) that allows
// rendering some common shapes and effects in complementary ways to the official
// github.com/hajimehoshi/ebiten/v2/vector package. Unlike vector, shapes relies more on
// Kage shaders instead of raw triangles rasterization. This means rendering tends to be
// smoother, but extra care has to be taken as changing shaders or some of its parameters will
// break batching.
package shapes
