package profile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/A-Murchison/dot-plasma/internal/paths"
)

func TestApplyPreflightRejectsCorruptLaterFile(t *testing.T) {
	t.Parallel()
	root, out := t.TempDir(), t.TempDir()
	roots := paths.LiveRoots{Config: root, LocalShare: t.TempDir()}
	writeFile(t, filepath.Join(root, "kdeglobals"), "[General]\nvalue=saved\n")
	writeFile(t, filepath.Join(root, "kwinrc"), "[General]\nvalue=saved\n")
	if _, err := Save(context.Background(), SaveOptions{Profile: "work", Out: out, Roots: roots}); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(root, "kdeglobals"), "[General]\nvalue=live\n")
	writeFile(t, filepath.Join(out, "profiles", "work", "files", "config", "kwinrc"), "corrupt")
	if _, err := Apply(context.Background(), ApplyOptions{Profile: "work", Out: out, Roots: roots}); err == nil {
		t.Fatal("Apply() succeeded with corrupted later file")
	}
	if got := readFile(t, filepath.Join(root, "kdeglobals")); got != "[General]\nvalue=live\n" {
		t.Fatalf("live file changed during preflight: %q", got)
	}
	if _, err := os.Stat(filepath.Join(out, "backups")); !os.IsNotExist(err) {
		t.Fatalf("backup created during failed preflight: %v", err)
	}
}

func TestApplyRollsBackEarlierFileOnBackupFailure(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		missing bool
	}{
		{name: "restore updated file"},
		{name: "remove created file", missing: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root, out := t.TempDir(), t.TempDir()
			roots := paths.LiveRoots{Config: root, LocalShare: t.TempDir()}
			first := filepath.Join(root, "kdeglobals")
			writeFile(t, first, "[General]\nvalue=saved\n")
			writeFile(t, filepath.Join(root, "kwinrc"), "[General]\nvalue=saved\n")
			if _, err := Save(context.Background(), SaveOptions{Profile: "work", Out: out, Roots: roots, Now: fixedNow}); err != nil {
				t.Fatal(err)
			}
			if tc.missing {
				if err := os.Remove(first); err != nil {
					t.Fatal(err)
				}
			} else {
				writeFile(t, first, "[General]\nvalue=live\n")
			}
			writeFile(t, filepath.Join(root, "kwinrc"), "[General]\nvalue=live\n")
			block := filepath.Join(out, "backups", "work", fixedNow().UTC().Format("20060102T150405Z"), "files", "config", "kwinrc")
			if err := os.MkdirAll(block, 0o700); err != nil {
				t.Fatal(err)
			}
			if _, err := Apply(context.Background(), ApplyOptions{Profile: "work", Out: out, Roots: roots, Now: fixedNow}); err == nil {
				t.Fatal("Apply succeeded despite blocked backup")
			}
			if tc.missing {
				if _, err := os.Stat(first); !os.IsNotExist(err) {
					t.Fatalf("created file remains: %v", err)
				}
			} else if got := readFile(t, first); got != "[General]\nvalue=live\n" {
				t.Fatalf("rollback = %q, want original live content", got)
			}
		})
	}
}

func TestRollbackApplyReportsMissingBackup(t *testing.T) {
	t.Parallel()
	profileDir, backupDir := t.TempDir(), t.TempDir()
	target := filepath.Join(t.TempDir(), "kwinrc")
	writeFile(t, target, "applied")
	actions := []Action{{Status: "update", Source: filepath.Join(profileDir, "files", "config", "kwinrc"), Target: target, written: []byte("applied")}}
	if err := rollbackApply(actions, profileDir, backupDir); err == nil || !strings.Contains(err.Error(), "rollback live file") {
		t.Fatalf("rollbackApply error = %v, want rollback failure", err)
	}
}

func TestRollbackApplyPreservesChangedTargets(t *testing.T) {
	for _, status := range []string{"create", "update"} {
		t.Run(status, func(t *testing.T) {
			profileDir, backupDir := t.TempDir(), t.TempDir()
			target := filepath.Join(t.TempDir(), "kwinrc")
			writeFile(t, target, "concurrent edit")
			backup := filepath.Join(backupDir, "files", "config", "kwinrc")
			writeFile(t, backup, "original")
			actions := []Action{{Status: status, Source: filepath.Join(profileDir, "files", "config", "kwinrc"), Target: target, written: []byte("applied")}}
			err := rollbackApply(actions, profileDir, backupDir)
			if err == nil || !strings.Contains(err.Error(), "rollback conflict") {
				t.Fatalf("rollbackApply error = %v, want conflict", err)
			}
			if got := readFile(t, target); got != "concurrent edit" {
				t.Errorf("target = %q", got)
			}
			if got := readFile(t, backup); got != "original" {
				t.Errorf("backup = %q", got)
			}
		})
	}
}

