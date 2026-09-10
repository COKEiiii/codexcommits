# Similar-project research

Research date: 2026-09-10. Searches covered GitHub repositories and public web
results for combinations of `Codex`, `commit message`, `staged diff`,
`codex-commit`, and `codexcommits`.

## Conclusion

The overall idea already exists, including CLI tools that call Codex to produce
commit messages. `codexcommits` should therefore be evaluated on its narrow
workflow and safety properties, not as the first tool in this category.

## Closest projects found

| Project | Similarity | Main difference from this project |
|---|---|---|
| [takai/git-ai-commit](https://github.com/takai/git-ai-commit) | CLI reads staged diffs, calls existing LLM CLIs including Codex, and commits | Mature multi-provider tool with configuration, filtering, and Homebrew support; this project is intentionally Codex/ChatGPT-only and checks snapshot identity before confirmation |
| [spotdemo4/codex-commit](https://github.com/spotdemo4/codex-commit) | Rust CLI calls Codex, shows an editable message, and commits | Automatically runs `git add .`; this project requires explicit staging and never stages files |
| [rajatxs/codex-commit](https://github.com/rajatxs/codex-commit) | Codex CLI generates Conventional Commits from staged changes | VS Code extension that fills the SCM input and never commits |
| [jiying2007/codex-commit](https://github.com/jiying2007/codex-commit) | Strong staged-snapshot and structured-output safety | Much larger VS Code safety suite; this project is a small terminal command with no runtime package dependencies |
| [kholmatov/codex-commit-generator](https://github.com/kholmatov/codex-commit-generator) | Codex CLI, staged diff, review, optional commit | VS Code extension with optional staging and push workflows |
| [Commit Freeloader](https://marketplace.visualstudio.com/items?itemName=dazzatronus.commit-freeloader) | Reuses paid CLI subscriptions, including Codex | VS Code extension and multi-provider integration |

OpenAI's Codex repository also contains user requests and examples for piping a
staged diff into `codex exec`, including [issue #1123](https://github.com/openai/codex/issues/1123)
and [issue #3481](https://github.com/openai/codex/issues/3481). These further
confirm that the basic workflow is established.

## Defensible project scope

`codexcommits` combines the following in a dependency-free terminal tool:

1. ChatGPT-authenticated Codex only, with API-key environment fallback removed.
2. An immutable staged Git tree used as the generation input.
3. HEAD/index identity checks before showing and accepting the result.
4. Isolated ephemeral read-only generation in a temporary directory.
5. JSON Schema output plus local Conventional Commit validation.
6. Explicit commit confirmation while preserving Git hooks and never pushing.

This combination is useful, but individual elements also appear in existing
projects. Future documentation should continue to avoid novelty claims.
