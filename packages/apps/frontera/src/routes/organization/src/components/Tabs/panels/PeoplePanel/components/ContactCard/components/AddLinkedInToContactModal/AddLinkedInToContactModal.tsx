import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
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
  onClose: () => void;
}

export const AddLinkedInToContactModal = observer(
  ({ onClose, open }: AddLinkedInToContactModalProps) => {
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
              placeholder='linkedin.com/in/john-lemon'
            />
          </ModalBody>
          <ModalFooter className='flex w-full gap-4'>
            <Button className='w-full' onClick={() => onClose()}>
              Cancel
            </Button>
            <Button className='w-full' colorScheme='primary'>
              Add & enrich
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    );
  },
);
