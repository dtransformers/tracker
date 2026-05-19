---
trigger: always_on
---

# Minimal-Oriented Development Rule

## Objective

The system MUST prioritize the smallest, safest, and least disruptive implementation capable of achieving the requested result.

The goal is to minimize:

- Code changes
- Architectural impact
- Dependency introduction
- Refactoring scope
- Cognitive complexity
- Runtime risk
- Migration surface
- Maintenance overhead

while still delivering:

- Correctness
- Readability
- Extensibility when justified
- Compatibility with the existing architecture

---

# Core Principles

## 1. Prefer Modification Over Replacement

Before introducing new modules, services, abstractions, or patterns, the system MUST evaluate whether the requirement can be fulfilled by:

- Extending existing code
- Reusing existing utilities
- Adding small isolated functions
- Improving current interfaces
- Injecting configuration instead of redesigning logic

Avoid rewriting working systems unless explicitly required.

---

## 2. Minimize Surface Area

Every implementation MUST attempt to:

- Touch the fewest files possible
- Modify the fewest lines possible
- Introduce the fewest new concepts possible
- Avoid unnecessary exports/public APIs
- Avoid unnecessary abstractions

Large-scale refactors MUST NOT be introduced for small feature requests.

---

## 3. Search for Existing Capabilities First

Before implementing anything new, the system MUST:

- Search the current codebase for similar logic
- Reuse existing helpers/services/patterns
- Verify whether the framework already provides the capability
- Prefer native/runtime features over external libraries

Do not introduce dependencies when the standard library or current stack already solves the problem adequately.

---

## 4. Delay Abstraction

The system MUST avoid premature abstraction.

Avoid introducing:

- Generic factories
- Complex interfaces
- Plugin systems
- Event buses
- Middleware layers
- Multi-strategy architectures

unless:

- The requirement already demands them
- There are at least multiple proven use cases
- Existing duplication creates real maintenance cost

---

## 5. Prefer Localized Changes

Changes SHOULD remain as close as possible to the affected domain.

Examples:

- Validation changes stay in validation layer
- Transport changes stay in transport layer
- Database changes stay in repository layer
- UI changes stay in UI layer

Avoid cascading modifications across unrelated modules.

---

## 6. Preserve Existing Architecture

The implementation MUST align with:

- Existing project structure
- Existing naming conventions
- Existing dependency injection patterns
- Existing error handling conventions
- Existing logging strategy
- Existing testing approach

Do not introduce new architectural paradigms unless explicitly requested.

---

## 7. Minimize Dependency Growth

Before adding a dependency:

- Verify if existing dependencies already provide the capability
- Verify if the standard library is sufficient
- Evaluate bundle/runtime/build impact
- Evaluate security and maintenance implications

A new dependency MUST be justified.

---

## 8. Prefer Explicitness Over Cleverness

The implementation SHOULD prioritize:

- Single Responsibility Principle
- Readability
- Predictability
- Maintainability

Avoid:

- Meta-programming
- Excessive generics
- Hidden side effects
- Magic automation
- Over-engineered optimizations

---

## 9. Maintain Backward Compatibility

Changes MUST avoid:

- Breaking public contracts
- Changing existing API behavior unexpectedly
- Modifying database schemas unnecessarily
- Renaming existing interfaces without reason

Backward compatibility is preferred unless breaking change is explicitly approved.

---

## 10. Evaluate the Cost of Every Change

Before implementation, the system MUST evaluate:

- Is this change necessary?
- Is there a smaller alternative?
- Can the same result be achieved with configuration?
- Can existing logic be extended safely?
- Is the abstraction justified now?

The simplest valid solution SHOULD be preferred.

---

---

## 11. Test Only The affected module

We should not re-test every module everytime, the test commands are allowed only for the modules affected with the changes

---

# Implementation Decision Order

The system MUST evaluate solutions in the following order:

1. Configuration change
2. Existing code extension
3. Small isolated function
4. Local refactor
5. Shared utility extraction
6. New module/service
7. Architectural redesign

Higher-cost approaches require justification.

---

# Anti-Patterns

The system MUST avoid:

- Refactoring unrelated code during feature implementation
- Introducing frameworks for small problems
- Creating abstractions for hypothetical future needs
- Rewriting stable logic without measurable benefit
- Expanding scope beyond the requested requirement
- Large formatting-only commits mixed with logic changes
- Creating helper layers with single usage

---

# Expected Outcome

Every implementation SHOULD aim for:

- Small diff size
- High clarity
- Low regression risk
- Fast reviewability
- Fast rollback capability
- Predictable behavior
- Minimal operational impact

The best implementation is the smallest implementation that correctly solves the problem.
