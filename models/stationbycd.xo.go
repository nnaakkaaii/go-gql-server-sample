package models

// Code generated by xo. DO NOT EDIT.

import (
	"context"
)

// StationByCD represents a row from 'public.station_by_c_d'.
type StationByCD struct {
	LineCd      int    `json:"line_cd"`      // line_cd
	LineName    string `json:"line_name"`    // line_name
	StationCd   int    `json:"station_cd"`   // station_cd
	StationGCd  int    `json:"station_g_cd"` // station_g_cd
	StationName string `json:"station_name"` // station_name
	Address     string `json:"address"`      // address
}

// StationByCDsByStationCD runs a custom query, returning results as StationByCD.
func StationByCDsByStationCD(ctx context.Context, db DB, stationCD int) ([]*StationByCD, error) {
	// query
	const sqlstr = `select l.line_cd, l.line_name, s.station_cd, station_g_cd, s.station_name, s.address ` +
		`from station s ` +
		`inner join line l on s.line_cd = l.line_cd ` +
		`where s.station_cd = $1 ` +
		`and s.e_status = 0`
	// run
	logf(sqlstr, stationCD)
	rows, err := db.QueryContext(ctx, sqlstr, stationCD)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// load results
	var res []*StationByCD
	for rows.Next() {
		var sbcd StationByCD
		// scan
		if err := rows.Scan(&sbcd.LineCd, &sbcd.LineName, &sbcd.StationCd, &sbcd.StationGCd, &sbcd.StationName, &sbcd.Address); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &sbcd)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}
