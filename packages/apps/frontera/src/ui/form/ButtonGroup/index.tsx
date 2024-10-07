import { twMerge } from 'tailwind-merge';

export const ButtonGroup = ({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) => {
  return (
    <div
      {...props}
      className={twMerge(
        '[&>*:not(:last-child):not(:first-child)]:rounded-none [&>*:first-child]:rounded-r-none [&>*:first-child]:border-r-0 [&>*:last-child]:rounded-l-none [&>*:last-child]:border-l-0',
        className,
      )}
    />
  );
};
