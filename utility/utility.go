package utility

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/ShrewdSpirit/credman/data"
)

func Kbmbgb(value int64) string {
	if value/1000000000 > 0 {
		return fmt.Sprintf("%dgb", value/1000000000)
	} else if value/1000000 > 0 {
		return fmt.Sprintf("%dmb", value/1000000)
	} else if value/1000 > 0 {
		return fmt.Sprintf("%dkb", value/1000)
	} else {
		return fmt.Sprintf("%db", value)
	}
}

func ForkSelf(args ...string) error {
	execFile, err := os.Executable()
	if err != nil {
		return err
	}

	execDir := filepath.Dir(execFile)
	pidFile := path.Join(data.DataDir, ".pid")

	if _, err := os.Stat(pidFile); err == nil {
		pidBytes, err := ioutil.ReadFile(pidFile)
		if err != nil {
			return err
		}

		pid64, err := strconv.ParseInt(string(pidBytes), 10, 32)
		if err != nil {
			return err
		}

		proc, err := os.FindProcess(int(pid64))
		if err != nil {
			return err
		}

		proc.Kill()
		RemovePidFile()
	}

	argv := []string{
		execFile,
	}
	argv = append(argv, args...)

	proc, err := os.StartProcess(execFile, argv, &os.ProcAttr{
		Dir: execDir,
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	})

	if err != nil {
		return err
	}

	ioutil.WriteFile(pidFile, []byte(strconv.FormatInt(int64(proc.Pid), 10)), os.ModePerm)

	proc.Release()

	return nil
}

func RemovePidFile() {
	pidFile := path.Join(data.DataDir, ".pid")

	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		return
	}

	os.Remove(pidFile)
}

func RunClearClipboardDelayed() {
	err := ForkSelf("clsclip")

	if err != nil {
		fmt.Printf("Failed to fork self: %s\n", err)
		return
	}

	fmt.Printf("Clipboard will be cleared after 10 seconds\n")
}
