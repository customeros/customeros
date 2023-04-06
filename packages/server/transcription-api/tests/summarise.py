import unittest
import transcribe.summary as summary
import json

class MyTestCase(unittest.TestCase):
    def test_summary(self):
        with open('data/transcription.json', 'r') as file:
            transcript = json.load(file)
        result = summary.summarise(transcript)
        print(result)
        self.assertEqual(True, False)  # add assertion here


if __name__ == '__main__':
    unittest.main()