func TestApplyDryRunDoesNotWriteLiveFilesOrBackups(t *testing.T) {
	t.Parallel()

	configRoot := t.TempDir()
	out := t.TempDir()
	writeFile(t, filepath.Join(configRoot, "kdeglobals"), "[General]\nColorScheme=Old\n")
	_, err := Save(context.Background(), SaveOptions{
		Profile: "work",
		Out:     out,
		Roots:   paths.LiveRoots{Config: configRoot, LocalShare: t.TempDir()},
		Now:     fixedNow,
	})
	if err != nil {
		t.Fatalf("Save returned error: %s", err.Error())
	}
	writeFile(t, filepath.Join(configRoot, "kdeglobals"), "[General]\nColorScheme=New\n")

	result, err := Apply(context.Background(), ApplyOptions{
		Profile: "work",
		Out:     out,
		DryRun:  true,
		Roots:   paths.LiveRoots{Config: configRoot, LocalShare: t.TempDir()},
		Now:     fixedNow,
	})
	if err != nil {
		t.Fatalf("Apply returned error: %s", err.Error())
	}
	if !hasActionStatus(result.Actions, "would-update") {
		t.Fatalf("Apply actions = %#v, want would-update", result.Actions)
	}
	if got := readFile(t, filepath.Join(configRoot, "kdeglobals")); !strings.Contains(got, "ColorScheme=New") {
		t.Fatalf("dry run changed live file: %q", got)
	}
	if _, err := os.Stat(filepath.Join(out, "backups")); !os.IsNotExist(err) {
		t.Fatalf("dry run created backups or unexpected stat error: %v", err)
	}
}

func TestApplyRestoresSavedFilesAndBacksUpLiveFiles(t *testing.T) {
	t.Parallel()

	configRoot := t.TempDir()
	out := t.TempDir()
	writeFile(t, filepath.Join(configRoot, "kdeglobals"), "[General]\nColorScheme=Saved\n")
	_, err := Save(context.Background(), SaveOptions{
		Profile: "work",
		Out:     out,
		Roots:   paths.LiveRoots{Config: configRoot, LocalShare: t.TempDir()},
		Now:     fixedNow,
	})
	if err != nil {
		t.Fatalf("Save returned error: %s", err.Error())
	}
	writeFile(t, filepath.Join(configRoot, "kdeglobals"), "[General]\nColorScheme=Live\n")

	result, err := Apply(context.Background(), ApplyOptions{
		Profile: "work",
		Out:     out,
		Roots:   paths.LiveRoots{Config: configRoot, LocalShare: t.TempDir()},
		Now:     fixedNow,
	})
	if err != nil {
		t.Fatalf("Apply returned error: %s", err.Error())
	}
	if !hasActionStatus(result.Actions, "update") {
		t.Fatalf("Apply actions = %#v, want update", result.Actions)
	}
	if got := readFile(t, filepath.Join(configRoot, "kdeglobals")); !strings.Contains(got, "ColorScheme=Saved") {
		t.Fatalf("live file was not restored: %q", got)
	}
	backup := readFile(t, filepath.Join(result.BackupDir, "files", "config", "kdeglobals"))
	if !strings.Contains(backup, "ColorScheme=Live") {
		t.Fatalf("backup does not contain previous live data: %q", backup)
	}
}

