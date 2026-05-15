// Package main is INTENTIONALLY BAD. Each numbered comment points to a tool
// that catches the issue. Compare with ../good/main.go.
//
// Run `make demo-bad` from the parent directory to see the tools fire.
package main

import (
	"crypto/md5" // (2) gosec G401: weak crypto
	"database/sql"
	"fmt"
	"io"
	"math/rand" // (3) gosec G404: math/rand for sensitive use
	"net/http"
	"os"
)

// (1) gosec G101: hardcoded credential.
const apiToken = "AKIAIOSFODNN7EXAMPLE"

// (4) gochecknoinits: avoid init() funcs. (5) predeclared: shadows `len`.
func init() {
	len := 10
	_ = len
}

// fetch leaks the response body and skips the error.
func fetch(url string) []byte {
	resp, _ := http.Get(url) // (6) errcheck + noctx
	body, _ := io.ReadAll(resp.Body)
	// (7) bodyclose: resp.Body never Close()d
	return body
}

// query builds SQL by concatenation.
func query(db *sql.DB, userID string) (*sql.Rows, error) {
	q := fmt.Sprintf("SELECT name FROM users WHERE id = '%s'", userID) // (8) gosec G201
	rows, err := db.Query(q)
	if err != nil {
		return nil, nil // (9) nilerr: returning nil error after err != nil
	}
	// (10) rowserrcheck / sqlclosecheck: rows.Err() not called, rows not Closed
	return rows, err
}

// hashSecret uses a broken hash function for a sensitive value.
func hashSecret(s string) []byte {
	h := md5.New() // (11) gosec G401
	_, _ = io.WriteString(h, s)
	return h.Sum(nil)
}

// pickToken seeds math/rand with a constant and uses it for an auth token.
func pickToken() int {
	r := rand.New(rand.NewSource(1)) // (12) gosec G404
	return r.Int()
}

// writeReport opens a file with world-writable perms.
func writeReport(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o777) // (13) gosec G306
}

// process has deep nesting / high cyclomatic complexity.
func process(items []string, mode int) string {
	out := ""
	for i := 0; i < len(items); i++ { // (15) intrange / rangeint
		if mode == 1 {
			if items[i] != "" {
				if i%2 == 0 {
					if len(items[i]) > 3 {
						out += items[i] // (14) gocognit / nestif / stringsbuilder
					} else {
						out += "?"
					}
				} else {
					out += "x"
				}
			}
		} else if mode == 2 {
			out += items[i]
		} else {
			out += "-"
		}
	}
	return out
}

type user struct{ Name string }

// describeUser dereferences a nil pointer when called with nil.
func describeUser(u *user) string {
	return u.Name // (16) nilaway: no nil guard on parameter
}

// compareErr uses == on errors instead of errors.Is.
func compareErr(err error) bool {
	return err == io.EOF // (17) errorlint
}

// unused is never referenced.
func unused() { fmt.Println("dead") } // (18) unusedfunc / unused

// ineffectiveAssign assigns then immediately overwrites.
func ineffectiveAssign() int {
	x := 1 // (19) ineffassign
	x = 2
	return x
}

func main() {
	_ = apiToken
	_, _ = query(nil, "1")
	_ = fetch("http://example.com")
	_ = hashSecret("hunter2")
	_ = pickToken()
	_ = writeReport("/tmp/out", "hi")
	_ = process([]string{"a", "bb", "ccc"}, 1)
	_ = describeUser(nil)
	_ = compareErr(io.EOF)
	_ = ineffectiveAssign()
	fmt.Println("done")
}
