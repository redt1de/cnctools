/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/redt1de/cnctools/util"
	"github.com/spf13/cobra"
)

// 2X2 L
/*
G21
G91
;
M3 S100
G01 Y2 F100
G01 X2 F100
G01 X-2 Y-2 F100
M5
*/

// focusCmd represents the focus command
var focusCmd = &cobra.Command{
	Use:   "focus",
	Short: "A brief description of your command",
	Long: `
1. touch off the laser and zero
2. run focus command
3. turn on laser to lowest power and jog X until its over the cleanest point on the line. 
4. run the focus command again with the --calculate [X] flag
`,
	Run: func(cmd *cobra.Command, args []string) {
		Zmin, _ := cmd.Flags().GetFloat64("z-min")
		Zmax, _ := cmd.Flags().GetFloat64("z-max")
		power, _ := cmd.Flags().GetInt("power")
		feed, _ := cmd.Flags().GetInt("feed")
		focal, _ := cmd.Flags().GetFloat64("focal-length")
		totalZtravel := Zmax - Zmin
		xtravel := float64(50)
		// slope := totalZtravel / xtravel
		g := util.Gcode("")

		// xtmp := 25.0
		// z := util.CalculateRise(slope, xtmp)
		// fmt.Println(Zmin + z)
		// out := util.G90Preamble()

		g.G91Preamble()
		g.Add("G0 Z%.2f ; Move to focal", focal)
		g.Add("G0 Z%.2f ; Move down", Zmin)

		g.Add("M3 S%d ; Set spindle", power)
		g.Add("G01 X%.2f Z%.2f F%d ; Move right and up\n", xtravel, totalZtravel, feed)
		g.Add(`M5; Turn off spindle`)

		g.Add("G0 X-%.2f Y-1 Z-%.2f ; return to start\n", xtravel, totalZtravel)
		g.Add("M3 S%d ; Set spindle", power)
		g.Add("G01 Y2 F%d ; draw tick", feed)
		g.Add("M5 ; Set spindle")
		g.Add("G0 Y-1 ")

		zCur := focal + Zmin
		chart := ""

		for i := 0; i < 10; i++ {

			g.Add("G0 X%.2f Y-1 Z%.2f; Move right\n", xtravel/10.0, totalZtravel/10.0)
			g.Add("M3 S%d ; Set spindle", power)
			g.Add("G01 Y2 F%d ; draw tick", feed)
			g.Add("M5 ; Set spindle")
			g.Add("G0 Y-1 ")
			chart += fmt.Sprintf("; %d: %.2f\n", i+1, zCur)
			zCur += totalZtravel / 10.0

		}

		g.Print()

		fmt.Println(chart)
	},
}

func init() {
	laserCmd.AddCommand(focusCmd)
	focusCmd.Flags().IntP("power", "p", 75, "power, should be minimal")
	focusCmd.Flags().IntP("feed", "f", 100, "feed rate")
	focusCmd.Flags().Float64P("focal-length", "F", 40, "rough focal length")
	focusCmd.Flags().Float64P("z-min", "z", -5, "lowest Z value")
	focusCmd.Flags().Float64P("z-max", "Z", 5, "highest Z value")
}
