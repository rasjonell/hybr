package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rasjonell/hybr/internal/services"
	"github.com/spf13/cobra"
)

var isRemote bool

func init() {
	servicesListCmd.Flags().BoolVarP(
		&isRemote,
		"remote", "r", false,
		"",
	)
	servicesCmd.AddCommand(servicesListCmd)
}

var servicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists hybr services",
	Long: `Lists hybr services.

If the "-r" or "--remote" flag is provided, this will list services in global hybr service registry.
`,
	Run: listServices,
}

func listServices(cmd *cobra.Command, args []string) {
	var err error
	var installPaths []string
	var serviceNames []string

	if isRemote {
		fmt.Println("Fetching from central registry...")
		return
	}

	installPaths, serviceNames, err = services.GetInstallableServicePaths()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reg := services.GetRegistry()
	bSet := make(map[string]struct{})
	for _, v := range reg.ListInstallations() {
		bSet[v.GetName()] = struct{}{}
	}

	common := []string{
		"Already Installed Services:",
	}

	unique := []string{
		"Installable Services:",
	}

	for i, v := range serviceNames {
		if _, exists := bSet[v]; exists {
			common = append(common, fmt.Sprintf(
				"  - %s", v,
			))
		} else {
			unique = append(unique, fmt.Sprintf(
				"  - %s:\n       hybr install %s", v, installPaths[i],
			))
		}
	}

	fmt.Printf(strings.Join(append(common, unique...), "\n") + "\n")
}
