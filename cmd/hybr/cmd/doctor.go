package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rasjonell/hybr/internal/services"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use:   "doctor [path]",
	Short: "Checks service definition validity",
	Long: `Checks service definition validity.
If [path] is provided, it will check the validity of the given service's definition.
Otherwise it will check the validity of services in $HYBR_DIR/services.`,
	Args: cobra.MaximumNArgs(1),
	Run:  hybrDoctor,
}

func hybrDoctor(cmd *cobra.Command, args []string) {
	var result []string
	if len(args) == 1 {
		result = checkServiceDefinitions(filepath.Join(".", args[0]))
	} else {
		paths, _, listErr := services.GetInstallableServicePaths()
		if listErr != nil {
			fmt.Println(listErr)
			os.Exit(1)
		}
		result = checkServiceDefinitions(paths...)
	}

	fmt.Printf(strings.Join(result, "\n") + "\n")
}

func checkServiceDefinitions(paths ...string) []string {
	statuses := make([]string, len(paths))

	for i, path := range paths {
		if _, err := os.Stat(filepath.Join(path, "service.json")); os.IsNotExist(err) {
			statuses[i] = fmt.Sprintf("[%s] 🛑 Invalid\nInvalid Service Directory\n\t - service.json doesn't exist!", filepath.Base(path))
			continue
		}

		if err := services.ValidateServiceJSON(filepath.Join(path, "service.json")); err != nil {
			statuses[i] = fmt.Sprintf("[%s] 🛑 Invalid\n%s", filepath.Base(path), err.Error())
			continue
		}

		statuses[i] = fmt.Sprintf("[%s] ✅ Valid", filepath.Base(path))
	}

	return statuses
}
