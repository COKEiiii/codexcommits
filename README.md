# codexcommits

[![CI](https://github.com/COKEiiii/codexcommits/actions/workflows/ci.yml/badge.svg)](https://github.com/COKEiiii/codexcommits/actions/workflows/ci.yml)
[![Python 3.10+](https://img.shields.io/badge/python-3.10%2B-blue.svg)](https://www.python.org/)
[![MIT License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Generate an accurate Conventional Commit from your staged snapshot using the
Codex CLI and the ChatGPT account you already use. Review, edit, regenerate, or
cancel before Git creates the commit.

```console
$ git add src/parser.py
$ codexcommits
Staged changes:
 src/parser.py | 12 +++++++++---

feat(parser): handle nested markdown tables

[y] commit  [e] edit  [r] regenerate  [n/Enter] cancel:
```

[简体中文](README.zh-CN.md)

## Why this project?

AI commit generators already exist. This project deliberately has a narrower
contract for people who use Codex through a ChatGPT subscription:

- requires `codex login` with ChatGPT and refuses silent API-key fallback;
- sends only an immutable staged-tree diff, never unstaged working-tree edits;
- never stages or pushes;
- invalidates the result if HEAD or the staged tree changes while generating;
- uses ephemeral, read-only `codex exec` with low reasoning;
- requires schema-valid, single-line Conventional Commit output;
- preserves normal Git hooks, signing, and configuration;
- uses only the Python standard library at runtime.

See [MARKET_RESEARCH.md](MARKET_RESEARCH.md) for similar projects and the
honest scope comparison.

## Requirements

- macOS or Linux
- Python 3.10+
- Git
- Codex CLI 0.149.0+ on `PATH`
- a working ChatGPT login: `codex login status`

Install Codex and sign in with ChatGPT by following the
[official authentication guide](https://learn.chatgpt.com/docs/auth).

## Install

The easiest isolated installation uses `pipx`:

```bash
pipx install git+https://github.com/COKEiiii/codexcommits.git
```

Or install from a clone:

```bash
git clone https://github.com/COKEiiii/codexcommits.git
cd codexcommits
python3 -m pip install .
```

## Usage

Stage exactly what you want to commit, then run the tool:

```bash
git add path/to/file
codexcommits
```

`git diff --cached` is optional and useful for your own review.
`codexcommits` reads the staged snapshot itself.

At the prompt:

- `y` commits with the displayed message;
- `e` lets you replace the complete subject and asks again;
- `r` calls Codex again and consumes additional allowance;
- `n`, Enter, or Ctrl-C cancels while preserving the index.

Generate without committing:

```bash
codexcommits --print
```

Use another available Codex model for one run:

```bash
codexcommits --model gpt-5.6-terra
```

The default is `gpt-5.6-luna` with low reasoning. You can set a persistent
shell-level default with `CODEXCOMMITS_MODEL`.

## Data and usage

The staged textual diff is sent through Codex to OpenAI under the policies of
the ChatGPT account used by Codex. Binary contents are not included in Git's
text diff. A request is rejected before Codex runs when the diff exceeds 100 KB.

Each generation consumes Codex allowance. Consumption varies with the model,
input size, reasoning, and other runtime factors. See the
[official Codex pricing and usage documentation](https://learn.chatgpt.com/docs/pricing).

`codexcommits` removes API-key environment variables from the Codex child
process and verifies that `codex login status` reports ChatGPT. It does not read
or store credentials.

## Development

```bash
python3 -m pip install -e .
python3 -m unittest discover -s tests -v
```

The tests use temporary Git repositories and a mocked generation boundary; CI
does not need Codex credentials. A maintainer can run a manual end-to-end check
with `codexcommits --print` in a repository containing staged changes.

## Limitations

- Generated text can still be inaccurate; review it before accepting.
- Only one-line Conventional Commit subjects are supported in v0.1.
- Diffs larger than 100 KB must be split into smaller commits.
- Windows support has not been tested.

## License

MIT
