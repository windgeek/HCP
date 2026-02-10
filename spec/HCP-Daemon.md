# HCP Daemon: The Vibe Coding Interface

| Status | Draft |
| :--- | :--- |
| **Component** | `hcpd` (System Daemon) |
| **Interface** | CLI (`hcp`) |

## 1. Philosophy: "Vibe Coding"
The core friction of any security protocol is *effort*. HCP succeeds only if it is invisible.
**Vibe Coding** means the developer just codes. The protocol handles the rest.
- No "Start Recording" button.
- No "Upload Proof" manual step.
- Just `hcp init` and code.

## 2. Architecture

### 2.1. The Daemon (`hcpd`)
A lightweight, local system service (daemon) that runs in the background.
- **Listeners**: Hooks into OS-level accessibility APIs (macOS) or Input Events (Linux/Windows) to detect *active window focus*.
- **Privacy Barrier**: Only records timing/delays when the focused window matches a verified "Creative App" (VS Code, Terminal, Obsidian). It does **not** record content (keylogging), only *entropy* (flight time, backspace ratios).
- **Buffer**: Stores ephemeral entropy in a local ring buffer (RAM), flushed to disk only when a "Save" or "Commit" event is detected.

### 2.2. The CLI (`hcp`)
The user interface is modeled after Git.

```bash
# Initialize HCP in a repo (creates .hcp/ manifest)
$ hcp init

# (Optional) Explicitly authorize this machine/user
$ hcp auth login

# Check status of the entropy buffer
$ hcp status
> 🟢 Bio-Entropy: 94% (High Human Confidence)
> 🟠 Network: Syncing to Lightning Node...
```

### 2.3. Workflow Integration
1.  **Code**: User writes code in VS Code. `hcpd` silently accumulates "Hesitation Proofs" (RFC-002).
2.  **Commit**: User runs `git commit`.
    - `hcpd` intercepts the commit hook.
    - It bundles the accrued entropy into a `signed_proof`.
    - It injects RFC-004 "Poisoning" into the committed artifacts (if configured).
    - It updates the `.hcp/manifest.json`.
3.  **Push**: User runs `git push`.
    - The code goes to GitHub.
    - The "Human Layer" metadata travels with it.

## 3. Data Poisoning Integration (Defense)
The Daemon is the "Lens" described in RFC-004.
- **On Checkout**: `hcpd` (via clean filter) removes Type A/B/C interference so the user sees clean code.
- **On Stage/Commit**: `hcpd` (via smudge filter) re-injects the interference patterns based on the file hash.

## 4. Implementation Strategy
- **Language**: Rust (for zero-overhead background operation).
- **Distribution**: `brew install hcp`, `apt install hcp`, `cargo install hcp`.
- **Update Channel**: Signed auto-updates via the Bitcoin blockchain (Release hashes anchored in OP_RETURN).
