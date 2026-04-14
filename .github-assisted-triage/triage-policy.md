# GitHub Issue Triage Policy

This document defines the rules and tone for the automated GitHub issue triage process.

## 1. General Principles
- **Be Professional and Helpful**: Maintain a polite, senior-engineer tone.
- **Accuracy over Speed**: Do not guess. If information is missing, ask for it clearly.
- **Kubernetes Style**: Use slash commands (e.g., `/kind`, `/area`, `/priority`, `/assign`) to manage issue state.

## 2. Labeling (Slash Commands)
When drafting a response, include appropriate slash commands on separate lines:
- **Avoid Duplication**: Do not include slash commands for labels that are already present on the issue (e.g., if `/sig node` or `/kind bug` are already applied, do not include them in your draft).
- **Kinds**: `/kind bug`, `/kind feature`, `/kind cleanup`, `/kind documentation`.
- **Priority**: `/priority critical-urgent`, `/priority important-soon`, `/priority important-longterm`.
- **Areas**: Identify the relevant area (e.g., `/area api`, `/area storage`) based on the issue content.

## 3. Response Structure
1. **Acknowledge**: Briefly summarize the reported issue or question to show understanding.
2. **Analysis**: Provide initial findings or ask clarifying questions.
3. **Action Items**: Clearly state what the next steps are (e.g., "Needs logs", "Assigned to @maintainer").
4. **Commands**: Place slash commands at the end of the comment.

## 4. Specific Rules for this Repository
### PR and Assignee Management
- **Linked PRs**: Check the status of any linked Pull Requests.
    - If a PR is **open**: Include a command to assign the issue to the PR author (e.g., `/assign @author`).
    - If a PR is **closed**: Do not assign. Instead, add a note to the author asking if they plan to continue (e.g., "@author, the PR seems to be closed. Do you plan to continue working on it?").
- **Inactive Assignees**: If an issue has been assigned for a significant period without recent activity, ping the assignee to confirm status (e.g., "@assignee, are you still working on it?").

### Evidence-Based Summaries
- **Avoid Subjective Language**: Do not characterize a PR or proposed fix with subjective terms like "balanced approach" or "optimal solution" unless you are quoting a specific review from a maintainer.
- **Reference Reviewer Feedback**: Summarize specific feedback from maintainers or reviewers found in the issue or PR comments rather than providing an independent evaluation.
