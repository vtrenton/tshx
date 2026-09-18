// Session status: checks whether a profile's Teleport certificate has
// expired, and offers to run `tsh login` when it has.
package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
)

// sessionStatus describes the local TLS certificate for a profile.
// known is false when no certificate could be found or parsed, which
// covers both "never logged in" and any unexpected key store state.
type sessionStatus struct {
	known     bool
	expired   bool
	expiresAt time.Time
}

func (s sessionStatus) needsLogin() bool {
	return !s.known || s.expired
}

// tag renders a short, colored annotation describing the session status,
// suitable for appending to a profile's display label.
func (s sessionStatus) tag() string {
	switch {
	case !s.known:
		return promptui.Styler(promptui.FGFaint)("(no session)")
	case s.expired:
		return promptui.Styler(promptui.FGRed)("(expired)")
	default:
		return promptui.Styler(promptui.FGGreen)(fmt.Sprintf("(expires in %s)", formatDuration(time.Until(s.expiresAt))))
	}
}

// profileUser extracts the "user:" field from a profile's yaml file.
func profileUser(tshHome, profile string) (string, error) {
	data, err := os.ReadFile(filepath.Join(tshHome, profile+profileSuffix))
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		rest, ok := strings.CutPrefix(line, "user:")
		if !ok {
			continue
		}
		return strings.Trim(strings.TrimSpace(rest), `"'`), nil
	}
	return "", nil
}

// checkSession reads the TLS certificate tsh stores for profile and reports
// whether it is still valid. Any failure to locate or parse the certificate
// is treated as an unknown (not-logged-in) session rather than an error,
// since that's the common case for a profile that has never been used.
func checkSession(tshHome, profile string) sessionStatus {
	user, err := profileUser(tshHome, profile)
	if err != nil || user == "" {
		return sessionStatus{}
	}

	certPath := filepath.Join(tshHome, "keys", profile, user+".crt")
	data, err := os.ReadFile(certPath)
	if err != nil {
		return sessionStatus{}
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return sessionStatus{}
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return sessionStatus{}
	}

	return sessionStatus{
		known:     true,
		expired:   time.Now().After(cert.NotAfter),
		expiresAt: cert.NotAfter,
	}
}

// formatDuration renders d as a short "1d2h" / "2h15m" / "15m" style string.
func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	d = d.Round(time.Minute)

	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	mins := d / time.Minute

	switch {
	case days > 0:
		return fmt.Sprintf("%dd%dh", days, hours)
	case hours > 0:
		return fmt.Sprintf("%dh%dm", hours, mins)
	default:
		return fmt.Sprintf("%dm", mins)
	}
}

// offerLogin checks profile's session and, if it's missing or expired, asks
// the user whether to run `tsh login` for it, defaulting to yes since it's
// a safe operation the user wants almost every time. Errors running tsh are
// reported but not treated as fatal, since the profile switch itself (which
// already happened by the time this runs) succeeded regardless.
func offerLogin(tshHome, profile string) {
	status := checkSession(tshHome, profile)
	if !status.needsLogin() {
		return
	}

	if status.known {
		fmt.Printf("Teleport session for %q has expired.\n", profile)
	} else {
		fmt.Printf("No active Teleport session found for %q.\n", profile)
	}

	prompt := promptui.Prompt{
		Label:     "Log in now",
		IsConfirm: true,
		Default:   "y",
	}
	if _, err := prompt.Run(); err != nil {
		// ErrAbort (explicit "no") or any other prompt error (e.g. ^C):
		// skip the login rather than treating it as fatal.
		return
	}

	cmd := exec.Command("tsh", "login", "--proxy="+profile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "tshx: tsh login failed: %v\n", err)
	}
}
