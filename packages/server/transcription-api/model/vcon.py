import json
from enum import Enum
from time import time
from typing import Any


class VConAppended:
    def __init__(self, uuid:str):
        self.uuid = uuid


class VConParty:
    def __init__(self, tel:str=None, stir:str=None, mailto:str=None, name:str=None, user_id:str=None, contact_id:str=None):
        if tel is not None:
            self.tel = tel
        if stir is not None:
            self.stir = stir
        if mailto is not None:
            self.mailto = mailto
        if name is not None:
            self.name = name
        if user_id is not None:
            self.user_id = user_id
        if contact_id is not None:
            self.contact_id = contact_id

class VConAnalysisType(Enum):
    SUMMARY = 'summary'
    TRANSCRIPT = 'transcript'
    TRANSLATION = 'translation'
    SENTIMENT = 'sentiment'
    TTS = 'tts'

class VConEncoding(Enum):
    NONE = 'none'
    BASE64 = 'base64'
    JSON = 'json'
class VConAnalysis:
    def __init__(self, type:VConAnalysisType, mimetype:str, body:str, encoding:VConEncoding=VConEncoding.NONE, dialog:[int]=[]):
        self.type = type
        self.mimetype = mimetype
        self.body = body
        self.encoding = encoding
        self.dialog = dialog

class VConDialogType(Enum):
    TEXT = 'text'
    RECORDING = 'recording'

class VConDialog:
    def __init__(self, type:VConDialogType, start:time, duration:int, mimetype:str, body:str, encoding:VConEncoding=VConEncoding.NONE, parties:[int]=[]):
        self.type = type
        self.start = start
        self.duration = duration
        self.mimetype = mimetype
        self.body = body
        self.encoding = encoding
        self.parties = parties

class VConAttachment:
    def __init__(self, mimetype:str, body:str, parties:[int]=[], encoding:VConEncoding=VConEncoding.NONE):
        self.mimetype = mimetype
        self.body = body
        self.parties = parties
        self.encoding = encoding

class VCon:

    def __init__(self, subject:str=None,  uuid:str=None, vcon:str='0.0.1', parties:[VConParty]=None, dialog:[VConDialog]=None, attachments:[VConAttachment]=None, analysis:[VConAnalysis]=None, appended:VConAppended=None):
        self.uuid = uuid
        self.parties = parties
        self.vcon = vcon

        if subject is not None:
            self.subject = subject
        if dialog is not None:
            self.dialog = dialog
        if attachments is not None:
            self.attachments = attachments
        if analysis is not None:
            self.analysis = analysis
        if appended is not None:
            self.appended = appended


    def encode(self):
        return json.dumps(self, cls=VConEncoder)


class VConEncoder(json.JSONEncoder):
    def default(self, o:VCon):
        if isinstance(o, VCon):
            return o.__dict__
        if isinstance(o, VConParty):
            return o.__dict__
        if isinstance(o, VConAnalysis):
            return o.__dict__
        if isinstance(o, VConDialog):
            return o.__dict__
        if isinstance(o, VConAttachment):
            return o.__dict__
        if isinstance(o, VConAppended):
            return o.__dict__
        if isinstance(o, VConAnalysisType):
            return o.value
        if isinstance(o, VConEncoding):
            return o.value
        if isinstance(o, VConDialogType):
            return o.value

        return super().default(o)