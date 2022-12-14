INSERT INTO conversations(contact_id, created_on, updated_on, state, last_message, last_sender_id, last_sender_type) VALUES ('echotest', now(), now(), 'NEW', 'test message', '1', 'CONTACT');
INSERT INTO conversation_items(type, sender_id, sender_type, message, channel, direction, "time", conversation_conversation_item)
    (SELECT 'MESSAGE','1', 'CONTACT','test message', 'CHAT', 'INBOUND', now(), id FROM conversations WHERE contact_id='echotest');
INSERT INTO conversation_items(type, sender_id, sender_type, message, channel, direction, "time", conversation_conversation_item)
    (SELECT 'MESSAGE','1', 'CONTACT','email message', 'MAIL', 'INBOUND', now(), id FROM conversations WHERE contact_id='echotest');
INSERT INTO conversation_items(type, sender_id, sender_type, message, channel, direction, "time", conversation_conversation_item)
    (SELECT 'MESSAGE','1', 'CONTACT', 'here would be the transcription for the phone call', 'VOICE', 'INBOUND', now(), id FROM conversations WHERE contact_id='echotest');