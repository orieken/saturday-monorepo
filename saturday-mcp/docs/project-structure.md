# Project Structure

```
saturday-mcp/
├── cmd/
│   └── saturday-mcp/
│       └── main.go
├── internal/
│   ├── server/
│   ├── generators/
│   ├── analyzers/
│   ├── validators/
│   ├── templates/
│   ├── models/
│   └── utils/
├── pkg/
├── templates/
├── docs/
├── test/
├── go.mod
├── Makefile
└── README.md
```

Notes:
- Follow clean architecture: no global state outside `main`
- Dependency injection for testability

