import { Command, useCommandState } from 'cmdk';

import { cn } from '@ui/utils/cn';
import { Tag, TagLabel } from '@ui/presentation/Tag/Tag';

interface CommandInputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  value?: string;
  onValueChange?: (search: string) => void;
  asChild?: boolean;
  placeholder: string;
  children?: React.ReactNode;
}

export const CommandInput = ({
  label,
  asChild,
  children,
  placeholder,
  ...rest
}: CommandInputProps) => {
  return (
    <div className='p-6 pb-2 flex flex-col gap-2 border-b border-b-gray-100'>
      {label && (
        <Tag size='md' variant='subtle' colorScheme='gray'>
          <TagLabel>{label}</TagLabel>
        </Tag>
      )}
      <div className='w-full h-10 flex items-center'>
        <Command.Input
          autoFocus
          asChild={asChild}
          children={children}
          placeholder={placeholder}
          {...rest}
        />
      </div>
    </div>
  );
};

interface CommandItemProps extends React.HTMLAttributes<HTMLDivElement> {
  onSelect?: () => void;
  children: React.ReactNode;
  leftAccessory?: React.ReactNode;
  rightAccessory?: React.ReactNode;
}

export const CommandItem = ({
  children,
  leftAccessory,
  rightAccessory,
  ...props
}: CommandItemProps) => {
  return (
    <Command.Item {...props}>
      {leftAccessory}
      {children}
      <div className='flex gap-1 items-center ml-auto'>{rightAccessory}</div>
    </Command.Item>
  );
};

export const StaticCommandItem = ({
  children,
  leftAccessory,
  rightAccessory,
  ...props
}: CommandItemProps) => {
  return (
    <div data-cmdk-item {...props}>
      {leftAccessory}
      {children}
      <div className='flex gap-1 items-center ml-auto'>{rightAccessory}</div>
    </div>
  );
};

interface KbdProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
}

export const Kbd = ({ children, className, ...props }: KbdProps) => {
  return (
    <kbd
      {...props}
      className={cn(
        'bg-gray-100 text-gray-700 px-2 py-1 rounded-md text-xs',
        className,
      )}
    >
      {children}
    </kbd>
  );
};

export { Command, useCommandState };
