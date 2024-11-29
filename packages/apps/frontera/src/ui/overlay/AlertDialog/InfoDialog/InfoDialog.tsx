import React, { MouseEventHandler } from 'react';

import { Button } from '@ui/form/Button/Button';

import {
  AlertDialog,
  AlertDialogBody,
  AlertDialogHeader,
  AlertDialogFooter,
  AlertDialogPortal,
  AlertDialogOverlay,
  AlertDialogContent,
  AlertDialogCloseIconButton,
} from '../AlertDialog';

interface InfoDialogProps {
  label?: string;
  isOpen: boolean;
  onClose: () => void;
  description?: string;
  body?: React.ReactNode;
  hideCloseButton?: boolean;
  confirmButtonLabel: string;
  children?: React.ReactNode;
  onConfirm: MouseEventHandler<HTMLButtonElement>;
}

export const InfoDialog = ({
  body,
  label,
  isOpen,
  onClose,
  children,
  onConfirm,
  description,
  hideCloseButton,
  confirmButtonLabel,
}: InfoDialogProps) => {
  return (
    <AlertDialog isOpen={isOpen} onClose={onClose}>
      <AlertDialogPortal>
        <AlertDialogOverlay>
          <AlertDialogContent className='top-[15%] rounded-md'>
            {!hideCloseButton && <AlertDialogCloseIconButton />}
            <AlertDialogHeader className='font-bold'>
              {label && <p className='pb-0 font-semibold truncate'>{label}</p>}
              {description && (
                <p className='mt-1 text-sm text-gray-700 font-normal'>
                  {description}
                </p>
              )}
            </AlertDialogHeader>

            {(body ?? children) && (
              <AlertDialogBody>{body ?? children}</AlertDialogBody>
            )}

            <AlertDialogFooter className='grid-cols-1'>
              <Button variant='outline' className='w-full' onClick={onConfirm}>
                {confirmButtonLabel}
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialogPortal>
    </AlertDialog>
  );
};
