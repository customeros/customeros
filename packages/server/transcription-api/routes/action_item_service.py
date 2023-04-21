import json
import datetime
import time


from flask import jsonify, request
import routes.routes as routes
from transcribe import action_items


def handle_action_item_post_request():
    error = routes.check_api_key()
    if error:
        return error

    current_time = time.time()

    try:
        if request.form.get('transcript') is not None:
            transcript = json.loads(request.form.get('transcript'))

        action_item_list = action_items.action_items(transcript)
    finally:
        print("Time taken: " + str(time.time() - current_time))
    return jsonify({
        'status': 'ok',
        'action_items': action_item_list
    }), 200