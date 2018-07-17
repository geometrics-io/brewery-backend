package postgres

import (
	"encoding/json"

	brewery "github.com/antschmidt/brewery-backend"
)

type ReportsService struct {
	client *Client
}

func (rs *ReportsService) Monthlies() ([]*brewery.MonthlyReport, error) {
	err := rs.client.Open()
	if err != nil {
		return nil, err
	}
	defer rs.client.db.Close()
	var monthlies []*brewery.MonthlyReport

	rows, err := rs.client.db.Query("SELECT monthlies FROM monthlies_report;")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var report brewery.MonthlyReport
		var reportBytes []byte

		err = rows.Scan(&reportBytes)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(reportBytes, &report)
		if err != nil {
			return nil, err
		}
		monthlies = append(monthlies, &report)
	}
	return monthlies, nil
}
