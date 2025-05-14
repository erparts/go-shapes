package shapes

import (
	"image"
)

func rectOriginSize(bounds image.Rectangle) (ox, oy, w, h int) {
	return bounds.Min.X, bounds.Min.Y, bounds.Dx(), bounds.Dy()
}

func topBorder(bounds image.Rectangle, borderSize int) image.Rectangle {
	bounds.Max.Y = bounds.Min.Y + borderSize
	return bounds
}

func rightBorder(bounds image.Rectangle, borderSize int) image.Rectangle {
	bounds.Min.X = bounds.Max.X - borderSize
	return bounds
}

func bottomBorder(bounds image.Rectangle, borderSize int) image.Rectangle {
	bounds.Min.Y = bounds.Max.Y - borderSize
	return bounds
}

func leftBorder(bounds image.Rectangle, borderSize int) image.Rectangle {
	bounds.Max.X = bounds.Min.X + borderSize
	return bounds
}

// top border without overlapping the right border
func clockwiseTopBorder(bounds image.Rectangle, borderSize int) image.Rectangle {
	bounds = topBorder(bounds, borderSize)
	bounds.Max.X -= borderSize
	return bounds
}

// right border without overlapping the bottom
func clockwiseRightBorder(bounds image.Rectangle, borderSize int) image.Rectangle {
	bounds = rightBorder(bounds, borderSize)
	bounds.Max.Y -= borderSize
	return bounds
}

// bottom border without overlapping the left border
func clockwiseBottomBorder(bounds image.Rectangle, borderSize int) image.Rectangle {
	bounds = bottomBorder(bounds, borderSize)
	bounds.Min.X += borderSize
	return bounds
}

// left border without overlapping the top border
func clockwiseLeftBorder(bounds image.Rectangle, borderSize int) image.Rectangle {
	bounds = bottomBorder(bounds, borderSize)
	bounds.Min.Y += borderSize
	return bounds
}
