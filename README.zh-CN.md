# codexcommits

面向 Codex 新手的 commit 小工具。

只要你的电脑已经能使用 Codex CLI，就可以暂存改动、运行一个命令，并在
Git 真正提交前检查 AI 生成的 Conventional Commit。无需申请 API Key、选择
模型提供商、编写提示词配置，也无需安装 Python 或 Node.js。

```console
$ git add src/parser.py
$ codexcommits
Staged changes:
 src/parser.py | 12 +++++++++---

feat(parser): handle nested markdown tables

[y] commit  [e] edit  [r] regenerate  [n/Enter] cancel:
```

[English](README.md)

## 使用前提

- 已安装 Git；
- 已安装 [Codex CLI](https://developers.openai.com/codex/cli/)，终端能运行 `codex`；
- Codex 已通过 ChatGPT 登录，`codex login status` 能正常显示状态。

## 安装

### macOS、Linux 和 WSL2

```bash
brew install COKEiiii/tap/codexcommits
```

以后更新：

```bash
brew upgrade codexcommits
```

### Windows PowerShell

```powershell
irm https://raw.githubusercontent.com/COKEiiii/codexcommits/main/install.ps1 | iex
```

安装脚本会下载最新的 Windows 可执行文件，放入当前用户目录并加入用户
`PATH`。执行前可以先阅读 [install.ps1](install.ps1)。macOS、Linux 和 Windows
的手动安装包也可以从 [Releases](https://github.com/COKEiiii/codexcommits/releases/latest)
下载。

## 日常使用

```bash
git add main.py
codexcommits
```

`git diff --cached` 不是必需步骤。它只用于你自己提前检查暂存差异；
`codexcommits` 会自动读取已经暂存的快照。

出现确认提示后：

| 输入 | 作用 |
|---|---|
| `y` | 使用当前说明创建 commit |
| `e` | 手动替换完整说明，然后再次确认 |
| `r` | 让 Codex 重新生成一次，会再次使用额度 |
| `n`、直接回车或 Ctrl-C | 取消，暂存内容保持不变 |

只生成说明，不创建 commit：

```bash
codexcommits --print
```

默认使用 Codex CLI 为你的账户选择的模型，并采用 low 推理强度。大多数用户
不需要修改任何配置。可选高级参数可以通过 `codexcommits --help` 查看。

## 这个命令做了什么

1. 检查当前位置是否为 Git 仓库，并确认存在暂存改动；
2. 记录准确的暂存 Git 树，只把它的文本 diff 交给 Codex；
3. 获得一条经过格式校验的 Conventional Commit；
4. 把建议显示给你，等待确认、编辑、重试或取消；
5. 确认提交前再次检查 HEAD 和暂存树没有变化，然后运行普通的
   `git commit`。

工具不会执行 `git add` 或 `git push`。你原有的 Git hooks、签名和配置仍然
生效。

## 隐私和额度

暂存的文本 diff 会通过 Codex 发送给 OpenAI，使用 Codex 当前登录的 ChatGPT
账户。不要暂存密钥。Git 文本 diff 不包含二进制文件内容；暂存 diff 超过
100 KB 时，工具会在调用 Codex 前要求拆分。

每次生成都会使用 Codex 额度，选择 `r` 会再请求一次。工具会移除 Codex 子
进程中的 API Key 环境变量，并确认当前登录状态包含 ChatGPT；它不会读取或
保存登录凭据。

## 与同类项目的关系

使用 AI 或 Codex 生成 commit message 的想法已经存在。本项目专注于让 Codex
新手获得一套小而明确的流程：自己暂存、使用现有 ChatGPT 登录、校验暂存快照、
提交前确认、不配置模型提供商。详细对比见
[MARKET_RESEARCH.md](MARKET_RESEARCH.md)。

## 开发

```bash
go test ./...
go vet ./...
go build .
```

MIT License
