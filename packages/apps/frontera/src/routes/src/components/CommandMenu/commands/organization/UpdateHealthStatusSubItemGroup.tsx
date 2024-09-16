import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { Activity } from '@ui/media/icons/Activity';
import { CommandSubItem } from '@ui/overlay/CommandMenu';
import { OpportunityRenewalLikelihood } from '@graphql/types';

import { organizationKeywords } from './keywords';

export const UpdateHealthStatusSubItemGroup = ({
  selectedIds,
  updateHealth,
  closeMenu,
}: {
  closeMenu: () => void;
  selectedIds: Array<string>;
  updateHealth: (
    ids: Array<string>,
    health: OpportunityRenewalLikelihood,
  ) => void;
}) => {
  const store = useStore();

  const isSelected = () => {
    if (selectedIds.length > 1) {
      return;
    } else {
      const organization = store.organizations.value.get(selectedIds[0]);

      return organization?.value.accountDetails?.renewalSummary
        ?.renewalLikelihood;
    }
  };

  return (
    <>
      <CommandSubItem
        rightLabel='High'
        icon={<Activity />}
        leftLabel='Change health status'
        keywords={organizationKeywords.change_health_status}
        onSelectAction={() => {
          updateHealth(selectedIds, OpportunityRenewalLikelihood.HighRenewal);
          closeMenu();
        }}
        rightAccessory={
          isSelected() === OpportunityRenewalLikelihood.HighRenewal ? (
            <Check />
          ) : null
        }
      />

      <CommandSubItem
        rightLabel='Medium'
        icon={<Activity />}
        leftLabel='Change health status'
        keywords={organizationKeywords.change_health_status}
        onSelectAction={() => {
          updateHealth(selectedIds, OpportunityRenewalLikelihood.MediumRenewal);
          closeMenu();
        }}
        rightAccessory={
          isSelected() === OpportunityRenewalLikelihood.MediumRenewal ? (
            <Check />
          ) : null
        }
      />

      <CommandSubItem
        rightLabel='Low'
        icon={<Activity />}
        leftLabel='Change health status'
        keywords={organizationKeywords.change_health_status}
        onSelectAction={() => {
          updateHealth(selectedIds, OpportunityRenewalLikelihood.LowRenewal);
          closeMenu();
        }}
        rightAccessory={
          isSelected() === OpportunityRenewalLikelihood.LowRenewal ? (
            <Check />
          ) : null
        }
      />
    </>
  );
};
