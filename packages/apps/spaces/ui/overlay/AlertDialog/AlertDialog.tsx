import React, { forwardRef, ElementRef } from 'react';

import { twMerge } from 'tailwind-merge';
import * as RadixAlertDialog from '@radix-ui/react-alert-dialog';

import { XClose } from '@ui/media/icons/XClose';
import { IconButton } from '@ui/form/IconButton/IconButton';

interface AlertDialogGenericProps extends React.HTMLAttributes<HTMLDivElement> {
  asChild?: boolean;
  className?: string;
  children: React.ReactNode;
}

interface AlertDialogProps {
  isOpen: boolean;
  className?: string;
  onClose: () => void;
  children: React.ReactNode;
}

export const AlertDialog = forwardRef<
  ElementRef<typeof RadixAlertDialog.AlertDialog>,
  AlertDialogProps
>(({ children, className, isOpen, onClose, ...props }, ref) => {
  return (
    <RadixAlertDialog.Root open={isOpen} onOpenChange={onClose} {...props}>
      {children}
    </RadixAlertDialog.Root>
  );
});

export const AlertDialogPortal = ({ children }: AlertDialogGenericProps) => {
  return (
    <RadixAlertDialog.Portal
      container={typeof window !== 'undefined' ? document?.body : null}
    >
      {children}
    </RadixAlertDialog.Portal>
  );
};
//TODO:remove z-[5000] from overlay after chakra-ui is removed
export const AlertDialogOverlay = forwardRef<
  ElementRef<typeof RadixAlertDialog.Overlay>,
  AlertDialogGenericProps
>(({ children, className }, ref) => {
  return (
    <RadixAlertDialog.Overlay
      ref={ref}
      className={twMerge(
        'z-[5000] backdrop-brightness-[.55] data-[state=open]:animate-overlayShow fixed inset-0',
        className,
      )}
    >
      {children}
    </RadixAlertDialog.Overlay>
  );
});

export const AlertDialogContent = forwardRef<
  ElementRef<typeof RadixAlertDialog.Content>,
  AlertDialogGenericProps
>(({ className, children, ...props }, ref) => {
  return (
    <RadixAlertDialog.Content
      ref={ref}
      {...props}
      className={twMerge(
        'data-[state=open]:animate-contentShow fixed top-[14%] left-[50%] max-h-[80vh] w-[100%] outline-offset-2 max-w-[448px] translate-x-[-50%] translate-y-[-50%] rounded-xl bg-white p-6 focus:outline-none',
        className,
      )}
    >
      {children}
    </RadixAlertDialog.Content>
  );
});

export const AlertDialogHeader = forwardRef<
  ElementRef<typeof RadixAlertDialog.Title>,
  AlertDialogGenericProps
>(({ children, className }, ref) => {
  return (
    <RadixAlertDialog.Title className={twMerge(className)} ref={ref}>
      {children}
    </RadixAlertDialog.Title>
  );
});

export const AlertDialogBody = forwardRef<
  ElementRef<typeof RadixAlertDialog.Description>,
  AlertDialogGenericProps
>(({ className, children, asChild }, ref) => {
  return (
    <RadixAlertDialog.Description
      ref={ref}
      asChild={asChild}
      className={twMerge(className, 'start-6 end-6 flex-1 py-2')}
    >
      {children}
    </RadixAlertDialog.Description>
  );
});

export const AlertDialogFooter = ({
  children,
  className,
}: AlertDialogGenericProps) => {
  return (
    <div className={twMerge('grid grid-cols-2 pt-4 gap-3', className)}>
      {children}
    </div>
  );
};
export const AlertDialogCloseButton = forwardRef<
  ElementRef<typeof RadixAlertDialog.AlertDialogCancel>,
  RadixAlertDialog.AlertDialogCancelProps
>(({ asChild, children, ...props }, ref) => {
  return (
    <RadixAlertDialog.Cancel ref={ref} asChild {...props}>
      {children}
    </RadixAlertDialog.Cancel>
  );
});

export const AlertDialogConfirmButton = forwardRef<
  ElementRef<typeof RadixAlertDialog.Action>,
  RadixAlertDialog.AlertDialogActionProps
>(({ children, ...props }, ref) => {
  return (
    <RadixAlertDialog.Action className='w-full' ref={ref} {...props}>
      {children}
    </RadixAlertDialog.Action>
  );
});

export const AlertDialogCloseIconButton = forwardRef<
  ElementRef<typeof RadixAlertDialog.Cancel>,
  RadixAlertDialog.AlertDialogCancelProps
>(({ asChild, children, className, ...props }, ref) => {
  return (
    <RadixAlertDialog.Cancel
      className={twMerge('absolute right-3 top-3', className)}
      ref={ref}
      {...props}
    >
      <IconButton
        variant='ghost'
        colorScheme='gray'
        icon={<XClose />}
        aria-label='Close dialog'
      />
    </RadixAlertDialog.Cancel>
  );
});
