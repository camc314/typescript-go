package vfs

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// MatchFilesNew provides a non-regex implementation of matchFiles.
// It follows the same behavior but uses string operations instead of regular expressions.
func MatchFilesNew(path string, extensions []string, excludes []string, includes []string, useCaseSensitiveFileNames bool, currentDirectory string, depth *int, host FS) []string {
	path = tspath.NormalizePath(path)
	currentDirectory = tspath.NormalizePath(currentDirectory)
	absolutePath := tspath.CombinePaths(currentDirectory, path)

	// For no includes, we should include everything under the path
	if len(includes) == 0 {
		includes = []string{"**/*"}
	}

	// Process includes - adding implicit "/**/*" where needed
	processedIncludes := processIncludes(includes)

	// Get base paths to start the search from
	basePaths := getBasePaths(path, includes, useCaseSensitiveFileNames)

	// Associate an array of results with each include pattern
	// If there are no includes, then just put everything in results[0]
	var results [][]string
	if len(processedIncludes) > 0 {
		results = make([][]string, len(processedIncludes))
		for i := range processedIncludes {
			results[i] = []string{}
		}
	} else {
		results = [][]string{{}}
	}

	// Create visitor for traversing file system
	v := &nonRegexVisitor{
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
		host:                      host,
		includePatterns:           processedIncludes,
		excludePatterns:           excludes,
		extensions:                extensions,
		results:                   results,
		basePath:                  absolutePath,
		depth:                     depth,
	}

	// Visit each base path
	for _, basePath := range basePaths {
		v.visitDirectory(basePath, tspath.CombinePaths(currentDirectory, basePath), depth)
	}

	return core.Flatten(results)
}

// processIncludes expands implicit globs in include patterns
func processIncludes(includes []string) []string {
	if len(includes) == 0 {
		return nil
	}

	result := make([]string, 0, len(includes))
	for _, include := range includes {
		// Check if this is an implicit glob (directory path)
		parts := strings.Split(include, "/")
		lastComponent := ""
		if len(parts) > 0 {
			lastComponent = parts[len(parts)-1]
		}

		if IsImplicitGlob(lastComponent) {
			// For directory includes, add "/**/*" to capture all files
			result = append(result, include+"/**/*")
		} else {
			result = append(result, include)
		}
	}
	return result
}

// getSearchBasePaths determines the base paths from which to start the search
func getSearchBasePaths(path string, includes []string, useCaseSensitiveFileNames bool) []string {
	// If no includes, search from the path itself
	if len(includes) == 0 {
		return []string{path}
	}

	// Otherwise, use the same base paths as the original implementation
	return getBasePaths(path, includes, useCaseSensitiveFileNames)
}

// Below is a skeleton of what the implementation would look like
// This is not complete but shows the general approach

type nonRegexVisitor struct {
	useCaseSensitiveFileNames bool
	host                      FS
	includePatterns           []string
	excludePatterns           []string
	extensions                []string
	results                   [][]string
	visited                   core.Set[string]
	basePath                  string
	depth                     *int
}

func (v *nonRegexVisitor) visitDirectory(path string, absolutePath string, depth *int) {
	// Check if already visited
	canonicalPath := tspath.GetCanonicalFileName(absolutePath, v.useCaseSensitiveFileNames)
	if v.visited.Has(canonicalPath) {
		return
	}
	v.visited.Add(canonicalPath)

	// Get directory entries
	entries := v.host.GetAccessibleEntries(absolutePath)

	// Process files
	for _, file := range entries.Files {
		filePath := tspath.CombinePaths(path, file)
		absoluteFilePath := tspath.CombinePaths(absolutePath, file)

		// Skip files not matching extensions
		if len(v.extensions) > 0 && !tspath.FileExtensionIsOneOf(filePath, v.extensions) {
			continue
		}

		// Skip excluded files
		if v.isExcluded(absoluteFilePath) {
			continue
		}

		// Add file to results if it matches include patterns
		if len(v.includePatterns) == 0 {
			// No include patterns means include everything
			v.results[0] = append(v.results[0], filePath)
		} else {
			// Check each include pattern
			for i, pattern := range v.includePatterns {
				// Normalize the pattern and path for matching
				absolutePattern := tspath.CombinePaths(v.basePath, pattern)
				if v.matchesGlobPattern(absolutePattern, absoluteFilePath) {
					v.results[i] = append(v.results[i], filePath)
					break
				}
			}
		}
	}

	// Process directories (unless depth limit reached)
	if depth != nil {
		newDepth := *depth - 1
		if newDepth == 0 {
			return
		}
		depth = &newDepth
	}

	for _, dir := range entries.Directories {
		dirPath := tspath.CombinePaths(path, dir)
		absoluteDirPath := tspath.CombinePaths(absolutePath, dir)

		// Skip excluded directories
		if v.isExcluded(absoluteDirPath) {
			continue
		}

		// Only visit directories that could match our include patterns
		if v.couldContainMatches(absoluteDirPath) {
			v.visitDirectory(dirPath, absoluteDirPath, depth)
		}
	}
}

