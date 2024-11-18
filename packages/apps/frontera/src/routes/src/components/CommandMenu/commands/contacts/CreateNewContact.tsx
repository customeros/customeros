import { useSearchParams } from 'react-router-dom';
import React, { useRef, useMemo, useState, useEffect } from 'react';

import Fuse from 'fuse.js';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import {
  Command,
  CommandItem,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

type FieldType = 'linkedin' | 'organizationId';

export const CreateNewContact = observer(() => {
  const hasSubmitedRef = useRef(false);
  const [linkedin, setLinkedin] = useState('');
  const [organizationId, setOrganizationId] = useState<string>('');
  const [validation, setValidation] = useState<Record<FieldType, boolean>>({
    linkedin: false,
    organizationId: false,
  });
  const [search, setSearch] = useState('');
  const [filteredOrganizations, setFilteredOrganizations] = useState<
    Array<{ id: string; name: string }>
  >([]);

  const [searchParams] = useSearchParams();

  const preset = searchParams.get('preset');

  const store = useStore();

  const contactsPreset = store.tableViewDefs.contactsPreset;
  const contactsTargetPreset = store.tableViewDefs.contactsTargetPreset;
  const contactsFlowPreset = store.tableViewDefs.contactsFlowsPreset;

  useEffect(() => {
    hasSubmitedRef?.current && validate();
  }, [linkedin, organizationId]);

  if (contactsFlowPreset === preset) {
    return null;
  }

  const organizationsList = useMemo(
    () =>
      store.organizations.toArray().map((org) => ({
        id: org.value.metadata.id,
        name: org.value.name,
      })),
    [store.organizations],
  );

  const fuse = useMemo(() => {
    return new Fuse(organizationsList, {
      keys: ['name'],
      threshold: 0.3,
    });
  }, [organizationsList]);

  const validate = () => {
    setValidation(() => ({
      linkedin: !linkedin,
      organizationId: !organizationId && contactsTargetPreset === preset,
    }));

    return linkedin && (organizationId || contactsTargetPreset !== preset);
  };

  const handleSubmit = () => {
    hasSubmitedRef.current = true;

    if (!validate()) return;

    if (contactsTargetPreset === preset || organizationId) {
      store.contacts.createWithSocial({
        organizationId,
        socialUrl: linkedin,
        options: {
          onSuccess: () => {
            handleClose();
          },
        },
      });
    } else {
      store.contacts.createWithoutOrg({
        socialUrl: linkedin,
        options: {
          onSuccess: () => {
            handleClose();
          },
        },
      });
    }
  };

  const handleClose = () => {
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.clearCallback();
  };

  const handleSearch = (v: string) => {
    const normalizedValue = v.normalize('NFD').replace(/[\u0300-\u036f]/g, '');

    setSearch(normalizedValue);

    if (normalizedValue.length > 0) {
      const results = fuse.search(normalizedValue, { limit: 20 });

      setFilteredOrganizations(results.map((v) => v.item));
    } else {
      setFilteredOrganizations([]);
    }
  };

  return (
    <Command shouldFilter={false}>
      <article className='w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold'>Create new contact</h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        <p className='text-sm mt-2 my-3'>
          We’ll auto-enrich this contact using its LinkedIn
        </p>

        <div className='flex flex-col mb-1'>
          <label htmlFor='linkedin' className='text-sm font-semibold'>
            Contact's LinkedIn URL
          </label>
          <Input
            id='linkedin'
            value={linkedin}
            placeholder='LinkedIn profile link'
            onChange={(e) => {
              setLinkedin(e.target.value);
            }}
            className={cn('border-b-0', {
              'border-b !border-error-500': validation.linkedin,
            })}
            onKeyDown={(e) => {
              if (e.key === 'Escape') {
                handleClose();
              }
              e.stopPropagation();
            }}
          />
          {validation.linkedin && (
            <p className='text-sm text-error-500 mt-1'>
              One does not simply skip LinkedIn
            </p>
          )}
        </div>

        <div className='flex flex-col'>
          <label
            htmlFor='organizationId'
            className='text-sm font-semibold mb-2'
          >
            Organization
            {contactsPreset === preset && ' (optional)'}
          </label>
          <Command.Input
            value={search}
            onValueChange={handleSearch}
            placeholder='Contact’s organization'
            onKeyDownCapture={(e) => {
              if (e.key === ' ') {
                e.stopPropagation();
              }

              if (e.metaKey && e.key === 'Enter') {
                store.ui.commandMenu.setOpen(false);
              } else {
                //handle select
              }
            }}
          />

          <Command.List className='-mx-4 '>
            {filteredOrganizations.map((orgList, idx) => {
              return (
                <CommandItem
                  key={idx}
                  onSelect={() => {
                    setOrganizationId(orgList.id);
                  }}
                >
                  {orgList.name}
                </CommandItem>
              );
            })}
          </Command.List>
          {validation.organizationId && contactsTargetPreset === preset && (
            <p className='text-sm text-error-500 mt-1'>
              Please select an organization
            </p>
          )}
        </div>
        <div className='flex justify-between gap-3 mt-6'>
          <CommandCancelButton onClose={handleClose} />

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
        </div>
      </article>
    </Command>
  );
});
