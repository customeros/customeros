INSERT INTO app_keys (app_id, key, active) VALUES ('customer-os-api', 'dd9d2474-b4a9-4799-b96f-73cd0a2917e4', true);
INSERT INTO app_keys (app_id, key, active) VALUES ('file-storage-api', '9eb87aa2-75e7-45b2-a1e6-53ed297d0ba8', true);

INSERT INTO public.user_to_tenant (username, tenant) VALUES ('dev@openline.ai', 'openline');

INSERT INTO public.conversation_event (event_uuid, tenant_name, conversation_id, type, subtype, initiator_username, sender_id,
                                       sender_type, sender_username, content, source, external_id, direction, created_at) VALUES
('481876a4-cd17-4b89-a67d-60660b9d9b87', 'openline', 'echotest',
 'WEB_CHAT', 'TEXT', 'echotest', '3b6087e8-a104-4d2d-acf7-4ee8ae5adf2c', 'CONTACT', 'echotest', 'Hello world!', 'openline', '',
 'INBOUND', '2023-01-13 12:35:07.333292 +00:00');