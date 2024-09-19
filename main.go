package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/viper"
)

func main() {
	log.Println("Reading configuration...")
	readViperConfig(getConfigPath(getConfigDirectory()))
	log.Println("Reading configuration successful.")

	log.Println("Creating GraphQL client...")
	ctx := context.Background()
	client := graphql.NewClient("https://api-prod.omnivore.app/api/graphql", &http.Client{Transport: &authedTransport{wrapped: http.DefaultTransport}})
	log.Println("Creating GraphQL client successful.")

	var updatesSince *getUpdatesUpdatesSinceUpdatesSinceSuccess
	for condition := true; condition; condition = updatesSince.PageInfo.HasNextPage {
		after := "0"
		if updatesSince != nil {
			after = updatesSince.PageInfo.EndCursor
		}
		resp, err := getUpdates(ctx, client, viper.GetString("timestamp"), after)
		if err != nil {
			log.Fatalf("Error fetching updates: %v \n", err)
		}
		updatesSince = resp.UpdatesSince.(*getUpdatesUpdatesSinceUpdatesSinceSuccess)

		for _, edge := range updatesSince.Edges {
			log.Printf("Processing %s...\n", edge.Node.Url)
			body := ""
			for _, highlight := range edge.Node.Highlights {
				if highlight.Quote != "" {
					body += fmt.Sprintf("#+begin_quote\n%s\n#+end_quote\n\n", highlight.Quote)
				}
				if highlight.Annotation != "" {
					body += fmt.Sprintf("%s\n\n", highlight.Annotation)
				}
				body += "\n"
			}

			data := url.Values{}
			data.Set("template", viper.GetString("template_key"))
			data.Set("url", edge.Node.Url)
			data.Set("title", edge.Node.Title)
			data.Set("body", body)

			err := exec.Command("xdg-open", fmt.Sprintf("org-protocol://capture?%s", data.Encode())).Run()
			if err != nil {
				log.Fatalf("Error running xdg-open: %v \n", err)
			}
			log.Printf("Processing %s successful.\n", edge.Node.Url)
		}
	}

	log.Println("Updating timestamp...")
	viper.Set("timestamp", time.Now().Format(time.RFC3339))
	err := viper.WriteConfig()
	if err != nil {
		log.Fatalf("Error writing config file: %v \n", err)
	}
	log.Println("Updating timestamp successful.")
	log.Println("Synced Omnivore to Org Mode successfully.")
}

func getConfigDirectory() string {
	configRoot, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error finding config directory: %v \n", err)
	}

	configDirectory := filepath.Join(configRoot, "omnivore-org-roam")

	err = os.MkdirAll(configDirectory, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating config directory: %v", err)
	}

	return configDirectory
}

func getConfigPath(configDirectory string) string {
	viper.SetDefault("template_key", "o")
	viper.SetDefault("api_key", "")
	viper.SetDefault("timestamp", "0")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDirectory)

	configPath := filepath.Join(configDirectory, "config.yaml")
	return configPath
}

func readViperConfig(configPath string) {
	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			err := viper.WriteConfigAs(configPath)
			if err != nil {
				log.Fatalf("Error writing config file: %v", err)
			}

			err = viper.ReadInConfig()
			if err != nil {
				log.Fatalf("Error reading config file: %v", err)
			}
		default:
			log.Fatalf("Error reading config file: %v \n", err)
		}
	}

	err = viper.WriteConfig()
	if err != nil {
		log.Fatalf("Error writing config file: %v", err)
	}
}

type authedTransport struct {
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	key := viper.GetString("api_key")
	req.Header.Set("Authorization", key)
	return t.wrapped.RoundTrip(req)
}
