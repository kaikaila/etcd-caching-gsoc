# GSoC Project Roadmap: etcd Caching Library

This document outlines the overall roadmap for the GSoC 2025 project "Develop a caching library for etcd".

The project is divided into several phases, grouped into a mainline ("Generic Proxy") and two optional branches ("Client Library" and "Adapter Integration").

---

## 🟩 Generic Proxy Mainline

| Phase | Name                      | Objective                                                     | Est. Duration  |
| ----- | ------------------------- | ------------------------------------------------------------- | -------------- |
| 0     | Initialization            | Set up local etcd + scaffold project structure                | 1 day ✅ done  |
| 1     | Etcd Watcher Integration  | Implement WatchKey with callback/channel support              | 2 days ✅ done |
| 2     | Cache Connection Layer    | Build `memoryCache`, wire WatchKey → Cache                    | 2 days ✅ done |
| 3     | WatchCache Layer          | Snapshot abstraction, dual-revision model, and compaction     | 4 days ✅ now  |
| 4     | Multi-Client Support      | Enable multiple clients + namespace isolation                 | 2 days         |
| 5     | Strategy Injection        | TTL / revision-based eviction strategies as pluggable modules | 2 days         |
| 6     | Testing and Delivery Prep | Unit tests, benchmark, documentation for upstream evaluation  | 3 days         |

---

## 🟦 Optional Branch: Client Library

| Phase | Name                | Objective                                                  | Est. Duration |
| ----- | ------------------- | ---------------------------------------------------------- | ------------- |
| A1    | Client API Skeleton | Define external interface for client consumers             | 2 days        |
| A2    | Query Interface     | Add Get/List with filters, consistent reads, field selects | 2 days        |

---

## 🟨 Optional Branch: Adapter Integration

| Phase | Name                    | Objective                                       | Est. Duration |
| ----- | ----------------------- | ----------------------------------------------- | ------------- |
| B1    | Adapter: Kubernetes     | Add adapter that mimics K8s-style WatchCache    | TBD           |
| B2    | Adapter: Istio / Cilium | Explore cache structure needs of downstream use | TBD           |

---

## 🔁 Notes

- This roadmap is milestone-driven and subject to adjustment based on community feedback and mentor review.
- The Generic Proxy must be completed to qualify for GSoC midterm evaluation.
- Bonus branches can be explored after core functionality is stable.
