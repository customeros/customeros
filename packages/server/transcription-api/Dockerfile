# Base image
FROM python:3.9-slim-buster

# Set the working directory to /app
WORKDIR /app

# Copy the requirements file into the container
COPY requirements.txt .
RUN apt-get update && apt-get install -y ffmpeg && rm -rf /var/lib/apt/lists/*
# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of the application code into the container
COPY . .

# Set environment variables for Flask and gunicorn
ENV FLASK_APP=main.py
ENV GUNICORN_CMD_ARGS="--bind=0.0.0.0:8014 --workers=2 --threads=4"

# Expose port 8014 for the container
EXPOSE 8014

# Start the application using gunicorn
CMD ["gunicorn", "service.transcription_service:app"]
