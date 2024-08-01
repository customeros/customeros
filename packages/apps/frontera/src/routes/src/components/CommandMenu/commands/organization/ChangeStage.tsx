import { observer } from 'mobx-react-lite';

import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { OrganizationStage } from '@graphql/types';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';
import {
  stageOptions,
  getStageOptions,
} from '@organization/components/Tabs/panels/AboutPanel/util.ts';

export const ChangeStage = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const organization = store.organizations.value.get(context.id as string);
  const label = `Organization - ${organization?.value?.name}`;
  const selectedStageOption = stageOptions.find(
    (option) => option.value === organization?.value.stage,
  );

  const applicableStageOptions = getStageOptions(
    organization?.value?.relationship,
  );

  const handleSelect = (value: OrganizationStage) => () => {
    if (!context.id) return;

    if (!organization) return;
    organization?.update((org) => {
      org.stage = value;

      return org;
    });
    store.ui.commandMenu.toggle('ChangeStage');
  };

  return (
    <Command label='Change Stage'>
      <CommandInput label={label} placeholder='Change stage...' />

      <Command.List>
        {applicableStageOptions.map((option) => (
          <CommandItem
            key={option.value}
            onSelect={handleSelect(option.value)}
            rightAccessory={
              selectedStageOption?.value === option.value ? <Check /> : null
            }
          >
            {option.label}
          </CommandItem>
        ))}
      </Command.List>
    </Command>
  );
});
