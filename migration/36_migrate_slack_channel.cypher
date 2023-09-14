match (i:InteractionEvent) where i.channel = 'SLACK' set i.channel = 'CHAT';
match (i:InteractionSession) where i.channel = 'SLACK' set i.channel = 'CHAT';