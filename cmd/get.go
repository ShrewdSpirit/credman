package cmd

import (
	"fmt"
	"strings"

	"github.com/ShrewdSpirit/credman/config"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Args:  cobra.ExactArgs(1),
	Short: "Gets specified fields of a site or all fields if none is specified",
	Run: func(cmd *cobra.Command, args []string) {
		siteName := args[0]

		profileName, _ := cmd.Flags().GetString("profile")
		if profileName == "" {
			profileName = config.AppConfig.DefaultProfile
		}
		if profileName == "" {
			fmt.Println("Create a profile first")
			return
		}

		profile, err := config.GetProfile(ProfilesDir, profileName)
		if profile == nil && err != nil {
			fmt.Printf("No such profile '%s'\n", profileName)
			return
		}

		profilePassword, err := PromptPassword("Profile password")
		if err != nil {
			fmt.Println(err)
			return
		}

		if Verbose {
			fmt.Println("Decrypting profile")
		}
		err = profile.Decrypt(profilePassword)
		if err != nil {
			fmt.Println(err)
			return
		}

		if Verbose {
			fmt.Println("Checking site's existence")
		}
		site := profile.Site(siteName)
		if site == nil {
			fmt.Printf("Site '%s' doesn't exist\n", siteName)
			return
		}

		copy, _ := cmd.Flags().GetBool("copy")

		getEmail, _ := cmd.Flags().GetBool("email")
		getUsername, _ := cmd.Flags().GetBool("username")
		getPassword, _ := cmd.Flags().GetBool("password")
		fields, _ := cmd.Flags().GetStringSlice("fields")

		if copy {
			if getEmail {
				if err := doCopy(siteName, "Email", site.Fields["email"]); err != nil {
					fmt.Println(err)
					return
				}
			} else if getUsername {
				if err := doCopy(siteName, "Username", site.Fields["username"]); err != nil {
					fmt.Println(err)
					return
				}
			} else if getPassword {
				if err := doCopy(siteName, "Password", site.Fields["password"]); err != nil {
					fmt.Println(err)
					return
				}
			} else if len(fields) > 0 {
				field := strings.TrimSpace(fields[0])
				value, ok := site.Fields[field]
				field = strings.Title(field)
				if !ok {
					fmt.Printf("Invalid field '%s'\n", field)
					return
				}
				if err := doCopy(siteName, field, value); err != nil {
					fmt.Println(err)
					return
				}
			} else {
				fmt.Println("No field is selected")
			}
		} else {
			fmt.Printf("Fields of site '%s':\n", siteName)
			printedOneField := false

			if getEmail {
				doPrint("Email", site.Fields["email"])
				printedOneField = true
			}
			if getUsername {
				doPrint("Username", site.Fields["username"])
				printedOneField = true
			}
			if getPassword {
				doPrint("Password", site.Fields["password"])
				printedOneField = true
			}

			for _, field := range fields {
				field = strings.TrimSpace(field)
				value, _ := site.Fields[field]
				doPrint(strings.Title(field), value)
				printedOneField = true
			}

			if !printedOneField {
				for field, value := range site.Fields {
					doPrint(strings.Title(field), value)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	AddFlagsProfileName(getCmd)
	AddFlagsSiteFields(getCmd, true)

	getCmd.Flags().BoolP("copy", "c", false, "Copy the single field get value to clipboard")
}

func doCopy(sitename, fieldname, value string) error {
	if len(value) == 0 {
		fmt.Printf("No value for field '%s' of '%s' is set\n", fieldname, sitename)
		return nil
	}

	if err := clipboard.WriteAll(value); err != nil {
		return err
	}

	fmt.Printf("%s of '%s' has been copied to clipboard\n", fieldname, sitename)

	return nil
}

func doPrint(fieldname, value string) {
	if len(value) == 0 {
		return
	}
	fmt.Printf("%s: %s\n", fieldname, value)
}
