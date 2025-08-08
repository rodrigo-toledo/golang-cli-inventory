# Active Context

Current work focus
- Initialize the project's Memory Bank to capture the high-level project context and make future sessions resumable.
- Provide concise, developer-focused documentation that mirrors the project's existing OVERVIEW and PROJECT_SPEC.

Recent changes (this session)
- Created the following memory-bank files:
  - `projectbrief.md` — project summary and goals.
  - `productContext.md` — users, core flows, UX goals.
  - `systemPatterns.md` — architecture, patterns, testing rules.
  - `techContext.md` — tech stack, tooling, common commands.

Important patterns & decisions
- Memory Bank mirrors key project documentation so that a fresh session can be resumed with minimal manual context loading.
- Core files follow the hierarchy described in the Memory Bank spec: `projectbrief.md` feeds into `productContext.md`, `systemPatterns.md`, and `techContext.md`.
- Keep the Memory Bank up-to-date whenever major decisions are made (architecture changes, new workflows, or CI/migration changes).

Next steps
1. Create `progress.md` describing current status and open tasks.
2. Encourage maintainers to update `activeContext.md` after implementing significant changes (new features, schema changes, or CI updates).
3. Optionally add:
   - Integration notes (how to run integration tests locally).
   - Feature-specific docs (e.g., report generation specifics) under `memory-bank/` as needed.

How to update
- Edit the relevant file(s) under `memory-bank/`.
- Include date and short rationale for major changes.
- Keep `activeContext.md` focused on the immediate next steps and decisions in flight.
