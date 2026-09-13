package activity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry represents a single CLI activity (clone, commit, push, edit, etc.)
type Entry struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	RepoURL   string `json:"repoUrl"`
	RepoName  string `json:"repoName"`
	RepoPath  string `json:"repoPath"`
	Action    string `json:"action"` // clone, commit, push, pull, edit, delete, view, status, branch
	Command   string `json:"command"`
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
}

var mu sync.Mutex

func getLogPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".commandspeak-activity.json"), nil
}

// Log appends a new activity entry to the JSON log file.
func Log(entry Entry) error {
	mu.Lock()
	defer mu.Unlock()

	if entry.Timestamp == "" {
		entry.Timestamp = time.Now().Format(time.RFC3339)
	}
	if entry.ID == "" {
		entry.ID = "act-" + time.Now().Format("20060102-150405-") + entry.Action
	}

	logPath, err := getLogPath()
	if err != nil {
		return err
	}

	entries, _ := readEntries(logPath) // ignore error if file doesn't exist yet
	entries = append(entries, entry)

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(logPath, data, 0644)
}

// GetAll returns all activity entries.
func GetAll() ([]Entry, error) {
	logPath, err := getLogPath()
	if err != nil {
		return nil, err
	}
	return readEntries(logPath)
}

// GetByRepo returns activity entries filtered by repo name.
func GetByRepo(repoName string) ([]Entry, error) {
	all, err := GetAll()
	if err != nil {
		return nil, err
	}

	var filtered []Entry
	for _, e := range all {
		if e.RepoName == repoName {
			filtered = append(filtered, e)
		}
	}
	return filtered, nil
}

// GetRepos returns a deduplicated list of repo names from activity.
func GetRepos() ([]map[string]string, error) {
	all, err := GetAll()
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var repos []map[string]string
	for _, e := range all {
		if e.RepoName != "" && !seen[e.RepoName] {
			seen[e.RepoName] = true
			repos = append(repos, map[string]string{
				"name": e.RepoName,
				"url":  e.RepoURL,
				"path": e.RepoPath,
			})
		}
	}
	return repos, nil
}

func readEntries(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return []Entry{}, nil
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return []Entry{}, nil
	}
	return entries, nil
}
