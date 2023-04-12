import replicate
import requests
import json
import os
import numpy
from pydub.utils import make_chunks
import datetime
from sklearn.metrics.pairwise import cosine_similarity
import concurrent.futures
from io import BytesIO
import traceback


AUDIO_SEGMENT_LENGTH = 5 * 60 * 1000  # 5 minutes in milliseconds
NUM_SEGMENTS_PARALLEL = 30
TASK_TIMEOUT = 10 * 60  # 10 minutes

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

    results = []

    # Run tasks asynchronously, with up to NUM_SEGMENTS_PARALLEL tasks running in parallel
    with concurrent.futures.ThreadPoolExecutor(max_workers=NUM_SEGMENTS_PARALLEL) as executor:
        future_to_result = {executor.submit(diarise_chunk, chunk): chunk for chunk in chunks}
        try:
            for future in concurrent.futures.as_completed(future_to_result, timeout=TASK_TIMEOUT):
                chunk = future_to_result[future]
                try:
                    result = future.result(timeout=TASK_TIMEOUT)
                    results.append(result)
                except Exception as exc:
                    print('%r generated an exception: %s' % (chunk, exc))
                    traceback.print_exc()
        except concurrent.futures.TimeoutError:
            print("Timed out waiting for diarise tasks to complete")

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

    print(sorted_diariastions[0])
    for segment in sorted_diariastions[0]['diariastion']['segments']:
        s = {}
        s['start'] = get_milliseconds(segment['start'])
        s['stop'] = get_milliseconds(segment['stop'])
        s['speaker'] = segment['speaker']
        print("Adding segment " + str(s['start']) + " to " + str(s['stop']))
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

def build_transcribe_prompt(participants, industries, descriptions, topic):
#    prompt = "Your Job is to transcribe an audio conversation into text.\n"
#    prompt += "(Only if Needed) Extra context is provided below.\n"
    prompt = "------------\n"

    prompt += "Description:\nThis is an audio conversation between " + " and ".join(participants) + "\n"

    if len(industries)  > 0:
        prompt += "The participants are from the following industries: " + ", ".join(industries) + "\n"


    if len(descriptions) > 0:
        prompt += "The participants are described as follows:\n"
        for description in descriptions:
            prompt += description + "\n"

    if topic:
             prompt += "The following is the topic of the discussion:\n" + topic + "\n"
    prompt += "------------\n\n"

    #prompt += "\nIf the context isn't useful return the original transcription\n"

    return prompt


def transcribe_audio_buffer(audio_segment, prompt, temperature):
    buffer = BytesIO()
    audio_segment.export(buffer, format="mp3")
    buffer.seek(0)
    segment_output = replicate.run(
        "openai/whisper:e39e354773466b955265e969568deb7da217804d8e771ea8c9cd0cef6591f8bc",
        input={"audio": buffer, "initial_prompt": prompt,
               "condition_on_previous_text": True,
               "compression_ratio_threshold": 2.4,
               "logprob_threshold": -1,
               "temperature": temperature, }
    )
    return segment_output

def transcribe_segment(segment):
    start = segment['start']
    stop = segment['stop']
    speaker = segment['speaker']
    audio_segment = segment['audio']
    prompt = segment['prompt']
    print(f"Transcribing segment {start} to {stop} for speaker {speaker}...")

    try:
        temperature = 0.0
        error_rate = 0.0
        while True:
            segment_output = transcribe_audio_buffer(audio_segment, prompt, temperature)

            errors = 0
            for chunk in segment_output['segments']:
                if chunk['avg_logprob'] < -1:
                    errors += 1
                    continue
                if chunk['compression_ratio'] > 2.4:
                    errors += 1
                    continue

            if len(segment_output['segments']) == 0:
                if prompt != "":
                    print(f"{start}: No output, retrying without prompt")
                    prompt = ""
                    continue
                else:
                    print(f"{start}: No output even without prompt, giving up")
                    break


            error_rate = float(errors) / float(len(segment_output['segments']))
            if error_rate < 0.25:
                break

            if prompt != "":
                print(f"{start}: Error rate too high {str(error_rate)}, retrying without prompt")
                prompt = ""
                continue
            else:
                print(f"{start}: Error rate too high {str(error_rate)} even without prompt, giving up")
                break



        if prompt == "" and error_rate < 0.25 and len(segment_output['segments']) > 0:
            print(f"*******{start}: Success by removing prompt: {str(temperature)}" )

        text = ""
        last_end = 0
        total_new_segment_outputs = []
        for chunk in segment_output['segments']:
            new_start = chunk['start'] * 1000
            new_end = chunk['end'] * 1000
            print( f"{start}: Chunk {new_start} to {new_end}: Gap: {new_start - last_end}")
            if new_start - last_end > 5000:
                print(f"{start}:**** gap detected")
                new_segment_output = transcribe_audio_buffer(audio_segment[last_end:new_start], "", temperature)
                for new_chunk in new_segment_output['segments']:
                    new_chunk['start'] = new_chunk['start'] + last_end / 1000
                    print( f"{start}: New chunk {new_chunk['start']} to {new_chunk['end']}: {new_chunk['text']}")
                    total_new_segment_outputs.append(new_chunk)
            last_end = new_end

        if len(total_new_segment_outputs) > 0:
            print(f"{start}: Merging {len(total_new_segment_outputs)} new chunks")
            segment_output['segments'] = segment_output['segments'] + total_new_segment_outputs
            segment_output['segments'] = sorted(segment_output['segments'], key=lambda k: k['start'])

        for chunk in segment_output['segments']:
            if chunk['avg_logprob'] < -1:
                print(f"{start}: Skipping chunk with low logprob: " + chunk['text'])
                continue
            if chunk['compression_ratio'] > 2.4:
                print(f"{start}: Skipping chunk with high compression ratio: " + chunk['text'])
                continue
            text += chunk['text']
            print(f"{start}: Chunk{chunk}")
    except Exception as e:
        print(f"Error transcribing segment {start} to {stop} for speaker {speaker}: error {e}")
        return None

    segment_info = {'speaker': speaker, 'text': text, 'start': start}
    return segment_info




def run_transcription_tasks(tasks):
    results = []
    # Run tasks asynchronously, with up to NUM_SEGMENTS_PARALLEL tasks running in parallel
    with concurrent.futures.ThreadPoolExecutor(max_workers=NUM_SEGMENTS_PARALLEL) as executor:
        future_to_task = {executor.submit(transcribe_segment, task): task for task in tasks}
        try:
            for future in concurrent.futures.as_completed(future_to_task, timeout=TASK_TIMEOUT):
                task = future_to_task[future]
                try:
                    data = future.result(timeout=TASK_TIMEOUT)
                except Exception as exc:
                    print('%r generated an exception: %s' % (task, exc))
                else:
                    results.append(data)
        except concurrent.futures.TimeoutError:
            print("Timed out waiting for transcription tasks to complete")

    return results


def transcribe(mp3_file, diarisation, participants=[], industries=[], descriptions=[], topic=""):
    tasks = []

    prompt = build_transcribe_prompt(participants, industries, descriptions, topic)
    print("Transcription prompt: " + prompt)

    for segment in diarisation['segments']:
        start = segment['start']
        stop = segment['stop']
        speaker = segment['speaker']
        tasks.append({'speaker': speaker, 'audio': mp3_file[start:stop], 'start': start, 'stop': stop,'prompt': prompt})

    result_transcript = list(filter(lambda x: isinstance(x, dict), run_transcription_tasks(tasks)))
    sorted_transcript = sorted(result_transcript, key=lambda k: k['start'])
    return sorted_transcript


