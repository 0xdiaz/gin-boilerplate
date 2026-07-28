package utils

import (
	"fmt"
	"gorm.io/gorm"
)

// Search creates a GORM scope for searching a specific field with LIKE pattern.
//
// Parameters:
//   - search: The search term to look for
//   - field: The database column name to search in (must be validated)
//
// Returns a GORM scope function that can be used with db.Scopes()
//
// Security: The field parameter should be validated against a whitelist
// to prevent SQL injection. Never pass user input directly as field.
func Search(search, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if search != "" && field != "" {
			// Use fmt.Sprintf to build the WHERE clause safely
			// field should be validated by caller against whitelist
			whereClause := fmt.Sprintf("%s LIKE ?", field)
			db = db.Where(whereClause, "%"+search+"%")
		}
		return db
	}
}
