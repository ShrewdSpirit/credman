package cmd

import (
	"fmt"

	"github.com/ShrewdSpirit/credman/config"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
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
		if err != nil {
			fmt.Println(err)
			return
		}

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
				if err := SiteGetCopy(siteName, "Email", site.Email); err != nil {
					fmt.Println(err)
					return
				}
			} else if getUsername {
				if err := SiteGetCopy(siteName, "Username", site.Username); err != nil {
					fmt.Println(err)
					return
				}
			} else if getPassword {
				if err := SiteGetCopy(siteName, "Password", site.Password); err != nil {
					fmt.Println(err)
					return
				}
			} else if getNotes {
				if err := SiteGetCopy(siteName, "Notes", site.Notes); err != nil {
					fmt.Println(err)
					return
				}
			} else if getSecq1 {
				if err := SiteGetCopy(siteName, "Security question 1", site.SecurityQuestions[0]); err != nil {
					fmt.Println(err)
					return
				}
			} else if getSecq2 {
				if err := SiteGetCopy(siteName, "Security question 2", site.SecurityQuestions[1]); err != nil {
					fmt.Println(err)
					return
				}
			} else if getSecq3 {
				if err := SiteGetCopy(siteName, "Security question 3", site.SecurityQuestions[2]); err != nil {
					fmt.Println(err)
					return
				}
			} else if getSecq4 {
				if err := SiteGetCopy(siteName, "Security question 4", site.SecurityQuestions[3]); err != nil {
					fmt.Println(err)
					return
				}
			} else if getSecq5 {
				if err := SiteGetCopy(siteName, "Security question 5", site.SecurityQuestions[4]); err != nil {
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
				fmt.Printf("Email: %s\n", site.Email)
				printedOneField = true
			}
			if getUsername {
				fmt.Printf("Username: %s\n", site.Username)
				printedOneField = true
			}
			if getPassword {
				fmt.Printf("Password: %s\n", site.Password)
				printedOneField = true
			}
			if getNotes {
				fmt.Printf("Notes: %s\n", site.Notes)
				printedOneField = true
			}
			if getSecq1 {
				fmt.Printf("Security question 1: %s\n", site.SecurityQuestions[0])
				printedOneField = true
			}
			if getSecq2 {
				fmt.Printf("Security question 2: %s\n", site.SecurityQuestions[1])
				printedOneField = true
			}
			if getSecq3 {
				fmt.Printf("Security question 3: %s\n", site.SecurityQuestions[2])
				printedOneField = true
			}
			if getSecq4 {
				fmt.Printf("Security question 4: %s\n", site.SecurityQuestions[3])
				printedOneField = true
			}
			if getSecq5 {
				fmt.Printf("Security question 5: %s\n", site.SecurityQuestions[4])
				printedOneField = true
			}

			if !printedOneField {
				fmt.Printf("Email: %s\n", site.Email)
				fmt.Printf("Username: %s\n", site.Username)
				fmt.Printf("Password: %s\n", site.Password)
				fmt.Printf("Notes: %s\n", site.Notes)
				fmt.Printf("Security question 1: %s\n", site.SecurityQuestions[0])
				fmt.Printf("Security question 2: %s\n", site.SecurityQuestions[1])
				fmt.Printf("Security question 3: %s\n", site.SecurityQuestions[2])
				fmt.Printf("Security question 4: %s\n", site.SecurityQuestions[3])
				fmt.Printf("Security question 5: %s\n", site.SecurityQuestions[4])
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	ProfileAddFlags(getCmd)
	SiteAddFieldFlags(getCmd, true)

	getCmd.Flags().BoolP("copy", "c", false, "Copy the single field get value to clipboard")
}
