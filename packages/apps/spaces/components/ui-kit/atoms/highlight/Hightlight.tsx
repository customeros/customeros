import React from 'react';

interface Props {
  text: string;
  highlight: string;
}

export const Highlight: React.FC<Props> = ({ text, highlight }) => {
  if (!text) return null;
  // Split text on higlight term, include term itself into parts, ignore case
  const parts = text.split(new RegExp(`(${highlight})`, 'gi'));

  return (
    <>
      {parts.map((part: string, index: number) => (
        <React.Fragment key={index}>
          {part.toLowerCase() === highlight.toLowerCase() ? (
            <mark>{part}</mark>
          ) : (
            part
          )}
        </React.Fragment>
      ))}
    </>
  );
};
