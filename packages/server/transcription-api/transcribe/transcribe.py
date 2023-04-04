#!/usr/bin/env python3
import time

import replicate
import requests
import json
import os
import numpy
from pydub import AudioSegment
from pydub.utils import make_chunks
import datetime
from sklearn.metrics.pairwise import cosine_similarity
import multiprocessing
from io import BytesIO


AUDIO_SEGMENT_LENGTH = 5 * 60 * 1000  # 5 minutes in milliseconds
NUM_SEGMENTS_PARALLEL = 30

#The splitting of the audio is just for diarisation, not for transcription
#so the splitting doesn't need to be prettily on silence boundaries
#we will compensate for this in the merge_diarisation step
def split_audio_file(audio):
    chunks = make_chunks(audio, AUDIO_SEGMENT_LENGTH)
    result = []
    last_start = 0
    for chunk in chunks:
        result.append({'offset': last_start, 'audio': chunk})
        last_start += len(chunk)

    # Split the audio into chunks of 5 minutes (ish)
#    silences = silence.detect_silence(audio, min_silence_len=500, silence_thresh=-16)
#    for start, end in silences:
#        if start - last_start > AUDIO_SEGMENT_LENGTH:
#            chunk = audio[last_start:end]
#            chunks.append({'offset': last_start, 'audio': chunk})
#            last_start = end


    return result

def diarise_chunk(chunk):
    buffer = BytesIO()
    chunk['audio'].export(buffer, format="mp3")
    buffer.seek(0)
    print("calling diariastion")
    output = replicate.run(
        "meronym/speaker-diarization:64b78c82f74d78164b49178443c819445f5dca2c51c8ec374783d49382342119",
        input={"audio": buffer}
    )
    print(output)
    response = requests.get(output)
    if response.status_code == 200:
        output_json = json.loads(response.content)
        # Process the data as needed
    else:
        print(f"Error fetching data: {response.status_code}")
        os.exit(1)

    return {'offset': chunk['offset'], 'diariastion': output_json}
def diarise_chunks(chunks):
    # Run tasks asynchronously, with up to NUM_SEGMENTS_PARALLEL tasks running in parallel
    pool = multiprocessing.Pool(processes=NUM_SEGMENTS_PARALLEL)

    # Start the processes
    result_iterator = pool.imap_unordered(diarise_chunk, chunks, chunksize=1)

    # Wait for the processes to complete
    pool.close()
    pool.join()

    # Collect the results
    results = []
    for result in result_iterator:
        results.append(result)

    return results

def add_new_speaker(speakers, embedding):
    for letter in 'ABCDEFGHIJKLMNOPQRSTUVWXYZ':
        if letter not in speakers:
            speakers[letter] = embedding
            return letter

def map_speakers(transcript_speaker_dict, known_speaker_dict):
    speaker_mapping = {}
    for transcript_speaker, transcript_vector in transcript_speaker_dict.items():
        found = False
        best_similarity = 0
        for known_speaker, known_vector in known_speaker_dict.items():
            similarity = cosine_similarity(numpy.array(known_vector).reshape(1,-1),
                                           numpy.array(transcript_vector).reshape(1,-1))
            #print("Comparison between " + label + " and " + speaker['name'] + " is " + str(similarity))
            if similarity > 0.5:
                found = True
                if similarity > best_similarity:
                    best_similarity = similarity
                    speaker_mapping[transcript_speaker] = known_speaker
                    print("Speaker " + transcript_speaker + " is " + known_speaker + " confidence " + str(similarity))

        if not found:
            print("Speaker " + transcript_speaker + " not found")
            new_speaker = add_new_speaker(known_speaker_dict, transcript_vector)
            speaker_mapping[transcript_speaker] = new_speaker
            print("Speaker " + transcript_speaker + " mapped to " + new_speaker)
    return speaker_mapping
