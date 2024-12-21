package autolevel

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// Boundaries struct to hold the min/max values for X, Y, and Z
type Boundaries struct {
	MinX, MaxX, MinY, MaxY, MinZ, MaxZ float64
}

// ParseGcodeBoundaries parses a Gcode file to find the boundaries
func ParseGcodeBoundaries(filePath string) (Boundaries, error) {
	// Initialize boundaries with extreme values
	boundaries := Boundaries{
		MinX: 1e9, MaxX: -1e9,
		MinY: 1e9, MaxY: -1e9,
		MinZ: 1e9, MaxZ: -1e9,
	}

	// Open the Gcode file
	file, err := os.Open(filePath)
	if err != nil {
		return boundaries, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Regular expressions to match X, Y, and Z coordinates
	reX := regexp.MustCompile(`X(-?\d+(\.\d+)?)`)
	reY := regexp.MustCompile(`Y(-?\d+(\.\d+)?)`)
	reZ := regexp.MustCompile(`Z(-?\d+(\.\d+)?)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Extract X coordinate
		if match := reX.FindStringSubmatch(line); match != nil {
			x, _ := strconv.ParseFloat(match[1], 64)
			if x < boundaries.MinX {
				boundaries.MinX = x
			}
			if x > boundaries.MaxX {
				boundaries.MaxX = x
			}
		}

		// Extract Y coordinate
		if match := reY.FindStringSubmatch(line); match != nil {
			y, _ := strconv.ParseFloat(match[1], 64)
			if y < boundaries.MinY {
				boundaries.MinY = y
			}
			if y > boundaries.MaxY {
				boundaries.MaxY = y
			}
		}

		// Extract Z coordinate
		if match := reZ.FindStringSubmatch(line); match != nil {
			z, _ := strconv.ParseFloat(match[1], 64)
			if z < boundaries.MinZ {
				boundaries.MinZ = z
			}
			if z > boundaries.MaxZ {
				boundaries.MaxZ = z
			}
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return boundaries, fmt.Errorf("error reading file: %v", err)
	}

	return boundaries, nil
}

// func main() {
// 	// Example usage
// 	filePath := "example.gcode"
// 	boundaries, err := ParseGcodeBoundaries(filePath)
// 	if err != nil {
// 		log.Fatalf("Error: %v", err)
// 	}

// 	fmt.Printf("Boundaries: %+v\n", boundaries)
// }
