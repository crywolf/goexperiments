package sub_process_tests

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func Crasher() {
	fmt.Println("Going down in flames")
	os.Exit(1)
}

func TestCrasher(t *testing.T) {
	if os.Getenv("CALL_CRASHER") == "1" {
		Crasher()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCrasher")
	cmd.Env = append(os.Environ(), "CALL_CRASHER=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		t.Logf("exited as expected with %v", e)
		return
	}
	t.Fatalf("program exited with error: %v, want: exit status 1", err)
}
