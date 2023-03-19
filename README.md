# GCP TTS Scripts

This is a set of scripts for generating voice using the GCP Studio Voice, which is designed to sound very similar to a human voice.

## Prerequisites
To use these scripts, you will need:
* A GCP account
* A project with the Text-to-Speech API enabled. You can enable the API by running the following command:
    ``` bash
    PROJECT_ID="your_project_id"
    # Enable Text-to-Speech API
    gcloud services enable texttospeech.googleapis.com --project=$PROJECT_ID --quiet
    ```

* A service account with access to the Text-to-Speech API
* The service account credentials file downloaded and stored in a safe location. You will need to set the GOOGLE_APPLICATION_CREDENTIALS environmental variable to the path of the credentials file.
* Go installed.


## Running the script
To use the script:
1. Copy the text that you want to convert to speech into a .txt files and place them in the /data folder or any other folder. Make sure that the following conditions are met:
    * All sentences must end with a period. 
    * The Text-to-Speech API cannot process lists, so remove any list marks.
    * Remove any non-text content such as emojis.
2. Run the Go script which will read the text files and send them in chunks to GCP. The script will try to break the text down into chunks, ensuring that each sentence is included in a single chunk. To run the script, execute the following commands: 

    ```bash
    export GOOGLE_APPLICATION_CREDENTIALS={path to credentials file}
    export DATA_PATH=data  
    go cmd/main.go
    ```
    The script will loop over all .txt files in the specified folder and generate a separate .wav file for each chunk of text.

If you need to regenerate some files, simply delete the corresponding .wav files and run the script again.