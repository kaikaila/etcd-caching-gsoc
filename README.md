# Project Documentation Overview

This `docs/` directory contains the architectural and design documentation for the etcd caching project.

Unlike `doc.go` files in Go packages, which are intended to help developers understand how to use exported types, interfaces, and methods (i.e., they serve **code-level API consumers**), the Markdown files in this directory are aimed at **design reviewers, maintainers, and system architects**.

These documents explain:

- Why the system is structured the way it is
- What tradeoffs were made
- How core components interact
- The rationale behind performance and architectural decisions

Use this folder to collect all long-form design notes that describe _why_ and _how_ the system is built — not just _what_ it exposes via Go interfaces.
