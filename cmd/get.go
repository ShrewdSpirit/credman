package cmd

import (
	"fmt"

	"github.com/ShrewdSpirit/credman/config"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var GetShort bool

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
		getNotes, _ := cmd.Flags().GetBool("notes")
		getSecq1, _ := cmd.Flags().GetBool("secq1")
		getSecq2, _ := cmd.Flags().GetBool("secq2")
		getSecq3, _ := cmd.Flags().GetBool("secq3")
		getSecq4, _ := cmd.Flags().GetBool("secq4")
		getSecq5, _ := cmd.Flags().GetBool("secq5")

		if copy {
			if getEmail {
				if err := doCopy(siteName, "Email", site.Email); err != nil {
					fmt.Println(err)
					return
				}
			} else if getUsername {
				if err := doCopy(siteName, "Username", site.Username); err != nil {
					fmt.Println(err)
					return
				}
			} else if getPassword {
				if err := doCopy(siteName, "Password", site.Password); err != nil {
					fmt.Println(err)
					return
				}
			} else if getNotes {
				if err := doCopy(siteName, "Notes", site.Notes); err != nil {
					fmt.Println(err)
					return
				}
			} else if getSecq1 {
				if err := doCopy(siteName, "Security question 1", site.SecurityQuestions[0]); err != nil {
					fmt.Println(err)
					return
				}
			} else if getSecq2 {
				if err := doCopy(siteName, "Security question 2", site.SecurityQuestions[1]); err != nil {
					fmt.Println(err)
					return
				}
			} else if getSecq3 {
				if err := doCopy(siteName, "Security question 3", site.SecurityQuestions[2]); err != nil {
					fmt.Println(err)
					return
				}
			} else if getSecq4 {
				if err := doCopy(siteName, "Security question 4", site.SecurityQuestions[3]); err != nil {
					fmt.Println(err)
					return
				}
			} else if getSecq5 {
				if err := doCopy(siteName, "Security question 5", site.SecurityQuestions[4]); err != nil {
					fmt.Println(err)
					return
				}
			} else {
				fmt.Println("No field is selected")
			}
		} else {
			if !GetShort {
				fmt.Printf("Fields of site '%s':\n", siteName)
			}
			printedOneField := false

			if getEmail {
				doPrint("Email", site.Email)
				printedOneField = true
			}
			if getUsername {
				doPrint("Username", site.Username)
				printedOneField = true
			}
			if getPassword {
				doPrint("Password", site.Password)
				printedOneField = true
			}
			if getNotes {
				doPrint("Notes", site.Notes)
				printedOneField = true
			}
			if getSecq1 {
				doPrint("Security question 1", site.SecurityQuestions[0])
				printedOneField = true
			}
			if getSecq2 {
				doPrint("Security question 2", site.SecurityQuestions[1])
				printedOneField = true
			}
			if getSecq3 {
				doPrint("Security question 3", site.SecurityQuestions[2])
				printedOneField = true
			}
			if getSecq4 {
				doPrint("Security question 4", site.SecurityQuestions[3])
				printedOneField = true
			}
			if getSecq5 {
				doPrint("Security question 5", site.SecurityQuestions[4])
				printedOneField = true
			}

			if !printedOneField {
				doPrint("Email", site.Email)
				doPrint("Username", site.Username)
				doPrint("Password", site.Password)
				doPrint("Notes", site.Notes)
				doPrint("Security question 1", site.SecurityQuestions[0])
				doPrint("Security question 2", site.SecurityQuestions[1])
				doPrint("Security question 3", site.SecurityQuestions[2])
				doPrint("Security question 4", site.SecurityQuestions[3])
				doPrint("Security question 5", site.SecurityQuestions[4])
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	AddFlagsProfileName(getCmd)
	AddFlagsSiteFields(getCmd, true)

	getCmd.Flags().BoolP("copy", "c", false, "Copy the single field get value to clipboard")
	getCmd.Flags().BoolVarP(&GetShort, "short", "s", false, "Omits all extra strings for output and only prints the field values")
}

func doCopy(sitename, fieldname, value string) error {
	if err := clipboard.WriteAll(value); err != nil {
		return err
	}

	if !GetShort {
		fmt.Printf("%s of '%s' has been copied to clipboard\n", fieldname, sitename)
	}

	return nil
}

func doPrint(fieldname, value string) {
	if len(value) == 0 {
		return
	}

	if GetShort {
		fmt.Println(value)
	} else {
		fmt.Printf("%s: %s\n", fieldname, value)
	}
}
