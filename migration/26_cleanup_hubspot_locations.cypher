#Execute after deploying to prod https://github.com/openline-ai/openline-customer-os/issues/1861

MATCH (loc:Location) WHERE loc.source = 'hubspot' AND loc.name = 'Default location' SET loc.name = '';