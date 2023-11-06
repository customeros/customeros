import React, { useRef, MouseEventHandler } from 'react';

import { AlertDialogBody } from '@chakra-ui/modal';

import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { Icons, FeaturedIcon } from '@ui/media/Icon';
import {
  AlertDialog,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogContent,
  AlertDialogOverlay,
  AlertDialogCloseButton,
} from '@ui/overlay/AlertDialog';

interface ConfirmDeleteDialogProps {
  label: string;
  isOpen: boolean;
  isLoading?: boolean;
  onClose: () => void;
  description?: string;
  colorScheme?: string;
  icon?: React.ReactNode;
  body?: React.ReactNode;
  confirmButtonLabel: string;
  cancelButtonLabel?: string;
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
  icon,
  colorScheme = 'red',
}: ConfirmDeleteDialogProps) => {
  const cancelRef = useRef<HTMLButtonElement>(null);

  return (
    <AlertDialog
      isOpen={isOpen}
      onClose={onClose}
      leastDestructiveRef={cancelRef}
      closeOnEsc
    >
      <AlertDialogOverlay>
        <AlertDialogContent
          borderRadius='xl'
          backgroundImage='/backgrounds/organization/circular-bg-pattern.png'
          backgroundRepeat='no-repeat'
        >
          <AlertDialogCloseButton color='gray.400' top={6} />
          <AlertDialogHeader fontSize='lg' fontWeight='bold' pt='6'>
            <FeaturedIcon size='lg' colorScheme={colorScheme}>
              {icon ? icon : <Icons.Trash1 />}
            </FeaturedIcon>
            <Text mt='4'>{label}</Text>
            {description && (
              <Text mt='4' fontSize='md' color='gray.600' fontWeight='normal'>
                {description}
              </Text>
            )}
          </AlertDialogHeader>

          {body && <AlertDialogBody>{body}</AlertDialogBody>}

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
              colorScheme={colorScheme}
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
