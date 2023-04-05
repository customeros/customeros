# Execute before release & after release

MATCH (tt:Ticket)--(n:Note) detach delete n;

MATCH (tt:Ticket) detach delete tt;

MATCH (:Organization)--(n:Note) where n.source='zendesk_support' detach delete n;
MATCH (:Contact)--(n:Note) where n.source='zendesk_support' detach delete n;
MATCH (:Contact)--(cf:CustomField) where cf.source='zendesk_support' detach delete cf;
