export const removeProtocolFromLink = (link: string): string => {
  const protocolIndex = link.indexOf('://');
  if (protocolIndex !== -1) {
    return link.slice(protocolIndex + 3);
  }
  return link;
};
export const getExternalUrl = (link: string) => {
  const linkWithoutProtocol = removeProtocolFromLink(link);
  return `https://${linkWithoutProtocol}`;
};

export const getFormattedLink = (url: string): string => {
  return url.replace(/^(https?:\/\/)?(www\.)?/i, '');
};
