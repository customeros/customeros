import React from 'react';

export const ExternalLink = ({ url }: { url: string }) => {
  const removeProtocolFromLink = (link: string): string => {
    const protocolIndex = link.indexOf('://');
    if (protocolIndex !== -1) {
      return link.slice(protocolIndex + 3);
    }
    return link;
  };
  const getExternalUrl = (link: string) => {
    const linkWithoutProtocol = removeProtocolFromLink(link);
    return `https://${linkWithoutProtocol}`;
  };

  const getFormattedLink = (url: string): string => {
    return url.replace(/^(https?:\/\/)?(www\.)?/i, '');
  };
  return (
    <a href={getExternalUrl(url)} rel='noopener noreferrer' target='_blank'>
      {getFormattedLink(url)}
    </a>
  );
};
