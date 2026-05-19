---
description: research and history of request
---

# WORKFLOW

the workflow rules of changes and user requests implementations :

## 1. Scope Compliance (Root README Gate)

Every code change MUST be validated against the project scope defined in `Root README.md`.

### Rules:

- The LLM must parse and interpret `Root README.md` before generating or modifying any code.
- Any new feature, module, or refactor must:
  - Explicitly map to an existing scope item, OR
  - Be justified as a scoped extension with rationale.

- If a change is outside scope:
  - It MUST be rejected or flagged for explicit user confirmation.
- No implicit expansion of project boundaries is allowed.

### MCP USAGE

- For every question you have about not included on the README.md you must ask using notion MCP, if not found you can add it to the FAQ.md with Q/A format.

### Enforcement Condition:

> If a change cannot be traced to a README-defined objective, it is invalid.

---

## 2. Structural Alignment (Architecture Consistency Rule)

All modifications MUST align with the existing repository structure.

### Rules:

- Follow current module boundaries, naming conventions, and layering (e.g., domain / service / repository / transport).
- Do NOT introduce new structural paradigms unless:
  - Existing architecture cannot support the requirement, AND
  - A migration plan is documented.

- Ensure:
  - No orphan modules
  - No duplicated responsibility layers
  - No bypassing of established abstractions

### Enforcement Condition:

> New code must “fit” the architecture, not reshape it by default.

---

## 3. Change History Enforcement (`./HISTORY.md`)

Every user-requested change MUST be recorded in:

```
./HISTORY.md
```

### Required Entry Format:

Each entry must include:

- Timestamp
- Request summary
- Implemented changes (technical breakdown)
- Reasoning / design decisions
- Impact analysis (what changed in system behavior)
- Linked files/modules affected

### Rules:

- No code change is valid without a HISTORY entry.
- Entries must be appended, never overwritten.
- Must be written in deterministic, structured format (not conversational).

### Enforcement Condition:

> HISTORY.md is the system of record for all evolution of the codebase.

---

## 4. Documentation Synchronization Rule

Any code change MUST be reflected in documentation.

### Rules:

- Update all relevant docs, including but not limited to:
  - README.md (if scope or usage changes)
  - module-level docs
  - API documentation
  - architecture notes

- Documentation must never lag behind implementation.
- Outdated documentation is treated as a critical defect.

### Enforcement Condition:

> Code and documentation must remain bidirectionally consistent.

---

## 5. Test Integrity Rule

All changes MUST preserve or improve the current test suite.

### Rules:

- Existing tests must continue to pass.
- New features require corresponding tests.
- If tests are broken:
  - Fix implementation first, not tests, unless tests are incorrect by design.

- No commit is valid if it degrades test coverage or correctness.

### Enforcement Condition:

> “Green tests” is a mandatory gate before change acceptance.

---

## 6. FAQ Knowledge Base Requirement

The system MUST maintain a living FAQ used by the LLM agent.

### Location:

```
./commands/FAQ.md
```

### Rules:

- FAQ must contain:
  - Known system behaviors
  - Common user questions
  - Architectural decisions and rationale
  - Operational constraints and deployment rules

- After each user request:
  - Add new discovered questions if they emerge
  - Or update existing answers if knowledge evolves

- FAQ must be structured for retrieval, not prose (Q/A format preferred)

### Agent Behavior Requirement:

- The LLM must consult FAQ BEFORE proposing solutions.
- If ambiguity exists, the LLM should:
  - Check FAQ first
  - Then Root README.md
  - Then codebase

### Enforcement Condition:

> FAQ is a first-class knowledge dependency, not optional documentation.

---

## 7. Global Consistency Rule (Meta Constraint)

All changes MUST satisfy all constraints simultaneously:

- README scope compliance
- Architectural alignment
- HISTORY.md traceability
- Documentation synchronization
- Test suite integrity
- FAQ consistency

If any constraint conflicts:

1. Halt generation
2. Report conflict explicitly
3. Request clarification or decision

---

## Summary Constraint Statement

> No backend change is valid unless it is scoped, structured, traceable, documented, tested, and knowledge-integrated.
