package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ShrewdSpirit/credman/utility/cmdutlitity"

	"github.com/atotto/clipboard"

	"github.com/ShrewdSpirit/credman/data"
	"github.com/ShrewdSpirit/credman/utility"
	"github.com/spf13/cobra"
)

var siteCmd = &cobra.Command{
	Use:     "site",
	Aliases: []string{"s"},
	Short:   "Site management",
	Run:     func(cmd *cobra.Command, args []string) {},
}

var siteAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a", "n", "new"},
	Args:    cobra.ExactArgs(1),
	Short:   "Adds a new site",
	Run: func(cmd *cobra.Command, args []string) {
		siteAdd(args[0])
	},
}

var siteRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "rem", "del", "delete"},
	Args:    cobra.ExactArgs(1),
	Short:   "Removes a site",
	Run: func(cmd *cobra.Command, args []string) {
		siteRemove(args[0])
	},
}

var siteRenameCmd = &cobra.Command{
	Use:     "rename",
	Aliases: []string{"rn", "ren"},
	Args:    cobra.ExactArgs(2),
	Short:   "Renames a site",
	Run: func(cmd *cobra.Command, args []string) {
		siteRename(args[0], args[1])
	},
}

var siteSetCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s", "u", "edit", "e"},
	Args:    cobra.ExactArgs(1),
	Short:   "Updates site fields",
	Run: func(cmd *cobra.Command, args []string) {
		siteSet(args[0])
	},
}

var siteListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Lists sites using optional pattern",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			siteList("")
		} else {
			siteList(args[0])
		}
	},
}

var siteGetCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g", "f"},
	Args:    cobra.ExactArgs(1),
	Short:   "Gets value(s) of specified field(s) or copy the first field into clipboard",
	Run: func(cmd *cobra.Command, args []string) {
		siteGet(args[0])
	},
}

var siteFieldsMap map[string]string
var siteFieldsList []string
var siteFieldsDelete []string
var siteGetCopy bool
var siteSetPassword bool
var siteAddNoPassword bool

func init() {
	rootCmd.AddCommand(siteCmd)
	cmdutility.FlagsAddProfileName(siteCmd)

	siteCmd.AddCommand(siteAddCmd)
	siteAddCmd.Flags().BoolVarP(&siteAddNoPassword, "no-password", "n", false, "Doesn't prompt for site password. Useful for sites that you don't want any password for.")
	siteFlagsFields(siteAddCmd, false)
	cmdutility.FlagsAddPasswordOptions(siteAddCmd)

	siteCmd.AddCommand(siteRemoveCmd)
	siteCmd.AddCommand(siteRenameCmd)
	siteCmd.AddCommand(siteSetCmd)
	siteSetCmd.Flags().BoolVarP(&siteSetPassword, "password", "w", false, "Change password. Can be used with password generator or it will prompt user")
	siteSetCmd.Flags().StringSliceVarP(&siteFieldsDelete, "delete", "d", []string{}, "Deletes specified fields")
	siteFlagsFields(siteSetCmd, false)
	cmdutility.FlagsAddPasswordOptions(siteSetCmd)

	siteCmd.AddCommand(siteListCmd)
	siteCmd.AddCommand(siteGetCmd)
	siteFlagsFields(siteGetCmd, true)
	siteGetCmd.Flags().BoolVarP(&siteGetCopy, "copy", "c", false, "Copy first selected field into clipboard")
}

func siteFlagsFields(cmd *cobra.Command, get bool) {
	if get {
		cmd.Flags().StringSliceVarP(&siteFieldsList, "fields", "f", make([]string, 0), "-f=Key1,Key2 ...")
	} else {
		cmd.Flags().StringToStringVarP(&siteFieldsMap, "field", "f", make(map[string]string), "-f=Key=Value ...")
	}
}

func siteAdd(siteName string) {
	profile, profilePassword := cmdutility.GetProfileCommandLine()
	if profile == nil {
		return
	}

	if profile.SiteExist(siteName) {
		utility.LogColor(utility.BoldHiYellow, "Site %s exists.", siteName)
		return
	}

	site := data.NewSite(siteName)
	var password string
	if !siteAddNoPassword {
		var err error
		password, err = cmdutility.ParsePasswordGenerationFlags("Site new password")
		if err != nil {
			utility.LogError("Site creation failed", err)
			return
		}
	}

	site.Fields["password"] = password
	for field, value := range siteFieldsMap {
		if field == "password" {
			continue
		}
		site.Fields[field] = value
	}
	profile.AddSite(site)
	if err := profile.Save(profilePassword); err != nil {
		utility.LogError("Failed saving profile", err)
		return
	}

	utility.LogColor(utility.Green, "Site %s has been added to profile %s", siteName, profile.Name)
}

