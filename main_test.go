package main

import (
	"bytes"
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
	err = exec.Command("go", "build", "-o", tmpBuildName).Run()
	if err != nil {
		t.Fatal(err)
	}

	// get free port
	l, err := net.ListenTCP("tcp", &net.TCPAddr{})
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	portString := strconv.Itoa(port)
	baseURL := "http://localhost:" + portString

	cmd := exec.Command(tmpBuildName, "-a", ":"+portString, "-g", "test/gallery")

	// read std{out,err} for debugging (logged in cleanup)
	output := &bytes.Buffer{}
	cmd.Stdout = output
	cmd.Stderr = output

	// keep port reserved until ready to bind
	l.Close()

	// start server
	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		t.Log(output)
		cmd.Process.Kill() // more reliable than a context
		os.Remove(tmpBuildName)
	})

	// wait for HTTP server to start
	<-time.After(10 * time.Millisecond) // more reliable that waiting for start message

	// Simple tests for identifiers protects from human error.
	// This is not a replacement for UI tests but that is far too complex for this simple project.
	tests := []struct {
		path       string
		identifier string
	}{
		{"/", "<!-- homem98y2r8 -->"},
		{"/gallery/image", "<!-- gallery72yr98mj -->"},
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

	t.Run("favicon", func(t *testing.T) {
		t.Parallel()

		response, err := http.Get(baseURL + "/favicon.ico")
		if err != nil {
			t.Fatal(err)
		}
		want := 404
		got := response.StatusCode
		if got != want {
			t.Fatalf("Want: %d, Got: %d", want, got)
		}
	})

	t.Run("gallery", func(t *testing.T) {
		t.Parallel()

		c := http.Client{
			// do not follow redirects
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		// make request
		response, err := c.Get(baseURL + "/gallery")
		if err != nil {
			t.Fatal(err)
		}

		// read and log body
		t.Cleanup(func() {
			if t.Failed() {
				body := make([]byte, 50)
				response.Body.Read(body)
				t.Log(string(body))
			}
		})

		t.Run("statusCode", func(t *testing.T) {
			want := 307
			got := response.StatusCode
			if got != want {
				t.Fatalf("Want: %d, Got: %d", want, got)
			}
		})

		t.Run("locationHeader", func(t *testing.T) {
			want := "/gallery/image"
			got := response.Header.Get("Location")
			if got != want {
				t.Fatalf("Want: %s, Got: %s", want, got)
			}
		})
	})

	t.Run("galleryImage/image", func(t *testing.T) {
		t.Parallel()

		response, err := http.Get(baseURL + "/galleryImage/image")
		if err != nil {
			t.Fatal(err)
		}
		want := 200
		got := response.StatusCode
		if got != want {
			t.Fatalf("Want: %d, Got: %d", want, got)
		}
	})

	t.Run("galleryImage/doesNotExist", func(t *testing.T) {
		// This test is necessary because the mock fs used in the unit tests
		// (testfs.MapFS) behaves differently to the real fs used in main() (os.DirFS())
		// Specifically, when a file doesn't exist:
		// fstest.MapFS.Open() returns fs.ErrNotExist
		// os.DirFS().Open() returns syscall.ENOENT
		// This is not documented
		t.Parallel()

		response, err := http.Get(baseURL + "/galleryImage/doesNotExist")
		if err != nil {
			t.Fatal(err)
		}
		want := 404
		got := response.StatusCode
		if got != want {
			t.Fatalf("Want: %d, Got: %d", want, got)
		}
	})
}
