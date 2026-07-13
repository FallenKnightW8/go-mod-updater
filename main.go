package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	repoURL    string
	directOnly bool
	jsonOutput bool
)

type ModuleInfo struct {
	ModuleName   string           `json:"module_name"`
	GoVersion    string           `json:"go_version"`
	Dependencies []DependencyInfo `json:"dependencies"`
}

type DependencyInfo struct {
	Name           string `json:"name"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version,omitempty"`
	IsDirect       bool   `json:"is_direct"`
	CanUpdate      bool   `json:"can_update"`
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "go-mod-updater",
		Short: "Analyze Go module dependencies for available updates",
		Long: `A CLI tool that clones a Git repository, analyzes its Go module,
and lists dependencies that have available updates.`,
		RunE: runAnalyzer,
	}

	rootCmd.Flags().StringVarP(&repoURL, "repo", "r", "", "Git repository URL (required)")
	rootCmd.Flags().BoolVarP(&directOnly, "direct", "d", false, "Show only direct dependencies")
	rootCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
	rootCmd.MarkFlagRequired("repo")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runAnalyzer(cmd *cobra.Command, args []string) error {
	fmt.Printf("Analyzing repository: %s\n\n", repoURL)

	// Создаем временную папку
	tmpDir, err := os.MkdirTemp("", "go-mod-updater-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Клонируем репозиторий
	repoDir := filepath.Join(tmpDir, "repo")
	fmt.Println("Cloning repository...")
	
	cloneCmd := exec.Command("git", "clone", "--depth", "1", repoURL, repoDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	
	if err := cloneCmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}
	fmt.Println("Repository cloned successfully\n")

	// Проверяем go.mod
	goModPath := filepath.Join(repoDir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found in repository. Is this a Go project?")
	}

	// Парсим go.mod
	moduleName, goVersion, err := parseGoMod(goModPath)
	if err != nil {
		return fmt.Errorf("failed to parse go.mod: %w", err)
	}

	// Получаем зависимости
	deps, err := getDependencies(repoDir, directOnly)
	if err != nil {
		return fmt.Errorf("failed to get dependencies: %w", err)
	}

	// Формируем результат
	info := &ModuleInfo{
		ModuleName:   moduleName,
		GoVersion:    goVersion,
		Dependencies: deps,
	}

	// Выводим результат
	if jsonOutput {
		return outputJSON(info)
	}
	return outputTable(info)
}

func parseGoMod(path string) (string, string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", "", err
	}

	lines := strings.Split(string(content), "\n")
	var moduleName, goVersion string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			moduleName = strings.TrimSpace(strings.TrimPrefix(line, "module "))
		} else if strings.HasPrefix(line, "go ") {
			goVersion = strings.TrimSpace(strings.TrimPrefix(line, "go "))
		}
	}

	if moduleName == "" {
		return "", "", fmt.Errorf("module name not found in go.mod")
	}

	return moduleName, goVersion, nil
}

func getDependencies(repoDir string, directOnly bool) ([]DependencyInfo, error) {
	// Получаем текущие зависимости
	fmt.Println("Analyzing current dependencies...")
	
	cmd := exec.Command("go", "list", "-m", "-json", "all")
	cmd.Dir = repoDir
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("go list failed: %w", err)
	}

	var currentDeps []DependencyInfo
	decoder := json.NewDecoder(strings.NewReader(string(output)))
	
	for decoder.More() {
		var mod struct {
			Path     string
			Version  string
			Main     bool
			Indirect bool
		}
		
		if err := decoder.Decode(&mod); err != nil {
			continue
		}
		
		if mod.Main {
			continue
		}

		if directOnly && mod.Indirect {
			continue
		}
		
		currentDeps = append(currentDeps, DependencyInfo{
			Name:           mod.Path,
			CurrentVersion: mod.Version,
			IsDirect:       !mod.Indirect,
		})
	}

	// Проверяем обновления
	fmt.Println("Checking for available updates...")
	
	updateCmd := exec.Command("go", "list", "-m", "-u", "-json", "all")
	updateCmd.Dir = repoDir
	
	updateOutput, err := updateCmd.Output()
	if err != nil {
		fmt.Printf("Warning: Could not check for updates: %v\n", err)
		return currentDeps, nil
	}

	updates := make(map[string]string)
	updateDecoder := json.NewDecoder(strings.NewReader(string(updateOutput)))
	
	for updateDecoder.More() {
		var mod struct {
			Path   string
			Update *struct {
				Version string
			}
		}
		
		if err := updateDecoder.Decode(&mod); err != nil {
			continue
		}
		
		if mod.Update != nil {
			updates[mod.Path] = mod.Update.Version
		}
	}

	// Объединяем информацию
	for i, dep := range currentDeps {
		if latest, exists := updates[dep.Name]; exists {
			currentDeps[i].LatestVersion = latest
			currentDeps[i].CanUpdate = true
		}
	}

	fmt.Printf("Found %d dependencies\n\n", len(currentDeps))
	return currentDeps, nil
}

func outputJSON(info *ModuleInfo) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(info)
}

func outputTable(info *ModuleInfo) error {
	fmt.Printf("=======================================\n")
	fmt.Printf("Module: %s\n", info.ModuleName)
	if info.GoVersion != "" {
		fmt.Printf("Go Version: %s\n", info.GoVersion)
	}
	fmt.Printf("=======================================\n\n")

	updatableCount := 0
	for _, dep := range info.Dependencies {
		if dep.CanUpdate {
			updatableCount++
		}
	}

	if len(info.Dependencies) == 0 {
		fmt.Println("No dependencies found.")
		return nil
	}

	fmt.Printf("Total dependencies: %d | Updates available: %d\n\n", 
		len(info.Dependencies), updatableCount)
	
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "DEPENDENCY\tCURRENT\tLATEST\tSTATUS\tTYPE")
	fmt.Fprintln(w, "----------\t-------\t------\t------\t----")

	for _, dep := range info.Dependencies {
		status := "Up to date"
		latest := dep.CurrentVersion
		depType := "indirect"
		
		if dep.CanUpdate {
			status = "Update available"
			latest = dep.LatestVersion
		}
		
		if dep.IsDirect {
			depType = "direct"
		}
		
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", 
			dep.Name, 
			dep.CurrentVersion, 
			latest, 
			status,
			depType)
	}

	fmt.Fprintln(w)
	return w.Flush()
}