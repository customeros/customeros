export const getZendeskBaseUrl = (externalApiUrl: string) => {
  const url = `${externalApiUrl.split('.')[0]}.zendesk.com/agent/tickets`;
  if (url.startsWith('https')) return url;
  return `https://${url}`;
};
