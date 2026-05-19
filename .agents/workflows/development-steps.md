---
description: AGENT EXECUTION PROTOCOL
---

# AGENT EXECUTION PROTOCOL

You are an autonomous AI coding agent. You MUST execute every user request following this strict, step-by-step operational pipeline. You are required to use your available tools (File Reader, File Writer, Terminal, MCP) to complete these steps. Do NOT skip any phases.

## PHASE 1: Discovery & Scope Validation

_Objective: Understand constraints before writing code._

1. **Scan Scope:** Use your file reading tool to scan `Root README.md` to validate the project boundaries.
2. **Query Knowledge Base:** Search `FAQ.md` or use the Notion MCP to search for existing architectural rules or similar past problems.
3. **SCOPE GATE:** If the user request is outside the scope of `README.md`, you MUST halt and ask the user for clarification or permission to expand the scope. Do not assume implicit expansion.

## PHASE 2: Code Implementation

_Objective: Write structurally aligned code._

1. **Apply Changes:** Write or modify the code to fulfill the user request.
2. **Align Architecture:** Ensure new code fits existing domain, service, or repository boundaries. Do not introduce new architectural paradigms without explicitly asking the user first.

## PHASE 3: Verification

\*Objective: Ensure system integrity.

1. **Execute Tests:** You MUST use your terminal execution tool to run the relevant test suite (e.g., `npm test`, `go test ./...`, or `cargo test`).
2. **TEST GATE:** If the test suite fails, you must debug and fix the code. Do NOT proceed to Phase 4 until you receive a successful exit code from the test runner.

## PHASE 4: State Synchronization (MANDATORY)

_Objective: Update the system of record. You MUST use your file writing tool for these steps._

1. **Append to History:** Open `./HISTORY.md` and append a new log entry. You must explicitly save this file. Use this exact format:
   - **Timestamp:** [Insert current date/time]
   - **Request:** [Summary of user prompt]
   - **Changes:** [Technical breakdown of modules/files edited]
   - **Reasoning:** [Why you chose this approach]
   - **Impact:** [How system behavior changed]
2. **Sync Docs:** Update any relevant inline documentation, API docs, or `README.md` sections that were affected by your code changes.
3. **Update FAQ:** If you discovered a new edge case, or if the Notion MCP lacked an answer, use your file writing tool to add a new Q/A entry to `./commands/FAQ.md`.

## PHASE 5: Completion Handoff

_Objective: Report status to the user._

1. Respond to the user confirming the task is complete.
2. Explicitly state: "I have updated HISTORY.md and run the test suite."
