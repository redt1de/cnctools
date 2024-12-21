/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	spindleSpeed = 10000
	cornerSize   = 10
	Conventional = 0
	Climb        = 1
)

var (
	xMax     float64
	yMax     float64
	cutDepth float64
	stepOver float64
	feedRate float64
	climb    bool
	corner   bool
)

// spoilboardCmd represents the spoilboard command
var spoilboardCmd = &cobra.Command{
	Use:     "spoilboard",
	Aliases: []string{"spoil"},
	Short:   "generate surfacing gcode",
	Long:    `Zero on corner and define X and Y limits, optionally set cut depth, step over and feed rate.`,
	Run: func(cmd *cobra.Command, args []string) {
		startX, startY := 0.0, 0.0
		xMax, _ = cmd.Flags().GetFloat64("x-max")
		yMax, _ = cmd.Flags().GetFloat64("y-max")
		cutDepth, _ = cmd.Flags().GetFloat64("depth")
		stepOver, _ = cmd.Flags().GetFloat64("step-over")
		feedRate, _ = cmd.Flags().GetFloat64("feed-rate")
		corner, _ = cmd.Flags().GetBool("corner")
		mode := Conventional
		if climb {
			mode = Climb
		}
		gcode := GenerateRectangularSpiralSurfacing(startX, startY, xMax, yMax, cutDepth, feedRate, stepOver, mode)
		fmt.Println(gcode)

	},
}

func init() {
	rootCmd.AddCommand(spoilboardCmd)
	spoilboardCmd.Flags().Float64P("x-max", "x", 200, "X max position")
	spoilboardCmd.Flags().Float64P("y-max", "y", 150, "Y max position")
	spoilboardCmd.Flags().Float64P("depth", "d", 2, "depth of cut")
	spoilboardCmd.Flags().Float64P("step-over", "s", 8, "stepover, should be roughly half of the tool diameter")
	spoilboardCmd.Flags().Float64P("feed-rate", "f", 100, "feed rate")
	spoilboardCmd.Flags().BoolP("corner", "c", false, "leave an alignement corner")
}

// GenerateRectangularSpiralSurfacing generates G-code for a rectangular spiral surfacing operation.
// Arguments: xMax, yMax -> size of the spoilboard; cutDepth -> depth of the cut; feedRate -> feed rate in units/min; stepover -> distance between successive spiral loops; cuttingMode -> conventional or climb cutting.
func GenerateRectangularSpiralSurfacing(startX, startY, xMax, yMax, cutDepth, feedRate, stepover float64, cuttingMode int) string {
	z := -cutDepth
	gcode := ""

	gcode += "G21 ; Set units to mm\n"
	gcode += "G90 ; Absolute positioning\n"

	gcode += fmt.Sprintf("G0 X%.2f Y%.2f ; Move to starting corner\n", startX, startY)
	gcode += fmt.Sprintf("M3 S%d ; Start spindle\n", spindleSpeed)
	gcode += fmt.Sprintf("G1 Z%.2f F25 ; Move down to fixed cutting depth\n", z)

	if corner {
		gcode += fmt.Sprintf("G1 X%.2f Y%.2f F%.2f\n", startX+cornerSize, startY+cornerSize, feedRate)
		startX += cornerSize
		startY += cornerSize
		gcode += fmt.Sprintf("G1 Y%.2fF%.2f\n", yMax, feedRate)
		gcode += fmt.Sprintf("G1 Y%.2fF%.2f\n", startY, feedRate)

	}

	gcode += fmt.Sprintf("G1 F%.2f ; Set feed rate\n", feedRate)
	xMin, yMin := startX, startY
	xMaxCurrent, yMaxCurrent := xMax, yMax

	// Loop to generate rectangular spiral moves
	for xMin < xMaxCurrent && yMin < yMaxCurrent {
		gcode += fmt.Sprintf("G1 X%.2f Y%.2f\n", xMaxCurrent, yMin)
		gcode += fmt.Sprintf("G1 X%.2f Y%.2f\n", xMaxCurrent, yMaxCurrent)
		gcode += fmt.Sprintf("G1 X%.2f Y%.2f\n", xMin, yMaxCurrent)
		gcode += fmt.Sprintf("G1 X%.2f Y%.2f\n", xMin, yMin+stepover)

		xMin += stepover
		yMin += stepover
		xMaxCurrent -= stepover
		yMaxCurrent -= stepover
	}

	gcode += "G0 Z5.00 ; Retract tool\n"
	gcode += "M5 ; Stop spindle\n"
	gcode += "M30 ; End program\n"

	return gcode
}
