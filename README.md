# gcp-tts-scripts

This is a scripts for using 

## Prerequisites
* GCP account 
* GCloud cli installed

## Step 1 - preparing project, enabling API
Check you logged in in cli
``` bash
 gcloud auth list
 ```

Set variables with project name, output bucket name
``` bash
PROJECT_ID="tts-20230318"
REGION="us-central1"
BUCKET_NAME="tts-20230318-out"

# Create new project
gcloud projects create $PROJECT_ID --set-as-default --quiet

# Please go to the Cloud Console and attach a billing account to this project.

# Enable Text-to-Speech API
gcloud services enable texttospeech.googleapis.com --project=$PROJECT_ID --quiet

# Create Cloud Storage bucket for output
gsutil mb -l $REGION gs://$BUCKET_NAME
```


## Step 2 - prepare file texts 

Copy text into .txt file and place it into /data folder.  
* all sentences must finish with dot. 
* It doesn't understand lists , so remove any list marks
* remove othe non-text - emojis, etc.

once files are ready - 
run python script which would prepare JSON objects for 
'''bash 
python .\scripts\text2json.py
'''

Troubleshoot the file at https://cloud.google.com/text-to-speech if required

## Step 3 - start encoding

gcloud text-to-speech synthesize \
  --voice "en-US-Studio-O" \
  --input "text/plain" \
  --input-file "data/data_31.json" \
  --output-file "data_31.mp3"