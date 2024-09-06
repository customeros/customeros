import React, { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useDidMount, useKeyBindings } from 'rooks';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Command } from '@ui/overlay/CommandMenu';
import { OrganizationStage, OrganizationRelationship } from '@graphql/types';

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

export const AddNewOrganization = observer(() => {
  const store = useStore();
  const [allowSubmit, setAllowSubmit] = useState(false);
  const { organizations, tableViewDefs, ui } = useStore();
  const [searchParams] = useSearchParams();

  const [website, setWebsite] = useState('');
  const [name, setName] = useState<string>('');
  const [validation, setValidation] = useState<boolean>(false);

  const preset = searchParams?.get('preset');

  const tableViewName = tableViewDefs.getById(`${preset}`)?.value.name;

  useEffect(() => {
    if (ui.searchCount === 0) {
      setName(searchParams.get('search') ?? '');
    }
  }, []);

  useDidMount(() => {
    setTimeout(() => {
      setAllowSubmit(true);
    }, 100);
  });

  const handleConfirm = () => {
    if (!allowSubmit) return;
    setAllowSubmit(false);

    setValidation(false);

    if (website && !isValidURL(website)) {
      setValidation(true);
      setAllowSubmit(true);

      return;
    }
    const payload = defaultValuesNewOrganization(tableViewName ?? '');

    organizations.create(
      {
        ...payload,
        website,
        name,
      },
      {},
    );

    store.ui.commandMenu.toggle('AddNewOrganization');
  };

  useKeyBindings(
    {
      Enter: handleConfirm,
      Escape: () => (store.ui.commandMenu.isOpen = false),
    },
    { when: allowSubmit },
  );

  return (
    <Command label={`Rename `}>
      <div className='p-6 pb-4 flex flex-col gap-1 '>
        <p className='font-semibold'>Create new organization</p>
        <p className='text-sm'>
          We’ll auto-enrich this organization using its website
        </p>
      </div>

      <div className='pr-6 pl-6 pb-6 flex flex-col gap-2 '>
        <div className='flex flex-col mb-0'>
          <label htmlFor='website' className='text-sm font-semibold'>
            Organization's website (optional)
          </label>
          <Input
            autoFocus
            id='website'
            value={website}
            variant='unstyled'
            placeholder='Website link'
            onChange={(e) => {
              setWebsite(e.target.value);
            }}
            onKeyDown={(e) => {
              if (e.key === 'Escape') {
                store.ui.commandMenu.setOpen(false);
              }
            }}
            onKeyUp={(e) => {
              if (e.key === 'Backspace' && website.length === 0) {
                setValidation(false);
              }
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
            variant='unstyled'
            placeholder='Organization name'
            defaultValue={searchParams.get('name') ?? ''}
            data-test='organizations-create-new-org-org-name'
            onChange={(e) => {
              setName(e.target.value);
            }}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }

              if (e.key === 'Escape') {
                store.ui.commandMenu.setOpen(false);
              }
            }}
          />
        </div>
      </div>

      <div className='flex w-full gap-2 pl-6 pr-6 pb-6'>
        <Button
          className='w-full'
          onClick={() => {
            store.ui.commandMenu.setOpen(false);
            store.ui.commandMenu.setType('OrganizationCommands');
          }}
        >
          Cancel
        </Button>
        <Button
          className='w-full'
          colorScheme='primary'
          onClick={handleConfirm}
        >
          Create
        </Button>
      </div>
    </Command>
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
