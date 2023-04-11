import unittest
import json

from pydub import AudioSegment

import transcribe.transcribe as transcribe


class MyTranscription(unittest.TestCase):
    def test_transcription(self):
        with open('data/diarisation.json', 'r') as file:
            diarisation = json.load(file)
        mp3_file = AudioSegment.from_file("/Users/torreysearle/Documents/vuy-wxso-sik (2022-06-27 02_12 GMT-7).mp3", format="mp3")

        parties = ['John Doe',
                   'Jane Doe',
                   'Joe Bloggs',
                   'Jane Bloggs']

        transcript = transcribe.transcribe(mp3_file, diarisation, participants=parties,
                                           industries=None, descriptions=None ,topic=None)
        print(transcript)
        self.assertEqual(True, False)  # add assertion here


if __name__ == '__main__':
    unittest.main()
