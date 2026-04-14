# portwatch

> CLI tool to monitor and alert on unexpected open ports on a local or remote host.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

```bash
# Scan localhost and alert on any ports outside the expected set
portwatch --host localhost --expected 22,80,443

# Monitor a remote host on a schedule (every 60 seconds)
portwatch --host 192.168.1.10 --expected 22,8080 --interval 60

# Output results as JSON
portwatch --host localhost --expected 22,443 --format json
```

When an unexpected port is detected, `portwatch` prints an alert to stderr and exits with a non-zero status code, making it easy to integrate into scripts or CI pipelines.

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--host` | Target host to scan | `localhost` |
| `--expected` | Comma-separated list of allowed ports | _(none)_ |
| `--interval` | Repeat scan every N seconds (0 = run once) | `0` |
| `--format` | Output format: `text` or `json` | `text` |
| `--timeout` | Connection timeout per port in milliseconds | `500` |

---

## Example Output

```
[ALERT] Unexpected open port detected on localhost:
  → Port 3306 (mysql) is open but not in the expected list

[OK] Scan complete. 1 unexpected port(s) found.
```

---

## License

MIT © 2024 yourusername