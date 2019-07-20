package commands

import (
	"net"
	"os"

	"github.com/ShrewdSpirit/credman/interfaces/commands/cmdutility"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Connect to ssh",
	Long:  `Connects to ssh using fields specified in given site.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		siteName := args[0]

		profile, _ := cmdutility.GetProfileCommandLine(true)
		if profile == nil {
			return
		}

		site := profile.GetSite(siteName)
		if site == nil {
			cmdutility.LogColor(cmdutility.BoldHiYellow, "Site %s doesn't exist.", siteName)
			return
		}

		siteUser, ok := site["user"]
		if !ok {
			cmdutility.LogColor(cmdutility.BoldHiYellow, "Site %s doesn't have %s field set.", siteName, "user")
			return
		}

		sitePass, ok := site["password"]
		if !ok {
			cmdutility.LogColor(cmdutility.BoldHiYellow, "Site %s doesn't have %s field set.", siteName, "password")
			return
		}

		siteAddr, ok := site["address"]
		if !ok {
			cmdutility.LogColor(cmdutility.BoldHiYellow, "Site %s doesn't have %s field set.", siteName, "address")
			return
		}

		sshshell(siteName, siteUser, sitePass, siteAddr)
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
	cmdutility.FlagsAddProfileName(sshCmd)
}

func sshshell(siteName, user, pw, addr string) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pw),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		addr = addr + ":22"
	}

	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		cmdutility.LogError("Failed to connect", err)
		return
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		cmdutility.LogError("Failed to create a session", err)
		return
	}
	defer session.Close()

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	fileDescriptor := int(os.Stdout.Fd())

	if terminal.IsTerminal(fileDescriptor) {
		originalState, err := terminal.MakeRaw(fileDescriptor)
		if err != nil {
			cmdutility.LogError("Failed to put terminal in raw mode", err)
			return
		}
		defer terminal.Restore(fileDescriptor, originalState)

		termWidth, termHeight, err := terminal.GetSize(fileDescriptor)
		if err != nil {
			cmdutility.LogError("Failed to get terminal size", err)
			return
		}

		err = session.RequestPty("xterm", termHeight, termWidth, modes)
		if err != nil {
			cmdutility.LogError("Failed to request remote pty", err)
			return
		}
	}

	if err := session.Shell(); err != nil {
		cmdutility.LogError("Failed to request remote shell", err)
		return
	}

	session.Wait()
	// if err := session.Wait(); err != nil {
	// 	cmdutility.LogError("SSH exited", err)
	// 	return
	// }
}
