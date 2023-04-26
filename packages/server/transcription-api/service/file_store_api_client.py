import tempfile

import requests

class FileStoreApiClient:
    def __init__(self, base_url:str, api_key:str, openline_username:str):
        self.base_url = base_url
        self.api_key = api_key
        self.openline_username = openline_username

    def upload_file(self, file_name:str):
        url = f"{self.base_url}/file"
        print(f"Uploading file {file_name} to {url} and username {self.openline_username}")

        headers = {
            "X-Openline-API-KEY": self.api_key,
            "X-Openline-USERNAME": self.openline_username,
            "Accept": "application/json"
        }

        print(f"Headers: {headers}")
        with open(file_name, 'rb') as f:
            form = {
                'file': (file_name, f)
            }
            response = requests.post(url, headers=headers, files=form)
        if response.status_code == 200:
            print("File uploaded successfully")
            result = response.json()
            return result

        return {'error': 'Unable to upload file', 'msg': response.text, 'status': response.status_code}

    def download_file(self, file_name:dict):
        url = f"{self.base_url}/file/{file_name['id']}/download"
        with tempfile.NamedTemporaryFile(delete=False, suffix="." + file_name['extension']) as temp_file:
            print(f"Downloading file {file_name} to {temp_file.name}")
            headers = {
                "X-Openline-API-KEY": self.api_key,
                "X-Openline-USERNAME": self.openline_username,
                "Accept": "application/json"
            }
            print(f"Headers: {headers}")
            response = requests.get(url, headers=headers, stream=True)
            if response.status_code == 200:
                for chunk in response.iter_content(chunk_size=8192):
                    temp_file.write(chunk)
                return temp_file.name
            else:
                print(f"Error downloading file: {response.text}")
            return None

