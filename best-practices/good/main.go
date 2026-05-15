// Package main mirrors ../bad/main.go but every issue is fixed.
// Numbers in comments match the planted issues in the bad version so the
// two files can be diffed side-by-side.
package main

import (
	"context"
	"crypto/rand"   // (3, 12) crypto/rand for security-sensitive randomness
	"crypto/sha256" // (2, 11) strong hash
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// (1) Read the token from the environment rather than hard-coding it.
func apiToken() string { return os.Getenv("API_TOKEN") }

// (5) No shadowing of `len`. Initialisation pulled out of init() (4).

// fetch uses a context, checks the error, and always closes the body.
func fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) // (6)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	if resp == nil { // nilaway: stdlib has no nilability annotations
		return nil, errors.New("nil response")
	}
	defer func() { _ = resp.Body.Close() }() // (7)
	return io.ReadAll(resp.Body)
}

// query uses a parameterised statement, propagates errors, closes rows.
func query(ctx context.Context, db *sql.DB, userID string) ([]string, error) {
	rows, err := db.QueryContext(ctx,
		"SELECT name FROM users WHERE id = ?", userID) // (8) parameterised
	if err != nil {
		return nil, fmt.Errorf("query: %w", err) // (9) propagate, don't swallow
	}
	defer func() { _ = rows.Close() }() // (10)

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil { // (10)
		return nil, fmt.Errorf("rows: %w", err)
	}
	return names, nil
}

// hashSecret uses SHA-256 instead of MD5.
func hashSecret(s string) []byte {
	h := sha256.Sum256([]byte(s)) // (11)
	return h[:]
}

// pickToken uses crypto/rand for unpredictable values.
func pickToken() (uint64, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil { // (12)
		return 0, fmt.Errorf("rand: %w", err)
	}
	return binary.BigEndian.Uint64(b[:]), nil
}

// writeReport uses restrictive permissions.
func writeReport(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o600) // (13)
}

// process is split into small helpers — each has low cyclomatic complexity.
func process(items []string, mode int) string {
	switch mode {
	case 1:
		return processMode1(items)
	case 2:
		return processMode2(items)
	default:
		return processDefault(items)
	}
}

func processMode1(items []string) string {
	var b strings.Builder
	for i, item := range items {
		if item == "" {
			continue
		}
		b.WriteString(pickFragment(i, item))
	}
	return b.String()
}

func pickFragment(i int, item string) string {
	if i%2 != 0 {
		return "x"
	}
	if len(item) > 3 {
		return item
	}
	return "?"
}

func processMode2(items []string) string {
	var b strings.Builder
	for _, item := range items {
		b.WriteString(item)
	}
	return b.String()
}

func processDefault(items []string) string {
	var b strings.Builder
	for range items {
		b.WriteString("-")
	}
	return b.String()
}

// user is exported and guarded against nil.
type user struct{ Name string }

func describeUser(u *user) string {
	if u == nil { // (16) nilaway happy
		return ""
	}
	return u.Name
}

// compareErr uses errors.Is — supports wrapped errors.
func compareErr(err error) bool {
	return errors.Is(err, io.EOF) // (17)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := fetch(ctx, "https://example.com"); err != nil {
		fmt.Fprintln(os.Stderr, "fetch:", err)
	}
	if _, err := query(ctx, nil, "1"); err != nil {
		fmt.Fprintln(os.Stderr, "query:", err)
	}
	_ = hashSecret("hunter2")
	if _, err := pickToken(); err != nil {
		fmt.Fprintln(os.Stderr, "token:", err)
	}
	if err := writeReport("/tmp/out", "hi"); err != nil {
		fmt.Fprintln(os.Stderr, "write:", err)
	}
	_ = process([]string{"a", "bb", "ccc"}, 1)
	_ = describeUser(nil)
	_ = compareErr(io.EOF)
	_ = apiToken()
	fmt.Println("done")
}
