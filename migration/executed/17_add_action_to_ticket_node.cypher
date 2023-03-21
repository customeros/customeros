# replace <tenant> with the tenant name

MATCH (tt:Ticket) SET tt:TimelineEvent;
MATCH (tt:Ticket_<tenant>) SET tt:Action_<tenant>;