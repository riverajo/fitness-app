---
description: Create a Pull Request for the current changes
---

1. **Branching**: Checkout `main`, pull latest, and create a new branch (`git checkout -b agent/topic-name`).

2. **Committing**: Stage all changes and commit with a conventional commit message.

3. **Pushing**: Push the new branch to origin.

4. **Creating PR**: Use the `tea` CLI to create the PR.
   ```bash
   tea pr create --title "Title" --description "Body"
   ```

5. **Auto-Merge**: Forgejo Actions/Tea doesn't support the same auto-merge flag universally. Just notify the user the PR is created.
