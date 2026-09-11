# Changelog

## 0.3.1

- Added explicit `--push` mode to run `git add .`, review and create the commit,
  then run `git push` after confirmation.

## 0.3.0

- Added a numbered model menu for one-time selection with `--choose-model`.
- Added saved model preferences with `--set-model` and `--reset-model`.
- Kept direct `--model`/`-m` overrides for advanced users.

## 0.2.2

- Added explicit per-run Codex model selection with `--model`/`-m`.
- Added the `CODEXCOMMITS_MODEL` persistent default documentation.

## 0.2.1

- Avoided unauthenticated GitHub API rate limits in the Windows installer.

## 0.2.0

- Replaced the Python package with standalone macOS, Linux, and Windows binaries.
- Added Homebrew installation for macOS, Linux, and WSL2.
- Added a one-command Windows PowerShell installer with checksum verification.
- Made the configured Codex CLI model the zero-configuration default.
- Added CI on macOS, Linux, and Windows plus automated release archives.

## 0.1.0

- Initial Python CLI with staged-snapshot validation and interactive review.
