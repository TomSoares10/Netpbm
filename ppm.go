package Netpbm

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint8
}

type Pixel struct {
	R, G, B uint8
}

func ReadPPM(filename string) (*PPM, error) {
	var err error
	var magicNumber string = ""
	var width int
	var height int
	var maxval int
	var counter int
	var headersize int
	var splitfile []string
	file, err := os.ReadFile(filename)
	if err != nil {
	}

	if strings.Contains(string(file), "\r") {
		splitfile = strings.SplitN(string(file), "\r\n", -1)
	} else {
		splitfile = strings.SplitN(string(file), "\n", -1)
	}
	for i, _ := range splitfile {
		if strings.Contains(splitfile[i], "P3") {
			magicNumber = "P3"
		} else if strings.Contains(splitfile[i], "P6") {
			magicNumber = "P6"
		}
		if strings.HasPrefix(splitfile[i], "#") && maxval != 0 {
			headersize = counter
		}
		splitl := strings.SplitN(splitfile[i], " ", -1)
		if width == 0 && height == 0 && len(splitl) >= 2 {
			width, err = strconv.Atoi(splitl[0])
			height, err = strconv.Atoi(splitl[1])
			headersize = counter
		}

		if maxval == 0 && width != 0 {
			maxval, err = strconv.Atoi(splitfile[i])
			headersize = counter
		}
		counter++
	}

	data := make([][]Pixel, height)

	for j := 0; j < height; j++ {
		data[j] = make([]Pixel, width)
	}
	var splitdata []string

	if counter > headersize {
		for i := 0; i < height; i++ {
			splitdata = strings.SplitN(splitfile[headersize+1+i], " ", -1)
			for j := 0; j < width*3; j += 3 {
				r, _ := strconv.Atoi(splitdata[j])
				if r > maxval {
					r = maxval
				}
				g, _ := strconv.Atoi(splitdata[j+1])
				if g > maxval {
					g = maxval
				}
				b, _ := strconv.Atoi(splitdata[j+2])
				if b > maxval {
					b = maxval
				}
				data[i][j/3] = Pixel{R: uint8(r), G: uint8(g), B: uint8(b)}
			}
		}
	}
	return &PPM{data: data, width: width, height: height, magicNumber: magicNumber, max: uint8(maxval)}, err
}

func display(data [][]Pixel) {
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[0]); j++ {
			fmt.Print(data[i][j], " ")
		}
		fmt.Println()
	}
}

// Size returns the width and height of the image.
func (ppm *PPM) Size() (int, int) {
	return ppm.height, ppm.width
}

// At returns the value of the pixel at (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[x][y] = value
}

// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert() {
	for i := 0; i < len(ppm.data); i++ {
		for j := 0; j < len(ppm.data[0]); j++ {
			ppm.data[i][j].R = ppm.max - ppm.data[i][j].R
			ppm.data[i][j].G = ppm.max - ppm.data[i][j].G
			ppm.data[i][j].B = ppm.max - ppm.data[i][j].B
		}
	}
}

// Flip flips the PPM image horizontally.
func (ppm *PPM) Flip() {
	// Height = Colums
	// Width = Rows
	NumRows := ppm.width
	NumColums := ppm.height
	for i := 0; i < NumRows; i++ {
		for j := 0; j < NumColums/2; j++ {
			ppm.data[i][j], ppm.data[i][NumColums-j-1] = ppm.data[i][NumColums-j-1], ppm.data[i][j]
		}
	}
}

// Flop flops the PPM image vertically.
func (ppm *PPM) Flop() {
	// Height = Colums
	// Width = Rows
	NumRows := ppm.width
	for i := 0; i < NumRows/2; i++ {
		ppm.data[i], ppm.data[NumRows-i-1] = ppm.data[NumRows-i-1], ppm.data[i]
	}
}

func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	oldMax := ppm.max
	ppm.max = maxValue
	for i := 0; i < len(ppm.data); i++ {
		for j := 0; j < len(ppm.data[0]); j++ {
			ppm.data[i][j].R = uint8(float64(ppm.data[i][j].R) * float64(ppm.max) / float64(oldMax))
			ppm.data[i][j].G = uint8(float64(ppm.data[i][j].G) * float64(ppm.max) / float64(oldMax))
			ppm.data[i][j].B = uint8(float64(ppm.data[i][j].B) * float64(ppm.max) / float64(oldMax))
		}
	}
}

