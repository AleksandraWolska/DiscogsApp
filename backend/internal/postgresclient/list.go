package postgresclient

import (
	"context"
	"fmt"
)

// its not needed, but I tried out things
func (pc *PostgresClient) ListTablesAndRelations(ctx context.Context) (map[string][]string, error) {
	query := `
        SELECT
            tc.table_name,
            kcu.column_name,
            ccu.table_name AS foreign_table_name,
            ccu.column_name AS foreign_column_name
        FROM
            information_schema.table_constraints AS tc
            JOIN information_schema.key_column_usage AS kcu
              ON tc.constraint_name = kcu.constraint_name
              AND tc.table_schema = kcu.table_schema
            JOIN information_schema.constraint_column_usage AS ccu
              ON ccu.constraint_name = tc.constraint_name
        WHERE tc.constraint_type = 'FOREIGN KEY';
    `

	rows, err := pc.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	relations := make(map[string][]string)
	for rows.Next() {
		var tableName, columnName, foreignTableName, foreignColumnName string
		if err := rows.Scan(&tableName, &columnName, &foreignTableName, &foreignColumnName); err != nil {
			return nil, err
		}
		relation := fmt.Sprintf("%s(%s) -> %s(%s)", tableName, columnName, foreignTableName, foreignColumnName)
		relations[tableName] = append(relations[tableName], relation)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return relations, nil
}
