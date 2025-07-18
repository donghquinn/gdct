package gdct

import (
	"testing"
)

func TestBuildSelect(t *testing.T) {
	tests := []struct {
		name     string
		dbType   DBType
		table    string
		columns  []string
		expected string
		hasError bool
	}{
		{
			name:     "Simple select all",
			dbType:   PostgreSQL,
			table:    "users",
			columns:  nil,
			expected: "SELECT * FROM users",
		},
		{
			name:     "Select specific columns",
			dbType:   PostgreSQL,
			table:    "users",
			columns:  []string{"id", "name", "email"},
			expected: "SELECT id, name, email FROM users",
		},
		{
			name:     "Empty table name should error",
			dbType:   PostgreSQL,
			table:    "",
			columns:  nil,
			hasError: true,
		},
		{
			name:     "Invalid DB type should error",
			dbType:   DBType("invalid"),
			table:    "users",
			columns:  nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := BuildSelect(tt.dbType, tt.table, tt.columns...)

			if tt.hasError {
				if qb.err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if qb.err != nil {
				t.Errorf("Unexpected error: %v", qb.err)
				return
			}

			query, _, err := qb.Build()
			if err != nil {
				t.Errorf("Build error: %v", err)
				return
			}

			if query != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, query)
			}
		})
	}
}

func TestQueryBuilderWhere(t *testing.T) {
	tests := []struct {
		name     string
		dbType   DBType
		table    string
		where    string
		args     []interface{}
		expected string
		hasError bool
	}{
		{
			name:     "PostgreSQL where with placeholder",
			dbType:   PostgreSQL,
			table:    "users",
			where:    "age > ?",
			args:     []interface{}{18},
			expected: "SELECT * FROM users WHERE age > $1",
		},
		{
			name:     "MySQL where with placeholder",
			dbType:   MariaDB,
			table:    "users",
			where:    "age > ?",
			args:     []interface{}{18},
			expected: "SELECT * FROM users WHERE age > ?",
		},
		{
			name:     "Empty condition should error",
			dbType:   PostgreSQL,
			table:    "users",
			where:    "",
			args:     nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := BuildSelect(tt.dbType, tt.table).Where(tt.where, tt.args...)

			if tt.hasError {
				if qb.err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if qb.err != nil {
				t.Errorf("Unexpected error: %v", qb.err)
				return
			}

			query, args, err := qb.Build()
			if err != nil {
				t.Errorf("Build error: %v", err)
				return
			}

			if query != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, query)
			}

			if len(args) != len(tt.args) {
				t.Errorf("Expected %d args, got %d", len(tt.args), len(args))
			}
		})
	}
}

func TestQueryBuilderInsert(t *testing.T) {
	data := map[string]interface{}{
		"name":  "John",
		"email": "john@example.com",
		"age":   30,
	}

	tests := []struct {
		name     string
		dbType   DBType
		table    string
		data     map[string]interface{}
		hasError bool
	}{
		{
			name:   "PostgreSQL insert",
			dbType: PostgreSQL,
			table:  "users",
			data:   data,
		},
		{
			name:   "MySQL insert",
			dbType: MariaDB,
			table:  "users",
			data:   data,
		},
		{
			name:     "Insert without data should error",
			dbType:   PostgreSQL,
			table:    "users",
			data:     nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := BuildInsert(tt.dbType, tt.table)
			if tt.data != nil {
				qb = qb.Values(tt.data)
			}

			query, args, err := qb.Build()

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if query == "" {
				t.Errorf("Query should not be empty")
			}

			if len(args) != len(tt.data) {
				t.Errorf("Expected %d args, got %d", len(tt.data), len(args))
			}
		})
	}
}

func TestDBTypeValidation(t *testing.T) {
	validTypes := []DBType{PostgreSQL, MariaDB, Mysql, Sqlite}
	invalidTypes := []DBType{"invalid", "oracle", ""}

	for _, dbType := range validTypes {
		if !dbType.IsValid() {
			t.Errorf("DBType %s should be valid", dbType)
		}
	}

	for _, dbType := range invalidTypes {
		if dbType.IsValid() {
			t.Errorf("DBType %s should be invalid", dbType)
		}
	}
}

func BenchmarkBuildSelect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		qb := BuildSelect(PostgreSQL, "users", "id", "name", "email").
			Where("age > ?", 18).
			Where("status = ?", "active").
			OrderBy("created_at", "DESC", nil).
			Limit(10)

		_, _, err := qb.Build()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBuildComplexQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		qb := BuildSelect(PostgreSQL, "users u", "u.id", "u.name", "p.title").
			LeftJoin("posts p", "p.user_id = u.id").
			Where("u.age > ?", 18).
			Where("u.status = ?", "active").
			Where("p.published = ?", true).
			GroupBy("u.id", "u.name").
			Having("COUNT(p.id) > ?", 5).
			OrderBy("u.created_at", "DESC", nil).
			Limit(20).
			Offset(40)

		_, _, err := qb.Build()
		if err != nil {
			b.Fatal(err)
		}
	}
}
