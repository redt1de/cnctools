/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// stepsCmd represents the steps command
var stepsCmd = &cobra.Command{
	Use:   "steps",
	Short: "calibrate steps per millimeter",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		current, _ := cmd.Flags().GetFloat64("current")
		target, _ := cmd.Flags().GetFloat64("target")
		actual, _ := cmd.Flags().GetFloat64("actual")
		steps := current / actual * target
		fmt.Printf("New steps per mm: %.3f\n", steps)
	},
}

func init() {
	rootCmd.AddCommand(stepsCmd)
	stepsCmd.Flags().Float64P("current", "c", 0.0, "current steps")
	stepsCmd.Flags().Float64P("target", "t", 25.0, "target distance")
	stepsCmd.Flags().Float64P("actual", "a", 0.0, "actual distance traveled")
}
