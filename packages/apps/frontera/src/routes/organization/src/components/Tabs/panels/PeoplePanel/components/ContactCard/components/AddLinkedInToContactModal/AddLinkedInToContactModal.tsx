import { useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { LinkedIn } from '@domain/usecases/people-contact-card/add-linkedin.usecase';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
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
    const contactStore = store.contacts.value.get(contactId);

    useEffect(() => {
      if (contactStore) {
        linkedInUseCase.setEntity(contactStore);
      }
    }, [contactId]);

    return (
      <Modal open={open}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>
            <p className='font-medium'>LinkedIn profile URL</p>
            <ModalCloseButton asChild />
          </ModalHeader>
          <ModalBody className='flex flex-col gap-4'>
            <p>We'll auto-enrich this contact using their LinkedIn profile</p>
            <Input
              variant='unstyled'
              value={linkedInUseCase.inputValue}
              placeholder='linkedin.com/in/john-lemon'
              onChange={(e) => linkedInUseCase.setInputValue(e.target.value)}
            />
          </ModalBody>
          <ModalFooter className='flex w-full gap-4'>
            <Button className='w-full' onClick={() => onClose()}>
              Cancel
            </Button>
            <Button
              className='w-full'
              colorScheme='primary'
              onClick={() => {
                linkedInUseCase.setLinkedInUrl();
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
