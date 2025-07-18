package gdct

import (
	"strings"
	"testing"
)

func TestSelectMethod(t *testing.T) {
	// Test the new Select method
	query, args, err := BuildSelect(PostgreSQL, "users", "id", "name").
		Select("email", "age").
		Where("active = ?", true).
		Build()

	if err != nil {
		t.Errorf("Select method failed: %v", err)
	}

	expected := "SELECT id, name, email, age FROM users WHERE active = $1"
	if query != expected {
		t.Errorf("Expected %q, got %q", expected, query)
	}

	if len(args) != 1 || args[0] != true {
		t.Errorf("Expected args [true], got %v", args)
	}
}

func TestOrWhereMethod(t *testing.T) {
	// Test the new OrWhere method
	query, args, err := BuildSelect(PostgreSQL, "users").
		Where("age > ?", 18).
		OrWhere("status = ?", "premium").
		Build()

	if err != nil {
		t.Errorf("OrWhere method failed: %v", err)
	}

	expected := "SELECT * FROM users WHERE (age > $1 OR status = $2)"
	if query != expected {
		t.Errorf("Expected %q, got %q", expected, query)
	}

	if len(args) != 2 || args[0] != 18 || args[1] != "premium" {
		t.Errorf("Expected args [18, 'premium'], got %v", args)
	}
}

func TestComplexQueryBuilding(t *testing.T) {
	// Test complex query with multiple methods
	query, args, err := BuildSelect(PostgreSQL, "users u").
		Select("u.name", "p.title").
		LeftJoin("posts p", "p.user_id = u.id").
		Where("u.age > ?", 18).
		Where("u.status = ?", "active").
		OrWhere("u.role = ?", "admin").
		GroupBy("u.id").
		Having("COUNT(p.id) > ?", 5).
		OrderBy("u.created_at", "DESC", nil).
		Limit(10).
		Offset(20).
		Build()

	if err != nil {
		t.Errorf("Complex query building failed: %v", err)
	}

	// Verify query structure (note: SELECT includes both * and additional columns)
	if !strings.Contains(query, "SELECT *, u.name, p.title FROM users u") {
		t.Errorf("Query missing SELECT clause: %s", query)
	}
	if !strings.Contains(query, "LEFT JOIN posts p ON p.user_id = u.id") {
		t.Errorf("Query missing JOIN clause: %s", query)
	}
	if !strings.Contains(query, "WHERE u.age > $1 AND (u.status = $2 OR u.role = $3)") {
		t.Errorf("Query missing WHERE clause: %s", query)
	}
	if !strings.Contains(query, "GROUP BY u.id") {
		t.Errorf("Query missing GROUP BY clause: %s", query)
	}
	if !strings.Contains(query, "HAVING COUNT(p.id) > $4") {
		t.Errorf("Query missing HAVING clause: %s", query)
	}
	if !strings.Contains(query, "ORDER BY u.created_at DESC") {
		t.Errorf("Query missing ORDER BY clause: %s", query)
	}
	if !strings.Contains(query, "LIMIT $5") {
		t.Errorf("Query missing LIMIT clause: %s", query)
	}
	if !strings.Contains(query, "OFFSET $6") {
		t.Errorf("Query missing OFFSET clause: %s", query)
	}

	// Verify args count
	if len(args) != 6 {
		t.Errorf("Expected 6 args, got %d: %v", len(args), args)
	}
}

func TestErrorHandling(t *testing.T) {
	// Test error accumulation
	qb := BuildSelect(PostgreSQL, "users").
		Where("", "invalid").  // Empty condition should cause error
		Select("name")         // Should not execute due to previous error

	query, args, err := qb.Build()
	if err == nil {
		t.Errorf("Expected error but got none")
	}
	if query != "" {
		t.Errorf("Expected empty query on error, got %q", query)
	}
	if len(args) != 0 {
		t.Errorf("Expected no args on error, got %v", args)
	}
}

func TestDataValidation(t *testing.T) {
	// Test INSERT with empty data
	_, _, err := BuildInsert(PostgreSQL, "users").
		Values(map[string]interface{}{}).
		Build()

	if err == nil {
		t.Errorf("Expected error for empty data")
	}

	// Test UPDATE with empty data
	_, _, err = BuildUpdate(PostgreSQL, "users").
		Set(map[string]interface{}{}).
		Where("id = ?", 1).
		Build()

	if err == nil {
		t.Errorf("Expected error for empty data")
	}
}

func TestDifferentDatabaseTypes(t *testing.T) {
	testCases := []struct {
		name   string
		dbType DBType
		expectedPlaceholder string
	}{
		{"PostgreSQL", PostgreSQL, "$1"},
		{"MariaDB", MariaDB, "?"},
		{"MySQL", Mysql, "?"},
		{"SQLite", Sqlite, "?"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, args, err := BuildSelect(tc.dbType, "users").
				Where("age > ?", 18).
				Build()

			if err != nil {
				t.Errorf("Query building failed for %s: %v", tc.name, err)
			}

			if !strings.Contains(query, tc.expectedPlaceholder) {
				t.Errorf("Expected placeholder %s in query for %s, got: %s", 
					tc.expectedPlaceholder, tc.name, query)
			}

			if len(args) != 1 || args[0] != 18 {
				t.Errorf("Expected args [18] for %s, got %v", tc.name, args)
			}
		})
	}
}

func BenchmarkNewSelectMethod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, err := BuildSelect(PostgreSQL, "users", "id").
			Select("name", "email").
			Where("age > ?", 18).
			Build()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOrWhereMethod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, err := BuildSelect(PostgreSQL, "users").
			Where("age > ?", 18).
			OrWhere("status = ?", "active").
			Build()
		if err != nil {
			b.Fatal(err)
		}
	}
}