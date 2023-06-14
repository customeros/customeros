import React from 'react';

export const ExternalLink = ({ url }: { url: string }) => {
  const createExternalLink = (link: string) => {
    if (link.includes('http')) {
      return link;
    }
    return `https://${link}`;
  };
  return (
    <a href={createExternalLink(url)} rel='noopener noreferrer' target='_blank'>
      {url}
    </a>
  );
};
