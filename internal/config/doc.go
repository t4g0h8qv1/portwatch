// Package config provides loading and validation of portwatch configuration
// files written in YAML.
//
// A minimal configuration requires at least a target host and a port
// expression understood by the portrange package:
//
//	# portwatch.yaml
//	target: localhost
//	ports: "22,80,443,8000-9000"
//	baseline: baseline.json
//
// Optional sections:
//
//	alert:
//	  stdout: true          # print changes to stdout (default)
//	  webhook: "https://…"  # future: POST alert payload to a URL
//
//	scan:
//	  timeout_ms: 500       # per-port TCP dial timeout (default 500 ms)
//	  concurrency: 100      # parallel probes (default 100)
package config
