-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Create hypertable for traces if it doesn't exist
CREATE TABLE IF NOT EXISTS traces (
    id SERIAL,
    trace_id TEXT NOT NULL,
    span_id TEXT NOT NULL,
    parent_span_id TEXT,
    name TEXT,
    kind INTEGER,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    attributes JSONB,
    resource_attributes JSONB,
    service_name TEXT,
    status_code INTEGER,
    status_message TEXT,
    PRIMARY KEY (id, start_time)
);

-- Convert traces table to a hypertable if not already
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM timescaledb_information.hypertables 
        WHERE hypertable_name = 'traces'
    ) THEN
        PERFORM create_hypertable('traces', 'start_time');
    END IF;
END $$;

-- Create indexes for traces if they don't exist
CREATE INDEX IF NOT EXISTS idx_traces_trace_id ON traces(trace_id);
CREATE INDEX IF NOT EXISTS idx_traces_span_id ON traces(span_id);
CREATE INDEX IF NOT EXISTS idx_traces_service_name ON traces(service_name);

-- Create hypertable for logs if it doesn't exist
CREATE TABLE IF NOT EXISTS logs (
    id SERIAL,
    timestamp TIMESTAMPTZ NOT NULL,
    trace_id TEXT,
    span_id TEXT,
    severity TEXT,
    body TEXT,
    attributes JSONB,
    resource_attributes JSONB,
    service_name TEXT,
    PRIMARY KEY (id, timestamp)
);

-- Convert logs table to a hypertable if not already
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM timescaledb_information.hypertables 
        WHERE hypertable_name = 'logs'
    ) THEN
        PERFORM create_hypertable('logs', 'timestamp');
    END IF;
END $$;

-- Create indexes for logs if they don't exist
CREATE INDEX IF NOT EXISTS idx_logs_trace_id ON logs(trace_id);
CREATE INDEX IF NOT EXISTS idx_logs_span_id ON logs(span_id);
CREATE INDEX IF NOT EXISTS idx_logs_service_name ON logs(service_name);
CREATE INDEX IF NOT EXISTS idx_logs_severity ON logs(severity); 