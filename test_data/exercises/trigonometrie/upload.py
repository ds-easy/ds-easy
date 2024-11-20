import os
import requests

# Directory containing the PDF files
directory = '.'  # Change this to the directory containing your PDF files

# URL to send the POST request to
url = 'http://localhost:8080/exercises'

# Fields to include in the multipart body
lesson_name = 'Trigonometrie'
uploaded_by = 'admin@dseasy.com'

print("UPLOADING")
# Loop through each file in the directory
for filename in os.listdir(directory):
    if filename.endswith('.typ'):
        print(filename)
        # Full path to the PDF file
        file_path = os.path.join(directory, filename)

        # Fields for the multipart body
        files = {
            'exo_file': (filename, open(file_path, 'rb')),
        }
        data = {
            'exercise_name': filename,
            'lesson_name': lesson_name,
            'uploadedBy': uploaded_by,
        }

        # Send the POST request
        response = requests.post(url, files=files, data=data)

        # Print the response status code and content
        print(f'Response for {filename}: {response.status_code}')
        print(response.content)

        # Close the file
        files['exo_file'][1].close()
