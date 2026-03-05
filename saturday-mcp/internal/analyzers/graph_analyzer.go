package analyzers

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/orieken/saturday-mcp/internal/logging"
)

// NodeType represents the type of a node in the dependency graph
type NodeType string

const (
	NodeTypePage    NodeType = "page"
	NodeTypeFlow    NodeType = "flow"
	NodeTypeStep    NodeType = "step"
	NodeTypeElement NodeType = "element"
	NodeTypeUnknown NodeType = "unknown"
)

// Node represents a file or component in the graph
type Node struct {
	ID       string   // Unique identifier (usually relative file path)
	Name     string   // Component name (e.g., "LoginPage")
	Type     NodeType // Type of component
	Path     string   // Absolute file path
	Outgoing []string // IDs of nodes this node depends on
	Incoming []string // IDs of nodes that depend on this node
}

// DependencyGraph represents the web of dependencies
type DependencyGraph struct {
	Nodes map[string]*Node
}

// GraphAnalyzer builds and analyzes the dependency graph
type GraphAnalyzer struct {
	logger *logging.Logger
}

// NewGraphAnalyzer creates a new graph analyzer
func NewGraphAnalyzer(logger *logging.Logger) *GraphAnalyzer {
	return &GraphAnalyzer{
		logger: logger,
	}
}

// Build creates a dependency graph from the project root
func (a *GraphAnalyzer) Build(rootPath string) (*DependencyGraph, error) {
	graph := &DependencyGraph{
		Nodes: make(map[string]*Node),
	}

	a.logger.Info("Building dependency graph", "root", rootPath)

	// 1. Identify all nodes (files)
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == "node_modules" || strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".ts") && !strings.HasSuffix(path, ".js") {
			return nil
		}

		relPath, _ := filepath.Rel(rootPath, path)
		nodeID := relPath

		node := &Node{
			ID:       nodeID,
			Path:     path,
			Outgoing: []string{},
			Incoming: []string{},
			Type:     determineNodeType(path),
			Name:     determineNodeName(path),
		}
		graph.Nodes[nodeID] = node
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk project: %w", err)
	}

	// 2. Parse files to find dependencies (Edges)
	for id, node := range graph.Nodes {
		deps, err := a.parseDependencies(node.Path, rootPath)
		if err != nil {
			a.logger.Error("Failed to parse file dependencies", "file", node.Path, "error", err)
			continue
		}

		for _, depPath := range deps {
			// Find the node for this dependency
			// depPath is relative to project root
			if targetNode, exists := graph.Nodes[depPath]; exists {
				// Add Edge
				node.Outgoing = append(node.Outgoing, targetNode.ID)
				targetNode.Incoming = append(targetNode.Incoming, id)
			}
		}
	}

	a.logger.Info("Graph build complete", "nodes", len(graph.Nodes))
	return graph, nil
}

// AnalyzeImpact returns a list of files affected by changes to the given file
func (a *GraphAnalyzer) AnalyzeImpact(graph *DependencyGraph, changedFileRelPath string) ([]string, error) {
	startNode, exists := graph.Nodes[changedFileRelPath]
	if !exists {
		// Try to find by fuzzy match if exact match fails
		for id, node := range graph.Nodes {
			if strings.HasSuffix(id, changedFileRelPath) {
				startNode = node
				break
			}
		}
		if startNode == nil {
			return nil, fmt.Errorf("node not found: %s", changedFileRelPath)
		}
	}

	visited := make(map[string]bool)
	var impacted []string

	// BFS to find all upstream dependencies (Incoming edges)
	queue := []string{startNode.ID}
	visited[startNode.ID] = true

	for len(queue) > 0 {
		currentID := queue[0]
		queue = queue[1:]

		// If it's not the start node, add it to impacted list
		if currentID != startNode.ID {
			impacted = append(impacted, currentID)
		}

		node := graph.Nodes[currentID]
		for _, incomingID := range node.Incoming {
			if !visited[incomingID] {
				visited[incomingID] = true
				queue = append(queue, incomingID)
			}
		}
	}

	return impacted, nil
}

// Helpers

func determineNodeType(path string) NodeType {
	lower := strings.ToLower(path)
	if strings.Contains(lower, "/pages/") || strings.HasSuffix(lower, "page.ts") {
		return NodeTypePage
	}
	if strings.Contains(lower, "/flows/") || strings.HasSuffix(lower, "flow.ts") {
		return NodeTypeFlow
	}
	if strings.Contains(lower, "/steps/") || strings.Contains(lower, ".features") { // checking for step definitions usually in steps folder
		return NodeTypeStep
	}
	if strings.Contains(lower, "/elements/") {
		return NodeTypeElement
	}
	return NodeTypeUnknown
}

func determineNodeName(path string) string {
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext)
}

// parseDependencies scans a file for import statements and resolves them to file paths relative to root
func (a *GraphAnalyzer) parseDependencies(filePath string, rootPath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var deps []string
	scanner := bufio.NewScanner(file)
	
	// Regex for detecting imports: import ... from 'release/path';
	importRegex := regexp.MustCompile(`import\s+.*from\s+['"](.+)['"]`)

	fileDir := filepath.Dir(filePath)

	for scanner.Scan() {
		line := scanner.Text()
		matches := importRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			importPath := matches[1]

			// Resolve path
			var resolvedPath string
			if strings.HasPrefix(importPath, ".") {
				// Relative import
				absPath := filepath.Join(fileDir, importPath)
				
				// Handle extensions if missing
				if !strings.HasSuffix(absPath, ".ts") && !strings.HasSuffix(absPath, ".js") {
					if _, err := os.Stat(absPath + ".ts"); err == nil {
						absPath += ".ts"
					} else if _, err := os.Stat(absPath + ".js"); err == nil {
						absPath += ".js"
					} else if _, err := os.Stat(absPath + "/index.ts"); err == nil {
						absPath += "/index.ts"
					}
				}
				
				rel, _ := filepath.Rel(rootPath, absPath)
				resolvedPath = rel
			} else {
				// TODO: Handle alias imports if configured in tsconfig (skipping for MVP)
				continue 
			}

			deps = append(deps, resolvedPath)
		}
	}

	return deps, scanner.Err()
}
