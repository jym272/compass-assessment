# `score app`

---

>Evaluate potential duplicates with a similarity score. 0 - 1000

---

### Start
Build and execute
```bash
go mod download
make build-app
./score.app -f input.xlsx -o out.xlsx
```

### Test
```bash
make test
```
### CI
You need `golangci-lint` installed
```bash
make all
```

### Docker utility container
To be implemented...