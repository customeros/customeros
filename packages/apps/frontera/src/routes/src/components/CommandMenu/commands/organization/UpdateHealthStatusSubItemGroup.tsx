import { Activity } from '@ui/media/icons/Activity';
import { CommandSubItem } from '@ui/overlay/CommandMenu';
import { OpportunityRenewalLikelihood } from '@graphql/types';
import { organizationKeywords } from '@shared/components/CommandMenu/commands';

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
      />
    </>
  );
};
