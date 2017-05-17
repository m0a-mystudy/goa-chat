package models

func AllRooms(db XODB, limit int) ([]*Room, error) {

	// sql query
	const sqlstr = `SELECT ` +
		`id, name, description, created ` +
		`FROM goa_chat.rooms LIMIT ?`
	// run query
	XOLog(sqlstr, limit)
	q, err := db.Query(sqlstr, limit)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	var res []*Room
	for q.Next() {
		r := Room{}
		err = q.Scan(&r.ID, &r.Name, &r.Description, &r.Created)
		if err != nil {
			return nil, err
		}
		res = append(res, &r)
	}
	return res, nil
}
