package tortuneai

import (
	"context"
	"fmt"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

const Version = "0.0.1"

// HitMe proxies a text prompt to Google Cloud's text-bison model hosted on
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
	// create Vertex AI client, but override with us-central1, due to
	// availability of hosted model
	client, err := aiplatform.NewPredictionClient(
		ctx,
		option.WithEndpoint("us-central1-aiplatform.googleapis.com:443"),
	)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// use default parameters, as found in the web console of Vertex AI
	parameters, err := structpb.NewValue(map[string]interface{}{
		"temperature":     0.2,
		"maxOutputTokens": 256,
		"topK":            40,
		"topP":            0.95,
	})
	if err != nil {
		return "", err
	}

	// this will only need the actual prompt
	instances, err := structpb.NewValue(map[string]interface{}{
		"prompt": prompt,
	})

	// construct the full resource name of the hosted model
	endpoint := fmt.Sprintf(
		"projects/%s/locations/us-central1/publishers/google/models/text-bison",
		project,
	)

	// instantiate the request object and execute it
	req := &aiplatformpb.PredictRequest{
		Endpoint:   endpoint,
		Instances:  []*structpb.Value{instances},
		Parameters: parameters,
	}
	resp, err := client.Predict(ctx, req)
	if err != nil {
		return "", err
	}

	// extract model response and return it
	predictionFields := resp.Predictions[0].GetStructValue().AsMap()
	return predictionFields["content"].(string), nil
}
