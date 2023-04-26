import uuid
import requests

from model.vcon import *

class Analysis:
    def __init__(self, type:VConAnalysisType, content_type:str, content:str):
        self.type = type
        self.content_type = content_type
        self.content = content

class VConPublisher:
    def __init__(self, url:str, api_key:str, openline_username:str, parties:[VConParty], type:str=None, uuid=None):
        self.url = url
        self.api_key = api_key
        self.openline_username = openline_username
        if uuid is not None:
            self.uuid = uuid
            self.first = False
        else:
            self.first = True
            self.uuid = None
        self.parties = parties
        self.type = type


    def publish_vcon(self, analysis:Analysis=None, attachments:[str]=None, dialog:VConDialog=None):
        print("Parties: " + str(self.parties))
        vcon = VCon(parties=self.parties)
        if self.first:
            self.uuid = str(uuid.uuid4())
            vcon.uuid = self.uuid
            self.first = False
        else:
            vcon.uuid = str(uuid.uuid4())
            vcon.appended = VConAppended(uuid=self.uuid)

        if self.type is not None:
            vcon.type = self.type

        if analysis is not None:
            vcon.analysis = [VConAnalysis(type=analysis.type, mimetype=analysis.content_type, body=analysis.content)]
        if dialog is not None:
            vcon.dialog = [dialog]
        if attachments is not None:
            vcon.attachments = [VConAttachment(mimetype="application/x-openline-file-store-id", body=attachment) for attachment in attachments]

        headers = {
            "X-Openline-VCon-Api-Key": self.api_key,
            "X-Openline-USERNAME": self.openline_username,
            "Content-Type": "application/json"
        }
        print(headers)
        print(vcon.encode())
        response = requests.post(f"{self.url}/vcon", headers=headers, data=vcon.encode())
        print(response.text)
