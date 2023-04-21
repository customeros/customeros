import json
import datetime
import time


from flask import jsonify, request
import routes.routes as routes
from transcribe import summary


def handle_summary_post_request():
    error = routes.check_api_key()
    if error:
        return error

    current_time = time.time()

    try:
        if request.form.get('transcript') is not None:
            transcript = json.loads(request.form.get('transcript'))

        sum_content = summary.summarise(transcript)
    finally:
        print("Time taken: " + str(time.time() - current_time))
    return jsonify({
        'status': 'ok',
        'summary': sum_content
    }), 200