// Rotate90CW rotates the PPM image 90Â° clockwise.
func (ppm *PPM) Rotate90CW() {
	// Height = Colums = Colonne vers le bas
	// Width = Rows = Ligne vers la droite
	NumRows := ppm.width
	NumColums := ppm.height
	for i := 0; i < NumColums; i++ {
		for j := i + 1; j < NumRows; j++ {
			vartemp := ppm.data[i][j]
			ppm.data[i][j] = ppm.data[j][i]
			ppm.data[j][i] = vartemp
		}
	}
	for i := 0; i < NumColums; i++ {
		for j := 0; j < NumRows/2; j++ {
			vartemp := ppm.data[i][j]
			ppm.data[i][j] = ppm.data[i][NumRows-j-1]
			ppm.data[i][NumRows-j-1] = vartemp
		}
	}
}

// ToPGM converts the PPM image to PGM.
func (ppm *PPM) ToPGM() *PGM {
	// Height = Colums = Colonne vers le bas
	// Width = Rows = Ligne vers la droite
	var newmagicnumber string

	if ppm.magicNumber == "P3" {
		newmagicnumber = "P2"
	} else if ppm.magicNumber == "P6" {
		newmagicnumber = "P5"
	}
	Numrows := ppm.width
	NumColumns := ppm.height
	var newdata = make([][]uint8, NumColumns)
	for i := 0; i < NumColumns; i++ {
		newdata[i] = make([]uint8, Numrows)
		for j := 0; j < Numrows; j++ {
			{
				newdata[i][j] = uint8((int(ppm.data[i][j].R) + int(ppm.data[i][j].G) + int(ppm.data[i][j].B)) / 3)
			}
		}
	}
	return &PGM{data: newdata, width: Numrows, height: NumColumns, max: ppm.max, magicNumber: newmagicnumber}
}

// ToPBM converts the PPM image to PBM.
func (ppm *PPM) ToPBM() *PBM {
	var newmagicnumber string

	if ppm.magicNumber == "P3" {
		newmagicnumber = "P1"
	} else if ppm.magicNumber == "P6" {
		newmagicnumber = "P4"
	}
	Numrows := ppm.width
	NumColumns := ppm.height
	var newdata = make([][]bool, NumColumns)
	for i := 0; i < NumColumns; i++ {
		newdata[i] = make([]bool, Numrows)
		for j := 0; j < Numrows; j++ {
			newdata[i][j] = uint8((int(ppm.data[i][j].R)+int(ppm.data[i][j].G)+int(ppm.data[i][j].B))/3) < ppm.max/2
		}
	}
	return &PBM{data: newdata, width: Numrows, height: NumColumns, magicNumber: newmagicnumber}
}

type Point struct {
	X, Y int
}

func maxAbs(a, b float64) float64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	if a > b {
		return a
	}
	return b
}

// DrawLine draws a line between two points.
func (ppm *PPM) SetPixel(p Point, color Pixel) {
	// Check if the point is within the PPM dimensions.
	if p.X >= 0 && p.X < ppm.width && p.Y >= 0 && p.Y < ppm.height {
		ppm.data[p.Y][p.X] = color
	}
}

// DrawLine draws a line between two points.

// DrawLine uses Bresenham's line algorithm to draw a line between two points.
// Bresenham's algorithm efficiently rasterizes a line on a grid of pixels.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	// Bresenham's line algorithm

	// Extract coordinates of the two points.
	x1, y1 := p1.X, p1.Y
	x2, y2 := p2.X, p2.Y

	// Calculate differences in x and y coordinates.
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)

	// Determine the direction of the line along the x-axis.
	var sx, sy int
	if x1 < x2 {
		sx = 1
	} else {
		sx = -1
	}

	// Determine the direction of the line along the y-axis.
	if y1 < y2 {
		sy = 1
	} else {
		sy = -1
	}

	// Initialize the error term.
	err := dx - dy

	// Iterate through the points along the line using Bresenham's algorithm.
	for {
		// Set the pixel at the current point on the line.
		ppm.SetPixel(Point{x1, y1}, color)

		// Check if the end point of the line is reached.
		if x1 == x2 && y1 == y2 {
			break
		}

		// Calculate the doubled error term.
		e2 := 2 * err

		// Update the error term based on the decision parameter.
		if e2 > -dy {
			err -= dy
			x1 += sx
		}

		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	// Draw the four sides of the rectangle using DrawLine.
	p2 := Point{p1.X + width, p1.Y}
	p3 := Point{p1.X + width, p1.Y + height}
	p4 := Point{p1.X, p1.Y + height}

	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p4, color)
	ppm.DrawLine(p4, p1, color)
}

