#!/usr/bin/env python3
import threading
import os
from http.server import BaseHTTPRequestHandler, HTTPServer
import cgi
import transcribe.transcribe as transcribe
import tempfile

class RequestHandler(BaseHTTPRequestHandler):
    def __init__(self, pool, *args, **kwargs):
        self.pool = pool
        super().__init__(*args, **kwargs)
    def do_POST(self):
        # Get content length from header
        content_length = int(self.headers['Content-Length'])

        # Get form data from the request body
        form = cgi.FieldStorage(fp=self.rfile, headers=self.headers, environ={'REQUEST_METHOD': 'POST'})

        # Get file data from the form data
        file_item = form['file']
        file_name = file_item.filename

        temp_file = tempfile.NamedTemporaryFile(delete=False)
        temp_file.write(file_item.file.read())
        temp_file.close()

        # Send a response to the client
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        self.wfile.write(bytes("Received the following form data and file:<br>File: " + file_name, 'utf-8'))
        #self.wfile.close()
        #pool = multiprocessing.Pool(processes=30)

        t = threading.Thread(target=transcribe.process_file,args=(temp_file.name,))
        t.start()


