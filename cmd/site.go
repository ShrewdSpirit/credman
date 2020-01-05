package commands

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/ShrewdSpirit/credman/cmd/cmdutility"
	"github.com/ShrewdSpirit/credman/data"
	"github.com/ShrewdSpirit/credman/utils"
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
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
	Run:     siteAdd,
}

var siteRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "rem", "del", "delete"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Removes a site",
	Run:     siteRemove,
}

var siteRenameCmd = &cobra.Command{
	Use:     "rename",
	Aliases: []string{"rn", "ren"},
	Args:    cobra.ExactArgs(2),
	Short:   "Renames a site",
	Run:     siteRename,
}

var siteSetCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s", "u", "edit", "e"},
	Args:    cobra.ExactArgs(1),
	Short:   "Updates site fields",
	Run:     siteSet,
}

var siteListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Lists sites using optional pattern",
	Run:     siteList,
}

var siteGetCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g", "f"},
	Args:    cobra.ExactArgs(1),
	Short:   "Gets value(s) of specified field(s) or copy the first field into clipboard",
	Run:     siteGet,
}

var (
	siteFieldsMap       map[string]string
	siteFieldsList      []string
	siteFieldsDelete    []string
	siteGetCopy         bool
	siteSetPassword     bool
	siteAddNoPassword   bool
	siteTags            []string
	siteDeleteTags      []string
	siteGroup           bool
	siteGetShowPassword bool
)

func init() {
	rootCmd.AddCommand(siteCmd)
	cmdutility.FlagsAddProfileName(siteCmd)

	siteCmd.AddCommand(siteAddCmd)
	siteAddCmd.Flags().BoolVarP(&siteAddNoPassword, "no-password", "n", false, "Doesn't prompt for site password. Useful for sites that you don't want any password for.")
	siteFlagsFields(siteAddCmd, false)
	siteFlagsTags(siteAddCmd)
	cmdutility.FlagsAddPasswordOptions(siteAddCmd)

	siteCmd.AddCommand(siteRemoveCmd)
	siteFlagsTags(siteRemoveCmd)

	siteCmd.AddCommand(siteRenameCmd)

	siteCmd.AddCommand(siteSetCmd)
	siteSetCmd.Flags().BoolVarP(&siteSetPassword, "password", "w", false, "Change password. Can be used with password generator or it will prompt user")
	siteSetCmd.Flags().StringSliceVarP(&siteFieldsDelete, "delete", "d", []string{}, "Deletes specified fields")
	siteSetCmd.Flags().StringSliceVar(&siteDeleteTags, "delete-tags", nil, "Delte tags")
	siteFlagsTags(siteSetCmd)
	siteFlagsFields(siteSetCmd, false)
	cmdutility.FlagsAddPasswordOptions(siteSetCmd)

	siteCmd.AddCommand(siteListCmd)
	siteFlagsTags(siteListCmd)
	siteListCmd.Flags().BoolVarP(&siteGroup, "group", "g", false, "Groups sites by tags")

	siteCmd.AddCommand(siteGetCmd)
	siteFlagsFields(siteGetCmd, true)
	siteGetCmd.Flags().BoolVarP(&siteGetCopy, "copy", "c", false, "Copy first selected field into clipboard. If no field is specified, will copy password")
	siteGetCmd.Flags().BoolVarP(&siteGetShowPassword, "show-password", "s", false, "By default credman won't show password unless this flag is set")
}

func siteFlagsTags(cmd *cobra.Command) {
	cmd.Flags().StringSliceVarP(&siteTags, "tags", "t", nil, "Site tags")
}

func siteFlagsFields(cmd *cobra.Command, get bool) {
	if get {
		cmd.Flags().StringSliceVarP(&siteFieldsList, "fields", "f", make([]string, 0), "-f=Key1,Key2 ...")
	} else {
		cmd.Flags().StringToStringVarP(&siteFieldsMap, "field", "f", make(map[string]string), "-f=Key=Value ...")
	}
}

