package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var (
	profileName    string
	profileUser    string
	profileProject string
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage local profiles",
	Long:  "Manage local YAML profiles stored in the current directory.",
}

var profileCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a profile",
	Long:  "Create a local YAML profile with user and project fields.",
	Example: `mws profile create --name=test --user=example --project=new-project
go run . profile create --name=test --user=example --project=new-project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented")
	},
}

var profileGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a profile",
	Long:  "Read and print a local YAML profile by name.",
	Example: `mws profile get --name=test
go run . profile get --name=test`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented")
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List profiles",
	Long:  "List local YAML profiles stored in the current directory.",
	Example: `mws profile list
go run . profile list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented")
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a profile",
	Long:  "Delete a local YAML profile by name.",
	Example: `mws profile delete --name=test
go run . profile delete --name=test`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented")
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileCreateCmd, profileGetCmd, profileListCmd, profileDeleteCmd)

	profileCreateCmd.Flags().StringVar(&profileName, "name", "", "profile name")
	profileCreateCmd.Flags().StringVar(&profileUser, "user", "", "profile user")
	profileCreateCmd.Flags().StringVar(&profileProject, "project", "", "profile project")

	profileGetCmd.Flags().StringVar(&profileName, "name", "", "profile name")
	profileDeleteCmd.Flags().StringVar(&profileName, "name", "", "profile name")
}
