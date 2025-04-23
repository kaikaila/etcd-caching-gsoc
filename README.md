# etcd-watchcache-prototype

This repository contains a prototype implementation for the [GSoC 2025 proposal](https://github.com/kaikaila/etcd-caching-gsoc/tree/main/docs/Proposal-Develop%20a%20caching%20library for%20etcd%20-%20YunkaiLi.pdf) titled **"Develop a Caching Library for etcd"**, focused on building a generic WatchCache proxy that can serve multiple downstream consumers such as Istio, Calico, and Cilium.

> ⚠️ **Note:** This is a mid-stage prototype. While the core components (e.g., `WatchCache`, `Snapshot`, `Compact`, `ClientSession`) have been designed and partially implemented, the interface boundaries and client-side abstractions are actively being refactored. Expect bugs and partial implementations.

---

## ✅ Current Progress

- ✅ Implemented `WatchCache` with memory-backed in-memory store
- ✅ Designed and tested `Snapshot()` and `Compact()` behaviors
- ✅ Integrated `EventLog` interface for replayable event pipeline
- ✅ Dockerized etcd setup and CLI harness for demo
- 🧩 WIP: Refactoring `ClientLibrary` interface and session isolation
- 🧩 WIP: Modular separation between `proxy`, `watcher`, and `client view`

---

## 📁 Project Structure

```plaintext
.
├── README.md                                                  # You're here
├── cmd/                                                       # Entry points and CLI demos
│   ├── demo/                                                  # Minimal usage demos (WIP)
│   └── proxy/                                                 # Main proxy startup logic
├── default.etcd/                                              # Local etcd volume mount
├── docs/                                                      # Architecture, roadmap, proposal, design docs
│   ├── Proposal-Develop a caching library for etcd - YunkaiLi.md   # Full GSoC proposal (markdown)
│   ├── Proposal-Develop a caching library for etcd - YunkaiLi.pdf  # Same proposal (PDF version)
│   ├── roadmap.md                                             # GSoC milestones and deliverables
│   ├── performance_decision.md                                # Trade-offs and performance notes
│   ├── clientli.md                                            # Client-side session design
│   ├── architecture.xmind                                     # High-level WatchCache architecture
│   ├── ClientLibrary Architecture.xmind                       # ClientSession internal structure
│   ├── file_structure_xmind.md                                # Codebase structure planning
│   └── watchcache_xmind.md                                    # Event flow / WatchCache logic
├── go.mod / go.sum                                            # Go module definition
├── pkg/                                                       # Core modules (under active development)
│   ├── adapter/                                               # Optional protocol-specific adapters
│   ├── api/                                                   # Shared interfaces and types
│   ├── clientlibrary/                                         # ClientSession logic and cache views
│   ├── eventlog/                                              # Append/ListSince/Compact abstraction
│   ├── proxy/                                                 # WatchCache implementation + API fanout
│   └── watcher/                                               # etcd Watcher and restart logic
├── run_etcd_docker.sh                                         # Dev script to launch etcd in Docker
└── stop_etcd_docker.sh                                        # Cleanup script
```

⸻

🧪 How to Try It

⚠️ Development is still ongoing, so components may not be fully wired.

# Step 1: Launch local etcd (Docker)

./run_etcd_docker.sh

# Step 2: Run a minimal proxy or demo (under cmd/)

cd cmd/proxy && go run main.go

⸻

📚 Docs & Design Notes
• 📄 GSoC Proposal (PDF-style)
• 🧠 Architecture Diagrams
• 🧱 Client Library Planning
• 🧭 Roadmap

⸻

## 🔍 Downstream Use Survey (Preliminary Thoughts)

While developing this prototype, I also began designing a downstream user survey — intended to help prioritize cache behaviors based on real-world needs from projects like Istio, Calico, and Cilium.

To be clear:

- This survey has **not been officially published** (pending GSoC selection outcome).
- It is only intended as an **early prototype of a design validation tool**.
- I welcome any feedback from community members on whether such an effort would be valuable or redundant.

📄 [Preview of Survey Form](https://docs.google.com/forms/d/e/1FAIpQLSfwlpMmI-Wjm08O5tw_MJFfKYq1o7DKBXXcOxtxlTmlf-c68g/viewform?usp=sharing)  
📊 [Preliminary Dashboard (fake data)](https://docs.google.com/document/d/1W8V0wBtyWBifgspzcviE17EAFLlzCjZAFBHSLx7KOcY/edit?usp=sharing)

This idea was inspired by the realization that not all downstream projects have the same consistency vs latency tradeoffs, and that making cache behavior configurable might depend on validating actual use cases.

---

🤝 Feedback & Collaboration

This project is intended to serve the broader etcd ecosystem, and is being actively refined based on community needs. If you’re from a project like Istio, Cilium, Calico, or others and have opinions about cache behavior, List-watch usage, or restart consistency, feedback is very welcome.

Please feel free to open issues or comment in etcd#19371, or reach out via Slack.

⸻

Author

Yunkai Li
MIMS @ UC Berkeley
GitHub: @kaikaila
Email: yunkai_li@berkeley.edu

Note: Commits prior to public release were made in a private repository context and may lack `--signoff`. All future contributions will follow the DCO signing standard.

```

```
