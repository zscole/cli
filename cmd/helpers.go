package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"bitbucket.org/whiteblockio/back_end/cli/project"
)

func Fatal(v ...interface{}) {
	fmt.Println("Error:", v)
	os.Exit(1)
}

// This function from github.com/spf13/cobra used under Apache 2.0 license
func isEmpty(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		Fatal(err)
	}

	if !fi.IsDir() {
		return fi.Size() == 0
	}

	f, err := os.Open(path)
	if err != nil {
		Fatal(err)
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil && err != io.EOF {
		Fatal(err)
	}

	for _, name := range names {
		if len(name) > 0 && name[0] != '.' {
			return false
		}
	}
	return true
}

func ExecWithOutput(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func ExecWithPipes(command string, in []byte, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	go func() {
		defer stdin.Close()
		stdin.Write(in)
	}()

	return cmd.CombinedOutput()
}

func RunInRoot(f func() error) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	prj, err := project.FindProject()
	if err != nil {
		return err
	}

	if err := os.Chdir(prj.AbsPath()); err != nil {
		return err
	}
	defer os.Chdir(wd)

	return f()
}

func runStub(stubcommand string, stubargs ...string) error {
	return RunInRoot(func() error {
		command := runtime.GOROOT() + "/bin/go"
		args := append([]string{"run", "stub/main.go", stubcommand}, stubargs...)
		return ExecWithOutput(command, args...)
	})
}
