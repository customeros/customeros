import React, { useState } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { Input } from '@ui/form/Input';
import { Spinner } from '@ui/feedback/Spinner';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { Building07 } from '@ui/media/icons/Building07.tsx';
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

interface CreateNewOrganizationModalProps {
  isOpen: boolean;
  setIsOpen: (open: boolean) => void;
}

export const CreateNewOrganizationModal: React.FC<CreateNewOrganizationModalProps> =
  observer(({ isOpen, setIsOpen }) => {
    const { organizations } = useStore();
    const [searchParams] = useSearchParams();

    const [website, setWebsite] = useState('');
    const [name, setName] = useState<string>('');
    const [validation, setValidation] = useState<Record<string, boolean>>({
      website: false,
      name: false,
    });

    const handleReset = () => {
      setWebsite('');
      setName('');
      setValidation({
        website: false,
        name: false,
      });
    };

    const handleSubmit = () => {
      organizations.create({
        website,
        name,
      });
    };

    return (
      <Modal
        open={isOpen}
        onOpenChange={(open) => {
          setIsOpen(open);
          if (!open) handleReset();
        }}
      >
        <ModalPortal>
          <ModalOverlay />
          <ModalFeaturedContent>
            <ModalFeaturedHeader featuredIcon={<Building07 />}>
              <p className='text-lg font-semibold mb-1'>
                Create new organization
              </p>
              <p className='text-sm'>
                We’ll auto-enrich this contact using its website
              </p>
            </ModalFeaturedHeader>
            <ModalCloseButton />
            <ModalBody className='flex flex-col gap-4'>
              <div className='flex flex-col'>
                <label className='text-sm font-semibold' htmlFor='website'>
                  Organization's website
                </label>
                <Input
                  id='website'
                  value={website}
                  placeholder='Website link'
                  className={cn(validation.linkedin && 'border-error-500')}
                  onChange={(e) => {
                    setWebsite(e.target.value);
                  }}
                />
                {validation.linkedin && (
                  <p className='text-sm text-error-500 mt-1'>
                    Please insert a valid URL
                  </p>
                )}
              </div>

              <div className='flex flex-col'>
                <label className='text-sm font-semibold' htmlFor='name'>
                  Organization name
                </label>
                <Input
                  id='name'
                  value={name}
                  defaultValue={searchParams.get('name') ?? ''}
                  placeholder='Orgnaization Name'
                  onChange={(e) => {
                    setName(e.target.value);
                  }}
                />
              </div>
            </ModalBody>
            <ModalFooter className='flex gap-3'>
              <ModalClose className='w-full'>
                <Button className='w-full'>Close</Button>
              </ModalClose>

              <Button
                className='w-full'
                colorScheme='primary'
                onClick={handleSubmit}
                isLoading={organizations.isLoading}
                loadingText='Creating contact'
                spinner={
                  <Spinner
                    label='loading'
                    size='sm'
                    className='text-primary-500 fill-primary-200'
                  />
                }
              >
                Create org
              </Button>
            </ModalFooter>
          </ModalFeaturedContent>
        </ModalPortal>
      </Modal>
    );
  });
