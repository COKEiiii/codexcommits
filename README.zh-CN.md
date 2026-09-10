# codexcommits

使用你已经登录的 Codex CLI 和 ChatGPT 订阅，根据 Git 暂存快照生成
Conventional Commit。提交前可以确认、编辑、重新生成或取消。

## 定位

这不是首个 AI commit 工具。它专注于一套小而明确的安全流程：

- 只接受 `codex login` 的 ChatGPT 登录，不会静默改用 API Key 计费；
- 只把已暂存快照交给 Codex，不读取未暂存版本；
- 不自动暂存、不自动 push；
- 生成期间 HEAD 或暂存树变化时，作废旧结果；
- 使用临时、只读、低推理强度的 `codex exec`；
- 用 JSON Schema 和本地规则校验单行 Conventional Commit；
- 保留正常的 Git hooks、签名和配置；
- 运行时只使用 Python 标准库。

同类项目和差异见 [MARKET_RESEARCH.md](MARKET_RESEARCH.md)。

## 要求与安装

需要 macOS/Linux、Python 3.10+、Git、Codex CLI 0.149.0+，并已使用
ChatGPT 登录：

```bash
codex login status
pipx install git+https://github.com/COKEiiii/codexcommits.git
```

## 使用

```bash
git add main.py
codexcommits
```

`git diff --cached` 只供你自己检查，可以省略。工具会自动读取暂存快照。

- `y`：使用当前说明提交；
- `e`：编辑完整说明并重新确认；
- `r`：重新生成，会再次消耗 Codex 额度；
- `n`、直接回车或 Ctrl-C：取消并保留暂存内容。

只生成、不提交：

```bash
codexcommits --print
```

默认使用 `gpt-5.6-luna` 和 low 推理强度。暂存 diff 超过 100 KB 时，会在
调用 Codex 前要求拆分。每次生成都会使用 Codex 额度，具体消耗取决于模型、
输入、推理和运行时因素。

暂存文本 diff 会通过 Codex 发送给 OpenAI。工具不读取或保存登录凭据，也不包含
Git 文本 diff 中不可见的二进制文件内容。生成结果仍可能不准确，请在确认前检查。

## 开发

```bash
python3 -m pip install -e .
python3 -m unittest discover -s tests -v
```

MIT License
