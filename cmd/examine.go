package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"justcompile/licenses/pkg/lookup"
	"justcompile/licenses/pkg/parser"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var examineCmd = &cobra.Command{
	Use:  "examine [SOURCE]",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		source := args[0]
		var r io.Reader
		if source == "-" {
			r = os.Stdin
		} else {
			f, err := os.Open(source)
			if err != nil {
				log.Fatalf("error reading file %s: %s", source, err.Error())
			}

			r = f
		}

		packages, err := parser.Parse(r)
		if err != nil {
			log.Fatalf("unable to extract packages: %s", err.Error())
		}

		var wg sync.WaitGroup
		wg.Add(len(packages))

		for _, p := range packages {
			go func(pkg *parser.Package) {
				defer wg.Done()

				license, err := lookup.Search(pkg.Name)
				if err != nil {
					fmt.Println(err)
					return
				}

				pkg.License = license
			}(p)
		}

		wg.Wait()

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(packages)

	}}

func init() {
	rootCmd.AddCommand(examineCmd)
}
