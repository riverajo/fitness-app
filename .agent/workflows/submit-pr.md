---
description: Create a Pull Request for the current changes
---

1. **Branching**: Checkout `main`, pull latest, and create a new branch (`git checkout -b agent/topic-name`).

2. **Committing**: Stage all changes and commit with a conventional commit message.

3. **Pushing**: Push the new branch to origin.

4. **Creating PR**: Use the `create_pull_request` tool with a descriptive title and body.

5. **Auto-Merge**: Attempt to enable auto-merge using `gh pr merge --auto --squash`. If this fails, just notify the user that the PR is ready for review.
