import React, { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { Input } from '@ui/form/Input';
import { Spinner } from '@ui/feedback/Spinner';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { Building07 } from '@ui/media/icons/Building07.tsx';
import { OrganizationStage, OrganizationRelationship } from '@graphql/types';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalPortal,
  ModalOverlay,
  ModalFeaturedHeader,
  ModalFeaturedContent,
} from '@ui/overlay/Modal';

interface CreateNewOrganizationModalProps {
  isOpen: boolean;
  setIsOpen: (open: boolean) => void;
}

function isValidURL(url: string) {
  const urlPattern =
    /^(https?:\/\/)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,6}(\/[a-zA-Z0-9#]+)*\/?$/;

  if (urlPattern.test(url)) {
    try {
      const parsedURL = new URL(url, 'http://example.com');

      return parsedURL.hostname.length > 0;
    } catch (e) {
      return false;
    }
  }

  return false;
}

export const CreateNewOrganizationModal: React.FC<CreateNewOrganizationModalProps> =
  observer(({ isOpen, setIsOpen }) => {
    const { organizations, tableViewDefs, ui } = useStore();
    const [searchParams] = useSearchParams();

    const [website, setWebsite] = useState('');
    const [name, setName] = useState<string>('');
    const [validation, setValidation] = useState<boolean>(false);

    const preset = searchParams?.get('preset');

    const tableViewName = tableViewDefs.getById(`${preset}`)?.value.name;

    useEffect(() => {
      if (isOpen && ui.searchCount === 0) {
        setName(searchParams.get('search') ?? '');
      }
    }, [isOpen]);

    const handleReset = () => {
      ui.setIsEditingTableCell(false);

      setWebsite('');
      setName('');
      setValidation(false);
    };

    const handleSubmit = () => {
      setValidation(false);

      if (website && !isValidURL(website)) {
        setValidation(true);

        return;
      }
      const payload = defaultValuesNewOrganization(tableViewName ?? '');

      setIsOpen(false);
      handleReset();

      organizations.create({
        ...payload,
        website,
        name,
      });
    };

    const handleClose = () => {
      handleReset();
      setIsOpen(false);
    };

    useKeyBindings(
      {
        Enter: handleSubmit,
        Escape: handleClose,
      },
      { when: isOpen },
    );

    return (
      <Modal open={isOpen}>
        <ModalPortal>
          <ModalOverlay />
          <ModalFeaturedContent>
            <ModalFeaturedHeader featuredIcon={<Building07 />}>
              <p className='text-lg font-semibold mb-1'>
                Create new organization
              </p>
              <p className='text-sm'>
                We’ll auto-enrich this organization using its website
              </p>
            </ModalFeaturedHeader>
            <ModalBody className='flex flex-col gap-4'>
              <div className='flex flex-col'>
                <label htmlFor='website' className='text-sm font-semibold'>
                  Organization's website
                </label>
                <Input
                  autoFocus
                  id='website'
                  value={website}
                  placeholder='Website link'
                  className={cn(validation && 'border-error-500')}
                  onChange={(e) => {
                    setWebsite(e.target.value);
                  }}
                />
                {validation && (
                  <p className='text-sm text-error-500 mt-1'>
                    Please insert a valid URL
                  </p>
                )}
              </div>

              <div className='flex flex-col'>
                <label htmlFor='name' className='text-sm font-semibold'>
                  Organization name
                </label>
                <Input
                  id='name'
                  value={name}
                  placeholder='Organization name'
                  defaultValue={searchParams.get('name') ?? ''}
                  onChange={(e) => {
                    setName(e.target.value);
                  }}
                />
              </div>
            </ModalBody>
            <ModalFooter className='flex gap-3'>
              <Button className='w-full' onClick={handleClose}>
                Close
              </Button>

              <Button
                className='w-full'
                colorScheme='primary'
                onClick={handleSubmit}
                isLoading={organizations.isLoading}
                loadingText='Creating organization'
                spinner={
                  <Spinner
                    size='sm'
                    label='loading'
                    className='text-primary-500 fill-primary-200'
                  />
                }
              >
                Create organization
              </Button>
            </ModalFooter>
          </ModalFeaturedContent>
        </ModalPortal>
      </Modal>
    );
  });

const defaultValuesNewOrganization = (organizationName: string) => {
  switch (organizationName) {
    case 'Customers':
      return {
        relationship: OrganizationRelationship.Customer,
        stage: OrganizationStage.Onboarding,
      };
    case 'Leads':
      return {
        relationship: OrganizationRelationship.Prospect,
        stage: OrganizationStage.Lead,
      };
    case 'Nurture':
      return {
        relationship: OrganizationRelationship.Prospect,
        stage: OrganizationStage.Target,
      };
    case 'All orgs':
      return {
        relationship: OrganizationRelationship.Prospect,
        stage: OrganizationStage.Target,
      };

    case 'Churn':
      return {
        relationship: OrganizationRelationship.FormerCustomer,
        stage: OrganizationStage.PendingChurn,
      };
    default:
      return {};
  }
};