func siteAdd(cmd *cobra.Command, args []string) {
	siteName := args[0]
	profile, profilePassword := cmdutility.GetProfileCommandLine(true)
	if profile == nil {
		return
	}

	var password string
	if !siteAddNoPassword {
		var err error
		password, err = cmdutility.ParsePasswordGenerationFlags("Site new password")
		if err != nil {
			cmdutility.LogError("Site creation failed", err)
			return
		}
	}

	if profile.SiteExist(siteName) {
		cmdutility.LogColor(cmdutility.BoldHiYellow, "Site %s exists.", siteName)
		return
	}

	site := data.NewSite(siteName, password, siteFieldsMap, siteTags)

	profile.AddSite(siteName, site)

	if err := profile.Save(profilePassword); err != nil {
		cmdutility.LogError("Failed saving profile", err)
		return
	}

	cmdutility.LogColor(cmdutility.Green, "Site %s has been added to profile %s", siteName, profile.Name)
}

func siteGet(cmd *cobra.Command, args []string) {
	siteName := args[0]
	profile, _ := cmdutility.GetProfileCommandLine(true)
	if profile == nil {
		return
	}

	if !profile.SiteExist(siteName) {
		cmdutility.LogColor(cmdutility.BoldHiYellow, "Site %s doesn't exist.", siteName)
		return
	}

	site := profile.GetSite(siteName)
	fields := site.GetFields(siteFieldsList)
	tags := site.GetTags()

	if siteGetCopy {
		if len(siteFieldsList) == 0 {
			if err := clipboard.WriteAll(site["password"]); err != nil {
				cmdutility.LogError("Failed write to clipboard", err)
				return
			}

			fmt.Println("Password copied to clipboard.")
			utils.RunClearClipboardDelayed()
			return
		} else {
			field := siteFieldsList[0]
			_, ok := site[field]
			if !ok || data.IsSpecialField(field) {
				cmdutility.LogColor(cmdutility.BoldRed, "No field %s was found.", siteFieldsList[0])
				return
			}

			if err := clipboard.WriteAll(site[field]); err != nil {
				cmdutility.LogError("Failed write to clipboard", err)
				return
			}

			cmdutility.LogColor(cmdutility.Green, "%s copied to clipboard.", siteFieldsList[0])

			if field == "password" {
				utils.RunClearClipboardDelayed()
			}
			return
		}
	}

	cmdutility.LogColor(cmdutility.Green, "Site %s fields:", siteName)
	tw := tabwriter.NewWriter(os.Stdout, 10, 0, 1, ' ', 0)
	for field, value := range fields {
		if field == "password" && !siteGetShowPassword {
			value = "*****"
		}

		if runtime.GOOS != "windows" {
			value = color.HiGreenString(value)
		}
		fmt.Fprintf(tw, " %s:\t%s\n", strings.Title(field), value)
	}

	if tags != nil && len(tags) > 0 {
		tagsString := "#" + strings.Replace(strings.Join(tags, " "), " ", " #", -1)
		if runtime.GOOS != "windows" {
			tagsString = color.CyanString(tagsString)
		}
		fmt.Fprintf(tw, " Tags:\t%s\n", tagsString)
	}

	tw.Flush()
}

func siteList(cmd *cobra.Command, args []string) {
	pattern := ""
	if len(args) == 1 {
		pattern = args[0]
	}

	profile, _ := cmdutility.GetProfileCommandLine(true)
	if profile == nil {
		return
	}

	sites, err := profile.GetSites(pattern, siteTags)
	if err != nil {
		cmdutility.LogError("Failed to list sites", err)
		return
	}

	// TODO: implement grouping
	tw := tabwriter.NewWriter(os.Stdout, 10, 0, 1, ' ', 0)
	for _, site := range sites {
		tagsString := "#" + strings.Replace(strings.Join(site.Tags, " "), " ", " #", -1)
		if site.Tags == nil || len(site.Tags) == 0 {
			tagsString = ""
		}

		if runtime.GOOS != "windows" {
			tagsString = color.CyanString(tagsString)
		}

		if site.Name == site.MatchParts[0] {
			fmt.Fprintf(tw, "%s\t%s\n", site.Name, tagsString)
		} else {
			matchPart := color.HiRedString(site.MatchParts[1])
			fmt.Fprintf(tw, "%s%s%s\t%s\n", site.MatchParts[0], matchPart, site.MatchParts[2], tagsString)
		}
	}

	tw.Flush()
}

