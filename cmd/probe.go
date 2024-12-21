/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// probeCmd represents the probe command
var probeCmd = &cobra.Command{
	Use:   "probe",
	Short: "generate gcode for probing",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		xStart, _ := cmd.Flags().GetFloat64("x-start")
		xEnd, _ := cmd.Flags().GetFloat64("x-end")
		yStart, _ := cmd.Flags().GetFloat64("y-start")
		yEnd, _ := cmd.Flags().GetFloat64("y-end")
		gridSize, _ := cmd.Flags().GetIntSlice("grid-size")
		depth, _ := cmd.Flags().GetFloat64("probe-depth")
		safe, _ := cmd.Flags().GetFloat64("safe-height")
		feedZ, _ := cmd.Flags().GetInt("probe-feed")

		gridSizeX := gridSize[0]
		gridSizeY := gridSize[0]
		if len(gridSize) > 1 {
			gridSizeY = gridSize[1]
		}
		xInterval := (xEnd - xStart) / float64(gridSizeX)
		yInterval := (yEnd - yStart) / float64(gridSizeY)
		curX := xStart
		curY := yStart

		//
		p("G21\nG90")
		p("G0 Z%.1f", safe)
		// p("G90 G0 X%.3f Y%.3f", curX, curY)

		for i := 0; i < gridSizeX; i++ {
			for j := 0; j < gridSizeY; j++ {
				p("\nG90 G0 X%.3f Y%.3f", curX, curY)
				p("G38.2 Z%.1f F%d", depth, feedZ)
				p("G0 Z%.1f", safe)

				curX += xInterval
			}
			curX = xStart
			curY += yInterval
		}

	},
}

func init() {
	autolevelCmd.AddCommand(probeCmd)
	probeCmd.Flags().Float64P("x-start", "x", 0.0, "x start postion")
	probeCmd.Flags().Float64P("x-end", "X", 0.0, "x end postion")
	probeCmd.Flags().Float64P("y-start", "y", 0.0, "y start postion")
	probeCmd.Flags().Float64P("y-end", "Y", 0.0, "y end postion")
	probeCmd.Flags().IntSliceP("grid-size", "g", []int{10}, "rows and columns")
	probeCmd.Flags().Float64P("probe-depth", "d", -5.0, "z min, probe depth")
	probeCmd.Flags().IntP("probe-feed", "F", 25, "feed rate for probing")
	probeCmd.Flags().Float64P("safe-height", "s", 5.0, "z max, safe height")

}

// (AL: probing initial point)
// G0 Z3
// G90 G0 X0.000 Y0.000 Z3
// G38.2 Z-5 F25
// G0 Z3
// G90 G0 X5.000 Y0.000 F600
// G38.2 Z-5 F50
// G0 Z3
// G90 G0 X10.000 Y0.000 F600
// G38.2 Z-5 F50
// G0 Z3
// G90 G0 X0.000 Y5.000 F600
// G38.2 Z-5 F50
// G0 Z3
// G90 G0 X5.000 Y5.000 F600
// G38.2 Z-5 F50
// G0 Z3
// G90 G0 X10.000 Y5.000 F600
// G38.2 Z-5 F50
// G0 Z3
// G90 G0 X0.000 Y10.000 F600
// G38.2 Z-5 F50
// G0 Z3
// G90 G0 X5.000 Y10.000 F600
// G38.2 Z-5 F50
// G0 Z3
// G90 G0 X10.000 Y10.000 F600
// G38.2 Z-5 F50
// G0 Z3

func p(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Printf(format, args...)
}