// DrawFilledRectangle draws a filled rectangle.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	// Ensure positive width and height.
	if width <= 0 || height <= 0 {
		return
	}

	for w := width; w > 0; w-- {
		// Draw a rectangle with reduced width.
		ppm.DrawRectangle(p1, w, height, color)

		// Move the starting point for the next iteration.
		p1.X++
	}
}

// DrawTriangle draws a triangle.
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p1, p3, color)
	ppm.DrawLine(p3, p2, color)
}
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	var Corners [3]Point
	// Sort the vertices based on Y-coordinate.
	if p1.Y <= p2.Y && p1.Y <= p3.Y {
		Corners[0], Corners[1], Corners[2] = p1, p2, p3
	} else if p2.Y <= p1.Y && p2.Y <= p3.Y {
		Corners[0], Corners[1], Corners[2] = p2, p1, p3
	} else {
		Corners[0], Corners[1], Corners[2] = p3, p1, p2
	}
	// Calculate slopes for the two edges of the triangle.
	slope1 := float64(Corners[2].X-Corners[0].X) / float64(Corners[2].Y-Corners[0].Y)
	slope2 := float64(Corners[2].X-Corners[1].X) / float64(Corners[2].Y-Corners[1].Y)

	x1 := float64(Corners[0].X)
	x2 := float64(Corners[1].X)

	for y := Corners[0].Y; y <= Corners[1].Y; y++ {
		ppm.DrawLine(Point{int(x1 + 0.5), y}, Point{int(x2 + 0.5), y}, color)
		x1 += slope1
		x2 += slope2
	}
	x2 = float64(Corners[1].X)
	for y := Corners[1].Y + 1; y <= Corners[2].Y; y++ {
		ppm.DrawLine(Point{int(x1 + 0.5), y}, Point{int(x2 + 0.5), y}, color)
		x1 += slope1
		x2 += slope2
	}
}
func (ppm *PPM) setPixel(x, y int, color Pixel) {
	if x >= 0 && x < ppm.width && y >= 0 && y < ppm.height {
		ppm.data[y][x] = color
	}
}

func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	// converts polar coordinates to Cartesian coordinates

	// Ensure non-negative radius.
	if radius < 0 {
		return
	}
	for theta := -0.01; theta <= 1.99*math.Pi; theta += (1.0 / float64(radius)) {
		x := center.X + int(float64(radius)*math.Cos(theta))
		y := center.Y + int(float64(radius)*math.Sin(theta))

		ppm.setPixel(x, y, color)
	}
}

func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	// Make sure the radius is non-negative.
	if radius < 0 {
		return
	}

	// Fill the center point of the circle
	ppm.setPixel(center.X, center.Y, color)

	for theta := -0.01; theta <= 1.99*math.Pi; theta += (1.0 / float64(radius)) {
		x := center.X + int(float64(radius)*math.Cos(theta))
		y := center.Y + int(float64(radius)*math.Sin(theta))

		// Fill the horizontal line from the center to the edges of the circle
		for xi := x; xi < center.X; xi++ {
			ppm.setPixel(xi, y, color)
			ppm.setPixel(center.X*2-xi, y, color)
		}

		// Fill the vertical line from the center to the edges of the circle
		for yi := y; yi < center.Y; yi++ {
			ppm.setPixel(x, yi, color)
			ppm.setPixel(x, center.Y*2-yi, color)
		}
	}
}
func (ppm *PPM) drawHorizontalLine(x1, x2, y int, color Pixel) {
	// Ensure valid y-coordinate.
	if y < 0 || y >= ppm.height {
		return
	}

	// Ensure x1 is less than or equal to x2.
	if x1 > x2 {
		x1, x2 = x2, x1
	}

	// Clip x-coordinates to the image bounds.
	x1 = clamp(x1, 0, ppm.width-1)
	x2 = clamp(x2, 0, ppm.width-1)

	for x := x1; x <= x2; x++ {
		ppm.setPixel(x, y, color)
	}
}
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