func siteRemove(cmd *cobra.Command, args []string) {
	profile, profilePassword := cmdutility.GetProfileCommandLine(true)
	if profile == nil {
		return
	}

	if len(siteTags) != 0 {
		sitesToRemove := make([]string, 0)
		for sn, s := range profile.Sites {
			hasTag, _ := s.HasTags(siteTags)
			if hasTag {
				sitesToRemove = append(sitesToRemove, sn)
			}
		}

		for _, sn := range sitesToRemove {
			helperSiteRemove(profile, profilePassword, sn)
		}

		return
	}

	if len(args) == 0 {
		cmdutility.LogColor(cmdutility.BoldHiYellow, "No input site")
		return
	}

	siteName := args[0]

	if !profile.SiteExist(siteName) {
		cmdutility.LogColor(cmdutility.BoldHiYellow, "Site %s doesn't exist.", siteName)
		return
	}

	helperSiteRemove(profile, profilePassword, siteName)
}

func helperSiteRemove(profile *data.Profile, profilePassword, siteName string) {
	remove, err := cmdutility.YesNoPrompt(fmt.Sprintf("Are you sure to delete site %s?", siteName))
	if err != nil {
		cmdutility.LogError("Reading input failed", err)
		return
	}

	if !remove {
		return
	}

	profile.DeleteSite(siteName)
	if err := profile.Save(profilePassword); err != nil {
		cmdutility.LogError("Failed saving profile", err)
		return
	}

	cmdutility.LogColor(cmdutility.Green, "Site %s has been removed.", siteName)
}

func siteRename(cmd *cobra.Command, args []string) {
	siteName := args[0]
	newName := args[1]

	profile, profilePassword := cmdutility.GetProfileCommandLine(true)
	if profile == nil {
		return
	}

	if !profile.SiteExist(siteName) {
		cmdutility.LogColor(cmdutility.BoldHiYellow, "Site %s doesn't exist.", siteName)
		return
	}

	if profile.GetSite(siteName).IsFile() {
		cmdutility.LogColor(cmdutility.BoldHiYellow, "Site %s is a file.", siteName)
		return
	}

	profile.RenameSite(siteName, newName)
	if err := profile.Save(profilePassword); err != nil {
		cmdutility.LogError("Failed saving profile", err)
		return
	}

	cmdutility.LogColor(cmdutility.Green, "Site %s has been renamed to %s.", siteName, newName)
}

func siteSet(cmd *cobra.Command, args []string) {
	siteName := args[0]
	profile, profilePassword := cmdutility.GetProfileCommandLine(true)
	if profile == nil {
		return
	}

	var password string
	if siteSetPassword {
		var err error
		password, err = cmdutility.ParsePasswordGenerationFlags("Site new password")
		if err != nil {
			cmdutility.LogError("Site creation failed", err)
			return
		}
	}

	if !profile.SiteExist(siteName) {
		cmdutility.LogColor(cmdutility.BoldHiYellow, "Site %s doesn't exist.", siteName)
		return
	}

	site := profile.GetSite(siteName)

	if site.IsFile() {
		cmdutility.LogColor(cmdutility.BoldHiYellow, "Site %s is a file.", siteName)
		return
	}

	if len(password) != 0 {
		site["password"] = password
	}

	for field, value := range siteFieldsMap {
		if data.IsSpecialField(field) {
			continue
		}

		site[strings.ToLower(field)] = value
	}

	for _, field := range siteFieldsDelete {
		delete(site, strings.ToLower(field))
	}

	if siteTags != nil && len(siteTags) > 0 {
		site.AddTags(siteTags)
	}

	if siteDeleteTags != nil && len(siteDeleteTags) > 0 {
		site.RemoveTags(siteDeleteTags)
	}

	if err := profile.Save(profilePassword); err != nil {
		cmdutility.LogError("Failed saving profile", err)
		return
	}

	cmdutility.LogColor(cmdutility.Green, "Site %s has been updated.", siteName)
}
