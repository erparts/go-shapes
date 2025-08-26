package shapes

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Noise draws pseudo-random, hash-based white noise with the current renderer color over
// the given target.
//
// The cycle parameter controls the noise animation. Progressively increasing the cycle
// value from 0 to 1 and looping back to zero will create a continuous, looping animation
// with an organic feel. If you don't need animation, leave cycle to zero to reduce shader
// calculations.
//
// Seed must be in [0, 1].
func (r *Renderer) Noise(target *ebiten.Image, intensity float32, seed, cycle float32) {
	if seed < 0.0 || seed > 1.0 {
		panic("seed must be in [0..1]")
	}
	ensureShaderNoiseLoaded()
	r.setFlatCustomVAs(intensity, seed, cycle, 0.0)
	r.DrawShader(target, 0, 0, shaderNoise)
}

// NoiseGolden draws a grid geometric noise with the current renderer color over the
// given target. This noise is based on the golden ratio and it's highly sensitive
// to the scale, producing very different results at different levels. Some interesting
// scales are 0.06, 1.0, 64.0, 93.0 and upwards (patterns start to darken and vanish
// afterwards).
//
// The param t controls the animation pace. Increase t at a rate of 1.0 per second
// for a natural animation rate.
func (r *Renderer) NoiseGolden(target *ebiten.Image, scale, intensity, t float32) {
	ensureShaderNoiseGoldenLoaded()
	r.setFlatCustomVAs(scale, intensity, t, 0)
	r.DrawShader(target, 0, 0, shaderNoiseGolden)
}
