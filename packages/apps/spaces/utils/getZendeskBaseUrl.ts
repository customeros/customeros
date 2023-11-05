export const getZendeskIssueBaseUrl = (externalApiUrl: string) => {
  const url = `${externalApiUrl.split('.')[0]}.zendesk.com/agent/tickets`;
  if (url.startsWith('https')) return url;

  return `https://${url}`;
};

export const getZendeskIssuesBaseUrl = (
  externalApiUrl: string,
  externalSource: string,
) => {
  const url = `${
    externalApiUrl.split('.')[0]
  }.zendesk.com/agent/${externalSource}`;
  if (url.startsWith('https')) return url;

  return `https://${url}`;
};
