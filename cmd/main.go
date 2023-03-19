package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"google.golang.org/api/option"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

const maxChunkSize = 1000

func splitTextIntoChunks(text string) []string {
	var chunks []string
	paragraphs := strings.Split(text, "\n")
	for _, paragraph := range paragraphs {
		if len(paragraph) <= maxChunkSize {
			if len(chunks) == 0 {
				chunks = append(chunks, paragraph)
			} else if len(chunks[len(chunks)-1])+len(paragraph)+2 <= maxChunkSize {
				chunks[len(chunks)-1] += "\n" + paragraph
			} else {
				chunks = append(chunks, paragraph)
			}
		} else {
			sentences := strings.Split(paragraph, ". ")
			var temp string
			for _, sentence := range sentences {
				if len(temp)+len(sentence)+2 > maxChunkSize {
					if len(chunks) == 0 {
						chunks = append(chunks, temp+".")
					} else if len(chunks[len(chunks)-1])+len(temp)+2 <= maxChunkSize {
						chunks[len(chunks)-1] += "\n" + temp + "."
					} else {
						chunks = append(chunks, temp+".")
					}
					temp = sentence
				} else {
					if temp == "" {
						temp = sentence
					} else {
						temp += ". " + sentence
					}
				}
			}
			if len(chunks) == 0 {
				chunks = append(chunks, temp+".")
			} else if len(chunks[len(chunks)-1])+len(temp)+2 <= maxChunkSize {
				chunks[len(chunks)-1] += "\n" + temp + "."
			} else {
				chunks = append(chunks, temp+".")
			}
		}
	}
	return chunks
}

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

	// Loop over all .txt files in the data folder
	err = filepath.Walk(dataPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".txt" {
			return nil
		}

		fileName := info.Name()[:len(info.Name())-4]

		// Read the text from the file
		textBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		text := string(textBytes)

		chunks := splitTextIntoChunks(text)
		// Generate speech for each chunk
		for i, chunk := range chunks {
			// Check if output file already exists
			audioPath := fmt.Sprintf("%s/%s_%d.wav", dataPath, fileName, i+1)
			if _, err := os.Stat(audioPath); err == nil {
				// File already exists, skip this chunk
				fmt.Printf("Skipping file %s chunk %d - output file %s already exists\n", fileName, i+1, audioPath)
				continue
			}

			// Set the input text and voice parameters
			req := &texttospeechpb.SynthesizeSpeechRequest{
				Input: &texttospeechpb.SynthesisInput{
					InputSource: &texttospeechpb.SynthesisInput_Text{
						Text: chunk,
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

			// Perform the text-to-speech request
			resp, err := client.SynthesizeSpeech(ctx, req)
			if err != nil {
				return fmt.Errorf("failed to synthesize speech: %v", err)
			}

			// Save the output to a file
			err = ioutil.WriteFile(audioPath, resp.AudioContent, 0644)
			if err != nil {
				return fmt.Errorf("failed to save output to file: %v", err)
			}
			fmt.Printf("Saved file %s chunk %d to file %s\n", fileName, i+1, audioPath)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
