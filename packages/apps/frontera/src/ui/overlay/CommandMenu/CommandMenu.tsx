import { Command, useCommandState } from 'cmdk';

import { cn } from '@ui/utils/cn';
import { Tag, TagLabel } from '@ui/presentation/Tag/Tag';
import { ChevronRight } from '@ui/media/icons/ChevronRight';

interface CommandInputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  value?: string;
  asChild?: boolean;
  placeholder: string;
  children?: React.ReactNode;
  onValueChange?: (value: string) => void;
}

export const CommandInput = ({
  label,
  asChild,
  children,
  placeholder,
  onValueChange,
  ...rest
}: CommandInputProps) => {
  return (
    <div className='relative w-full p-6 pb-2 flex flex-col gap-2 border-b border-b-gray-100'>
      {label && (
        <Tag size='md' variant='subtle' colorScheme='gray'>
          <TagLabel>{label}</TagLabel>
        </Tag>
      )}
      <div className='w-full min-h-10 flex items-center'>
        <Command.Input
          autoFocus
          asChild={asChild}
          children={children}
          placeholder={placeholder}
          onValueChange={onValueChange}
          {...rest}
        />
      </div>
    </div>
  );
};

interface CommandItemProps extends React.HTMLAttributes<HTMLDivElement> {
  disabled?: boolean;
  keywords?: string[];
  onSelect?: () => void;
  children: React.ReactNode;
  leftAccessory?: React.ReactNode;
  rightAccessory?: React.ReactNode;
}

export const CommandItem = ({
  children,
  disabled,
  leftAccessory,
  rightAccessory,
  ...props
}: CommandItemProps) => {
  return (
    <Command.Item disabled={disabled} {...props}>
      {leftAccessory}
      {children}
      <div className='flex gap-1 items-center ml-auto'>{rightAccessory}</div>
    </Command.Item>
  );
};

interface CommandSubItemProps {
  leftLabel: string;
  rightLabel: string;
  keywords?: string[];
  icon: React.ReactNode;
  onSelectAction: () => void;
}

export const CommandSubItem: React.FC<CommandSubItemProps> = ({
  icon,
  onSelectAction,
  leftLabel,
  rightLabel,
  ...rest
}) => {
  const search = useCommandState((state) => state.search);

  return (
    <CommandItem
      leftAccessory={icon}
      onSelect={onSelectAction}
      disabled={search.length <= 3}
      className={cn(search.length <= 3 && 'hidden')}
      {...rest}
    >
      <span className='text-gray-500'>{leftLabel}</span>
      <ChevronRight className='mx-1' />
      <span>{rightLabel}</span>
    </CommandItem>
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
