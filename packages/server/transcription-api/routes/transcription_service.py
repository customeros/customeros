import time

from flask import Flask, request, jsonify
import threading
import tempfile
import json
import os
from pydub import AudioSegment
import service.customer_os_api as customer_os_api
from service.vcon_service import VConPublisher, Analysis, VConAnalysisType
import service.vcon_service

import transcribe.transcribe as transcribe
import transcribe.summary as summary
from model.vcon import VConParty, VConEncoder

app = Flask(__name__)

def make_transcript(raw_transcript):
    result = []
    for line in raw_transcript:
        result.append({
            'party': VConParty(name=line['speaker']),
            'text': line['text']
        })
    return result



def process_file(filename, participants, vcon_api:VConPublisher):
    print("Processing file " + filename)
    current_time = time.time()

    mp3_file = AudioSegment.from_file(filename, format="mp3")
    print("File loaded in " + str(time.time() - current_time) + " seconds")
    diarisation = transcribe.diarise(mp3_file)

    print(diarisation)
    transcript = transcribe.transcribe(mp3_file, diarisation, participants=[t['firstName'] + " " + t['lastName'] for t in participants],
                                       topic=None)


    print(transcript)
    openline_transcript = make_transcript(transcript)
    vcon_api.publish_analysis(Analysis(content_type="application/x-openline-transcript", content=json.dumps(openline_transcript, cls=VConEncoder), type=VConAnalysisType.TRANSCRIPT))
    sum = summary.summarise(transcript)
    print(sum)
    vcon_api.publish_analysis(Analysis(content_type="text/plain", content=sum, type=VConAnalysisType.SUMMARY))
    print("Time taken: " + str(time.time() - current_time))
    os.unlink(filename)

def check_api_key():
    if request.headers.get('X-Openline-API-KEY') is None or os.environ.get('TRANSCRIPTION_KEY') != request.headers.get('X-OPENLINE-API-KEY'):
        return jsonify({
            'status': 'error',
            'message': 'Invalid API key'
        }), 401
    return None
@app.route('/transcribe', methods=['POST'])
def handle_post_request():
    # Check the API key
    error = check_api_key()
    if error:
        return error

    # Get file data from the request
    file_item = request.files['file']
    file_name = file_item.filename

    users = []
    contacts = []

    if request.form.get('users') is not None:
        users = json.loads(request.form.get('users'))
    if request.form.get('contacts') is not None:
        contacts = json.loads(request.form.get('contacts'))

    cos_api = customer_os_api.CustomerOsApi(os.environ.get('CUSTOMER_OS_API_URL'), os.environ.get('CUSTOMER_OS_API_KEY'), request.headers.get('X-Openline-USERNAME'))

    participants = []
    for user in users:
        info = cos_api.get_user(user)
        participants.append(info)

    for contact in contacts:
        info = cos_api.get_contact(contact)
        participants.append(info)

    print(participants)
    temp_file = tempfile.NamedTemporaryFile(delete=False)
    temp_file.write(file_item.read())
    temp_file.close()


    print("Users: " + str(users))
    print("Contacts: " + str(contacts))

    parties = [VConParty(user_id=u) for u in users] + [VConParty(contact_id=c) for c in contacts]
    print("Parties: " + str(parties))

    vcon_api = VConPublisher(os.environ.get('VCON_API_URL'), os.environ.get('VCON_API_KEY'), request.headers.get('X-Openline-USERNAME'), parties)

    # Start a new thread to process the file
    t = threading.Thread(target=process_file, args=(temp_file.name, participants, vcon_api))
    t.start()

    # Send a JSON response to the client
    return jsonify({
        'status': 'success',
        'message': f'Received file: {file_name}'
    })




