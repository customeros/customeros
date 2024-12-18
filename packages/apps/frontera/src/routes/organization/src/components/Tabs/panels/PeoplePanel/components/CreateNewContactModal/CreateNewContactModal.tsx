import { useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { CreateContact } from '@domain/Contacts/CreateContact.useCase';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { Mail01 } from '@ui/media/icons/Mail01';
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
  orgId: string;
}
const nameUseCase = new CreateContact();

export const CreateNewContactModal = observer(
  ({ orgId }: CreateNewContactModalProps) => {
    const store = useStore();
    const [createType, setCreateType] = useState<'linkedin' | 'email' | 'name'>(
      'linkedin',
    );

    const inputPlaceholder =
      createType === 'linkedin'
        ? 'linkedin.com/in/johnlemon'
        : createType === 'email'
        ? 'john@heyjude.band'
        : 'First and last name';

    const confirmButtonPlaceholder =
      createType === 'name' ? 'Add contact' : 'Add & enrich';

    console.log(nameUseCase.inputValue);

    return (
      <Modal open={false}>
        <ModalPortal>
          <ModalOverlay />
          <ModalContent>
            <ModalHeader>
              <p>Add a contact using their...</p>
              <ModalCloseButton asChild />
            </ModalHeader>
            <ModalBody className='flex flex-col w-full gap-4'>
              <ButtonGroup className=' flex w-full'>
                <Button
                  size='xs'
                  leftIcon={<LinkedinOutline />}
                  data-inactive={createType !== 'linkedin'}
                  onClick={() => setCreateType('linkedin')}
                  className='w-full data-[inactive=true]:bg-gray-50 focus:bg-white'
                >
                  LinkedIn
                </Button>
                <Button
                  size='xs'
                  leftIcon={<Mail01 />}
                  data-inactive={createType !== 'email'}
                  onClick={() => setCreateType('email')}
                  className='w-full data-[inactive=true]:bg-gray-50 focus:bg-white'
                >
                  Email
                </Button>
                <Button
                  size='xs'
                  leftIcon={<Signature />}
                  data-inactive={createType !== 'name'}
                  onClick={() => setCreateType('name')}
                  className='w-full data-[inactive=true]:bg-gray-50 focus:bg-white'
                >
                  Name
                </Button>
              </ButtonGroup>

              <Input
                variant='unstyled'
                placeholder={inputPlaceholder}
                onChange={(e) => nameUseCase.setInputValue(e.target.value)}
              />
            </ModalBody>
            <ModalFooter className='w-full flex gap-3'>
              <ModalCloseButton asChild>
                <Button size='sm' className='w-full'>
                  Cancel
                </Button>
              </ModalCloseButton>
              <Button
                size='sm'
                className='w-full'
                colorScheme='primary'
                onClick={() => {
                  if (createType === 'name') {
                    store.contacts.create(
                      orgId,
                      {},
                      { name: nameUseCase.inputValue },
                    );
                  }
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
