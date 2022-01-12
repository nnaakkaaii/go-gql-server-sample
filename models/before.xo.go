package models

// Code generated by xo. DO NOT EDIT.

import (
	"context"
)

// Before represents a row from 'public.before'.
type Before struct {
	LineCd               int    `json:"line_cd"`                // line_cd
	LineName             string `json:"line_name"`              // line_name
	StationCd            int    `json:"station_cd"`             // station_cd
	StationName          string `json:"station_name"`           // station_name
	Address              string `json:"address"`                // address
	BeforeStationCd      int    `json:"before_station_cd"`      // before_station_cd
	BeforeStationName    string `json:"before_station_name"`    // before_station_name
	BeforeStationGCd     int    `json:"before_station_g_cd"`    // before_station_g_cd
	BeforeStationAddress string `json:"before_station_address"` // before_station_address
}

// BeforesByStationCD runs a custom query, returning results as Before.
func BeforesByStationCD(ctx context.Context, db DB, stationCD int) ([]*Before, error) {
	// query
	const sqlstr = `select sl.line_cd, ` +
		`sl.line_name, ` +
		`s.station_cd, ` +
		`s.station_name, ` +
		`s.address, ` +
		`COALESCE(js.station_cd, 0)    as before_station_cd, ` +
		`COALESCE(js.station_name, '') as before_station_name, ` +
		`COALESCE(js.station_g_cd, 0)  as before_station_g_cd, ` +
		`COALESCE(js.address, '')      as before_station_address ` +
		`from station s ` +
		`left outer join line sl on s.line_cd = sl.line_cd ` +
		`left outer join station_join j on s.line_cd = j.line_cd and s.station_cd = j.station_cd1 ` +
		`left outer join station js on j.station_cd2 = js.station_cd ` +
		`where s.e_status = 0 ` +
		`and s.station_cd = $1`
	// run
	logf(sqlstr, stationCD)
	rows, err := db.QueryContext(ctx, sqlstr, stationCD)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// load results
	var res []*Before
	for rows.Next() {
		var b Before
		// scan
		if err := rows.Scan(&b.LineCd, &b.LineName, &b.StationCd, &b.StationName, &b.Address, &b.BeforeStationCd, &b.BeforeStationName, &b.BeforeStationGCd, &b.BeforeStationAddress); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &b)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}
