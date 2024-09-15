package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type Events struct {
	EventType string `json:"type"`
	Repo      struct {
		Name string `json:"name,omitempty"`
	} `json:"repo"`
	NumCommits struct {
		Commits int `json:"size,omitempty"`
	} `json:"payload"`
}

var rootCmd = &cobra.Command{
	Use:   "github-activity",
	Short: "github-Activity summarizes the user activity",
	Long:  `fast flexible way to access github-activity`,
	Run: func(cmd *cobra.Command, args []string) {
		url := "https://api.github.com/users/" + args[0] + "/events"

		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Authorization", "Bearer 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")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		if resp.StatusCode != 200 {
			fmt.Println("Request unsuccessful", resp.StatusCode)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		data := []Events{}
		json.Unmarshal(body, &data)
		Display(data[:3])
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Display(data []Events) {
	for _, v := range data {
		if v.EventType == "PushEvent" {
			if v.NumCommits.Commits == 1 {
				fmt.Printf("- Pushed %d commit to %s \n", v.NumCommits.Commits, v.Repo.Name)
			} else {
				fmt.Printf("- Pushed %d commits to %s \n", v.NumCommits.Commits, v.Repo.Name)
			}
		} else if v.EventType == "WatchEvent" {
			fmt.Printf("- Starred %s\n", v.Repo.Name)
		} else {
			fmt.Printf(" - %s %s\n", v.EventType, v.Repo.Name)
		}
	}
}
