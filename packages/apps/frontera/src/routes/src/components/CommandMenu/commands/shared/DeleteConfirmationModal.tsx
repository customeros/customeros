import { useNavigate } from 'react-router-dom';
import React, { useRef, useEffect } from 'react';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { ContactStore } from '@store/Contacts/Contact.store.ts';
import { OpportunityStore } from '@store/Opportunities/Opportunity.store.ts';
import { TableViewDefStore } from '@store/TableViewDefs/TableViewDef.store.ts';
import { OrganizationStore } from '@store/Organizations/Organization.store.ts';

import { TableIdType } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const DeleteConfirmationModal = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const navigate = useNavigate();

  const confirmButtonRef = useRef<HTMLButtonElement>(null);
  const closeButtonRef = useRef<HTMLButtonElement>(null);

  const entity = match(context.entity)
    .returnType<
      | OpportunityStore
      | OrganizationStore
      | TableViewDefStore
      | ContactStore
      | undefined
    >()
    .with('Opportunity', () => store.opportunities.value.get(context.ids?.[0]))
    .with('Organization', () => store.organizations.value.get(context.ids?.[0]))
    .with('Contact', () => store.contacts.value.get(context.ids?.[0]))
    .with('TableViewDef', () => store.tableViewDefs.getById(context.ids?.[0]))
    .otherwise(() => undefined);

  const handleClose = () => {
    store.ui.commandMenu.toggle('DeleteConfirmationModal');
    store.ui.commandMenu.clearCallback();
  };

  const handleConfirm = () => {
    match(context.entity)
      .with('Organization', () => {
        const oppoortunityOfOrgSelected = store.opportunities
          .toArray()
          .filter((o) => o.value.organization?.metadata.id === context.ids[0]);

        const oppotunityIdOfOrgSelected = oppoortunityOfOrgSelected.map(
          (o) => o.value.id,
        );

        store.organizations.hide(context.ids as string[]);
        store.opportunities.value.delete(oppotunityIdOfOrgSelected[0]);

        context.callback?.();
      })
      .with('Organizations', () => {
        store.organizations.hide(context.ids as string[]);
        context.callback?.();
      })
      .with('Contact', () => {
        store.contacts.archive(context.ids);
        context.callback?.();
      })
      .with('Opportunity', () => {
        store.opportunities.archive(context.ids?.[0]);
        context.callback?.();
      })
      .with('Opportunities', () => {
        store.opportunities.archiveMany(context.ids);
        context.callback?.();
      })
      .with('Sequence', () => {
        store.flowSequences.archive(context.ids?.[0]);
        context.callback?.();
      })
      .with('Sequences', () => {
        store.flowSequences.archiveMany(context.ids);
        context.callback?.();
      })
      .with('TableViewDef', () => {
        store.tableViewDefs.archive(context.ids?.[0], {
          onSuccess: () => {
            const allOrgsViewId = store.tableViewDefs
              ?.toArray()
              .find((e) => e.value.tableId === TableIdType.Organizations)
              ?.value.id;

            navigate(`/finder?preset=${allOrgsViewId}`);
          },
        });
      })
      .otherwise(() => {});

    handleClose();
  };

  const title = match(context.entity)
    .with(
      'Organization',
      () =>
        `Archive ${(entity as OrganizationStore)?.value.name || 'Unnamed'}?`,
    )
    .with(
      'Organizations',
      () => `Archive ${context.ids?.length} organizations?`,
    )
    .with(
      'Opportunities',
      () => `Archive ${context.ids?.length} opportunities?`,
    )
    .with(
      'Opportunity',
      () => `Archive ${(entity as OpportunityStore)?.value.name}?`,
    )
    .with('Sequences', () => `Archive ${context.ids.length} sequences?`)
    .with(
      'Sequence',
      () =>
        `Archive ${
          store.flowSequences.value.get(context.ids?.[0])?.value.name
        }?`,
    )

    .with('Contact', () =>
      context.ids?.length > 1
        ? `Archive ${context.ids?.length} contacts?`
        : `Archive ${(entity as ContactStore)?.value.name}?`,
    )
    .with(
      'TableViewDef',
      () => `Archive ${(entity as TableViewDefStore)?.value.name}?`,
    )
    .otherwise(() => `Archive selected ${context.entity?.toLowerCase()}`);
  const description = match(context.entity)
    .with(
      'Sequences',
      () =>
        `Archiving these sequences will remove all contacts currently enrolled to it`,
    )
    .with(
      'Sequence',
      () =>
        `Archiving this sequence will remove all contacts currently enrolled to it`,
    )

    .otherwise(() => null);

  useEffect(() => {
    closeButtonRef.current?.focus();
  }, []);

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold'>{title}</h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>
        {description && <p className='mt-1 text-sm'>{description}</p>}

        <div className='flex justify-between gap-3 mt-6'>
          <CommandCancelButton ref={closeButtonRef} onClose={handleClose} />

          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='error'
            ref={confirmButtonRef}
            onClick={handleConfirm}
            data-test='org-actions-confirm-archive'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Archive
          </Button>
        </div>
      </article>
    </Command>
  );
});
