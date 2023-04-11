from flask import Flask, request, jsonify

import routes.transcription_service as transcription_service
import routes.healthcheck as healthcheck

app = Flask(__name__)

@app.route('/transcribe', methods=['POST'])
def handle_transcribe_post_request():
    return transcription_service.handle_transcribe_post_request()

@app.route('/health', methods=['GET'])
def handle_health_get_request():
    return healthcheck.handle_health_get_request()