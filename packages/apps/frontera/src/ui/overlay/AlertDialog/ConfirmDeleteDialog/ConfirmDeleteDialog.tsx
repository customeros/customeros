import { useRef, MouseEventHandler } from 'react';

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
  label: string;
  isOpen: boolean;
  dataTest?: string;
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

export const ConfirmDeleteDialog = ({
  body,
  label,
  isOpen,
  onClose,
  dataTest,
  isLoading,
  onConfirm,
  description,
  hideCloseButton,
  confirmButtonLabel,
  colorScheme = 'error',
  cancelButtonLabel = 'Cancel',
  loadingButtonLabel = 'Deleting',
}: ConfirmDeleteDialogProps) => {
  const cancelRef = useRef<HTMLButtonElement>(null);

  return (
    <AlertDialog isOpen={isOpen} onClose={onClose} className='z-[99999]'>
      <AlertDialogPortal>
        <AlertDialogOverlay>
          <AlertDialogContent className='rounded-xl'>
            {!hideCloseButton && <AlertDialogCloseIconButton />}
            <AlertDialogHeader className='text-lg font-bold'>
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
                  data-test={dataTest}
                  isLoading={isLoading}
                  loadingText={loadingButtonLabel}
                  colorScheme={colorScheme || 'error'}
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
