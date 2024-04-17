import React, { MouseEventHandler } from 'react';

import { Icons } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon2';
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
        <AlertDialogContent className='top-[25%] rounded-xl bg-[url(/backgrounds/organization/circular-bg-pattern.png)] bg-no-repeat'>
          <AlertDialogHeader className='text-lg font-bold pt-6'>
            <FeaturedIcon
              size='lg'
              colorScheme='primary'
              className='translate-y-[-10px] translate-x-[10px]'
            >
              <Icons.InfoCircle />
            </FeaturedIcon>
            {label && <p className='mt-4'>{label}</p>}
            {children ??
              (description && (
                <p className='mt-4 text-base text-gray-600 font-normal'>
                  {description}
                </p>
              ))}
          </AlertDialogHeader>

          <AlertDialogFooter className='pb-6 grid-cols-1'>
            <Button className='w-full' variant='outline' onClick={onConfirm}>
              {confirmButtonLabel}
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialogOverlay>
    </AlertDialog>
  );
};
