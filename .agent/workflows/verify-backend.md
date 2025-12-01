---
description: Verify backend code quality
---

1. Run unit tests
   - Action: Run `go test ./...`
   - Directory: `backend`

2. Run linter
   - Action: Run `golangci-lint run ./...`
   - Directory: `backend`

3. If any fail, FIX them before notifying the user.
