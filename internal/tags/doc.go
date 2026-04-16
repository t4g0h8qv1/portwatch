// Package tags maps port numbers to human-readable labels for use in
// portwatch output, alerts, and reports.
//
// A Registry is initialised with common well-known service names (ssh, http,
// postgres, …) and can be extended at runtime via Set or by loading a JSON
// file with LoadFile.
//
// Example JSON tag file:
//
//	{
//	  "9200": "elasticsearch",
//	  "5601": "kibana"
//	}
package tags
