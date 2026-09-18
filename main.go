// tshx switches between Teleport (tsh) cluster profiles, kubectx-style.
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
)

const (
	currentProfileFile = "current-profile"
	previousMarkerFile = "tshx-previous"
	profileSuffix      = ".yaml"
)

// version is set at build time via -ldflags "-X main.version=...".
var version = "dev"

func tshHomeDir() (string, error) {
	if dir := os.Getenv("TSH_HOME"); dir != "" {
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, ".tsh"), nil
}

// listProfiles returns the sorted list of proxy domains that have a
// profile (<domain>.yaml) in the tsh home directory.
func listProfiles(tshHome string) ([]string, error) {
	entries, err := os.ReadDir(tshHome)
	if err != nil {
		return nil, fmt.Errorf("could not read tsh directory %s: %w", tshHome, err)
	}

	var profiles []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, profileSuffix) {
			continue
		}
		profiles = append(profiles, strings.TrimSuffix(name, profileSuffix))
	}
	sort.Strings(profiles)
	return profiles, nil
}

func readCurrentProfile(tshHome string) (string, error) {
	data, err := os.ReadFile(filepath.Join(tshHome, currentProfileFile))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("could not read current profile: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}

func writeCurrentProfile(tshHome, profile string) error {
	path := filepath.Join(tshHome, currentProfileFile)
	return os.WriteFile(path, []byte(profile+"\n"), 0o640)
}

func readPrevious(tshHome string) (string, error) {
	data, err := os.ReadFile(filepath.Join(tshHome, previousMarkerFile))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func writePrevious(tshHome, profile string) error {
	if profile == "" {
		return nil
	}
	path := filepath.Join(tshHome, previousMarkerFile)
	return os.WriteFile(path, []byte(profile+"\n"), 0o640)
}

// switchTo sets profile as the current profile, recording the prior
// current profile so `tshx -` can swap back to it.
func switchTo(tshHome, profile, previousCurrent string) error {
	if profile != previousCurrent {
		if err := writePrevious(tshHome, previousCurrent); err != nil {
			return fmt.Errorf("could not record previous profile: %w", err)
		}
	}
	if err := writeCurrentProfile(tshHome, profile); err != nil {
		return fmt.Errorf("could not switch profile: %w", err)
	}
	fmt.Printf("Switched to Teleport cluster %q.\n", profile)
	return nil
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func runInteractive(tshHome string, profiles []string, current string) error {
	items := make([]string, len(profiles))
	for i, p := range profiles {
		tag := checkSession(tshHome, p).tag()
		if p == current {
			items[i] = fmt.Sprintf("%s (current) %s", p, tag)
		} else {
			items[i] = fmt.Sprintf("%s %s", p, tag)
		}
	}

	startPos := 0
	for i, p := range profiles {
		if p == current {
			startPos = i
			break
		}
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F449 {{ . | cyan }}",
		Inactive: "  {{ . }}",
		Selected: "✔ {{ . | green }}",
	}

	prompt := promptui.Select{
		Label:     "Select a Teleport cluster",
		Items:     items,
		Templates: templates,
		CursorPos: startPos,
		HideHelp:  true,
		Size:      10,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) || errors.Is(err, promptui.ErrEOF) {
			return nil
		}
		return fmt.Errorf("selection failed: %w", err)
	}

	selected := profiles[idx]
	if selected == current {
		fmt.Printf("Already on Teleport cluster %q.\n", selected)
	} else if err := switchTo(tshHome, selected, current); err != nil {
		return err
	}
	offerLogin(tshHome, selected)
	return nil
}

func run(args []string) error {
	if len(args) > 0 && (args[0] == "--version" || args[0] == "-v") {
		fmt.Println("tshx", version)
		return nil
	}

	tshHome, err := tshHomeDir()
	if err != nil {
		return err
	}

	profiles, err := listProfiles(tshHome)
	if err != nil {
		return err
	}
	if len(profiles) == 0 {
		return fmt.Errorf("no Teleport profiles found in %s (run `tsh login` first)", tshHome)
	}

	current, err := readCurrentProfile(tshHome)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return runInteractive(tshHome, profiles, current)
	}

	target := args[0]

	if target == "-" {
		previous, err := readPrevious(tshHome)
		if err != nil {
			return err
		}
		if previous == "" {
			return errors.New("no previous Teleport cluster to switch to")
		}
		if !contains(profiles, previous) {
			return fmt.Errorf("previous cluster %q no longer has a profile", previous)
		}
		if err := switchTo(tshHome, previous, current); err != nil {
			return err
		}
		offerLogin(tshHome, previous)
		return nil
	}

	if !contains(profiles, target) {
		return fmt.Errorf("no Teleport profile named %q found in %s", target, tshHome)
	}
	if target == current {
		fmt.Printf("Already on Teleport cluster %q.\n", target)
		offerLogin(tshHome, target)
		return nil
	}
	if err := switchTo(tshHome, target, current); err != nil {
		return err
	}
	offerLogin(tshHome, target)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "tshx:", err)
		os.Exit(1)
	}
}
