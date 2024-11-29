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

  // Regex to match all occurrences of the term (case-insensitive)
  const regex = new RegExp(`(${term})`, 'gi');
  const parts = children.split(regex);

  return (
    <>
      {parts.map((part, index) =>
        part.toLowerCase() === term.toLowerCase() ? (
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
