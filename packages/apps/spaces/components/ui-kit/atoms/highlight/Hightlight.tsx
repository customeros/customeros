import escapeStringRegexp from 'escape-string-regexp';
import React from 'react';
interface Props {
  text: string;
  highlight: string;
}

export const Highlight: React.FC<Props> = ({ text, highlight }) => {
  if (!text) return null;
  if (!highlight) return <>{text}</>;
  // Split text on higlight term, include term itself into parts, ignore case
  const regex = new RegExp(`(${escapeStringRegexp(highlight)})`, 'gi');
  const parts = text.split(regex);

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
