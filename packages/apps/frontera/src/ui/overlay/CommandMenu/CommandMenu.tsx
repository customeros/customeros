import { Command, useCommandState } from 'cmdk';

import { cn } from '@ui/utils/cn';
import { Tag, TagLabel } from '@ui/presentation/Tag/Tag';

interface CommandInputProps {
  label?: string;
  placeholder: string;
}

export const CommandInput = ({ label, placeholder }: CommandInputProps) => {
  return (
    <div className='p-6 pb-2 flex flex-col gap-2 border-b border-b-gray-100'>
      {label && (
        <Tag variant='subtle' colorScheme='gray' size='lg'>
          <TagLabel>{label}</TagLabel>
        </Tag>
      )}
      <div className='w-full h-10 flex items-center'>
        <Command.Input placeholder={placeholder} />
      </div>
    </div>
  );
};

interface CommandItemProps {
  onSelect: () => void;
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
