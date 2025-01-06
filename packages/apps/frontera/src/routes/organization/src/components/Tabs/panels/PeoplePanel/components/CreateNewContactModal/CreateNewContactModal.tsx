import { useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { CreateContact } from '@domain/usecases/people-contact-card/create-contact.usecase';

import { cn } from '@ui/utils/cn.ts';
import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { Signature } from '@ui/media/icons/Signature';
import { LinkedinOutline } from '@ui/media/icons/LinkedinOutline';
import {
  Modal,
  ModalBody,
  ModalPortal,
  ModalHeader,
  ModalFooter,
  ModalContent,
  ModalOverlay,
  ModalCloseButton,
} from '@ui/overlay/Modal';

interface CreateNewContactModalProps {
  open: boolean;
  orgId: string;
  onClose: () => void;
}
const contactCreate = new CreateContact();

export const CreateNewContactModal = observer(
  ({ orgId, open, onClose }: CreateNewContactModalProps) => {
    const store = useStore();
    const inputPlaceholder =
      contactCreate.getType === 'linkedin'
        ? 'linkedin.com/in/johnlemon'
        : contactCreate.getType === 'email'
        ? 'john@heyjude.band'
        : 'First and last name';
    const org = store.organizations.getById(orgId);
    const confirmButtonPlaceholder =
      contactCreate.getType === 'name' ? 'Add contact' : 'Add & enrich';

    useEffect(() => {
      if (orgId && org) {
        contactCreate.setEntity(org);
      }
    }, [orgId]);

    return (
      <Modal open={open}>
        <ModalPortal>
          <ModalOverlay />
          <ModalContent>
            <ModalHeader>
              <p>Add a contact using their...</p>
              <ModalCloseButton asChild />
            </ModalHeader>
            <ModalBody className='flex flex-col w-full gap-4'>
              <ButtonGroup className='flex items-center w-full'>
                <Button
                  size='xs'
                  leftIcon={<LinkedinOutline />}
                  onClick={() => contactCreate.setType('linkedin')}
                  data-inactive={contactCreate.getType !== 'linkedin'}
                  className={cn('w-full', {
                    selected: contactCreate.getType === 'linkedin',
                  })}
                >
                  LinkedIn
                </Button>
                {/* <Button
                  size='xs'
                  leftIcon={<Mail01 />}
                  onClick={() => contactCreate.setType('email')}
                  data-inactive={contactCreate.getType !== 'email'}
                  className='w-full data-[inactive=true]:bg-gray-50 focus:bg-white'
                >
                  Email
                </Button> */}
                <Button
                  size='xs'
                  leftIcon={<Signature />}
                  dataTest='org-people-add-by-name'
                  onClick={() => contactCreate.setType('name')}
                  data-inactive={contactCreate.getType !== 'name'}
                  className={cn('w-full', {
                    selected: contactCreate.getType === 'name',
                  })}
                >
                  Name
                </Button>
              </ButtonGroup>

              <Input
                variant='unstyled'
                placeholder={inputPlaceholder}
                onChange={(e) => contactCreate.setInputValue(e.target.value)}
              />
            </ModalBody>
            <ModalFooter className='w-full flex gap-3'>
              <ModalCloseButton asChild>
                <Button size='sm' className='w-full' onClick={() => onClose()}>
                  Cancel
                </Button>
              </ModalCloseButton>
              <Button
                size='sm'
                className='w-full'
                colorScheme='primary'
                dataTest='org-people-add-contact'
                onClick={() => {
                  contactCreate.setOrganizationId(orgId);
                  contactCreate.submit();
                  onClose();
                }}
              >
                {confirmButtonPlaceholder}
              </Button>
            </ModalFooter>
          </ModalContent>
        </ModalPortal>
      </Modal>
    );
  },
);