def merge_diarisations(diarisations):
    result = {"segments": [], "speakers": {"labels": [], "embeddings": {}}}
    sorted_diariastions = sorted(diarisations, key=lambda k: k['offset'])
    if len(sorted_diariastions) == 0:
        return []

    speakers =  sorted_diariastions[0]['diariastion']['speakers']['embeddings']

    for segment in sorted_diariastions[0]['diariastion']['segments']:
        s = {}
        s['start'] = get_milliseconds(segment['start'])
        s['stop'] = get_milliseconds(segment['stop'])
        s['speaker'] = segment['speaker']
        result['segments'].append(s)

    for diariastion in sorted_diariastions[1:]:
        print("Merging diariastion for offset " + str(diariastion['offset']))
        speaker_mapping = map_speakers(diariastion['diariastion']['speakers']['embeddings'], speakers)
        for segment in diariastion['diariastion']['segments']:
            segment['speaker'] = speaker_mapping[segment['speaker']]

        last_segment = result['segments'][-1]
        first = True
        for segment in diariastion['diariastion']['segments']:
            #if the first segment is the same speaker as last segment, merge them for better transcription
            if first and segment['speaker'] == last_segment['speaker']:
                last_segment['stop'] = get_milliseconds(segment['stop']) + diariastion['offset']
                print("Merging segment " + str(last_segment['start']) + " to " + str(last_segment['stop']))
            else:
                s = {}
                s['start'] = get_milliseconds(segment['start']) + diariastion['offset']
                s['stop'] = get_milliseconds(segment['stop']) + diariastion['offset']
                s['speaker'] = segment['speaker']
                result['segments'].append(s)
            first = False

    for speaker, vector in speakers.items():
        result['speakers']['embeddings'][speaker] = vector
        result['speakers']['labels'].append(speaker)

    result['speakers']['embeddings']['count'] = len(speakers)
    return result


def diarise(mp3_file):
    chunks = split_audio_file(mp3_file)
    print("Split into " + str(len(chunks)) + " chunks")
    diarisations = diarise_chunks(chunks)
    final_diarisation = merge_diarisations(diarisations)


    return final_diarisation
def get_milliseconds(timestamp_str):
    timestamp_obj = datetime.datetime.strptime(timestamp_str, '%H:%M:%S.%f')
    milliseconds = (timestamp_obj - datetime.datetime(1900, 1, 1)).total_seconds() * 1000.0
    return milliseconds

def transcribe_segment(segment):
    start = segment['start']
    stop = segment['stop']
    speaker = segment['speaker']
    audio_segment = segment['audio']
    print(f"Transcribing segment {start} to {stop} for speaker {speaker}...")
    buffer = BytesIO()
    audio_segment.export(buffer, format="mp3")
    buffer.seek(0)

    segment_output = replicate.run(
        "openai/whisper:e39e354773466b955265e969568deb7da217804d8e771ea8c9cd0cef6591f8bc",
        input={"audio": buffer}
    )
    text = ""
    for chunk in segment_output['segments']:
        text += chunk['text']
        print(chunk['text'])
    segment_info = {'speaker': speaker, 'text': text, 'start': start}

    return segment_info

def run_transcription_tasks(tasks):

    # Run tasks asynchronously, with up to NUM_SEGMENTS_PARALLEL tasks running in parallel
    pool = multiprocessing.Pool(processes=NUM_SEGMENTS_PARALLEL)

    # Start the processes
    result_iterator = pool.imap_unordered(transcribe_segment, tasks, chunksize=1)

    # Wait for the processes to complete
    pool.close()
    pool.join()

    # Collect the results
    results = []
    for result in result_iterator:
        results.append(result)


    return results


def transcribe(mp3_file, diarisation):
    tasks = []

    for segment in diarisation['segments']:
        start = segment['start']
        stop = segment['stop']
        speaker = segment['speaker']
        tasks.append({'speaker': speaker, 'audio': mp3_file[start:stop], 'start': start, 'stop': stop})

    result_transcript = run_transcription_tasks(tasks)
    sorted_transcript = sorted(result_transcript, key=lambda k: k['start'])
    return sorted_transcript

def process_file(filename):
    print("Processing file " + filename)
    current_time = time.time()

    mp3_file = AudioSegment.from_file(filename, format="mp3")
    print("File loaded in " + str(time.time() - current_time) + " seconds")
    diarisation = diarise(mp3_file)

    print(diarisation)
    transcript = transcribe(mp3_file, diarisation)
    print(transcript)
    with open("result.json", "w") as file:
        json.dump(transcript, file)
    print("Time taken: " + str(time.time() - current_time))
    os.unlink(filename)

