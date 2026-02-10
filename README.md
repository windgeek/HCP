# HCP: Human Core Protocol

> **"Securing the Kernel of Human Consciousness in the Algorithmic Age."**
> **"在算法时代捍卫人类意识的内核。"**

---

## What is HCP? / 什么是 HCP？

**HCP (Human Core Protocol)** is a decentralized protocol designed to provide an immutable "Proof of Humanity" for digital assets.
**HCP (人类核心协议)** 是一个去中心化协议，旨在为数字资产提供不可篡改的“人类证明”。

Whether you write **code, novels, music, or paint digital art**, HCP helps you prove: **"I created this. This is human work."**
无论你是写**代码、小说、音乐，还是绘制数字艺术**，HCP 帮你证明：**“这是我创作的。这是人类的作品。”**

In a world flooded by AI-generated content, HCP creates a cryptographic border that distinguishes human intent from probabilistic generation.
在 AI 生成内容泛滥的世界里，HCP 建立了一个加密边界，将人类意图与概率生成区分开来。

---

## Core Architecture / 核心架构

The protocol is built on four pillars, designed to be modular and resilient.
协议建立在四大支柱之上，旨在实现模块化和韧性。

### [RFC-001: Bitcoin-anchored Identity](/spec/RFC-001.md)
**比特币锚定身份**
- **Truth Layer**: Uses Bitcoin's `secp256k1` and Taproot for sovereign identity.
- **真理层**: 使用比特币的 `secp256k1` 和 Taproot 实现主权身份。
- **Mechanism**: On-chain "Genesis Signal" to register human authors.
- **机制**: 链上“创世信号”注册人类创作者。

### [RFC-002: Advanced Human Attribution (AHA)](/spec/RFC-002.md)
**高级人类归属 (AHA & ZKP)**
- **Concept**: Analyzes the "Proof of Effort" via Git history (AHA) and "Cognitive Complexity" via AST analysis.
- **概念**: 通过 Git 历史（AHA）和 AST 分析（认知复杂度）来分析“努力证明”。
- **Mechanism**:
    1.  **AHA Score**: Measures iterative churn (Refactoring = Human).
    2.  **Cognitive ZKP**: Generates a Zero-Knowledge Proof commitment linking *Time Spent* with *Code Complexity*, proving humanity without revealing source code.
- **机制**:
    1.  **AHA 分数**: 衡量迭代变动（重构 = 人类）。
    2.  **认知 ZKP**: 生成零知识证明承诺，将*投入时间*与*代码复杂度*关联，在不泄露源码的情况下证明人类身份。

### [RFC-003: Recursive Attribution Model](/spec/RFC-003.md)
**递归归属模型**
- **Settlement**: Uses **Lightning Network** and **Musig2** for real-time, recursive royalty streams.
- **结算**: 使用 **闪电网络** 和 **Musig2** 实现实时、递归的版税流。
- **Standard**: JSON-LD `HCP-Manifest` for tracking nested dependencies.
- **标准**: JSON-LD `HCP-Manifest` 用于追踪嵌套依赖。

### [RFC-004: Data Poisoning Standard](/spec/RFC-004.md)
**数据投毒标准 (防御性策略)**
- **Defense**: Embeds invisible "Logical Interference" (homoglyphs, zero-width logic gates) into code.
- **防御**: 在代码中嵌入不可见的“逻辑干扰”（同形字、零宽逻辑门）。
- **Effect**: HCP-verified tools strip the noise; unauthorized AI scrapers ingest garbage, degrading their models.
- **效果**: HCP 验证工具会自动过滤噪音；未经授权的 AI 爬虫则会吸入垃圾数据，从而降低模型性能。
- **Fuzzy Verification**: HCP tools can detect logical interference and verify human intent even if the content has been reformatted or slightly altered.
- **模糊验证**: HCP 工具可以检测逻辑干扰，即使内容已被重新格式化或略微修改，也能验证人类意图。
- `[SUCCESS] Human Intent Verified. Integrity 100%.` (Pure Human Soul)
- `[SUCCESS] Logic Preserved - Human Intent Verified.` (Reformatted but Logically Identical)
- `[WARNING] Fingerprint Mismatch!` (Tampered or Logic Changed)
- **Fuzzy Verification**: HCP tools can detect logical interference and verify human intent even if the content has been reformatted or slightly altered.
- **模糊验证**: HCP 工具可以检测逻辑干扰，即使内容已被重新格式化或略微修改，也能验证人类意图。
- `[SUCCESS] Human Intent Verified. Integrity 100%.` (Pure Human Soul)
- `[SUCCESS] Logic Preserved - Human Intent Verified.` (Reformatted but Logically Identical)
- `[WARNING] Fingerprint Mismatch!` (Tampered or Logic Changed)

