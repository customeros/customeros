import { useRef, useEffect } from 'react';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';
import { CreateContact } from '@domain/usecases/people-contact-card/create-contact.usecase';

import { cn } from '@ui/utils/cn.ts';
import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { useModKey } from '@shared/hooks/useModKey';
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
    const modalRef = useRef<HTMLDivElement>(null);
    const inputRef = useRef<HTMLInputElement>(null);
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

    useModKey(
      'Enter',
      () => {
        handleSubmit();
      },
      { targetRef: modalRef, when: open },
    );
    useKey(
      'Escape',
      () => {
        onClose();
      },
      { target: modalRef, when: open },
    );

    const handleSubmit = () => {
      contactCreate.setOrganizationId(orgId);

      if (contactCreate.getType === 'name') {
        contactCreate.submit();
        !contactCreate.invalidName && onClose();
      }

      if (contactCreate.getType === 'linkedin') {
        contactCreate.submit();
        !contactCreate.emptyLinkedInUrl &&
          !contactCreate.invalidLinkedInUrl &&
          onClose();
      }
    };

    useEffect(() => {
      setTimeout(() => {
        if (inputRef.current) {
          inputRef.current.focus();
        }
      }, 100);
    }, [open, contactCreate.getType]);

    return (
      <Modal open={open} onOpenChange={onClose}>
        <ModalPortal>
          <ModalOverlay />
          <ModalContent ref={modalRef}>
            <ModalHeader>
              <p className='font-medium'>Add a contact using their...</p>
              <ModalCloseButton asChild />
            </ModalHeader>
            <ModalBody className='flex flex-col w-full'>
              <div className='flex flex-col gap-4'>
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
                  ref={inputRef}
                  variant='unstyled'
                  placeholder={inputPlaceholder}
                  dataTest='org-people-name-input'
                  onKeyDown={(e) => {
                    if (e.key === 'Escape') {
                      onClose();
                    }
                    e.stopPropagation();
                  }}
                  onChange={(e) => {
                    contactCreate.setInputValue(e.target.value);

                    if (contactCreate.inputValue) {
                      contactCreate.getType === 'name' &&
                        contactCreate.validateName();
                    }
                  }}
                />
              </div>

              {contactCreate.getType === 'name' && (
                <p
                  className={cn(
                    'text-error-500 text-[12px] mt-0 opacity-0',
                    contactCreate.invalidName && 'opacity-100',
                  )}
                >
                  Every hero needs a name
                </p>
              )}

              {contactCreate.getType === 'linkedin' && (
                <>
                  {!contactCreate.inputValue && (
                    <p
                      className={cn(
                        'text-error-500 text-[12px] mt-0 opacity-0',
                        contactCreate.emptyLinkedInUrl && 'opacity-100',
                      )}
                    >
                      Huston we have a blank...
                    </p>
                  )}

                  {contactCreate.inputValue && (
                    <p
                      className={cn(
                        'text-error-500 text-[12px] mt-0 opacity-0',
                        contactCreate.invalidLinkedInUrl && 'opacity-100',
                      )}
                    >
                      Invalid LinkedIn URL
                    </p>
                  )}
                </>
              )}
            </ModalBody>
            <ModalFooter className='w-full flex gap-3'>
              <ModalCloseButton asChild>
                <Button
                  size='sm'
                  className='w-full'
                  onClick={() => {
                    contactCreate.clearState();
                    onClose();
                  }}
                >
                  Cancel
                </Button>
              </ModalCloseButton>
              <Button
                size='sm'
                className='w-full'
                colorScheme='primary'
                dataTest='org-people-add-new-contact'
                onClick={() => {
                  handleSubmit();
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
