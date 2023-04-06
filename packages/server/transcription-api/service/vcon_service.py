import uuid
import requests

from model.vcon import *

class Analysis:
    def __init__(self, type:VConAnalysisType, content_type:str, content:str):
        self.type = type
        self.content_type = content_type
        self.content = content

class VConPublisher:
    def __init__(self, url:str, api_key:str, openline_username:str, parties:[VConParty]):
        self.url = url
        self.api_key = api_key
        self.openline_username = openline_username
        self.uuid = None
        self.first = True
        self.parties = parties

    def publish_analysis(self, analysis:Analysis):
        print("Parties: " + str(self.parties))
        vcon = VCon(parties=self.parties)
        if self.first:
            self.uuid = str(uuid.uuid4())
            vcon.uuid = self.uuid
            self.first = False
        else:
            vcon.uuid = str(uuid.uuid4())
            vcon.appended = VConAppended(uuid=self.uuid)

        vcon.analysis = [VConAnalysis(type=analysis.type, mimetype=analysis.content_type, body=analysis.content)]
        headers = {
            "X-Openline-VCon-Api-Key": self.api_key,
            "X-Openline-USERNAME": self.openline_username,
            "Content-Type": "application/json"
        }
        print(headers)
        print(vcon.encode())
        response = requests.post(f"{self.url}/api/v1/vcon", headers=headers, data=vcon.encode())
        print(response.text)
