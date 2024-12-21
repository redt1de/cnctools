package autolevel

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"golang.org/x/exp/rand"
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

// randomFloat generates a random float64 between min and max
func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// roundToDecimal rounds a float to the specified number of decimal places
func roundToDecimal(value float64, places int) float64 {
	multiplier := math.Pow(10, float64(places))
	return math.Round(value*multiplier) / multiplier
}

// fakeProbeData generates a random float between -1 and 1
func fakeProbeData() float64 {
	tmp := randomFloat(-1, 1)
	return roundToDecimal(tmp, 3)
}

type FakeSerialPort struct {
}

func (f *FakeSerialPort) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
}

// Point represents a single probed point.
type Point struct {
	X, Y, Z float64
}

// HeightMap is a collection of probed points
type HeightMap []Point

func ProbeGrid(port *FakeSerialPort, Xmin, Xmax, Ymin, Ymax float64, rows, cols int) HeightMap {
	// Initialize a 2D slice to store the probing results
	ret := make(HeightMap, 0)

	// Calculate step size between probing points
	stepX := (Xmax - Xmin) / float64(rows-1)
	stepY := (Ymax - Ymin) / float64(cols-1)

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			// Calculate the current probing position
			x := Xmin + float64(i)*stepX
			y := Ymin + float64(j)*stepY

			// Move to the probing point
			moveCommand := fmt.Sprintf("G0 X%.3f Y%.3f\n", x, y)
			if _, err := port.Write([]byte(moveCommand)); err != nil {
				log.Fatalf("Error moving to position X%.3f Y%.3f: %v", x, y, err)
			}

			// Probe and record the Z position
			probeCommand := "G38.2 Z-10 F50\n"
			if _, err := port.Write([]byte(probeCommand)); err != nil {
				log.Fatalf("Error sending probe command: %v", err)
			}

			// Read the response and extract Z
			z := parseZResponse(port) // Custom function to parse the Z-coordinate from the response

			// Store the probing result
			ret = append(ret, Point{X: x, Y: y, Z: z})
		}
	}

	return ret
}

func parseZResponse(port *FakeSerialPort) float64 {
	// Fake implementation, return a random float between -1 and 1
	return fakeProbeData()

}

func (h *HeightMap) Json() string {
	out, err := json.Marshal(h)

	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func (h *HeightMap) CSV() string {
	var out string
	for _, p := range *h {
		out += fmt.Sprintf("%.3f,%.3f,%.3f\n", p.X, p.Y, p.Z)
	}
	return out
}

func (h *HeightMap) Pretty() string {
	var out string
	out += "Height Map:\n"
	for _, p := range *h {
		out += fmt.Sprintf("   X: %.3f, Y: %.3f, Z: %.3f\n", p.X, p.Y, p.Z)
	}
	return out
}

func (h *HeightMap) GoCode() string {
	var out string
	out += "var hm = HeightMap{\n"

	for _, p := range *h {
		out += fmt.Sprintf("Point{X: %.3f, Y: %.3f, Z: %.3f},\n", p.X, p.Y, p.Z)
	}
	out += "}\n"
	return out
}
