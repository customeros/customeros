export const getMetadata = (metadataString?: string | null) => {
  let metadata;
  try {
    metadata = metadataString && JSON.parse(metadataString);
  } catch (error) {
    metadata = '';
  }
  return metadata;
};

export const getLikelihoodDisplayData = (text: string) => {
  const match = text.match(/(.+? to )(.+?)(?: by )(.+)/);

  if (!match) {
    return { preText: '', likelihood: '', author: '' };
  }

  return {
    preText: match?.[1], // "Renewal likelihood set to "
    likelihood: match?.[2], // "Low"
    author: match?.[3], // "Olivia Rhye"
  };
};

export const getCurrencyString = (text: string) => {
  const match = text
    .split(/(\$[\d,]+(\.\d{2})?)/)
    .filter(Boolean)
    .map((item) => (item.endsWith(',') ? item.slice(0, -1) : item));
  return match?.[1];
};
