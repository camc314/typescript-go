package vfs_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

// setupTestFS creates a test file system with a specific structure for testing glob patterns
func setupTestFS(useCaseSensitiveFileNames bool) vfs.FS {
	return vfstest.FromMap(map[string]any{
		"/src/foo.ts":                   "export const foo = 1;",
		"/src/bar.ts":                   "export const bar = 2;",
		"/src/baz.tsx":                  "export const baz = 3;",
		"/src/subfolder/qux.ts":         "export const qux = 4;",
		"/src/subfolder/quux.tsx":       "export const quux = 5;",
		"/src/node_modules/lib.ts":      "export const lib = 6;",
		"/src/.hidden/secret.ts":        "export const secret = 7;",
		"/src/test.min.js":              "console.log('minified');",
		"/dist/output.js":               "console.log('output');",
		"/build/temp.ts":                "export const temp = 8;",
		"/test/test1.spec.ts":           "describe('test1', () => {});",
		"/test/test2.spec.tsx":          "describe('test2', () => {});",
		"/test/subfolder/test3.spec.ts": "describe('test3', () => {});",
	}, useCaseSensitiveFileNames)
}

func TestMatchFilesVsMatchFilesNew(t *testing.T) {
	fs := setupTestFS(true)

	testCases := []struct {
		name     string
		path     string
		exts     []string
		excludes []string
		includes []string
		depth    *int
	}{
		{
			name:     "Simple include",
			path:     "/src",
			exts:     []string{".ts", ".tsx"},
			excludes: nil,
			includes: []string{"**/*.ts"},
			depth:    nil,
		},
		{
			name:     "Multiple extensions",
			path:     "/",
			exts:     []string{".ts", ".tsx"},
			excludes: nil,
			includes: []string{"**/*.ts", "**/*.tsx"},
			depth:    nil,
		},
		{
			name:     "With excludes",
			path:     "/",
			exts:     []string{".ts", ".tsx"},
			excludes: []string{"**/node_modules/**", "**/.hidden/**"},
			includes: []string{"**/*.ts"},
			depth:    nil,
		},
		{
			name:     "With depth",
			path:     "/",
			exts:     []string{".ts", ".tsx"},
			excludes: nil,
			includes: []string{"**/*.ts"},
			depth:    func() *int { d := 1; return &d }(),
		},
		{
			name:     "Implicit glob",
			path:     "/",
			exts:     []string{".ts", ".tsx"},
			excludes: nil,
			includes: []string{"src"},
			depth:    nil,
		},
		{
			name:     "Complex pattern",
			path:     "/",
			exts:     []string{".ts", ".tsx"},
			excludes: []string{"**/dist/**", "**/node_modules/**"},
			includes: []string{"src/**/*.ts", "test/**/*.spec.ts"},
			depth:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expected := vfs.ReadDirectory(fs, "/", tc.path, tc.exts, tc.excludes, tc.includes, tc.depth)

			// Directly call matchFilesNew with the same parameters
			actual := vfs.MatchFilesNew(tc.path, tc.exts, tc.excludes, tc.includes, fs.UseCaseSensitiveFileNames(), "/", tc.depth, fs)

			// Sort both slices for consistent comparison
			assert.DeepEqual(t, actual, expected)
		})
	}
}
