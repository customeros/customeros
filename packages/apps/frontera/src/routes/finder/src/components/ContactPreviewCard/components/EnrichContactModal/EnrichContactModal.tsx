import { useRef, useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { Spinner } from '@ui/feedback/Spinner';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import {
  Modal,
  ModalBody,
  ModalClose,
  ModalFooter,
  ModalPortal,
  ModalOverlay,
  ModalCloseButton,
  ModalFeaturedHeader,
  ModalFeaturedContent,
} from '@ui/overlay/Modal';

export const EnrichContactModal = observer(
  ({
    isModalOpen = false,
    onClose,
    contactId,
  }: {
    onClose: () => void;
    isModalOpen: boolean;
    contactId: string | number;
  }) => {
    const store = useStore();
    const hasSubmitedRef = useRef(false);
    const [linkedin, setLinkedin] = useState(
      () => store.contacts.value.get(String(contactId))?.value.socials[0]?.url,
    );
    const [validation, setValidation] = useState<Record<'linkedin', boolean>>({
      linkedin: false,
    });

    const contactStore = store.contacts.value.get(String(contactId));

    const validate = () => {
      setValidation(() => ({
        linkedin: !linkedin,
      }));

      return linkedin;
    };

    const reset = () => {
      setLinkedin('');
      setValidation({
        linkedin: false,
      });
      hasSubmitedRef.current = false;
    };

    const handleSubmit = () => {
      hasSubmitedRef.current = true;

      if (!validate()) return;

      contactStore?.addSocial(linkedin || '', {
        onSuccess: () => {
          onClose();
          reset();
        },
      });
    };

    return (
      <Modal
        open={isModalOpen}
        onOpenChange={(open) => {
          if (!open) {
            reset();
            onClose();
          }
        }}
      >
        <ModalPortal>
          <ModalOverlay className='z-[999]' />
          <ModalFeaturedContent className='z-[9999]'>
            <ModalFeaturedHeader>
              <p className='text-lg font-semibold mb-1'>
                What’s this contact’s LinkedIn?
              </p>
              <p className='text-sm'>
                To enrich this contact, we need their LinkedIn URL
              </p>
            </ModalFeaturedHeader>
            <ModalCloseButton />
            <ModalBody className='flex flex-col gap-4'>
              <div className='flex flex-col'>
                <Input
                  id='linkedin'
                  value={linkedin}
                  placeholder='LinkedIn profile link'
                  className={cn(validation.linkedin && 'border-error-500')}
                  onChange={(e) => {
                    setLinkedin(e.target.value);
                  }}
                  onKeyDown={(e) => {
                    if (e.key === 'Escape') {
                      onClose();
                    }
                    e.stopPropagation();
                  }}
                />
                {validation.linkedin && (
                  <p className='text-sm text-error-500 mt-1'>
                    One does not simply skip LinkedIn
                  </p>
                )}
              </div>
            </ModalBody>
            <ModalFooter className='flex gap-3'>
              <ModalClose className='w-full'>
                <Button className='w-full'>Cancel</Button>
              </ModalClose>

              <Button
                className='w-full'
                colorScheme='primary'
                onClick={handleSubmit}
                loadingText='Creating contact'
                isLoading={store.contacts.isLoading}
                rightSpinner={
                  <Spinner
                    size='sm'
                    label='loading'
                    className='text-primary-500 fill-primary-200'
                  />
                }
              >
                Enrich contact
              </Button>
            </ModalFooter>
          </ModalFeaturedContent>
        </ModalPortal>
      </Modal>
    );
  },
);