func siteRemove(siteName string) {
	profile, profilePassword := cmdutility.GetProfileCommandLine()
	if profile == nil {
		return
	}

	if !profile.SiteExist(siteName) {
		utility.LogColor(utility.BoldHiYellow, "Site %s doesn't exist.", siteName)
		return
	}

	remove, err := utility.YesNoPrompt(fmt.Sprintf("Are you sure to delete site %s?", siteName))
	if err != nil {
		utility.LogError("Reading input failed", err)
		return
	}

	if !remove {
		return
	}

	profile.DeleteSite(siteName)

	if err := profile.Save(profilePassword); err != nil {
		utility.LogError("Failed saving profile", err)
		return
	}

	utility.LogColor(utility.Green, "Site %s has been removed.", siteName)
}

func siteRename(siteName string, newName string) {
	profile, profilePassword := cmdutility.GetProfileCommandLine()
	if profile == nil {
		return
	}

	if !profile.SiteExist(siteName) {
		utility.LogColor(utility.BoldHiYellow, "Site %s doesn't exist.", siteName)
		return
	}

	site := profile.GetSite(siteName)
	site.Name = newName

	if err := profile.Save(profilePassword); err != nil {
		utility.LogError("Failed saving profile", err)
		return
	}

	utility.LogColor(utility.Green, "Site %s has been renamed to %s.", siteName, newName)
}

func siteSet(siteName string) {
	profile, profilePassword := cmdutility.GetProfileCommandLine()
	if profile == nil {
		return
	}

	if !profile.SiteExist(siteName) {
		utility.LogColor(utility.BoldHiYellow, "Site %s doesn't exist.", siteName)
		return
	}

	site := profile.GetSite(siteName)

	if siteSetPassword {
		password, err := cmdutility.ParsePasswordGenerationFlags("Site new password")
		if err != nil {
			utility.LogError("Site creation failed", err)
			return
		}
		site.Fields["password"] = password
	}

	updatedFields := make([]string, len(siteFieldsMap))
	for field, value := range siteFieldsMap {
		if field == "password" {
			continue
		}
		site.Fields[field] = value
		updatedFields = append(updatedFields, field)
	}

	deletedFields := make([]string, len(siteFieldsDelete))
	for _, field := range siteFieldsDelete {
		delete(site.Fields, field)
		deletedFields = append(deletedFields, field)
	}

	if err := profile.Save(profilePassword); err != nil {
		utility.LogError("Failed saving profile", err)
		return
	}

	utility.LogColor(utility.Green, "Site %s has been updated.", siteName)
}

func siteList(pattern string) {
	profile, _ := cmdutility.GetProfileCommandLine()
	if profile == nil {
		return
	}

	if len(pattern) == 0 {
		for _, site := range profile.Sites {
			fmt.Println(site.Name)
		}
	} else {
		rx, err := regexp.Compile(pattern)
		if err != nil {
			utility.LogError("Failed to compile pattern", err)
			return
		}

		for _, site := range profile.Sites {
			if rx.MatchString(site.Name) {
				idx := rx.FindStringIndex(site.Name)
				part1 := site.Name[:idx[0]]
				part2 := site.Name[idx[0]:idx[1]]
				part3 := site.Name[idx[1]:]
				utility.LogColor(utility.BoldRed, part1+"%s"+part3, part2)
			}
		}
	}
}

func siteGet(siteName string) {
	profile, _ := cmdutility.GetProfileCommandLine()
	if profile == nil {
		return
	}

	if !profile.SiteExist(siteName) {
		utility.LogColor(utility.BoldHiYellow, "Site %s doesn't exist.", siteName)
		return
	}

	site := profile.GetSite(siteName)
	if siteGetCopy {
		if len(siteFieldsList) == 0 {
			if err := clipboard.WriteAll(site.Fields["password"]); err != nil {
				utility.LogError("Failed write to clipboard", err)
				return
			}
			fmt.Println("Password copied to clipboard.")
		} else {
			field := siteFieldsList[0]
			_, ok := site.Fields[field]
			if !ok {
				utility.LogColor(utility.BoldRed, "No value for field %s is set.", field)
				return
			}
			if err := clipboard.WriteAll(site.Fields[field]); err != nil {
				utility.LogError("Failed write to clipboard", err)
				return
			}
			utility.LogColor(utility.Green, "%s copied to clipboard.", field)
		}
	} else {
		if len(siteFieldsList) == 0 {
			utility.LogColor(utility.Green, "Fields of site %s:", siteName)
			for field, value := range site.Fields {
				utility.LogColor(utility.HiGreen, strings.Title(field)+": %s", value)
			}
		} else {
			for _, field := range siteFieldsList {
				value, ok := site.Fields[field]
				if !ok {
					continue
				}
				utility.LogColor(utility.HiGreen, strings.Title(field)+": %s", value)
			}
		}
	}
}
