package main

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestTuneWebServerServesPageAndGenerateAPI(t *testing.T) {
	server := NewTuneWebServer()
	port := freeTuneWebTestPort(t)
	if err := server.Start(port); err != nil {
		t.Fatalf("start tune web server: %v", err)
	}
	defer func() { _ = server.Stop() }()

	status := server.Status()
	if !status.Running || status.URL == "" {
		t.Fatalf("status = %#v", status)
	}

	client := http.Client{Timeout: 5 * time.Second}
	pageResp, err := client.Get(status.URL)
	if err != nil {
		t.Fatalf("get tune page: %v", err)
	}
	defer pageResp.Body.Close()
	if pageResp.StatusCode != http.StatusOK {
		t.Fatalf("page status = %d", pageResp.StatusCode)
	}

	body := []byte(`{"useCase":"Drift","pi":700,"drivetrain":"RWD","weightKG":1400,"frontWeightPct":54,"tireCompound":"drift","balanceBias":100,"stiffnessBias":100,"speedBias":100}`)
	apiResp, err := client.Post(strings.Replace(status.URL, "/tune", "/api/tune/generate", 1), "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("post generate: %v", err)
	}
	defer apiResp.Body.Close()
	if apiResp.StatusCode != http.StatusOK {
		t.Fatalf("generate status = %d", apiResp.StatusCode)
	}
	var payload struct {
		ProfileDraft struct {
			UseCase       string   `json:"useCase"`
			RearDiffAccel *float64 `json:"rearDiffAccel"`
		} `json:"profileDraft"`
	}
	if err := json.NewDecoder(apiResp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode generate response: %v", err)
	}
	if payload.ProfileDraft.UseCase != "Drift" || payload.ProfileDraft.RearDiffAccel == nil || *payload.ProfileDraft.RearDiffAccel != 100 {
		t.Fatalf("unexpected generate payload: %#v", payload)
	}
}

func TestTuneWebServerGenerateRejectsInvalidInput(t *testing.T) {
	server := NewTuneWebServer()
	port := freeTuneWebTestPort(t)
	if err := server.Start(port); err != nil {
		t.Fatalf("start tune web server: %v", err)
	}
	defer func() { _ = server.Stop() }()

	url := strings.Replace(server.Status().URL, "/tune", "/api/tune/generate", 1)
	resp, err := http.Post(url, "application/json", strings.NewReader(`{"useCase":"TimeAttack","pi":700,"drivetrain":"AWD","weightKG":1400,"frontWeightPct":54}`))
	if err != nil {
		t.Fatalf("post invalid generate: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", resp.StatusCode)
	}
}

func TestTuneWebServerPortInUse(t *testing.T) {
	address := preferredTuneWebAddress()
	listener, err := net.Listen("tcp", net.JoinHostPort(address, "0"))
	if err != nil {
		t.Fatalf("reserve port: %v", err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port

	server := NewTuneWebServer()
	err = server.Start(port)
	if err == nil {
		_ = server.Stop()
		t.Fatal("expected port-in-use error")
	}
	status := server.Status()
	if status.Running || status.LastError == "" {
		t.Fatalf("status should record start failure: %#v", status)
	}
}

func freeTuneWebTestPort(t *testing.T) int {
	t.Helper()
	address := preferredTuneWebAddress()
	listener, err := net.Listen("tcp", net.JoinHostPort(address, "0"))
	if err != nil {
		t.Fatalf("reserve free port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatalf("close free port listener: %v", err)
	}
	if port <= 0 {
		t.Fatalf("invalid test port %s", strconv.Itoa(port))
	}
	return port
}
