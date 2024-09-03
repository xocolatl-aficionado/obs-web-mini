package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/rs/cors"
)

const obsPath = "/usr/bin/obs" // Adjust the path to match your system's OBS executable path

// Check if OBS is running
func checkOBS() (bool, error) {
	fmt.Println("Checking if OBS is running...")
	processes, err := process.Processes()
	if err != nil {
		return false, fmt.Errorf("error listing processes: %v", err)
	}

	for _, p := range processes {
		name, err := p.Name()
		if err == nil && name == "obs" { // Adjust "obs" if needed
			fmt.Println("OBS is running.")
			return true, nil
		}
	}
	fmt.Println("OBS is not running.")
	return false, nil
}

// Start OBS Studio
func startOBS() error {
	fmt.Println("Attempting to start OBS Studio...")
	cmd := exec.Command(obsPath)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("error starting OBS: %v", err)
	}

	// Wait for OBS to fully start
	time.Sleep(5 * time.Second) // Adjust the delay if needed
	return nil
}

// Ensure OBS is running
func ensureOBSRunning(w http.ResponseWriter, r *http.Request) {
	isRunning, err := checkOBS()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to check OBS status: %v", err), http.StatusInternalServerError)
		return
	}

	if !isRunning {
		fmt.Println("OBS Studio is not running, attempting to start it...")
		err = startOBS()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to start OBS: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Println("OBS Studio started successfully.")
	}

	// Respond with success JSON
	response := map[string]bool{"success": true}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Set up routes
	mux := http.NewServeMux()
	mux.HandleFunc("/start-obs", ensureOBSRunning)

	// Enable CORS
	handler := cors.Default().Handler(mux)

	// Start the server
	port := "3001"
	fmt.Printf("OBS controller backend running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
