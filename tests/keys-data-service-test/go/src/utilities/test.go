package utilities

import (
	"keys-data-service-test/libraries"
	"keys-data-service-test/settings"
	"net/http"
	"os/exec"
	"testing"
	"time"
)

type Test struct {
	command   *exec.Cmd
	process   int
	directory string
}

func SetupTest(tester *testing.T) *Test {
	tester.Helper()

	testDirectory, issue := ResolveTestDirectory()
	if issue != nil {
		tester.Fatalf("resolve test directory with issue, %v", issue)
	}

	command, issue := ConstructCommand("go", "run", ".")
	if issue != nil {
		tester.Fatalf("construct test command with issue, %v", issue)
	}

	if issue := command.Start(); issue != nil {
		tester.Fatalf("start test command with issue, %v", issue)
	}

	test := &Test{
		command:   command,
		process:   command.Process.Pid,
		directory: testDirectory,
	}

	tester.Cleanup(func() { DropTest(tester, test) })

	CheckTestHealthy(tester, test)

	return test
}

func DropTest(tester *testing.T, test *Test) {
	defer func() {
		if issue := DropTestDirectory(test.directory); issue != nil {
			tester.Errorf("drop test directory %s with issue, %v", test.directory, issue)
		}
	}()

	tester.Helper()

	if test.command.ProcessState != nil {
		return
	}

	if _, issue := SendPost(settings.TargetAddress + "/service/data/keys/stop"); issue != nil {
		if issue := KillProcess(test.process); issue != nil {
			tester.Errorf("drop test on kill process %d with issue, %v", test.process, issue)
		}

		return
	}

	if !CheckCommandExit(test.command) && test.command.ProcessState == nil {
		if issue := KillProcess(test.process); issue != nil {
			tester.Errorf("drop test on kill process %d with issue, %v", test.process, issue)
		}
	}
}

func SendTestPost(tester *testing.T, address string, body ...any) *http.Response {
	tester.Helper()

	var response *http.Response
	var issue error

	if len(body) > 0 && body[0] != nil {
		response, issue = SendPost(address, body[0])
	} else {
		response, issue = SendPost(address)
	}

	if issue != nil {
		tester.Fatalf("send test post at %s with issue, %v", address, issue)
	}

	return response
}

func SendTestGet(tester *testing.T, address string) *http.Response {
	tester.Helper()

	response, issue := SendGet(address)
	if issue != nil {
		tester.Fatalf("send test get at %s with issue, %v", address, issue)
	}

	return response
}

func SendTestDelete(tester *testing.T, address string) *http.Response {
	tester.Helper()

	response, issue := SendDelete(address)
	if issue != nil {
		tester.Fatalf("send test delete at %s with issue, %v", address, issue)
	}

	return response
}

func CheckTestHealthy(tester *testing.T, test *Test) {
	tester.Helper()

	deadline := time.Now().Add(settings.TestTimeout)
	for time.Now().Before(deadline) {
		response, issue := SendGet(settings.TargetAddress + "/service/data/keys/health")
		if issue == nil {
			if response.StatusCode == http.StatusOK {
				response.Body.Close()
				return
			}

			response.Body.Close()
		}

		libraries.Sleep(100)
	}

	tester.Fatal("check test healthy failed: service did not become healthy in time")
}

func CheckTestDone(tester *testing.T, test *Test) bool {
	tester.Helper()

	if CheckCommandExit(test.command) {
		return true
	}

	tester.Fatal("check test done failed")
	return false
}
