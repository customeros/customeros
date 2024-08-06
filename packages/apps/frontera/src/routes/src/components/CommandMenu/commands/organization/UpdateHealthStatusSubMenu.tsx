import { Activity } from '@ui/media/icons/Activity';
import { CommandSubItem } from '@ui/overlay/CommandMenu';
import { OpportunityRenewalLikelihood } from '@graphql/types';

export const UpdateHealthStatusSubMenu = ({
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
        keywords={['renewal likelihood']}
        onSelectAction={() => {
          updateHealth(selectedIds, OpportunityRenewalLikelihood.HighRenewal);
          closeMenu();
        }}
      />

      <CommandSubItem
        rightLabel='Medium'
        icon={<Activity />}
        leftLabel='Change health status'
        keywords={['renewal likelihood']}
        onSelectAction={() => {
          updateHealth(selectedIds, OpportunityRenewalLikelihood.MediumRenewal);
          closeMenu();
        }}
      />

      <CommandSubItem
        rightLabel='Low'
        icon={<Activity />}
        leftLabel='Change health status'
        keywords={['renewal likelihood']}
        onSelectAction={() => {
          updateHealth(selectedIds, OpportunityRenewalLikelihood.LowRenewal);
          closeMenu();
        }}
      />
    </>
  );
};
