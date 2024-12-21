package util

import (
	"fmt"
	"math"
	"strings"
)

type Gcode string

func (g *Gcode) Add(format string, a ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	*g += Gcode(fmt.Sprintf(format, a...))
}
func (g *Gcode) G90Preamble() {
	*g += "G21\nG90\n"
}

func (g *Gcode) G91Preamble() {
	*g += "G21\nG91\n"
}
func (g *Gcode) Print() {
	fmt.Println(*g)
}

// Function to generate G-code for square with spiral fill
func (g *Gcode) G91SpiralFill(boxSize, spindle, beamDiameter, overlap, feedrate float64) {
	// Initialize G-code with laser setup
	g.Add("M3 S%.2f ; Set spindle\n", spindle)

	// Adjust boxSize for beam diameter
	adjustedSize := boxSize - beamDiameter
	halfBeam := beamDiameter / 2

	// Starting at the outer edge, account for beam diameter
	currentX, currentY := -halfBeam, -halfBeam
	currentSize := adjustedSize
	totalX, totalY := 0.0, 0.0 // Track total displacement

	step := beamDiameter * (1 - overlap) // Step size based on beam diameter and overlap

	// Start spiral fill loop
	for currentSize > 0 {
		// Move right
		g.Add("G1 X%.2f F%.2f ; Right\n", currentSize, feedrate)
		currentX += currentSize
		totalX += currentSize

		// Move up
		g.Add("G1 Y%.2f ; Up\n", currentSize)
		currentY += currentSize
		totalY += currentSize

		// Reduce the size for the next inward square
		currentSize -= step
		if currentSize <= 0 {
			break
		}

		// Move left
		g.Add("G1 X%.2f ; Left\n", -currentSize)
		currentX -= currentSize
		totalX -= currentSize

		// Move down
		g.Add("G1 Y%.2f ; Down\n", -currentSize)
		currentY -= currentSize
		totalY -= currentSize

		// Reduce again for the next pass
		currentSize -= step
	}

	// Calculate return move by subtracting total displacement from current position
	g.Add("G0 X%.2f Y%.2f ; Return to start position\n", -totalX, -totalY)

	// Turn off the laser
	g.Add("M5 ; Turn off spindle\n")

}

// ///////////////////////////////
func G90Preamble() string {
	return "G21\nG90\n"
}

func G91Preamble() string {
	return "G21\nG91\n"
}

// Function to generate G-code for square with spiral fill
func G91SpiralFill(boxSize, spindle, beamDiameter, overlap, feedrate float64) string {
	// Initialize G-code with laser setup
	gcode := fmt.Sprintf("M3 S%.2f ; Set spindle\n", spindle)

	// Adjust boxSize for beam diameter
	adjustedSize := boxSize - beamDiameter
	halfBeam := beamDiameter / 2

	// Starting at the outer edge, account for beam diameter
	currentX, currentY := -halfBeam, -halfBeam
	currentSize := adjustedSize
	totalX, totalY := 0.0, 0.0 // Track total displacement

	step := beamDiameter * (1 - overlap) // Step size based on beam diameter and overlap

	// Start spiral fill loop
	for currentSize > 0 {
		// Move right
		gcode += fmt.Sprintf("G1 X%.2f F%.2f ; Right\n", currentSize, feedrate)
		currentX += currentSize
		totalX += currentSize

		// Move up
		gcode += fmt.Sprintf("G1 Y%.2f ; Up\n", currentSize)
		currentY += currentSize
		totalY += currentSize

		// Reduce the size for the next inward square
		currentSize -= step
		if currentSize <= 0 {
			break
		}

		// Move left
		gcode += fmt.Sprintf("G1 X%.2f ; Left\n", -currentSize)
		currentX -= currentSize
		totalX -= currentSize

		// Move down
		gcode += fmt.Sprintf("G1 Y%.2f ; Down\n", -currentSize)
		currentY -= currentSize
		totalY -= currentSize

		// Reduce again for the next pass
		currentSize -= step
	}

	// Calculate return move by subtracting total displacement from current position
	gcode += fmt.Sprintf("G0 X%.2f Y%.2f ; Return to start position\n", -totalX, -totalY)

	// Turn off the laser
	gcode += "M5 ; Turn off spindle\n"

	return gcode
}

func G91Triangle(size, spindle, feedrate float64) string {

	/*
		g0
	*/
	out := fmt.Sprintf("M3 S%.2f ; Set spindle\n", spindle)
	out += fmt.Sprintf("G1 Y%.2f F%.2f ; Move right\n", size, feedrate)
	out += fmt.Sprintf("G1 X%.2f F%.2f ; Move left\n", size, feedrate)

	out += fmt.Sprintf("G1 X%.2f Y%.2f F%.2f ; Move up\n", -size, -size, feedrate)
	out += "M5 ; Turn off spindle\n"
	return out
}

// Function to calculate the hypotenuse using the Pythagorean theorem
func CalculateHypotenuse(a, b float64) float64 {
	return math.Sqrt(a*a + b*b)
}

// Function to calculate the rise given the slope and run
func CalculateRise(slope, run float64) float64 {
	return slope * run
}
