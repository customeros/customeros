import { observer } from 'mobx-react-lite';

import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { CommandSubItem } from '@ui/overlay/CommandMenu';
import { SelectOption } from '@shared/types/SelectOptions';
import { Certificate02 } from '@ui/media/icons/Certificate02';

const roleOptions = [
  { value: 'Decision Maker', label: 'Decision Maker' },
  { value: 'Influencer', label: 'Influencer' },
  { value: 'User', label: 'User' },
  { value: 'Stakeholder', label: 'Stakeholder' },
  { value: 'Gatekeeper', label: 'Gatekeeper' },
  { value: 'Champion', label: 'Champion' },
  { value: 'Data Owner', label: 'Data Owner' },
];

export const AddJobRolesSubItemGroup = observer(() => {
  const store = useStore();

  const context = store.ui.commandMenu.context;

  const contact = store.contacts.value.get(context.ids?.[0] as string);
  const selectedIds = context.ids;

  const handleSelect = (opt: SelectOption[]) => {
    if (!context.ids?.[0] || !contact) return;

    if (selectedIds?.length === 1) {
      contact.value.jobRoles.map((jobRole) => {
        jobRole.description = opt.map((v) => v.value).join(',');
      });
    } else {
      selectedIds.forEach((id) => {
        const contact = store.contacts.value.get(id);

        if (contact) {
          contact.value.jobRoles.map((jobRole) => {
            jobRole.description = opt.map((v) => v.value).join(',');
          });
        }
      });
    }
    contact.commit();
  };

  useModKey('Enter', () => {
    store.ui.commandMenu.setOpen(false);
  });

  return (
    <>
      {roleOptions.map((role, idx) => {
        const selectedDescriptions =
          contact?.value?.jobRoles?.[0]?.description?.split(',') || [];

        const isSelected = selectedDescriptions.includes(role.value);

        return (
          <CommandSubItem
            key={idx}
            rightLabel={role.label}
            icon={<Certificate02 />}
            leftLabel='Change job role'
            rightAccessory={isSelected ? <Check /> : undefined}
            onSelectAction={() => {
              const newSelections = isSelected
                ? selectedDescriptions.filter((desc) => desc !== role.value)
                : [...selectedDescriptions, role.value];

              const newOptions = roleOptions.filter((r) =>
                newSelections.includes(r.value),
              );

              handleSelect(newOptions);
              store.ui.commandMenu.setOpen(false);
            }}
          />
        );
      })}
    </>
  );
});
