# Security

Please report security issues privately through GitHub's **Report a
vulnerability** feature rather than a public issue.

The tool sends staged textual diffs to Codex/OpenAI. Do not stage secrets. The
tool does not read or store Codex credentials and refuses API-key fallback, but
it cannot determine whether staged source code is confidential or licensed for
external processing.

Generated commit messages are untrusted suggestions. Review them before
confirmation. Git hooks and signing remain active during `git commit`.
