package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"

	"google.golang.org/api/option"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

func main() {
	ctx := context.Background()

	// Get the path to the service account key file from the GOOGLE_APPLICATION_CREDENTIALS environment variable
	keyPath, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if !ok {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS environment variable not set")
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		log.Fatalf("Service account key file not found at path: %s", keyPath)
	}

	// Get the path to the data folder from the DATA_PATH environment variable
	dataPath, ok := os.LookupEnv("DATA_PATH")
	if !ok {
		log.Fatal("DATA_PATH environment variable not set")
	}

	// Create a Text-to-Speech client using service account credentials
	client, err := texttospeech.NewClient(ctx, option.WithCredentialsFile(keyPath))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Set the input text and voice parameters
	req := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: "Hello, world!",
			},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			Name:         "en-US-Studio-O",
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_LINEAR16,
		},
	}

	// Synthesize speech from the input text
	resp, err := client.SynthesizeSpeech(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	// Save the output audio to a file in the data folder
	outputPath := fmt.Sprintf("%s/output.wav", dataPath)
	if err := ioutil.WriteFile(outputPath, resp.AudioContent, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Speech synthesized and saved to %s\n", outputPath)
}
