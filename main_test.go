package main

import (
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

func TestBlackBox(t *testing.T) {
	t.Parallel()
	const defaultPort = "8000"

	// build to temp file
	file, err := os.CreateTemp(os.TempDir(), "go-test-build")
	if err != nil {
		t.Fatal(err)
	}
	tmpBuildName := file.Name()
	exec.Command("go", "build", "-o", tmpBuildName).Run()
	defer os.Remove(tmpBuildName)

	// get free port
	l, err := net.ListenTCP("tcp", &net.TCPAddr{})
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	portString := strconv.Itoa(port)
	baseURL := "http://[::1]:" + portString

	// create cancellable command
	cmd := exec.Command(tmpBuildName, "-a", ":"+portString)

	// keep port reserved until ready to bind
	l.Close()

	// start server
	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cmd.Process.Kill() // more reliable than a context
	})

	// wait for HTTP server to start
	<-time.After(10 * time.Millisecond) // more reliable that waiting for start message

	// test various routes
	tests := []struct {
		path       string
		identifier string
	}{
		{"/", "<!-- 98x7y3m49t -->"},
		{"/gallery", "<!-- q98ny7g0sk -->"},
		{"/static/blackBoxTestFile", "This file exists"},
	}
	for _, test := range tests {
		// make a copy of test in this scope because tests do not run immediately or syncronously
		test := test
		t.Run(test.path, func(t *testing.T) {
			t.Parallel()

			// make request
			response, err := http.Get(baseURL + test.path)
			if err != nil {
				t.Fatal(err)
			}

			// check body for identifier
			b := make([]byte, len(test.identifier))
			_, err = response.Body.Read(b)
			if err != nil && err != io.EOF {
				t.Fatal(err)
			}
			response.Body.Close()
			if string(b) != test.identifier {
				t.Fatalf("Identifying characters not present at start of body. Want: %s, Got: %s", test.identifier, b)
			}
		})
	}
}