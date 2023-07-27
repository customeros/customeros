import { useRef } from 'react';

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

import { BackgroundPattern } from './BackgroundPattern';

interface ConfirmDeleteDialogProps {
  isOpen: boolean;
  isLoading?: boolean;
  onClose: () => void;
  onConfirm: () => void;
}

export const ConfirmDeleteDialog = ({
  isOpen,
  onClose,
  isLoading,
  onConfirm,
}: ConfirmDeleteDialogProps) => {
  const cancelRef = useRef<HTMLButtonElement>(null);

  return (
    <AlertDialog
      isOpen={isOpen}
      onClose={onClose}
      leastDestructiveRef={cancelRef}
    >
      <AlertDialogOverlay>
        <AlertDialogContent borderRadius='xl'>
          <AlertDialogCloseButton />
          <AlertDialogHeader fontSize='lg' fontWeight='bold' pt='6'>
            <BackgroundPattern />
            <FeaturedIcon size='lg' colorScheme='red'>
              <Icons.Trash1 />
            </FeaturedIcon>
            <Text mt='4'>Delete this contact?</Text>
          </AlertDialogHeader>

          <AlertDialogFooter pb='6'>
            <Button
              w='full'
              ref={cancelRef}
              onClick={onClose}
              isDisabled={isLoading}
              variant='outline'
            >
              Cancel
            </Button>
            <Button
              ml={3}
              w='full'
              variant='outline'
              colorScheme='red'
              onClick={onConfirm}
              isLoading={isLoading}
            >
              Delete contact
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialogOverlay>
    </AlertDialog>
  );
};
