import React, { Ref, forwardRef, MouseEvent, KeyboardEvent } from 'react';

import { Command, useCommandState } from 'cmdk';

import { cn } from '@ui/utils/cn';
import { IconButton } from '@ui/form/IconButton';
import { XClose } from '@ui/media/icons/XClose.tsx';
import { Button } from '@ui/form/Button/Button.tsx';
import { Tag, TagLabel } from '@ui/presentation/Tag/Tag';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { isUserPlatformMac } from '@utils/getUserPlatform.ts';
import { Command as CommandIcon } from '@ui/media/icons/Command';

interface CommandInputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  value?: string;
  asChild?: boolean;
  dataTest?: string;
  placeholder: string;
  wrapperClassName?: string;
  children?: React.ReactNode;
  inputWrapperClassName?: string;
  onValueChange?: (value: string) => void;
  onKeyDown?: (e: React.KeyboardEvent<HTMLInputElement>) => void;
}

export const CommandInput = ({
  label,
  asChild,
  children,
  dataTest,
  placeholder,
  onValueChange,
  onKeyDown,
  wrapperClassName,
  inputWrapperClassName,
  ...rest
}: CommandInputProps) => {
  return (
    <div
      className={cn(
        'relative w-full p-6 pb-2 flex flex-col gap-2 border-b border-b-gray-100',
        wrapperClassName,
      )}
    >
      {label && (
        <Tag size='md' variant='subtle' colorScheme='gray'>
          <TagLabel>{label}</TagLabel>
        </Tag>
      )}
      <div
        className={cn(
          'w-full min-h-10 flex items-center',
          inputWrapperClassName,
        )}
      >
        <Command.Input
          autoFocus
          asChild={asChild}
          children={children}
          data-test={dataTest}
          placeholder={placeholder}
          onValueChange={onValueChange}
          {...rest}
        />
      </div>
    </div>
  );
};

interface CommandItemProps extends React.HTMLAttributes<HTMLDivElement> {
  value?: string;
  dataTest?: string;
  disabled?: boolean;
  keywords?: string[];
  onSelect?: () => void;
  children: React.ReactNode;
  leftAccessory?: React.ReactNode;
  rightAccessory?: React.ReactNode;
}

export const CommandItem = forwardRef<HTMLDivElement, CommandItemProps>(
  (
    {
      children,
      dataTest,
      disabled = false,
      leftAccessory,
      rightAccessory,
      ...props
    }: CommandItemProps,
    ref,
  ) => {
    return (
      <Command.Item
        ref={ref}
        disabled={disabled}
        data-test={dataTest}
        {...props}
      >
        {leftAccessory && <span>{leftAccessory}</span>}
        {children}
        {rightAccessory && (
          <div className='flex gap-1 items-center ml-auto'>
            {rightAccessory}
          </div>
        )}
      </Command.Item>
    );
  },
);

interface CommandSubItemProps extends Partial<CommandItemProps> {
  leftLabel: string;
  rightLabel: string;
  keywords?: string[];
  icon: React.ReactNode;
  onSelectAction: () => void;
}

export const CommandSubItem = ({
  icon,
  onSelectAction,
  leftLabel,
  rightLabel,
  ...rest
}: CommandSubItemProps) => {
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
      <span className='overflow-hidden text-ellipsis whitespace-nowrap max-w-[250px]'>
        {rightLabel}
      </span>
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
        'bg-gray-100 text-gray-700 size-5 flex items-center justify-center rounded-[4px] text-xs',
        className,
      )}
    >
      {children}
    </kbd>
  );
};

export const CommandKbd = ({
  className,
}: React.HTMLAttributes<HTMLDivElement>) => {
  if (isUserPlatformMac()) {
    return (
      <Kbd className={className}>
        <CommandIcon className='size-3' />
      </Kbd>
    );
  }

  return (
    <kbd
      className={cn(
        'bg-gray-100 text-gray-700 flex p-1 py-0.5 items-center justify-center rounded-[4px] text-xs',
        className,
      )}
    >
      Ctrl
    </kbd>
  );
};

export const CommandCancelIconButton = ({
  onClose,
}: {
  onClose: (
    e: MouseEvent<HTMLButtonElement> | KeyboardEvent<HTMLButtonElement>,
  ) => void;
} & React.HTMLAttributes<HTMLButtonElement>) => {
  return (
    <IconButton
      size='xs'
      variant='ghost'
      icon={<XClose />}
      onClick={onClose}
      aria-label='cancel'
      onKeyDown={(e) => {
        if (e.key === 'Enter') {
          e.stopPropagation();
          onClose(e);
        }
      }}
    />
  );
};

export const CommandCancelButton = forwardRef(
  (
    {
      onClose,
      dataTest,
    }: {
      dataTest?: string;
      onClose: () => void;
    } & React.HTMLAttributes<HTMLButtonElement>,
    ref: Ref<HTMLButtonElement>,
  ) => {
    return (
      <Button
        size='sm'
        ref={ref}
        variant='outline'
        onClick={onClose}
        className='w-full'
        data-test={dataTest}
        onKeyDown={(e) => {
          if (e.key === 'Enter') {
            e.stopPropagation();
            onClose();
          }
        }}
      >
        Cancel
      </Button>
    );
  },
);

export { Command, useCommandState };
