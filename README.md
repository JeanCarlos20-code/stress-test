## 🛠️ Build and Run (Locally)

```bash
# Clone the repository
git clone https://github.com/JeanCarlos20-code/stress-test.git
cd stress-test

# Install dependencies
go mod tidy

# Run Build
go build -o stress-test cmd/main.go
./stress-test --url={string} --requests={number} --concurrency={number}

# Run Docker
docker build -t stress-test .
docker run stress-test --url={string} --requests={number} --concurrency={number}
```
