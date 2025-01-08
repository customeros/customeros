import { twMerge } from 'tailwind-merge';

export const ButtonGroup = ({
  className,
  variant = 'default',
  ...props
}: React.HTMLAttributes<HTMLDivElement> & { variant?: 'default' | 'old' }) => {
  return (
    <div
      {...props}
      className={twMerge(
        variant === 'default' ? defaultButtonGroup : '',
        variant === 'old' ? oldButtonGroup : '',
        className,
      )}
    />
  );
};
const defaultButtonGroup = `
  border border-grayModern-300 rounded-lg flex transition-all duration-200 ease-out 
  [&>*]:border-0
  [&>*.selected]:border
  [&>*.selected]:border-grayModern-300
  [&>*.selected]:!rounded-lg
  [&>*.selected]:-mt-[1px]
  [&>*.selected]:-mb-[1px]
  [&>*.selected]:relative
  [&>*.selected]:z-10
  [&>*.selected]:text-primary-700
  [&>*.selected *]:text-primary-700
  [&>*.selected]:bg-white
  [&>*.selected]:hover:text-primary-600
  [&>*:not(.selected)]:border-r
  [&>*:not(.selected)]:border-grayModern-300
  [&>*:not(.selected)]:bg-grayModern-50
  [&>*:not(.selected)]:text-gray-500
  [&>*:not(.selected)]:font-normal
  [&>*:not(.selected):has(+.selected)]:border-r-0
  [&>*:not(:last-child):not(:first-child)]:rounded-none
  [&>*:first-child:not(.selected)]:rounded-r-none
  [&>*:last-child:not(.selected)]:rounded-l-none
  [&>*:first-child]:border-l-0
  [&>*:last-child]:border-r-0
  [&>*:last-child:not(.selected)]:border-r-0
`;

const oldButtonGroup =
  '[&>*:not(:last-child):not(:first-child)]:rounded-none [&>*:first-child]:rounded-r-none [&>*:first-child]:border-r-0 [&>*:last-child]:rounded-l-none [&>*:last-child]:border-l-0';
