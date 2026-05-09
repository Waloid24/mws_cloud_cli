package cmd

import (
	"fmt"
	"strings"

	profilepkg "github.com/Waloid24/mws_cloud_cli/internal/profile"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newProfileCommand() *cobra.Command {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage local profiles",
		Long:  "Manage local YAML profiles stored in the current directory.",
	}

	profileCmd.AddCommand(
		newProfileCreateCommand(),
		newProfileGetCommand(),
		newProfileListCommand(),
		newProfileDeleteCommand(),
	)

	return profileCmd
}

func newProfileCreateCommand() *cobra.Command {
	var name string
	var user string
	var project string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a profile",
		Long:  "Create a local YAML profile with user and project fields.",
		Example: `mws_cloud_cli profile create --name=test --user=example --project=new-project
go run . profile create --name=test --user=example --project=new-project`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireFlag("name", name); err != nil {
				return err
			}
			if err := requireFlag("user", user); err != nil {
				return err
			}
			if err := requireFlag("project", project); err != nil {
				return err
			}

			p := profilepkg.Profile{
				User:    user,
				Project: project,
			}
			if err := profilepkg.Create(name, p); err != nil {
				return err
			}

			cmd.Printf("profile %q created\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "profile name")
	cmd.Flags().StringVar(&user, "user", "", "profile user")
	cmd.Flags().StringVar(&project, "project", "", "profile project")

	return cmd
}

func newProfileGetCommand() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a profile",
		Long:  "Read and print a local YAML profile by name.",
		Example: `mws_cloud_cli profile get --name=test
go run . profile get --name=test`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireFlag("name", name); err != nil {
				return err
			}

			p, err := profilepkg.Get(name)
			if err != nil {
				return err
			}

			data, err := yaml.Marshal(p)
			if err != nil {
				return err
			}

			cmd.Print(string(data))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "profile name")

	return cmd
}

func newProfileListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List profiles",
		Long:  "List local YAML profiles stored in the current directory.",
		Example: `mws_cloud_cli profile list
go run . profile list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := profilepkg.List()
			if err != nil {
				return err
			}

			for _, name := range names {
				cmd.Println(name)
			}

			return nil
		},
	}
}

func newProfileDeleteCommand() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a profile",
		Long:  "Delete a local YAML profile by name.",
		Example: `mws_cloud_cli profile delete --name=test
go run . profile delete --name=test`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireFlag("name", name); err != nil {
				return err
			}

			if err := profilepkg.Delete(name); err != nil {
				return err
			}

			cmd.Printf("profile %q deleted\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "profile name")

	return cmd
}

func requireFlag(name, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("required flag %q not set", name)
	}

	return nil
}
