package controller

import (
	"bytes"
	"encoding/json"
	"github.com/Meromen/mfc-backend/db"
	"io/ioutil"
	"net/http"
)

func (c controller) GetMfcs(w http.ResponseWriter, r *http.Request) {
	mfcs, err := c.mfcStorage.SelectAll()
	if err != nil {
		c.logger.Errorf(
			"Failed to get mfcs from database: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := response{}
	if len(mfcs) == 0 {
		res.Status = "ok"
		res.Code = http.StatusNotFound
		res.Body = "mfcs not found"
	} else {
		res.Status = "ok"
		res.Code = http.StatusOK
		res.Body = mfcs
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		c.logger.Errorf(
			"Failed to encode response to json: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c controller) UpdateMfcs() {
	client := http.Client{}

	resp, err := client.Get("http://mfc-25.ru/queue/statistics")
	if err != nil {
		c.logger.Errorf("Failed to GET mfc statistics: %v\n", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))

	mfcs := make([]db.Mfc, 0)
	err = json.Unmarshal(body, &mfcs)
	if err != nil {
		c.logger.Errorf("Failed to unmarshal body: %v\n", err)
		return
	}

	rows := make([]db.DBRow, 0)
	for _, mfc := range mfcs {
		rows = append(rows, &mfc)
	}

	err = c.mfcStorage.UpdateAll(rows)
	if err != nil {
		c.logger.Errorf("Failed to update statistics in database: %v\n", err)
	}
}