// isExcluded checks if a path matches any exclude pattern
func (v *nonRegexVisitor) isExcluded(path string) bool {
	if len(v.excludePatterns) == 0 {
		return false
	}

	for _, pattern := range v.excludePatterns {
		// For excludes, we need to normalize the pattern and path for matching
		absolutePattern := tspath.CombinePaths(v.basePath, pattern)

		// Special handling for exclude patterns: they match if the path starts with the pattern
		if v.matchesGlobPattern(absolutePattern, path) {
			return true
		}
	}
	return false
}

// couldContainMatches checks if a directory could potentially contain files that match our include patterns
func (v *nonRegexVisitor) couldContainMatches(dirPath string) bool {
	// No include patterns means include everything
	if len(v.includePatterns) == 0 {
		return true
	}

	// Handle common package folders and hidden directories
	dirName := tspath.GetBaseFileName(dirPath)
	isCommonPackage := dirName == "node_modules" || dirName == "bower_components" || dirName == "jspm_packages"
	isHidden := strings.HasPrefix(dirName, ".")

	if isCommonPackage || isHidden {
		// For these special directories, we need explicit include patterns that mention them
		hasExplicitInclude := false
		for _, pattern := range v.includePatterns {
			if strings.Contains(pattern, dirName) {
				hasExplicitInclude = true
				break
			}
		}
		if !hasExplicitInclude {
			return false
		}
	}

	// For each include pattern, check if the directory could match
	for _, pattern := range v.includePatterns {
		// Normalize the pattern for matching
		absolutePattern := tspath.CombinePaths(v.basePath, pattern)
		if v.directoryCouldMatchPattern(absolutePattern, dirPath) {
			return true
		}
	}
	return false
}

// directoryCouldMatchPattern checks if a directory path could match a glob pattern
// This is a performance optimization to avoid traversing directories that can't possibly match
func (v *nonRegexVisitor) directoryCouldMatchPattern(pattern, dirPath string) bool {
	// Common package folders should be skipped unless explicitly included
	dirName := tspath.GetBaseFileName(dirPath)

	// Skip common package folders (node_modules, etc) and hidden directories unless explicitly referenced
	if strings.HasPrefix(dirName, ".") || dirName == "node_modules" || dirName == "bower_components" || dirName == "jspm_packages" {
		// Only include if the pattern explicitly references this directory
		return strings.Contains(pattern, dirName)
	}

	// If pattern has ** wildcard, it could match anything underneath
	if strings.Contains(pattern, "**") {
		return true
	}

	// Split pattern and path into segments
	patternSegments := strings.Split(pattern, "/")
	pathSegments := strings.Split(dirPath, "/")

	// If pattern is shorter than path, it could only match if it ends with a wildcard
	if len(patternSegments) < len(pathSegments) {
		return strings.Contains(patternSegments[len(patternSegments)-1], "*")
	}

	// Check if each segment of the directory path matches the corresponding pattern segment
	for i := 0; i < len(pathSegments); i++ {
		if i >= len(patternSegments) {
			return false
		}
		if !v.matchesSegment(patternSegments[i], pathSegments[i]) {
			return false
		}
	}
	return true
}

