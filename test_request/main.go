package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/eunsukko/attendancebook-gcloud/models/team"
)

func main() {
	t := time.Now().Format(time.RFC3339)
	teamID := team.ID(3001)

	tmp := make(map[string]interface{})

	tmp["t"] = t
	tmp["team_id"] = teamID

	jsonData, _ := json.Marshal(tmp)

	fmt.Printf("before print\n")
	fmt.Printf("%s\n", jsonData)
}
