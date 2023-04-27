# Execute before release

MATCH (c:Contact) where c.prefix is null SET c.prefix = c.title REMOVE c.title;