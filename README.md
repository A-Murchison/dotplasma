# dotplasma

[![Release](https://img.shields.io/github/v/release/A-Murchison/dot-plasma?sort=semver)](https://github.com/A-Murchison/dot-plasma/releases)
[![Documentation](https://img.shields.io/badge/docs-usage%20reference-blue)](docs/usage.md)


`dotplasma` is a CLI for saving, comparing, importing, and restoring KDE Plasma desktop configuration as plain files. It is designed for Plasma users who want an inspectable backup of panels, widgets, shortcuts, desktop layout, KWin settings, and related config before making changes or moving a setup between machines.

<img width="960" height="540" alt="dotplasma demo" src="https://github.com/user-attachments/assets/7e330581-bcd3-4ced-ac8e-04b0b2206565" />

## Features

- **Profile snapshots:** Save your current Plasma configuration into named profiles under `profiles/<name>/`.
- **Readable diffs:** Compare a saved profile with your live desktop before deciding what to keep.
- **Safe imports:** Import a profile from another machine without applying it immediately.
- **Cautious restores:** Preview writes with `--dry-run`; real restores create a backup first.
- **Git-friendly output:** Profiles are normal files that you can review, copy, commit, or ignore yourself.
- **Conservative scope:** Only an [allow list](docs/usage.md#what-gets-captured) set of known Plasma configuration files is captured.
- **No root, no network:** `dotplasma` works in your user account and does not push or sync anything.

## Diff Demo



<img width="982" height="693" alt="dotplasma diff output" src="https://github.com/user-attachments/assets/49b5bb1c-dc49-4b61-a544-fb2bb1e0d14e" />

## Installation

### Install with Go

```sh
go install github.com/A-Murchison/dot-plasma/cmd/dotplasma@latest
dotplasma --help
```

### Download a release

Linux `.deb` and `.rpm` packages are available from the [GitHub Releases](https://github.com/A-Murchison/dot-plasma/releases) page. Release archives are also available for Linux, macOS, and Windows.

### Build from source

```sh
git clone https://github.com/A-Murchison/dot-plasma.git
cd dot-plasma
go build -o ~/.local/bin/dotplasma ./cmd/dotplasma
```

You can also run from a checkout without installing:

```sh
go run ./cmd/dotplasma --help
```

## Usage

By default, `dotplasma` stores profiles and restore backups under:

```text
${XDG_CONFIG_HOME:-~/.config}/dotplasma/
```

Use `--out DIR` with supported commands if you want to keep profiles somewhere else.

```sh
# Check that dotplasma can see your Plasma files and output directory.
dotplasma doctor

# Show which tracked Plasma config files exist on this machine.
dotplasma inspect-live

# Save your current desktop as a profile named "main".
dotplasma save main

# See what changed since that snapshot.
dotplasma diff main

# Preview a restore without writing anything.
dotplasma apply main --dry-run

# Apply the profile and optionally reload Plasma Shell.
dotplasma apply main --reload plasmashell
```

## Command overview

```text
dotplasma [--config PATH] doctor [--out DIR]
dotplasma [--config PATH] inspect-live
dotplasma [--config PATH] save <profile> [--out DIR] [--dry-run]
dotplasma [--config PATH] diff [profile] [--out DIR] [--format text] [--color auto|always|never] [--verbose]
dotplasma [--config PATH] list [--out DIR]
dotplasma [--config PATH] import <source-dir> <profile> [--out DIR] [--dry-run]
dotplasma [--config PATH] apply <profile> [--out DIR] [--dry-run] [--reload none|plasmashell]
dotplasma [--config PATH] version
```

## Safety model

`dotplasma` is intentionally conservative:

- It does not require root.
- It does not access the network.
- It does not make Git commits or push anywhere.
- It only tracks known Plasma files from its allowlist.
- Writing commands support `--dry-run`.
- Restore creates backups before updating tracked live files.
- Plasma Shell reload is opt-in.

> [!WARNING]
> Treat saved profiles as private. Plasma config can include local paths, widget settings, wallpaper paths, monitor layout details, and other machine-specific data. Review files before sharing or publishing them.

## Documentation

- [Usage and reference](docs/usage.md) — install details, commands, configuration, backups, captured files, workflows, Plasma Shell reloads, output layout, and safety notes.
- [Contributing](CONTRIBUTING.md) — development commands and project guidelines.
- [Changelog](CHANGELOG.md) — release history.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development commands and project guidelines.

## Credits

- Thanks to [EliverLara/Nordic](https://github.com/EliverLara/Nordic) for the Nordic theme shown in the demos.
- Thanks to [Prayag2/konsave](https://github.com/Prayag2/konsave) for inspiration.
