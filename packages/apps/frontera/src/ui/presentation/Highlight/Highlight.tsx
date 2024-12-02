import React from 'react';

interface HighlightProps {
  term: string;
  className?: string;
  children?: React.ReactNode;
}

export const Highlight = ({ term, className, children }: HighlightProps) => {
  if (typeof children !== 'string') return <>{children}</>;

  if (
    !term ||
    term.length === 0 ||
    !children.toLowerCase().includes(term.toLowerCase())
  ) {
    return <>{children}</>;
  }

  // Escape special regex characters in the term
  const escapeRegex = (str: string) =>
    str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');

  const safeTerm = escapeRegex(term);
  const regex = new RegExp(`(${safeTerm})`, 'gi');
  const parts = children.split(regex);

  return (
    <>
      {parts.map((part, index) =>
        regex.test(part) ? (
          <span key={index} className={className}>
            {part}
          </span>
        ) : (
          part
        ),
      )}
    </>
  );
};
