package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ageoffron/arubacentral/centralrest"

	"github.com/spf13/cobra"
)

var swarmsCmd = &cobra.Command{
	Use:   "swarms",
	Short: "get list of swarms",
	Long:  ``,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		readTokenConfig()
		output := centralrest.Getswarms(centralrest.TokenStruct{AccessToken: ViperCentralAuth.GetString("access_token")}, Verbose)
		e, err := json.Marshal(output)
		if err != nil {
			panic(err)
		}
		fmt.Printf(string(e))
	},
}
