import os
import json

# Set path to folder containing .txt files
folder_path = "../data"
bucket_name = os.environ.get("BUCKET_NAME","tts-20230318-out")

# Get path to script directory
script_dir = os.path.dirname(os.path.abspath(__file__))
# Load template JSON
with open(os.path.join(script_dir, "template.json")) as f:
    template = json.load(f)

# Loop over all .txt files in folder
for file_name in os.listdir(folder_path):
    if file_name.endswith(".txt"):
        # Extract filename and text from file
        with open(os.path.join(folder_path, file_name), "r") as f:
            file_text = f.read()
        # Replace new line symbols with space
        file_text = file_text.replace("\n", " ")
        # Replace text in template JSON and save to new file
        json_name = os.path.splitext(file_name)[0] + ".json"
        json_path = os.path.join(folder_path, json_name)
        template["input"]["text"] = file_text
        template["output_gcs_uri"] = "gs://" + bucket_name + "/" + os.path.splitext(file_name)[0] + ".wav"
        with open(json_path, "w") as f:
            json.dump(template, f, indent=2)