/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/redt1de/cnctools/util"
	"github.com/spf13/cobra"
)

const (
	beamDia = 0.2
	overlap = 0.5
	boxSize = 2
)

// powerCmd represents the power command
var powerCmd = &cobra.Command{
	Use:   "power",
	Short: "generate a test pattern for speed and power",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		boxSize := 5.0
		power := 1000.0
		beamDiameter := 0.2
		overlap := 0.5
		feedrate := 500.0

		gcode := `G21\nG91\n`
		gcode += util.G91SpiralFill(boxSize, power, beamDiameter, overlap, feedrate)
		fmt.Println(gcode)
	},
}

func init() {
	laserCmd.AddCommand(powerCmd)
	powerCmd.Flags().IntP("power-max", "P", 1000, "max value for power")
	powerCmd.Flags().IntP("power-min", "p", 100, "min value for power")
	powerCmd.Flags().IntP("feed-max", "F", 600, "max value for feed")
	powerCmd.Flags().IntP("feed-min", "f", 100, "max value for feed")
	powerCmd.Flags().IntP("iterations", "i", 10, "number of step iterations")
}
