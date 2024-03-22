import ResizeTextarea, { TextareaAutosizeProps } from 'react-textarea-autosize';

import { cn } from '@ui/utils/cn';

interface TextareaProps extends TextareaAutosizeProps {
  border?: boolean;
  className?: string;
}

export const Textarea = ({ className, border }: TextareaProps) => {
  return (
    <ResizeTextarea
      className={cn(
        'flex items-center w-full border-b border-transparent gap-3 hover:broder-b hover:border-gray-300 bg-transparent focus-within:outline-none focus-within:hover:border-primary-500 focus-within:border-primary-500 focus-within:border-b hover:transition ease-in-out delay-200  ',
        className,
      )}
    >
      <textarea />
    </ResizeTextarea>
  );
};
