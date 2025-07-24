package shapes

import "math"

// Angle constants in radians, for use with drawing functions.
//
// The general conventions for working with angles in this package
// are the following:
//   - 0 is right, pi/2 is bottom, and so on (positive x axis goes left,
//     positive y axis goes down).
//   - All angles are expected in [0, 2*pi) range, and will be normalized
//     to this range if they are negative, using a mod/wrap operation.
//   - When working with ranges (start, end), the range is interpreted to
//     go in clockwise direction.
//   - A range of start == end will be considered empty.
//   - A range of end >= start + 2*pi will be considered full.
const (
	RadsRight  = 0.0
	RadsLeft   = math.Pi
	RadsBottom = math.Pi / 2
	RadsTop    = 3 * math.Pi / 2

	RadsBottomRight = math.Pi / 4
	RadsTopRight    = 7 * math.Pi / 4
	RadsBottomLeft  = 3 * math.Pi / 4
	RadsTopLeft     = 5 * math.Pi / 4
)

// Direction constants for use with gradient generation functions.
// See [RadsRight] constants for angle conventions and docs.
const (
	DirRadsLTR = 0.0             // left to right
	DirRadsRTL = math.Pi         // right to left
	DirRadsTTB = math.Pi / 2     // top to bottom
	DirRadsBTT = 3 * math.Pi / 2 // bottom to top

	DirRadsTLBR = math.Pi / 4     // top-left to bottom-right
	DirRadsBLTR = 7 * math.Pi / 4 // bottom-left to top-right
	DirRadsTRBL = 3 * math.Pi / 4 // top-right to bottom-left
	DirRadsBRTL = 5 * math.Pi / 4 // bottom-right to top-left
)
