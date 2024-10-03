import { useRef, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { Plus } from '@ui/media/icons/Plus';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { Select, getContainerClassNames } from '@ui/form/Select';
import { OrganizationStage } from '@shared/types/__generated__/graphql.types';
import {
  Modal,
  ModalBody,
  ModalClose,
  ModalPortal,
  ModalFooter,
  ModalOverlay,
  ModalTrigger,
  ModalCloseButton,
  ModalFeaturedHeader,
  ModalFeaturedContent,
} from '@ui/overlay/Modal/Modal';

type FieldType = 'linkedin' | 'organizationId';

export const ContactAvatarHeader = observer(() => {
  const hasSubmitedRef = useRef(false);
  const [isOpen, setIsOpen] = useState(false);
  const [linkedin, setLinkedin] = useState('');
  const [searchValue, setSearchValue] = useState('');
  const [organizationId, setOrganizationId] = useState<string>('');
  const [validation, setValidation] = useState<Record<FieldType, boolean>>({
    linkedin: false,
    organizationId: false,
  });

  const enableFeature = useFeatureIsOn('gp-dedicated-1');
  const store = useStore();

  const options = store?.organizations
    ?.toComputedArray((arr) => {
      const targets = arr.filter(
        (item) => item.value.stage === OrganizationStage.Target,
      );

      if (!searchValue) return targets;

      return targets.filter((item) => item.value.name.includes(searchValue));
    })
    .map((item) => ({
      label: item.value.name,
      value: item.getId(),
    }));

  const validate = () => {
    setValidation(() => ({
      linkedin: !linkedin,
      organizationId: !organizationId,
    }));

    return linkedin && organizationId;
  };

  const handleSubmit = () => {
    hasSubmitedRef.current = true;

    if (!validate()) return;

    store.contacts.createWithSocial({
      organizationId,
      socialUrl: linkedin,
      options: {
        onSuccess: () => {
          setIsOpen(false);
          reset();
        },
      },
    });
  };

  const reset = () => {
    setLinkedin('');
    setOrganizationId('');
    setValidation({
      linkedin: false,
      organizationId: false,
    });
    hasSubmitedRef.current = false;
  };

  useEffect(() => {
    hasSubmitedRef?.current && validate();
  }, [linkedin, organizationId]);

  return (
    <Modal
      open={isOpen}
      onOpenChange={(open) => {
        setIsOpen(open);

        if (!open) reset();
      }}
    >
      <ModalTrigger asChild>
        <div className='flex w-[26px] items-center justify-center'>
          <Tooltip
            asChild
            side='bottom'
            align='center'
            label='Create contact'
            className={cn(enableFeature ? 'visible' : 'hidden')}
          >
            <IconButton
              size='xxs'
              variant='ghost'
              aria-label='create contact'
              data-test='create-contact-from-table'
              icon={<Plus className='text-gray-400 size-5' />}
              className={cn('size-6', enableFeature ? 'visible' : 'hidden')}
            />
          </Tooltip>
        </div>
      </ModalTrigger>

      <ModalPortal>
        <ModalOverlay />
        <ModalFeaturedContent>
          <ModalFeaturedHeader>
            <p className='text-lg font-semibold mb-1'>Create new contact</p>
            <p className='text-sm'>
              We’ll auto-enrich this contact using its LinkedIn
            </p>
          </ModalFeaturedHeader>
          <ModalCloseButton />
          <ModalBody className='flex flex-col gap-4'>
            <div className='flex flex-col'>
              <label htmlFor='linkedin' className='text-sm font-semibold'>
                Contact's LinkedIn URL
              </label>
              <Input
                id='linkedin'
                value={linkedin}
                placeholder='LinkedIn profile link'
                className={cn(validation.linkedin && 'border-error-500')}
                onChange={(e) => {
                  setLinkedin(e.target.value);
                }}
                onKeyDown={(e) => {
                  if (e.key === 'Escape') {
                    setIsOpen(false);
                  }
                  e.stopPropagation();
                }}
              />
              {validation.linkedin && (
                <p className='text-sm text-error-500 mt-1'>
                  Please insert a LinkedIn URL
                </p>
              )}
            </div>

            <div className='flex flex-col'>
              <label htmlFor='organizationId' className='text-sm font-semibold'>
                Organization
              </label>
              <Select
                isClearable
                options={options}
                id='organizationId'
                backspaceRemovesValue
                onInputChange={setSearchValue}
                placeholder='Contact’s organization'
                onChange={(value) => {
                  setOrganizationId(value?.value);
                }}
                onKeyDown={(e) => {
                  if (e.key === 'Escape') {
                    setIsOpen(false);
                  }
                }}
                noOptionsMessage={({ inputValue }) => {
                  if (!inputValue) return 'Type to search orgs';

                  return `No org found with name "${inputValue}"`;
                }}
                classNames={{
                  container: (props) =>
                    getContainerClassNames(
                      cn(validation.organizationId && 'border-error-500'),
                      'flushed',
                      props,
                    ),
                }}
              />
              {validation.organizationId && (
                <p className='text-sm text-error-500 mt-1'>
                  Please select an organization
                </p>
              )}
            </div>
          </ModalBody>
          <ModalFooter className='flex gap-3'>
            <ModalClose className='w-full'>
              <Button className='w-full'>Cancel</Button>
            </ModalClose>

            <Button
              className='w-full'
              colorScheme='primary'
              onClick={handleSubmit}
              loadingText='Creating contact'
              isLoading={store.contacts.isLoading}
              rightSpinner={
                <Spinner
                  size='sm'
                  label='loading'
                  className='text-primary-500 fill-primary-200'
                />
              }
            >
              Create
            </Button>
          </ModalFooter>
        </ModalFeaturedContent>
      </ModalPortal>
    </Modal>
  );
});
