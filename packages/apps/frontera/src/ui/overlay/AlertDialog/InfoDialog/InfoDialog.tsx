import React, { MouseEventHandler } from 'react';

import { Button } from '@ui/form/Button/Button';
import {
  AlertDialog,
  AlertDialogHeader,
  AlertDialogFooter,
  AlertDialogOverlay,
  AlertDialogContent,
} from '@ui/overlay/AlertDialog/AlertDialog';

interface InfoDialogProps {
  label?: string;
  isOpen: boolean;
  onClose: () => void;
  description?: string;
  confirmButtonLabel: string;
  children?: React.ReactNode;
  onConfirm: MouseEventHandler<HTMLButtonElement>;
}

export const InfoDialog = ({
  isOpen,
  onClose,
  onConfirm,
  label,
  description,
  children,
  confirmButtonLabel,
}: InfoDialogProps) => {
  return (
    <AlertDialog isOpen={isOpen} onClose={onClose}>
      <AlertDialogOverlay>
        <AlertDialogContent className='top-[25%] rounded-xl '>
          <AlertDialogHeader className='text-lg font-bold'>
            {label && <p className='font-semibold text-lg'>{label}</p>}
            {children ??
              (description && (
                <p className='mt-4 text-base text-gray-600 font-normal'>
                  {description}
                </p>
              ))}
          </AlertDialogHeader>

          <AlertDialogFooter className='grid-cols-1'>
            <Button variant='outline' className='w-full' onClick={onConfirm}>
              {confirmButtonLabel}
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialogOverlay>
    </AlertDialog>
  );
};
