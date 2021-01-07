package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ageoffron/arubacentral/centralrest"

	"github.com/spf13/cobra"
)

var apsCmd = &cobra.Command{
	Use:   "aps",
	Short: "get list of aps",
	Long:  ``,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		readTokenConfig()
		output := centralrest.Getaps(centralrest.TokenStruct{AccessToken: ViperCentralAuth.GetString("access_token")}, Verbose)
		e, err := json.Marshal(output)
		if err != nil {
			panic(err)
		}
		fmt.Printf(string(e))
	},
}
