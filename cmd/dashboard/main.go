package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "Dashboard",
		Short: "Dashboard Showcase using webx",
	}

	root.AddCommand(serveCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
