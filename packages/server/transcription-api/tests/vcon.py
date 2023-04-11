import unittest
from service.vcon_service import VConPublisher, Analysis
from model.vcon import *

class MyTestCase(unittest.TestCase):
    def test_vcon(self):
        parties = [VConParty(name='John Doe'),
                   VConParty(name='Jane Doe'),
                   VConParty(name='Joe Bloggs'),
                   VConParty(name='Jane Bloggs')]
        vcon_service = VConPublisher('https://api.openline.io/v1/vcon', 'API_KEY', 'OPENLINE_USERNAME', parties)
        vcon_service.publish_analysis(Analysis(VConAnalysisType.TRANSCRIPT, 'text/plain', 'Hello World'))

        vcon_service.publish_analysis(Analysis(VConAnalysisType.SUMMARY, 'text/plain', 'Summary of the conversation'))

        self.assertEqual(True, False)  # add assertion here


if __name__ == '__main__':
    unittest.main()
