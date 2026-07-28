package seeders

import "fmt"

// Run seeds convenience demo data for local development. It is idempotent
// (safe to run on every startup) and intended to be called only when
// config.IsDevelopment() is true.
//
// Note: seeders that must run in every environment (e.g. RBAC roles/lookup
// tables) don't exist yet because that feature isn't implemented in this
// codebase. When it is, add a sibling file (e.g. role_seeder.go) and call it
// unconditionally from main.go alongside this one.
func Run() error {
	if err := seedDemoUsers(); err != nil {
		return fmt.Errorf("seed demo users: %w", err)
	}
	if err := seedDemoExamples(); err != nil {
		return fmt.Errorf("seed demo examples: %w", err)
	}
	return nil
}
