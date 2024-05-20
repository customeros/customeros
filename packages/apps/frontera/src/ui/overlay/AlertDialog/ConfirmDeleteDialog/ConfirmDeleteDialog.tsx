import React, { useRef, MouseEventHandler } from 'react';

import { Trash01 } from '@ui/media/icons/Trash01';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
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
      <AlertDialogPortal>
        <AlertDialogOverlay>
          <AlertDialogContent className='rounded-xl bg-no-repeat bg-[url(/backgrounds/organization/circular-bg-pattern.png)]'>
            {!hideCloseButton && <AlertDialogCloseIconButton />}
            <FeaturedIcon
              size='lg'
              // eslint-disable-next-line @typescript-eslint/no-explicit-any
              colorScheme={colorScheme as any}
              className='mt-[13px] ml-[11px]'
            >
              {icon ? icon : <Trash01 />}
            </FeaturedIcon>
            <AlertDialogHeader className='text-lg font-bold mt-4'>
              <p className='pb-0 font-semibold'>{label}</p>
              {description && (
                <p className='mt-1 text-base text-gray-700 font-normal'>
                  {description}
                </p>
              )}
            </AlertDialogHeader>
            {body && <AlertDialogBody>{body}</AlertDialogBody>}
            <AlertDialogFooter>
              <AlertDialogCloseButton>
                <Button
                  ref={cancelRef}
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
      </AlertDialogPortal>
    </AlertDialog>
  );
};
