import { observer } from 'mobx-react-lite';

import { Button } from '@ui/form/Button/Button';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';

export const ChangeContactOrganizationModal = observer(() => {
  return (
    <Modal>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>
          <h1>Change organization</h1>
        </ModalHeader>
        <ModalBody>
          <div>
            <p>
              Are you sure you want to change the organization of this contact?
            </p>
          </div>
        </ModalBody>
        <ModalFooter>
          <Button variant='ghost'>Cancel</Button>
          <Button>Change organization</Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
});
