import React, { useRef, MouseEventHandler } from 'react';

import {
  AlertDialog,
  AlertDialogOverlay,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogFooter,
  AlertDialogCloseButton,
} from '@ui/overlay/AlertDialog';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon, Icons } from '@ui/media/Icon';

interface ConfirmDeleteDialogProps {
  isOpen: boolean;
  isLoading?: boolean;
  onClose: () => void;
  onConfirm: MouseEventHandler<HTMLButtonElement>;
  label: string;
  confirmButtonLabel: string;
  cancelButtonLabel?: string;
}

export const ConfirmDeleteDialog = ({
  isOpen,
  onClose,
  isLoading,
  onConfirm,
  label,
  confirmButtonLabel,
  cancelButtonLabel = 'Cancel',
}: ConfirmDeleteDialogProps) => {
  const cancelRef = useRef<HTMLButtonElement>(null);

  return (
    <AlertDialog
      isOpen={isOpen}
      onClose={onClose}
      leastDestructiveRef={cancelRef}
    >
      <AlertDialogOverlay>
        <AlertDialogContent
          borderRadius='xl'
          backgroundImage='/backgrounds/organization/circular-bg-pattern.png'
          backgroundRepeat='no-repeat'
        >
          <AlertDialogCloseButton color='gray.400' top={6} />
          <AlertDialogHeader fontSize='lg' fontWeight='bold' pt='6'>
            <FeaturedIcon size='lg' colorScheme='red'>
              <Icons.Trash1 />
            </FeaturedIcon>
            <Text mt='4'>{label}</Text>
          </AlertDialogHeader>

          <AlertDialogFooter pb='6'>
            <Button
              w='full'
              ref={cancelRef}
              onClick={onClose}
              isDisabled={isLoading}
              variant='outline'
              bg='white'
            >
              {cancelButtonLabel}
            </Button>
            <Button
              ml={3}
              w='full'
              variant='outline'
              colorScheme='red'
              onClick={onConfirm}
              isLoading={isLoading}
              loadingText='Deleting'
            >
              {confirmButtonLabel}
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialogOverlay>
    </AlertDialog>
  );
};
