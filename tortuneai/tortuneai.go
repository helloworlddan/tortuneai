package tortuneai

import (
	"context"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
)

const Version = "0.0.4"

// HitMe proxies a text prompt to Google Cloud's gemini-pro model hosted on
// Vertex AI and returns the generated response from the model. If no prompt
// is passed to the function, it will default to a prompt that is engineered
// to generate a terrible dad joke.
func HitMe(prompt string, project string) (string, error) {
	if project == "" {
		return "", fmt.Errorf("no project ID set")
	}

	// override empty prompt if not specified by user
	if prompt == "" {
		prompt = "tell me a terrible dad joke about tech"
	}

	ctx := context.Background()
	region := "us-central1"

	client, err := genai.NewClient(
		ctx,
		project,
		region,
	)
	if err != nil {
		return "", err
	}

	model := client.GenerativeModel("gemini-2.0-flash-001")

	response, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", response.Candidates[0].Content.Parts[0]), nil
}
