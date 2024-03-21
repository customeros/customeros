import { cn } from '@ui/utils/cn';

interface TextareaProps
  extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  border?: boolean;
  className?: string;
}

export const Textarea = ({ className, border }: TextareaProps) => {
  return (
    <>
      <textarea
        className={cn(
          className,
          border ? ' border-gray-200' : 'border-transparent',
          'flex items-center w-full border-b gap-3 hover:broder-b hover:border-gray-300 bg-transparent focus-within:outline-none focus-within:hover:border-primary-500 focus-within:border-primary-500 focus-within:border-b hover:transition ease-in-out delay-200  ',
        )}
      />
    </>
  );
};
