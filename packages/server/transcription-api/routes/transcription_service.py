import time

from flask import request, jsonify
import threading
import tempfile
import json
import os
import subprocess
from pydub import AudioSegment
import service.customer_os_api_client as customer_os_api_client
import service.file_store_api_client as file_store_api_client
from service.vcon_service import VConPublisher, Analysis, VConAnalysisType

import transcribe.transcribe as transcribe
import transcribe.summary as summary
from model.vcon import VConParty, VConEncoder


def make_transcript(raw_transcript, meeting_id):
    vcon_transcript = []
    for line in raw_transcript:
        vcon_transcript.append({
            'party': VConParty(name=line['speaker']),
            'text': line['text'],
            'file_id': line['file_id']
        })
    return {"file_id": meeting_id, "transcript": vcon_transcript}



def process_file(filename, participants, topic, vcon_api:VConPublisher, fs_api:file_store_api_client.FileStoreApiClient):
    print("Processing file " + filename)
    current_time = time.time()

    meeting_recording = fs_api.upload_file(filename)
    if 'error' in meeting_recording:
        print(f"Error uploading file to file store: {meeting_recording}")
        return
    if 'id' not in meeting_recording:
        print("Error uploading file to file store: no id returned")
        return
    meeting_id = meeting_recording['id']
    print("File uploaded to file store with id " + meeting_id)

    file_suffix = os.path.splitext(filename)[1]
    if file_suffix == '.mp4':
        print("Movie file detected, converting to MP3")
        new_file = os.path.splitext(filename)[0] + ".mp3"
        ret = subprocess.run(["ffmpeg", "-i", filename, "-acodec", "mp3",  new_file])
        if ret.returncode != 0:
            os.unlink(filename)
            print("*** Error converting movie file to MP3")
            return
        os.unlink(filename)
        filename = new_file

    try:
        readahead_buffer_size = 1000 * 60 * 5
        mp3_file = AudioSegment.from_file(filename, read_ahead_limit=readahead_buffer_size, format="mp3")
        print("File loaded in " + str(time.time() - current_time) + " seconds")
        diarisation = transcribe.diarise(mp3_file)

        print(diarisation)
        organizations = {}
        for participant in participants:
            if 'organizations' in participant:
                for org in participant['organizations']:
                    organizations[org['id']] = org

        industries = []
        descriptions = []
        for org in organizations.values():
            industries.append(org['industry'])
            descriptions.append(org['description'])
        transcript = transcribe.transcribe(mp3_file, diarisation, participants=[t['firstName'] + " " + t['lastName'] for t in participants],
                                           industries=industries, descriptions=descriptions ,topic=topic, fs_api=fs_api)


        print(transcript)
        openline_transcript = make_transcript(transcript, meeting_id)
        transcript_attachments = [f['file_id'] for f in openline_transcript['transcript'] if 'file_id' in f]
        transcript_attachments.append(meeting_id)
        vcon_api.publish_analysis(Analysis(content_type="application/x-openline-transcript-v2", content=json.dumps(openline_transcript, cls=VConEncoder), type=VConAnalysisType.TRANSCRIPT), attachments=transcript_attachments)
        sum_content = summary.summarise(transcript)
        print(sum_content)
        vcon_api.publish_analysis(Analysis(content_type="text/plain", content=sum_content, type=VConAnalysisType.SUMMARY))
    finally:
        print("Time taken: " + str(time.time() - current_time))
        os.unlink(filename)

def check_api_key():
    if request.headers.get('X-Openline-API-KEY') is None or os.environ.get('TRANSCRIPTION_KEY') != request.headers.get('X-OPENLINE-API-KEY'):
        return jsonify({
            'status': 'error',
            'message': 'Invalid API key'
        }), 401
    return None

def handle_transcribe_post_request():
    # Check the API key
    error = check_api_key()
    if error:
        return error

    # Get file data from the request
    file_item = request.files['file']
    file_name = file_item.filename

    file_suffix = os.path.splitext(file_name)[1]

    users = []
    contacts = []

    if request.form.get('users') is not None:
        users = json.loads(request.form.get('users'))
    if request.form.get('contacts') is not None:
        contacts = json.loads(request.form.get('contacts'))

    topic = request.form.get('topic')

    cos_api = customer_os_api_client.CustomerOsApiCient(os.environ.get('CUSTOMER_OS_API_URL'), os.environ.get('CUSTOMER_OS_API_KEY'), request.headers.get('X-Openline-USERNAME'))

    participants = []
    for user in users:
        info = cos_api.get_user(user)
        participants.append(info)

    for contact in contacts:
        info = cos_api.get_contact(contact)
        participants.append(info)

    print(participants)
    temp_file = tempfile.NamedTemporaryFile(delete=False, suffix=file_suffix)
    temp_file.write(file_item.read())
    temp_file.close()

    file_to_process = temp_file.name


    print("Users: " + str(users))
    print("Contacts: " + str(contacts))

    parties = [VConParty(user_id=u) for u in users] + [VConParty(contact_id=c) for c in contacts]
    print("Parties: " + str(parties))

    vcon_api = VConPublisher(os.environ.get('VCON_API_URL'), os.environ.get('VCON_API_KEY'), request.headers.get('X-Openline-USERNAME'), parties)
    fs_api = file_store_api_client.FileStoreApiClient(os.environ.get('FILE_STORE_API_URL'), os.environ.get('FILE_STORE_API_KEY'), request.headers.get('X-Openline-USERNAME'))

    # Start a new thread to process the file
    t = threading.Thread(target=process_file, args=(file_to_process, participants, topic, vcon_api, fs_api))
    t.start()

    # Send a JSON response to the client
    return jsonify({
        'status': 'success',
        'message': f'Received file: {file_name}'
    })




