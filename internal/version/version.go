package version

// Release Represents a single release version
type Release struct {
	Version string
	Date    string
	Changes []string
}

// CurrentVersion holds the current application version
const CurrentVersion = "v0.1.0"

// Changelog holds the history of changes
var Changelog = []Release{
	{
		Version: "v0.1.0",
		Date:    "18 Gen 2026",
		Changes: []string{
			"Implementazione delle funzioni di base CRUD con una GUI soddisfacente",
		},
	},
}
