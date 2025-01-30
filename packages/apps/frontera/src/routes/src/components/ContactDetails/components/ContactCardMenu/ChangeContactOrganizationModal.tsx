import { observer } from 'mobx-react-lite';

import { Combobox } from '@ui/form/Combobox';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Popover, PopoverContent, PopoverTrigger } from '@ui/overlay/Popover';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
  ModalCloseButton,
} from '@ui/overlay/Modal';

interface ChangeContactOrganizationModalProps {
  open: boolean;
  contactId: string;
  onClose: () => void;
}

export const ChangeContactOrganizationModal = observer(
  ({ onClose, open, contactId }: ChangeContactOrganizationModalProps) => {
    const store = useStore();

    const contactStore = store.contacts.value.get(contactId);

    if (!contactStore) return null;

    return (
      <Modal open={open}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>
            <p className='font-medium'>{`Change ${contactStore?.name}'s company`}</p>
            <ModalCloseButton asChild />
          </ModalHeader>
          <ModalBody>
            <div>
              <p>
                Changing this contact’s company will associate them with the
                newly selected company.
              </p>
            </div>

            <Popover modal>
              <PopoverTrigger className='w-full flex items-start justify-start mt-4 text-gray-400'>
                <p className='text-start w-full'>Contact's new company</p>
              </PopoverTrigger>
              <PopoverContent
                align='center'
                className='min-w-[264px] max-w-[320px] overflow-auto z-[99999]'
              >
                <Combobox
                  options={store.organizations.toArray().map((o) => ({
                    label: o.value.name,
                    value: o.value.id,
                  }))}
                  onChange={(value) => {
                    contactStore?.draft();
                    contactStore.value.primaryOrganizationName = value.label;

                    contactStore?.commit();

                    onClose();
                  }}
                />
              </PopoverContent>
            </Popover>
          </ModalBody>
          <ModalFooter className='flex gap-4'>
            <Button
              variant='outline'
              className='w-full'
              onClick={() => onClose()}
            >
              Cancel
            </Button>
            <Button className='w-full' colorScheme='primary'>
              Change company
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    );
  },
);
