package autolevel

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
)

// parseGcodeLine extracts X, Y, Z values from a Gcode line
func extractPos(gcodeLine string) (float64, float64, float64) {
	var x, y, z float64
	x, y, z = math.NaN(), math.NaN(), math.NaN() // Default to NaN if not found

	reX := regexp.MustCompile(`X(-?\d+(\.\d+)?)`)
	reY := regexp.MustCompile(`Y(-?\d+(\.\d+)?)`)
	reZ := regexp.MustCompile(`Z(-?\d+(\.\d+)?)`)

	if match := reX.FindStringSubmatch(gcodeLine); match != nil {
		x, _ = strconv.ParseFloat(match[1], 64)
	}
	if match := reY.FindStringSubmatch(gcodeLine); match != nil {
		y, _ = strconv.ParseFloat(match[1], 64)
	}
	if match := reZ.FindStringSubmatch(gcodeLine); match != nil {
		z, _ = strconv.ParseFloat(match[1], 64)
	}

	return x, y, z
}

func ApplyHeightMap(gcodeFile string, heightMap HeightMap) []string {
	file, err := os.Open(gcodeFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var curX, curY, curZ float64
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		curLine := scanner.Text()
		r := regexp.MustCompile(`(?i)^(G21|G22|G38|G10)`) // ignore commands with X,Y,Z that dont actually move
		if curLine == "" || r.MatchString(curLine) {
			fmt.Println(curLine)
			continue
		}

		x, y, z := extractPos(curLine)
		if !math.IsNaN(x) {
			curX = x
		}
		if !math.IsNaN(y) {
			curY = y
		}
		if !math.IsNaN(z) {
			curZ = z
		}

		zoffset, err := heightMap.FindZOffset(curX, curY)
		if err != nil {
			log.Fatal(err)
		}

		lineout := fmt.Sprintf("%s ; %.3f + %.3f -> %.3f", curLine, curZ, zoffset, curZ+zoffset)
		fmt.Println(lineout)

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return nil
}

// FindZOffset calculates the Z offset for a given (x, y) using plane fitting
func (hm HeightMap) FindZOffset(x, y float64) (float64, error) {
	zoff, ok := hm.matchExact(x, y)
	if ok {
		return zoff, nil
	}

	// Find three valid closest points
	points := findThreeValidPoints(x, y, hm)
	if len(points) < 3 {
		return 0, fmt.Errorf("not enough valid points for plane fitting")
	}

	// Perform plane fitting
	p1, p2, p3 := points[0], points[1], points[2]

	// Calculate two vectors on the plane
	u := Point{X: p2.X - p1.X, Y: p2.Y - p1.Y, Z: p2.Z - p1.Z}
	v := Point{X: p3.X - p1.X, Y: p3.Y - p1.Y, Z: p3.Z - p1.Z}

	// Compute the normal vector of the plane
	normal := Point{
		X: u.Y*v.Z - u.Z*v.Y,
		Y: u.Z*v.X - u.X*v.Z,
		Z: u.X*v.Y - u.Y*v.X,
	}

	// Check for invalid plane (normal.Z should not be 0)
	if math.Abs(normal.Z) < 1e-9 {
		return 0, fmt.Errorf("plane is invalid; normal vector Z component is zero")
	}

	// Solve the plane equation for Z at the given (x, y)
	z := p1.Z - (normal.X*(x-p1.X)+normal.Y*(y-p1.Y))/normal.Z

	return z, nil
}

// findThreeValidPoints finds three non-collinear points closest to (x, y)
func findThreeValidPoints(x, y float64, hm HeightMap) []Point {
	type distPoint struct {
		point    Point
		distance float64
	}

	// Calculate distances and store them alongside points
	distances := []distPoint{}
	for _, p := range hm {
		dist := math.Sqrt((p.X-x)*(p.X-x) + (p.Y-y)*(p.Y-y))
		distances = append(distances, distPoint{point: p, distance: dist})
	}

	// Sort by distance
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].distance < distances[j].distance
	})

	// Extract non-collinear points
	closest := []Point{}
	for i := 0; i < len(distances) && len(closest) < 3; i++ {
		if len(closest) < 2 {
			closest = append(closest, distances[i].point)
		} else {
			// Check if the new point is collinear with the last two points
			if !isCollinear(closest[0], closest[1], distances[i].point) {
				closest = append(closest, distances[i].point)
			}
		}
	}

	return closest
}

// isCollinear checks if three points are collinear in 3D space
func isCollinear(p1, p2, p3 Point) bool {
	u := Point{X: p2.X - p1.X, Y: p2.Y - p1.Y, Z: p2.Z - p1.Z}
	v := Point{X: p3.X - p1.X, Y: p3.Y - p1.Y, Z: p3.Z - p1.Z}
	// Calculate the cross product
	cross := Point{
		X: u.Y*v.Z - u.Z*v.Y,
		Y: u.Z*v.X - u.X*v.Z,
		Z: u.X*v.Y - u.Y*v.X,
	}
	// Check if the magnitude of the cross product is near zero
	return math.Abs(cross.X) < 1e-9 && math.Abs(cross.Y) < 1e-9 && math.Abs(cross.Z) < 1e-9
}

func (hm HeightMap) matchExact(x, y float64) (float64, bool) {
	for _, p := range hm {
		if p.X == x && p.Y == y {
			return p.Z, true
		}
	}
	return 0, false
}
