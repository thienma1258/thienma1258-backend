package repository

import (
	"database/sql"
	"github.com/sirupsen/logrus"
)

func ScanRowIDs(rows sql.Rows) []int {
	result := make([]int, 0)
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			logrus.Error("error ScanRowIDs %v ", err)
		}
		result = append(result, id)
	}

	return result
}
