# Memory Bank

The Memory Bank is a collection of concise, developer-focused documentation that captures the essential context of the project. It serves as a quick reference for developers to understand the project's purpose, architecture, and current state without having to read through all the source code.

## Purpose

- Provide a quick onboarding experience for new contributors
- Maintain a clear record of project decisions and context
- Enable resumable development sessions by capturing key information
- Mirror key project documentation for easy reference

## Files

- `projectbrief.md` - Project summary and goals
- `productContext.md` - Users, core flows, and UX goals
- `systemPatterns.md` - Architecture, design patterns, and technical rules
- `techContext.md` - Technical stack, tooling, and developer commands
- `activeContext.md` - Current focus, recent changes, and next steps
- `progress.md` - Progress status and open tasks

## For Maintainers

When making significant changes to the project, update the relevant memory-bank files:

- [ ] Update `activeContext.md` after major changes (schema migrations, CI updates, new features)
- [ ] Record progress in `progress.md` with short dated entries when milestones are reached
- [ ] Update technical context in `techContext.md` when tooling or stack changes
- [ ] Update system patterns in `systemPatterns.md` when architecture or design patterns change
- [ ] Update product context in `productContext.md` when users or core flows change