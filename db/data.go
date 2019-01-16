package db

import (
	"log"
)

type Data struct {
	Id   int
	Pgn  int
	Time uint64
	Blob []byte
}

func (d *db) WriteData(pgn int, blob []byte, time int64) {
	_, err := d.dbase.Exec("INSERT INTO data VALUES (?, ?, ?)", pgn, time, blob)

	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("Write OK")
	}
}

func (d *db) ReadData()[]Data {
	var outData []Data

	rows, err := d.dbase.Query("SELECT ROWID, pgn, time, data FROM data ORDER BY time LIMIT 1000")

	if err != nil {
		return nil
	}

	for rows.Next() {
		data := Data{}

		err := rows.Scan(&data.Id, &data.Pgn, &data.Time, &data.Blob)

		if err != nil {
			continue
		}

		outData = append(outData, data)
	}

	return outData
}

func (d *db) DeleteDataByID(id int) bool {
	_, err := d.dbase.Exec("DELETE FROM data WHERE ROWID = ?", id)

	return err == nil
}