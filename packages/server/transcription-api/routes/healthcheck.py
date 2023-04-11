from flask import jsonify


def handle_health_get_request():
    return jsonify({
        'status': 'ok',
    }), 200