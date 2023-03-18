# gcp-tts-scripts

This is a set of scripts for generating voice from files with GCP stidio voice (looks very similar to human)

## Prerequisites
* GCP account 
* project with service account created, Text To Speech API enabled
``` bash
PROJECT_ID="tts-20230318"
# Enable Text-to-Speech API
gcloud services enable texttospeech.googleapis.com --project=$PROJECT_ID --quiet

```

* service account should have an access to Text To Speech
* service credentials file should be exposed to GOOGLE_APPLICATION_CREDENTIALS environmental variable
* Go installed



## Step 1 - prepare file texts 

Copy text into .txt file and place it into /data folder.  
* all sentences must finish with dot. 
* It doesn't understand lists , so remove any list marks
* remove other non-text - emojis, etc.

once files are ready - 
run GO script which would read those files and send in chunks to GCP
'''bash 
go cmd/main.go
'''