func TestApplyRestoresPreviouslyVolatileKeys(t *testing.T) {
	t.Parallel()

	configRoot := t.TempDir()
	out := t.TempDir()
	writeFile(t, filepath.Join(configRoot, "plasma-org.kde.plasma.desktop-appletsrc"), "[Containments][1]\nlastScreen=0\nplugin=org.kde.plasma.panel\n")
	_, err := Save(context.Background(), SaveOptions{
		Profile: "work",
		Out:     out,
		Roots:   paths.LiveRoots{Config: configRoot, LocalShare: t.TempDir()},
		Now:     fixedNow,
	})
	if err != nil {
		t.Fatalf("Save returned error: %s", err.Error())
	}
	writeFile(t, filepath.Join(configRoot, "plasma-org.kde.plasma.desktop-appletsrc"), "[Containments][1]\nlastScreen=2\nplugin=org.kde.plasma.panel\n")

	result, err := Apply(context.Background(), ApplyOptions{
		Profile: "work",
		Out:     out,
		Roots:   paths.LiveRoots{Config: configRoot, LocalShare: t.TempDir()},
		Now:     fixedNow,
	})
	if err != nil {
		t.Fatalf("Apply returned error: %s", err.Error())
	}
	if !hasActionStatus(result.Actions, "update") {
		t.Fatalf("Apply actions = %#v, want update", result.Actions)
	}
	if got := readFile(t, filepath.Join(configRoot, "plasma-org.kde.plasma.desktop-appletsrc")); !strings.Contains(got, "lastScreen=0") {
		t.Fatalf("screen binding key was not restored: %q", got)
	}
}

func TestApplyRejectsOldPlasmaShellLayoutMissingScreenBindings(t *testing.T) {
	t.Parallel()

	out := t.TempDir()
	profileDir := filepath.Join(out, "profiles", "work")
	data := "[Containments][1]\nplugin=org.kde.plasma.folder\n\n[Containments][2]\nlastScreen=0\nplugin=org.kde.panel\n"
	writeFile(t, filepath.Join(profileDir, "files", "config", "plasma-org.kde.plasma.desktop-appletsrc"), data)
	writeFile(t, filepath.Join(profileDir, "manifests", "files.toml"), strings.Join([]string{
		"[[files]]",
		`root = "config"`,
		`path = "plasma-org.kde.plasma.desktop-appletsrc"`,
		`parser = "kconfig"`,
		fmt.Sprintf(`sha256 = %q`, sha256Hex([]byte(data))),
		fmt.Sprintf("bytes = %d", len(data)),
		"",
	}, "\n"))
	writeFile(t, filepath.Join(profileDir, "profile.toml"), `profile = "work"`)
	configRoot := t.TempDir()
	writeFile(t, filepath.Join(configRoot, "plasma-org.kde.plasma.desktop-appletsrc"), "[Containments][1]\nlastScreen=1\nplugin=org.kde.plasma.folder\n")

	_, err := Apply(context.Background(), ApplyOptions{
		Profile: "work",
		Out:     out,
		Roots:   paths.LiveRoots{Config: configRoot, LocalShare: t.TempDir()},
		Now:     fixedNow,
	})
	if err == nil || !strings.Contains(err.Error(), "missing lastScreen for containments 1") {
		t.Fatalf("Apply error = %v, want missing lastScreen", err)
	}
	if got := readFile(t, filepath.Join(configRoot, "plasma-org.kde.plasma.desktop-appletsrc")); !strings.Contains(got, "lastScreen=1") {
		t.Fatalf("live layout was changed despite validation failure: %q", got)
	}
}

func TestApplyRejectsManifestFileOutsideAllowlist(t *testing.T) {
	t.Parallel()

	out := t.TempDir()
	profileDir := filepath.Join(out, "profiles", "work")
	writeFile(t, filepath.Join(profileDir, "files", "config", "not-allowlisted"), "data")
	writeFile(t, filepath.Join(profileDir, "manifests", "files.toml"), strings.Join([]string{
		"[[files]]",
		`root = "config"`,
		`path = "not-allowlisted"`,
		`parser = "raw-copy"`,
		`sha256 = "3a6eb0790f39ac87c94f3856b2dd2c5d110e6811602261a9a923d3bb23adc8b7"`,
		"bytes = 4",
		"",
	}, "\n"))
	writeFile(t, filepath.Join(profileDir, "profile.toml"), `profile = "work"`)

	_, err := Apply(context.Background(), ApplyOptions{
		Profile: "work",
		Out:     out,
		DryRun:  true,
		Roots:   paths.LiveRoots{Config: t.TempDir(), LocalShare: t.TempDir()},
		Now:     fixedNow,
	})
	if err == nil || !strings.Contains(err.Error(), "not in the current allowlist") {
		t.Fatalf("Apply error = %v, want not allowlisted", err)
	}
}

func hasActionStatus(actions []Action, status string) bool {
	for _, action := range actions {
		if action.Status == status {
			return true
		}
	}
	return false
}
