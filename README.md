# HCP: Human Core Protocol

> **"Securing the Kernel of Human Consciousness in the Algorithmic Age."**
> **"在算法时代捍卫人类意识的内核。"**

---

## Vision / 愿景

**HCP (Human Core Protocol)** is a decentralized protocol designed to provide an immutable "Proof of Humanity" for digital assets. As Artificial General Intelligence (AGI) floods the world with zero-cost content, HCP creates a cryptographic border that distinguishes human intent from probabilistic generation.

**HCP (人类核心协议)** 是一个去中心化协议，旨在为数字资产提供不可篡改的“人类证明”。随着通用人工智能 (AGI) 以零成本内容淹没世界，HCP 建立了一个加密边界，将人类意图与概率生成区分开来。

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

### [RFC-002: Biological Entropy Engine](/spec/RFC-002.md)
**生物熵引擎 (犹豫证明)**
- **Concept**: Captures the "Proof of Hesitation"—the non-linear timing and refactoring patterns unique to human thought.
- **概念**: 捕捉“犹豫证明”——人类思维特有的非线性时间与重构模式。
- **Security**: Asymmetric work function. Fast to produce (just work), expensive to simulate (requires modeling human mistakes).
- **安全**: 非对称工作量证明。生成容易（只需工作），模拟昂贵（需建模人类错误）。

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

---

## Vibe Coding / 沉浸式编程

We believe security should be invisible. HCP integrates into your workflow with zero friction.
我们相信安全应该是隐形的。HCP 零摩擦地融入您的工作流。

### Installation / 安装
### Build from Source / 源码编译
```bash
go build -o hcp ./cmd/hcp
```


### Usage / 用法
Just like Git, but for your humanity.
像使用 Git 一样，但为了证明你的人类身份。

**1. Generate Identity / 生成身份**
```bash
./hcp keygen
# Generates ~/.hcp/identity.key (encrypted) and displays your Address
# 生成 ~/.hcp/identity.key（加密）并显示你的地址
```

**2. Sign a File / 签名文件**
```bash
./hcp sign <file>
# Creates <file>.hcp containing the signed manifest
# 创建包含签名清单的 <file>.hcp
```

**3. Anchor to Bitcoin / 锚定到比特币**
```bash
./hcp anchor <file>.hcp
# Mocks an OP_RETURN transaction with the content hash
# 模拟带有内容哈希的 OP_RETURN 交易
```

**4. Defense / 防御**
When you push to public repos, HCP injects **RFC-004** poison to deter AI scrapers. Your code remains readable to humans, but toxic to machines.
当你推送到公共仓库时，HCP 会注入 **RFC-004** 毒素以阻慑 AI 爬虫。你的代码对人类可读，但对机器有毒。

---

*Join the resistance. Protect the core.*
*加入抵抗。捍卫核心。*
