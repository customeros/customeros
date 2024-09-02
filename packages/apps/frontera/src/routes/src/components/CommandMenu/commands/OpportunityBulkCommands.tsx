import React from 'react';

import { observer } from 'mobx-react-lite';

import { User01 } from '@ui/media/icons/User01';
import { Archive } from '@ui/media/icons/Archive';
import { useStore } from '@shared/hooks/useStore';
import { Delete } from '@ui/media/icons/Delete.tsx';
import { Columns03 } from '@ui/media/icons/Columns03';
import { Calculator } from '@ui/media/icons/Calculator.tsx';
import { ArrowBlockUp } from '@ui/media/icons/ArrowBlockUp.tsx';
import { Kbd, CommandKbd, CommandItem } from '@ui/overlay/CommandMenu';
import { CurrencyDollarCircle } from '@ui/media/icons/CurrencyDollarCircle';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';
import { organizationKeywords } from '@shared/components/CommandMenu/commands/organization';

export const OpportunityBulkCommands = observer(() => {
  const store = useStore();
  const selectedIds = store.ui.commandMenu.context.ids;

  const label = `${selectedIds?.length} opportunities`;

  return (
    <CommandsContainer label={label}>
      <>
        <CommandItem
          leftAccessory={<Columns03 />}
          onSelect={() => {
            store.ui.commandMenu.setType('ChangeStage');
          }}
        >
          Change stage...
        </CommandItem>

        <CommandItem
          leftAccessory={<Calculator />}
          onSelect={() => {
            store.ui.commandMenu.setType('ChangeBulkArrEstimate');
          }}
          keywords={[
            'change',
            'arr',
            'estimate',
            'edit',
            'update',
            'annual',
            'recurring',
            'revenue',
            'forecast',
            'projection',
          ]}
        >
          Change ARR estimate
        </CommandItem>

        <CommandItem
          leftAccessory={<CurrencyDollarCircle />}
          onSelect={() => {
            store.ui.commandMenu.setType('ChangeCurrency');
          }}
          keywords={[
            'change',
            'arr',
            'currency',
            'edit',
            'update',
            'annual',
            'recurring',
            'revenue',
            'usd',
            'eur',
            'gbp',
            'dollar',
            'euro',
            'great british pound',
          ]}
        >
          Change ARR currency...
        </CommandItem>

        <CommandItem
          leftAccessory={<User01 />}
          keywords={organizationKeywords.assign_owner}
          onSelect={() => {
            store.ui.commandMenu.setType('AssignOwner');
          }}
          rightAccessory={
            <>
              <Kbd>
                <ArrowBlockUp className='size-3' />
              </Kbd>
              <Kbd>O</Kbd>
            </>
          }
        >
          Assign owner...
        </CommandItem>

        <CommandItem
          leftAccessory={<Archive />}
          onSelect={() => {
            store.ui.commandMenu.setType('DeleteConfirmationModal');
          }}
          keywords={[
            'archive',
            'opportunity',
            'delete',
            'remove',
            'hide',
            'deal',
          ]}
          rightAccessory={
            <>
              <CommandKbd />
              <Kbd>
                <Delete className='size-3' />
              </Kbd>
            </>
          }
        >
          Archive opportunities
        </CommandItem>
      </>
    </CommandsContainer>
  );
});
