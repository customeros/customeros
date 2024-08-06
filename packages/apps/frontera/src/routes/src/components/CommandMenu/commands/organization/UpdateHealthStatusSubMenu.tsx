import { CommandSubItem } from '@ui/overlay/CommandMenu';
import { OpportunityRenewalLikelihood } from '@graphql/types';
import { AlignHorizontalCentre02 } from '@ui/media/icons/AlignHorizontalCentre02';

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
        leftLabel='Change health status'
        keywords={['renewal likelihood']}
        icon={<AlignHorizontalCentre02 />}
        onSelectAction={() => {
          updateHealth(selectedIds, OpportunityRenewalLikelihood.HighRenewal);
          closeMenu();
        }}
      />

      <CommandSubItem
        rightLabel='Medium'
        leftLabel='Change health status'
        keywords={['renewal likelihood']}
        icon={<AlignHorizontalCentre02 />}
        onSelectAction={() => {
          updateHealth(selectedIds, OpportunityRenewalLikelihood.MediumRenewal);
          closeMenu();
        }}
      />

      <CommandSubItem
        rightLabel='Low'
        leftLabel='Change health status'
        keywords={['renewal likelihood']}
        icon={<AlignHorizontalCentre02 />}
        onSelectAction={() => {
          updateHealth(selectedIds, OpportunityRenewalLikelihood.LowRenewal);
          closeMenu();
        }}
      />
    </>
  );
};
