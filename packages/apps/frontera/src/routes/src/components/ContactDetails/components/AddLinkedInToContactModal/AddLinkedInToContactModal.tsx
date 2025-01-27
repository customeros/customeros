import { useRef, useEffect } from 'react';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';
import { LinkedIn } from '@domain/usecases/contact-details/add-linkedin.usecase';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import {
  Modal,
  ModalBody,
  ModalHeader,
  ModalFooter,
  ModalOverlay,
  ModalContent,
  ModalCloseButton,
} from '@ui/overlay/Modal';

interface AddLinkedInToContactModalProps {
  open: boolean;
  contactId: string;
  onClose: () => void;
}

const linkedInUseCase = new LinkedIn();

export const AddLinkedInToContactModal = observer(
  ({ onClose, open, contactId }: AddLinkedInToContactModalProps) => {
    const store = useStore();
    const modalRef = useRef<HTMLDivElement>(null);

    const contactStore = store.contacts.value.get(contactId);

    useEffect(() => {
      if (contactStore) {
        linkedInUseCase.setEntity(contactStore);
      }
    }, [contactId]);

    useModKey(
      'Enter',
      () => {
        linkedInUseCase.submitLinkedInUrl();

        if (
          linkedInUseCase.emptyLinkedInUrl ||
          linkedInUseCase.invalidLinkedInUrl
        )
          return;

        onClose();
      },
      {
        targetRef: modalRef,
        when: open,
      },
    );

    useKey(
      'Escape',
      () => {
        onClose();
      },
      { target: modalRef, when: open },
    );

    return (
      <Modal open={open}>
        <ModalOverlay />
        <ModalContent ref={modalRef}>
          <ModalHeader>
            <p className='font-medium'>LinkedIn profile URL</p>
            <ModalCloseButton asChild />
          </ModalHeader>
          <ModalBody className='flex flex-col  text-sm'>
            <div className='flex flex-col gap-4'>
              <p>We'll auto-enrich this contact using their LinkedIn profile</p>
              <Input
                variant='unstyled'
                dataTest='linkedin-url-input'
                value={linkedInUseCase.inputValue}
                placeholder='linkedin.com/in/john-lemon'
                onChange={(e) => {
                  linkedInUseCase.setInputValue(e.target.value);

                  linkedInUseCase.resetErrors();
                }}
                onKeyDown={(e) => {
                  if (e.key === 'Escape') {
                    onClose();
                  }
                  e.stopPropagation();
                }}
              />
            </div>
            {linkedInUseCase.emptyLinkedInUrl && (
              <p className='text-error-500 text-[12px] mt-0'>
                Huston we have a blank...
              </p>
            )}
            {linkedInUseCase.invalidLinkedInUrl &&
              !linkedInUseCase.emptyLinkedInUrl && (
                <p className='text-error-500 text-[12px] mt-0'>
                  This LinkedIn appears to be invalid
                </p>
              )}
            {linkedInUseCase.error && (
              <p className='text-error-500 text-[12px] mt-0'>
                {linkedInUseCase.error}
              </p>
            )}
          </ModalBody>
          <ModalFooter className='flex w-full gap-4'>
            <Button className='w-full' onClick={() => onClose()}>
              Cancel
            </Button>
            <Button
              className='w-full'
              colorScheme='primary'
              dataTest='add-linkedin-url'
              onClick={() => {
                linkedInUseCase.submitLinkedInUrl();

                if (
                  linkedInUseCase.emptyLinkedInUrl ||
                  linkedInUseCase.invalidLinkedInUrl
                )
                  return;

                onClose();
              }}
            >
              Add & enrich
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    );
  },
);
