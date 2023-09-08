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
import { SvgIcon } from '@ui/media/icons/SvgIcon';

interface ConfirmDeleteDialogProps {
  isOpen: boolean;
  isLoading?: boolean;
  onClose: () => void;
  onConfirm: MouseEventHandler<HTMLButtonElement>;
  label: string;
  description?: string;
  confirmButtonLabel: string;
  cancelButtonLabel?: string;
  icon?: any;
}

export const ConfirmDeleteDialog = ({
  isOpen,
  onClose,
  isLoading,
  onConfirm,
  label,
  description,
  confirmButtonLabel,
  icon,
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
              {icon ? <SvgIcon>{icon}</SvgIcon> : <Icons.Trash1 />}
            </FeaturedIcon>
            <Text mt='4'>{label}</Text>
            {description && (
              <Text mt='4' fontSize='md' color='gray.600' fontWeight='normal'>
                {description}
              </Text>
            )}
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
