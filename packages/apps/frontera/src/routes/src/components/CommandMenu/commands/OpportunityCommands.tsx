import React from 'react';

import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { User01 } from '@ui/media/icons/User01';
import { Archive } from '@ui/media/icons/Archive';
import { useStore } from '@shared/hooks/useStore';
import { Delete } from '@ui/media/icons/Delete.tsx';
import { Command } from '@ui/media/icons/Command.tsx';
import { Columns03 } from '@ui/media/icons/Columns03';
import { Calculator } from '@ui/media/icons/Calculator';
import { ArrowsRight } from '@ui/media/icons/ArrowsRight';
import { Kbd, CommandItem } from '@ui/overlay/CommandMenu';
import { ArrowBlockUp } from '@ui/media/icons/ArrowBlockUp.tsx';
import { CurrencyDollarCircle } from '@ui/media/icons/CurrencyDollarCircle';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';

export const OpportunityCommands = observer(() => {
  const store = useStore();
  const opportunity = store.opportunities.value.get(
    store.ui.commandMenu.context.ids?.[0] as string,
  );
  const label = `Opportunity - ${opportunity?.value.name}`;

  return (
    <CommandsContainer label={label}>
      <>
        <CommandItem
          leftAccessory={<Columns03 />}
          onSelect={() => {
            store.ui.commandMenu.setType('ChangeStage');
          }}
          keywords={[
            'change',
            'stage',
            'edit',
            'update',
            'status',
            'pipeline',
            'phase',
          ]}
          rightAccessory={
            <>
              <Kbd>
                <ArrowBlockUp className='text-inherit size-3' />
              </Kbd>
              <Kbd>S</Kbd>
            </>
          }
        >
          Change stage...
        </CommandItem>

        <CommandItem
          leftAccessory={<ArrowsRight />}
          onSelect={() => {
            store.ui.commandMenu.setType('SetOpportunityNextSteps');
          }}
          keywords={[
            'set',
            'next',
            'step',
            'edit',
            'update',
            'action',
            'reminder',
            'follow up',
            'task',
          ]}
        >
          Set next step
        </CommandItem>

        <CommandItem
          leftAccessory={<Calculator />}
          onSelect={() => {
            store.ui.commandMenu.setType('ChangeArrEstimate');
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
          leftAccessory={<Edit03 />}
          onSelect={() => {
            store.ui.commandMenu.setType('RenameOpportunityName');
          }}
          rightAccessory={
            <>
              <Kbd>
                <ArrowBlockUp className='text-inherit size-3' />
              </Kbd>
              <Kbd>R</Kbd>
            </>
          }
          keywords={[
            'rename',
            'opportunity',
            'edit',
            'change',
            'update',
            'deal',
            'name',
            'title',
            'label',
          ]}
        >
          Rename opportunity
        </CommandItem>

        <CommandItem
          leftAccessory={<User01 />}
          keywords={['assign', 'owner', 'change', 'update', 'edit']}
          onSelect={() => {
            store.ui.commandMenu.setType('AssignOwner');
          }}
          rightAccessory={
            <>
              <Kbd>
                <ArrowBlockUp className='text-inherit size-3' />
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
              <Kbd>
                <Command className='size-3' />
              </Kbd>
              <Kbd>
                <Delete className='size-3' />
              </Kbd>
            </>
          }
        >
          Archive opportunity
        </CommandItem>
      </>
    </CommandsContainer>
  );
});
