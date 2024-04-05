import React, { useRef, MouseEventHandler } from 'react';

import { Icons, FeaturedIcon } from '@ui/media/Icon';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { Button, ButtonProps } from '@ui/form/Button/Button';

import {
  AlertDialog,
  AlertDialogBody,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogContent,
  AlertDialogOverlay,
  AlertDialogCloseButton,
  AlertDialogConfirmButton,
} from '../AlertDialog';

interface ConfirmDeleteDialogProps {
  label: string;
  isOpen: boolean;
  isLoading?: boolean;
  onClose: () => void;
  description?: string;
  icon?: React.ReactNode;
  body?: React.ReactNode;
  hideCloseButton?: boolean;
  confirmButtonLabel: string;
  cancelButtonLabel?: string;
  loadingButtonLabel?: string;
  colorScheme?: ButtonProps['colorScheme'];
  onConfirm: MouseEventHandler<HTMLButtonElement>;
}

export const ConfirmDeleteDialog = ({
  isOpen,
  onClose,
  isLoading,
  onConfirm,
  label,
  description,
  body,
  confirmButtonLabel,
  cancelButtonLabel = 'Cancel',
  loadingButtonLabel = 'Deleting',
  icon,
  colorScheme = 'error',
  hideCloseButton,
}: ConfirmDeleteDialogProps) => {
  const cancelRef = useRef<HTMLButtonElement>(null);

  return (
    <AlertDialog isOpen={isOpen} onClose={onClose} className='z-[99999]'>
      <AlertDialogOverlay>
        <AlertDialogContent className='rounded-xl bg-no-repeat bg-[url(/backgrounds/organization/circular-bg-pattern.png)]'>
          {!hideCloseButton && (
            <AlertDialogCloseButton>
              <Button variant='solid' colorScheme='gray' onClick={onClose}>
                Close
              </Button>
            </AlertDialogCloseButton>
          )}
          <FeaturedIcon size='lg' colorScheme={colorScheme as string}>
            {icon ? icon : <Icons.Trash1 />}
          </FeaturedIcon>
          <AlertDialogHeader className='text-lg font-bold mt-4'>
            <span className='mt-4 pt-6 pb-0'>{label}</span>
            {description && (
              <span className='mt-4 text-base text-gray-600 '>
                {description}
              </span>
            )}
          </AlertDialogHeader>
          {body && <AlertDialogBody>{body}</AlertDialogBody>}

          <AlertDialogFooter>
            <AlertDialogCloseButton>
              <Button
                ref={cancelRef}
                onClick={onClose}
                isDisabled={isLoading}
                variant='outline'
                colorScheme={'gray'}
                size='md'
                className='bg-white w-full'
              >
                {cancelButtonLabel}
              </Button>
            </AlertDialogCloseButton>
            <AlertDialogConfirmButton>
              <Button
                className='w-full'
                variant='outline'
                size='md'
                colorScheme={colorScheme || 'error'}
                onClick={onConfirm}
                isLoading={isLoading}
                loadingText={loadingButtonLabel}
                spinner={
                  <Spinner
                    size={'sm'}
                    label='deleting'
                    className='text-error-300 fill-error-700'
                  />
                }
              >
                {confirmButtonLabel}
              </Button>
            </AlertDialogConfirmButton>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialogOverlay>
    </AlertDialog>
  );
};
