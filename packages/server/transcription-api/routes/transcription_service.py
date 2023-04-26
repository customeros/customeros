import time
import datetime

from flask import request, jsonify
import threading
import tempfile
import json
import os
import subprocess
from pydub import AudioSegment
import service.customer_os_api_client as customer_os_api_client
import service.file_store_api_client as file_store_api_client
from service.vcon_service import VConPublisher, Analysis, VConAnalysisType, VConDialog

import transcribe.transcribe as transcribe
import transcribe.summary as summary
import transcribe.action_items as action_items
import routes.routes as routes
from model.vcon import VConParty, VConEncoder, VConDialogType


def make_transcript(raw_transcript, start: datetime):
    vcon_transcript = []
    for line in raw_transcript:
        vcon_transcript.append({
            'party': VConParty(name=line['speaker']),
            'text': line['text'],
            'file_id': line['file_id'],
            'start': (start + datetime.timedelta(milliseconds=line['start'])).isoformat(),
            'duration': int((line['stop'] - line['start'])/1000)
        })
    return vcon_transcript



def process_file(filename, participants, topic, start:datetime, vcon_api:VConPublisher, fs_api:file_store_api_client.FileStoreApiClient):
    print("Processing file " + filename)
    current_time = time.time()

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
        openline_transcript = make_transcript(transcript, start)
        for line in openline_transcript:
            element = {'text': line['text'], 'party': line['party']}
            print(f'element: {element}')

            dialog = VConDialog(start=line['start'],  mimetype="x-openline-transcript-element", body=json.dumps(element, cls=VConEncoder), type=VConDialogType.TEXT, duration=line['duration'])

            print(line['start'] + " " + line['party'].name + ": " + line['text'])
            attachments = []
            if line['file_id'] is not None:
                attachments.append(line['file_id'])
            vcon_api.publish_vcon(dialog=dialog, attachments=attachments)

        sum_content = summary.summarise(transcript)
        print(sum_content)
        vcon_api.publish_vcon(analysis=Analysis(content_type="text/plain", content=sum_content, type=VConAnalysisType.SUMMARY))

        action_list = action_items.action_items(transcript)
        print(action_list)
        vcon_api.publish_vcon(analysis=Analysis(content_type="application/x-openline-action_items", content=json.dumps({"action_list": action_list}, cls=VConEncoder), type=VConAnalysisType.ACTION_ITEMS))

    finally:
        print("Time taken: " + str(time.time() - current_time))
        os.unlink(filename)


def handle_transcribe_post_request():
    # Check the API key
    error = routes.check_api_key()
    if error:
        return error

    start = datetime.datetime.now()

    if request.form.get("start") is not None:
        try:
            start = datetime.datetime.fromisoformat(request.form.get("start"))
        except ValueError:
            return jsonify({
                'status': 'error',
                'message': 'Invalid start time'
            }), 400


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

    if request.form.get('file_id') is not None:
        file_id = request.form.get('file_id')
        print("File ID: " + file_id)
        file_info = cos_api.get_attachment_info(file_id)
        if file_info is None:
            return jsonify({
                'status': 'error',
                'message': f'File {file_id} not found'
            }), 404
    else:
        return jsonify({
            'status': 'error',
            'message': 'No file ID provided'
        }), 400



    print("Users: " + str(users))
    print("Contacts: " + str(contacts))

    parties = [VConParty(user_id=u) for u in users] + [VConParty(contact_id=c) for c in contacts]
    print("Parties: " + str(parties))

    vcon_api = VConPublisher(os.environ.get('VCON_API_URL'), os.environ.get('VCON_API_KEY'), request.headers.get('X-Openline-USERNAME'), parties, type=request.form.get('type'), uuid=request.form.get('group_id'))
    fs_api = file_store_api_client.FileStoreApiClient(os.environ.get('FILE_STORE_API_URL'), os.environ.get('FILE_STORE_API_KEY'), request.headers.get('X-Openline-USERNAME'))

    file_to_process = fs_api.download_file(file_info)

    # Start a new thread to process the file
    thread = threading.Thread(target=process_file, args=(file_to_process, participants, topic, start, vcon_api, fs_api))
    thread.start()

    # Send a JSON response to the client
    return jsonify({
        'status': 'success',
        'message': f'Received file: {file_info["name"]}'
    })




