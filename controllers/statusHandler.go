package controllers

import (
	"encoding/json"
	"net/http"
	"time"
)

type appStatus struct {
	Status      string        `json:"status"`
	Environment string        `json:"environment"`
	Version     string        `json:"version"`
	Uptime      time.Duration `json:"uptime"`
}

func (m *Repository) StatusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := appStatus{
		Status:      "Available",
		Environment: m.App.Config.Env,
		Version:     m.App.Version,
		Uptime:      time.Duration(time.Since(m.App.StartTime).Minutes()),
	}

	js, err := json.Marshal(currentStatus)
	if err != nil {
		m.App.Logger.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
