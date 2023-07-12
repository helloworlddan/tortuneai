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

func HitMe(prompt string, project string) (string, error) {
	if project == "" {
		return "", fmt.Errorf("no project ID set")
	}

	if prompt == "" {
		prompt = "tell me a terribe dad joke about tech"
	}

	ctx := context.Background()
	client, err := aiplatform.NewPredictionClient(
		ctx,
		option.WithEndpoint("us-central1-aiplatform.googleapis.com:443"),
	)
	if err != nil {
		return "", err
	}
	defer client.Close()

	parameters, err := structpb.NewValue(map[string]interface{}{
		"temperature":     0.2,
		"maxOutputTokens": 256,
		"topK":            40,
		"topP":            0.95,
	})
	if err != nil {
		return "", err
	}

	instances, err := structpb.NewValue(map[string]interface{}{
		"prompt": prompt,
	})

	endpoint := fmt.Sprintf(
		"projects/%s/locations/us-central1/publishers/google/models/text-bison",
		project,
	)

	req := &aiplatformpb.PredictRequest{
		Endpoint:   endpoint,
		Instances:  []*structpb.Value{instances},
		Parameters: parameters,
	}
	resp, err := client.Predict(ctx, req)
	if err != nil {
		return "", err
	}

	predictionFields := resp.Predictions[0].GetStructValue().AsMap()

	return predictionFields["content"].(string), nil
}
