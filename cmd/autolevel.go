/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"log"

	"github.com/redt1de/cnctools/autolevel"
	"github.com/spf13/cobra"
)

// autolevelCmd represents the autolevel command
var autolevelCmd = &cobra.Command{
	Use:   "autolevel",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		bounds, err := autolevel.ParseGcodeBoundaries(file)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Gcode boundaries:")
		fmt.Printf("\tXmin: %.3f, XMax: %.3f\n", bounds.MinX, bounds.MaxX)
		fmt.Printf("\tYmin: %.3f, YMax: %.3f\n", bounds.MinY, bounds.MaxY)
		fmt.Printf("\tZmin: %.3f, ZMax: %.3f\n", bounds.MinZ, bounds.MaxZ)
		hm := autolevel.ProbeGrid(&autolevel.FakeSerialPort{}, bounds.MinX, bounds.MaxX, bounds.MinY, bounds.MaxY, 5, 4)

		fmt.Println(hm.GoCode())

		autolevel.ApplyHeightMap(file, hm)

	},
}

func init() {
	rootCmd.AddCommand(autolevelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// autolevelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	autolevelCmd.Flags().StringP("file", "f", "", "Gcode file to analyze")
}
