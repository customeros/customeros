# Execute before release & after release

MATCH (c:Contact) SET c.prefix = c.title REMOVE c.title;