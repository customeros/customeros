import React, { useRef, useState, useEffect, MouseEventHandler } from 'react';

import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { Spinner } from '@ui/feedback/Spinner/Spinner.tsx';
import { getExternalUrl } from '@utils/getExternalLink.ts';
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
} from '@ui/overlay/AlertDialog/AlertDialog.tsx';

interface ConfirmDeleteDialogProps {
  isOpen: boolean;
  isLoading?: boolean;
  onClose: () => void;
  organizationId: string;
  onConfirm: (e: string) => void;
}
function validateLinkedInProfileUrl(url: string): boolean {
  const linkedInProfileRegex =
    /^(https:\/\/)?(www\.)?linkedin\.com\/in\/[a-zA-Z0-9-]{3,100}\/?$/;

  return linkedInProfileRegex.test(url);
}

export const CreateContactFromLinkedInModal = observer(
  ({
    isOpen,
    onClose,
    isLoading,
    onConfirm,
    organizationId,
  }: ConfirmDeleteDialogProps) => {
    const store = useStore();
    const organizationStore = store.organizations.value.get(organizationId);
    const cancelRef = useRef<HTMLButtonElement>(null);
    const [url, setUrl] = useState('');
    const [validationError, setValidationError] = useState(false);

    useEffect(() => {
      store.ui.setIsEditingTableCell(isOpen);
      if (isOpen && !url.includes('linkedin.com')) {
        setUrl('');
      }
    }, [isOpen]);
    const handleClose = () => {
      setValidationError(false);

      setUrl('');
      onClose();
    };

    const handleConfirm: MouseEventHandler<HTMLButtonElement> = (e) => {
      e.stopPropagation();

      setValidationError(false);
      const isValidUrl = validateLinkedInProfileUrl(url);
      if (isValidUrl) {
        const formattedUrl = getExternalUrl(url);
        onConfirm(formattedUrl);
        setUrl('');
        onClose();

        return;
      }
      setValidationError(true);
    };

    return (
      <AlertDialog isOpen={isOpen} onClose={handleClose} className='z-[99999]'>
        <AlertDialogPortal>
          <AlertDialogOverlay>
            <AlertDialogContent className='rounded-xl bg-no-repeat bg-[url(/backgrounds/organization/circular-bg-pattern.png)]'>
              <AlertDialogHeader className='text-lg font-bold'>
                <p className='pb-0 font-semibold'>
                  Create new contact for {organizationStore?.value?.name}
                </p>
                <p className='mt-1 mb-2 text-sm text-gray-700 font-normal'>
                  We will automatically enrich this contact when you create it
                </p>
              </AlertDialogHeader>
              <AlertDialogBody>
                <Input
                  autoComplete='off'
                  autoFocus
                  size='sm'
                  name='linkedin-input'
                  value={url}
                  placeholder='Contact`s LinkedIn URL'
                  onChange={(e) => {
                    setUrl(e.target.value);
                  }}
                />
                {validationError && (
                  <p className='text-sm text-warning-600'>
                    The LinkedIn URL seems to be invalid.
                  </p>
                )}
              </AlertDialogBody>
              <AlertDialogFooter>
                <AlertDialogCloseButton>
                  <Button
                    ref={cancelRef}
                    variant='outline'
                    colorScheme={'gray'}
                    size='md'
                    className='bg-white w-full'
                  >
                    Cancel
                  </Button>
                </AlertDialogCloseButton>
                <AlertDialogConfirmButton>
                  <Button
                    className='w-full'
                    variant='outline'
                    size='md'
                    colorScheme={'primary'}
                    onClick={handleConfirm}
                    isLoading={isLoading}
                    loadingText='Creating contact'
                    spinner={
                      <Spinner
                        size={'sm'}
                        label='deleting'
                        className='text-error-300 fill-error-700'
                      />
                    }
                  >
                    Create contact
                  </Button>
                </AlertDialogConfirmButton>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialogOverlay>
        </AlertDialogPortal>
      </AlertDialog>
    );
  },
);
