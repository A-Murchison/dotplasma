# dotplasma — KDE Plasma config backup, diff, and restore CLI

Save, compare, and restore KDE Plasma desktop configuration as plain files.

`dotplasma` is an open-source CLI for people who want their KDE Plasma setup to be inspectable, backed up, and easy to move between machines. It snapshots an allowlisted set of Plasma desktop configuration files from your home directory into regular folders that you can review, copy, or commit to Git yourself.

<img width="960" height="540" alt="dotplasmademo" src="https://github.com/user-attachments/assets/7e330581-bcd3-4ced-ac8e-04b0b2206565" />

## Good fit

`dotplasma` is most useful when you want to:

- keep a known-good Plasma layout before experimenting with panels, widgets, shortcuts, or desktop settings;
- review what changed in Plasma's config files before committing dotfiles;
- back up and restore Plasma configuration on the same machine;
- move a profile between similar Plasma installations while reviewing it before applying;
- backup to git.

`dotplasma` focuses on that workflow: snapshot, diff, review, and only then restore if you choose to.

## Features

- **Save Plasma profiles** into `profiles/<name>/` as normal files.
- **Review changes** between a saved profile and your live desktop.
- **Import profiles safely** from another machine without applying them.
- **Restore cautiously** with `--dry-run` support and automatic pre-restore backups.
- **Git-friendly output** that you control; no automatic commits or remote sync.
- **Conservative file scope** using an allowlist of known Plasma config files.
- **No root and no network access** required.

## Quick start

By default, `dotplasma` stores profiles and restore backups under:

```text
${XDG_CONFIG_HOME:-~/.config}/dotplasma/
```

You can also choose a separate directory with `--out`, which is useful if you want to keep profiles in a Git repository.

```sh
# Check that dotplasma can see your Plasma files and output directory.
dotplasma doctor

# Show which known Plasma config files exist on this machine.
dotplasma inspect-live

# Save your current desktop as a profile named "main".
dotplasma save main

# See what changed since that snapshot.
dotplasma diff main

# Preview a restore without writing anything.
dotplasma apply main --dry-run
```

<img width="982" height="693" alt="dotplasma-diff" src="https://github.com/user-attachments/assets/49b5bb1c-dc49-4b61-a544-fb2bb1e0d14e" />

## Install

### Go install

If you have Go installed, this is the simplest option:

```sh
go install github.com/A-Murchison/dot-plasma/cmd/dotplasma@latest
dotplasma --help
```

### Packages and manual builds

RPM packages are available for RPM-based systems from the [GitHub Releases](https://github.com/A-Murchison/dot-plasma/releases) page. For releases that include Debian packages, download the `.deb` package for your architecture from the same page.

From a local checkout:

```sh
go build -o ~/.local/bin/dotplasma ./cmd/dotplasma
# or run without installing
go run ./cmd/dotplasma --help
```

## Documentation

For the full usage reference, see:

- [Usage and reference](docs/usage.md) — install details, commands, configuration, backups, captured files, workflows, Plasma Shell reloads, output layout, and safety notes.
- [Contributing](CONTRIBUTING.md) — development commands and project guidelines.
- [Changelog](CHANGELOG.md) — release history.

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

## What gets captured?

`dotplasma` tracks an allowlisted set of KDE Plasma configuration files, including files related to panels, widgets, desktop layout, shortcuts, Plasma Shell, KWin, and display/layout state where supported.

Run this to see exactly which tracked files exist on your machine:

```sh
dotplasma inspect-live
```

See the [usage reference](docs/usage.md#what-gets-captured) for the current built-in allowlist.

## Safety model

`dotplasma` is intentionally conservative:

- It does not require root.
- It does not access the network.
- It does not make Git commits or push anywhere.
- It only tracks a known allowlist of Plasma files.
- Writing commands support `--dry-run`.
- Restore creates backups before updating tracked live files.
- Plasma Shell reload is opt-in.

> [!WARNING]
> Treat saved profiles as private. Plasma config can include local paths, widget settings, wallpaper paths, monitor layout details, and other machine-specific data. Review files before sharing or publishing them.

## Contributing

Issues, bug reports, real-world Plasma edge cases, and focused pull requests are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for development commands and project guidelines.

---

Shout out to https://github.com/EliverLara/Nordic for the nordic theme in the gifs.

Shout out to https://github.com/Prayag2/konsave for the inspiration
