import os

from flask import Flask, request, jsonify
import routes.transcription_service as transcription_service
import routes.healthcheck as healthcheck
import routes.summary_service as summary_service
import routes.action_item_service as action_item_service


app = Flask(__name__)

def check_api_key():
    if request.headers.get('X-Openline-API-KEY') is None or os.environ.get('TRANSCRIPTION_KEY') != request.headers.get('X-OPENLINE-API-KEY'):
        return jsonify({
            'status': 'error',
            'message': 'Invalid API key'
        }), 401
    return None

@app.route('/transcribe', methods=['POST'])
def handle_transcribe_post_request():
    return transcription_service.handle_transcribe_post_request()

@app.route('/health', methods=['GET'])
def handle_health_get_request():
    return healthcheck.handle_health_get_request()

@app.route('/summary', methods=['POST'])
def handle_summary_post_request():
    return summary_service.handle_summary_post_request()

@app.route('/action-items', methods=['POST'])
def handle_action_items_post_request():
    return action_item_service.handle_action_item_post_request()