// matchesGlobPattern checks if a path matches a glob pattern
func (v *nonRegexVisitor) matchesGlobPattern(pattern, path string) bool {
	// Normalize paths
	pattern = tspath.NormalizePath(pattern)
	path = tspath.NormalizePath(path)

	// Convert TypeScript-style glob patterns to segments for matching
	patternSegments := strings.Split(pattern, "/")
	pathSegments := strings.Split(path, "/")

	// Add terminal $ marker to pattern (like in the regex implementation)
	if !strings.HasSuffix(pattern, "$") {
		patternSegments = append(patternSegments, "$")
	}

	// Handle special cases like **
	return v.matchesSegments(patternSegments, pathSegments, 0, 0)
}

// matchesSegments recursively checks if path segments match pattern segments
func (v *nonRegexVisitor) matchesSegments(patternSegments, pathSegments []string, patternIndex, pathIndex int) bool {
	// Base cases
	if patternIndex == len(patternSegments) {
		return pathIndex == len(pathSegments)
	}

	// Terminal $ marker means we must be at the end of the path
	if patternSegments[patternIndex] == "$" {
		return pathIndex == len(pathSegments)
	}

	// Handle ** wildcard (matches zero or more path segments)
	if patternSegments[patternIndex] == "**" {
		// Skip this ** pattern segment and try to match the rest
		if v.matchesSegments(patternSegments, pathSegments, patternIndex+1, pathIndex) {
			return true
		}

		// Try to match the current ** with one path segment and continue matching
		if pathIndex < len(pathSegments) {
			// Skip hidden directories/files for ** wildcard
			dirName := pathSegments[pathIndex]
			if strings.HasPrefix(dirName, ".") {
				return false
			}

			// Skip common package folders for ** wildcard
			if dirName == "node_modules" || dirName == "bower_components" || dirName == "jspm_packages" {
				return false
			}

			return v.matchesSegments(patternSegments, pathSegments, patternIndex, pathIndex+1)
		}
		return false
	}

	// We've run out of path segments but still have pattern segments
	if pathIndex == len(pathSegments) {
		return false
	}

	// Check if current segments match
	if v.matchesSegment(patternSegments[patternIndex], pathSegments[pathIndex]) {
		return v.matchesSegments(patternSegments, pathSegments, patternIndex+1, pathIndex+1)
	}

	return false
}

// matchesSegment checks if a path segment matches a pattern segment
func (v *nonRegexVisitor) matchesSegment(pattern, segment string) bool {
	// Case insensitive comparison if needed
	if !v.useCaseSensitiveFileNames {
		pattern = strings.ToLower(pattern)
		segment = strings.ToLower(segment)
	}

	// Fast path for exact match
	if pattern == segment {
		return true
	}

	// Fast path for * wildcard
	if pattern == "*" {
		return true
	}

	// Contains wildcards?
	if !strings.ContainsAny(pattern, "*?") {
		return false
	}

	// Match with wildcards character by character
	return v.matchWithWildcards(pattern, segment)
}

// matchWithWildcards does character-by-character matching with * and ? wildcards
func (v *nonRegexVisitor) matchWithWildcards(pattern, segment string) bool {
	// Convert to a simple state machine that handles * and ? wildcards
	patternLen := len(pattern)
	segmentLen := len(segment)

	// Handle empty pattern
	if patternLen == 0 {
		return segmentLen == 0
	}

	// Simple dynamic programming approach to match wildcards
	// dp[i][j] represents if pattern[0...i-1] matches segment[0...j-1]
	dp := make([][]bool, patternLen+1)
	for i := range dp {
		dp[i] = make([]bool, segmentLen+1)
	}

	// Empty pattern matches empty segment
	dp[0][0] = true

	// Handle patterns that start with * (can match empty string)
	for i := 1; i <= patternLen; i++ {
		if pattern[i-1] == '*' {
			dp[i][0] = dp[i-1][0]
		}
	}

	// Fill the dp table
	for i := 1; i <= patternLen; i++ {
		for j := 1; j <= segmentLen; j++ {
			if pattern[i-1] == '*' {
				// * can match zero or more characters
				// Either ignore * (dp[i-1][j]) or use it to match current character (dp[i][j-1])
				dp[i][j] = dp[i-1][j] || dp[i][j-1]
			} else if pattern[i-1] == '?' || pattern[i-1] == segment[j-1] {
				// ? matches any single character, or exact character match
				dp[i][j] = dp[i-1][j-1]
			}
		}
	}

	return dp[patternLen][segmentLen]
}
