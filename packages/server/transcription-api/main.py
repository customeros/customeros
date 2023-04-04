#!/usr/bin/env python3
import multiprocessing
import service.transcription_service as transcription_service
from http.server import HTTPServer

NUM_SEGMENTS_PARALLEL = 30;

if __name__ == '__main__':
    # Create an HTTP server and listen for POST requests
    pool = multiprocessing.Pool(processes=NUM_SEGMENTS_PARALLEL)
    pool.daemon = False

    print("Starting server on port 8014")
    server_address = ('', 8014)
    httpd = HTTPServer(server_address, lambda *args, **kwargs: transcription_service.RequestHandler(pool, *args, **kwargs))
    httpd.serve_forever()

