import React, { useRef, MouseEventHandler } from 'react';

import {
  AlertDialog,
  AlertDialogOverlay,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogFooter,
} from '@ui/overlay/AlertDialog';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon, Icons } from '@ui/media/Icon';

interface InfoDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: MouseEventHandler<HTMLButtonElement>;
  label?: string;
  description?: string;
  confirmButtonLabel: string;
  children?: React.ReactNode;
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
          <AlertDialogHeader fontSize='lg' fontWeight='bold' pt='6'>
            <FeaturedIcon size='lg' colorScheme='primary'>
              <Icons.InfoCircle />
            </FeaturedIcon>
            {label && <Text mt='4'>{label}</Text>}
            {children ??
              (description && (
                <Text mt='4' fontSize='md' color='gray.600' fontWeight='normal'>
                  {description}
                </Text>
              ))}
          </AlertDialogHeader>

          <AlertDialogFooter pb='6'>
            <Button w='full' variant='outline' onClick={onConfirm}>
              {confirmButtonLabel}
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialogOverlay>
    </AlertDialog>
  );
};
