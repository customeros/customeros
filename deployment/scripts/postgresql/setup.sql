-- auto-generated definition
CREATE TABLE IF NOT EXISTS app_keys
(
    id     bigserial
    primary key,
    app_id varchar(255) not null,
    key    varchar(255) not null
    constraint idx_app_keys_key
    unique,
    active boolean      not null
    );

ALTER TABLE app_keys
    owner TO postgres;

CREATE UNIQUE INDEX idx_key
    ON app_keys (key);

INSERT INTO app_keys (app_id, key, active) VALUES ('customer-os-api', 'dd9d2474-b4a9-4799-b96f-73cd0a2917e4', true) ON CONFLICT DO NOTHING;
INSERT INTO app_keys (app_id, key, active) VALUES ('customer-os-webhooks', '25744e24-8a89-4f5d-aae4-d48ccfdbe1d6', true) ON CONFLICT DO NOTHING;
INSERT INTO app_keys (app_id, key, active) VALUES ('file-store-api', '9eb87aa2-75e7-45b2-a1e6-53ed297d0ba8', true) ON CONFLICT DO NOTHING;
INSERT INTO app_keys (app_id, key, active) VALUES ('settings-api', '8b010f38-e5ca-4923-a62e-9f073c5c7dbf', true) ON CONFLICT DO NOTHING;
INSERT INTO app_keys (app_id, key, active) VALUES ('oasis-api', '10a6747a-97cd-4a6c-bcf5-e4ee89a12567', true) ON CONFLICT DO NOTHING;
INSERT INTO app_keys (app_id, key, active) VALUES ('validation-api', 'c41192c4-08bc-4b4e-851c-10fc876b9ebb', true) ON CONFLICT DO NOTHING;
INSERT INTO app_keys (app_id, key, active) VALUES ('anthropic-api', '01bfb7a3-917b-4af6-87f9-2bdfb73bea97', true) ON CONFLICT DO NOTHING;
INSERT INTO app_keys (app_id, key, active) VALUES ('openai-api', 'b155a222-2a0b-11ee-be56-0242ac120002', true) ON CONFLICT DO NOTHING;
