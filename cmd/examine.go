package cmd

import (
	"context"
	"encoding/json"
	"io"
	"justcompile/licenses/pkg/examine"
	"justcompile/licenses/pkg/parser"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var examineCmd = &cobra.Command{
	Use:  "examine [SOURCE]",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workingDir := "."

		source := args[0]
		var r io.Reader
		if source == "-" {
			r = os.Stdin
		} else {
			workingDir = path.Dir(source)

			f, err := os.Open(source)
			if err != nil {
				log.Fatalf("error reading file %s: %s", source, err.Error())
			}

			r = f
		}

		ctx := context.WithValue(context.Background(), parser.ContextCurrentWorkingDir, workingDir)

		pkgMgr, _ := cmd.Flags().GetString("lang")

		examiner, err := examine.New(pkgMgr)
		if err != nil {
			log.Fatal(err)
		}

		packages, err := examiner.Process(ctx, r)
		if err != nil {
			log.Fatal(err)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(packages); err != nil {
			log.Fatal(err)
		}

	}}

func init() {
	rootCmd.AddCommand(examineCmd)
	examineCmd.Flags().String("lang", "py", "denotes the package manager for parsing")
}
