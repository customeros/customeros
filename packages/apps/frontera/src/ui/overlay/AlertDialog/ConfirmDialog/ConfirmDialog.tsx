import React, { useRef, MouseEventHandler } from 'react';

import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { Button, ButtonProps } from '@ui/form/Button/Button';

import {
  AlertDialog,
  AlertDialogBody,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogPortal,
  AlertDialogContent,
  AlertDialogOverlay,
  AlertDialogCloseButton,
  AlertDialogConfirmButton,
  AlertDialogCloseIconButton,
} from '../AlertDialog';

interface ConfirmDeleteDialogProps {
  title: string;
  isOpen: boolean;
  isLoading?: boolean;
  onClose: () => void;
  description?: string;
  body?: React.ReactNode;
  hideCloseButton?: boolean;
  confirmButtonLabel: string;
  cancelButtonLabel?: string;
  loadingButtonLabel?: string;
  colorScheme?: ButtonProps['colorScheme'];
  onConfirm: MouseEventHandler<HTMLButtonElement>;
}

export const ConfirmDialog = ({
  isOpen,
  onClose,
  isLoading,
  onConfirm,
  title,
  description,
  body,
  confirmButtonLabel,
  cancelButtonLabel = 'Cancel',
  loadingButtonLabel = 'Loading action...',
  colorScheme = 'primary',
  hideCloseButton,
}: ConfirmDeleteDialogProps) => {
  const cancelRef = useRef<HTMLButtonElement>(null);

  return (
    <AlertDialog isOpen={isOpen} onClose={onClose} className='z-[99999]'>
      <AlertDialogPortal>
        <AlertDialogOverlay>
          <AlertDialogContent className='rounded-xl '>
            {!hideCloseButton && <AlertDialogCloseIconButton />}

            <AlertDialogHeader className='font-bold mt-4'>
              <p className='pb-0 font-semibold'>{title}</p>
              {description && (
                <p className='mt-1 text-sm text-gray-700 font-normal'>
                  {description}
                </p>
              )}
            </AlertDialogHeader>
            {body && <AlertDialogBody>{body}</AlertDialogBody>}
            <AlertDialogFooter>
              <AlertDialogCloseButton>
                <Button
                  size='md'
                  ref={cancelRef}
                  variant='outline'
                  colorScheme={'gray'}
                  isDisabled={isLoading}
                  className='bg-white w-full'
                >
                  {cancelButtonLabel}
                </Button>
              </AlertDialogCloseButton>
              <AlertDialogConfirmButton>
                <Button
                  size='md'
                  variant='outline'
                  className='w-full'
                  onClick={onConfirm}
                  isLoading={isLoading}
                  loadingText={loadingButtonLabel}
                  colorScheme={colorScheme || 'primary'}
                  spinner={
                    <Spinner
                      size={'sm'}
                      label='deleting'
                      className='text-primary-300 fill-primary-700'
                    />
                  }
                >
                  {confirmButtonLabel}
                </Button>
              </AlertDialogConfirmButton>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialogPortal>
    </AlertDialog>
  );
};