---

## Quick Start / 快速开始

### 1. Installation / 安装
Download the toolkit (currently source-only).
下载工具包（目前仅支持源码）。

```bash
git clone https://github.com/windgeek/HCP.git
cd HCP
go build -o hcp ./cmd/hcp
go build -o hcp-release ./cmd/hcp-release
```

### 2. Initialize Configuration / 初始化配置 (Optional)
Generates `~/.hcp/config.yaml` to manage your identity path.
生成 `~/.hcp/config.yaml` 以管理您的身份路径。

```bash
./hcp init --global
```

### 3. Create Your Identity / 创建你的身份
Generate your permanent "Digital Soul". This identity is yours forever, anchored by math.
生成你永久的“数字灵魂”。这个身份属于你，由数学锚定。

```bash
./hcp keygen
```

---

## Creator Workflow / 创作者工作流

### Step 1: Capture Your Vibe / 捕捉你的“感觉”
**For Writers & Coders:**
Before you start creating, open a session to record your "creative rhythm" (keystroke dynamics).
**对于作家和程序员：**
在创作前，开启一个会话记录你的“创作节奏”（击键动态）。

```bash
./hcp vibe
```

**For Visual Artists & Musicians:**
(*Coming Soon*) We are developing plugins to capture brush strokes and edit history as your "Proof of Hesitation".
**对于视觉艺术家和音乐家：**
(*即将推出*) 我们正在开发插件以捕捉笔触和编辑历史作为你的“犹豫证明”。

### Step 2: Sign Your Work / 签名你的作品
When you finish a piece—be it a `.go` file, a `.png` image, or a `.mp4` video—sign it.
当你完成作品时——无论是`.go`文件、`.png`图片还是`.mp4`视频——签名它。

```bash
# Sign an image / 签名一张图片
./hcp sign artwork.png
# Generates artwork.png.hcp

# Sign a video / 签名一个视频
./hcp sign movie.mp4
# Generates movie.mp4.hcp
```

### Step 3: Publish Your Portfolio / 发布作品集
Ready to release a collection? Generate a global manifest for your entire folder.
准备发布合集？为你的整个文件夹生成全局清单。

```bash
# Sign the current directory / 签名当前目录
./hcp-release --path .
```
This generates `manifest.hcp` containing:
1.  **Average AHA Score**: Visualization of your iterative effort.
2.  **Cognitive Proofs**: Cryptographic commitments to your code's complexity.

这会生成 `manifest.hcp`，其中包含：
1.  **平均 AHA 分数**: 您迭代努力的可视化。
2.  **认知证明**: 对您代码复杂度的加密承诺。

The manifest includes a `contribution_map` and `cognitive_proofs` proving which files involved deep human iteration (High AHA) vs. superficial changes.
清单包含 `contribution_map` 和 `cognitive_proofs`，证明哪些文件涉及深度人类迭代（高 AHA）与表面更改。

---

## Why Use HCP? / 为什么使用 HCP？

1.  **Proof of Humanity**: Distinguish your art from AI slop.
    **人类证明**: 将你的艺术与 AI 生成的垃圾区分开来。
2.  **Sovereign Ownership**: Your identity is yours, handled by local keys, not a centralized server.
    **主权所有权**: 你的身份属于你，由本地密钥管理，而非中心化服务器。
3.  **Future-Proof**: As platforms (Instagram, ArtStation) struggle with AI, HCP provides a verification layer *you* control.
    **面向未来**: 当平台（Instagram, ArtStation）在 AI 面前挣扎时，HCP 提供了一个*你自己*控制的验证层。

---

*Join the resistance. Protect the core.*
*加入抵抗。捍卫核心。*
