import { OrganizationStage } from '@graphql/types';
import { Columns03 } from '@ui/media/icons/Columns03';
import { Kbd, CommandSubItem } from '@ui/overlay/CommandMenu';
import { organizationKeywords } from '@shared/components/CommandMenu/commands';

export const StageSubItemGroup = ({
  selectedIds,
  updateStage,
  closeMenu,
}: {
  closeMenu: () => void;
  selectedIds: Array<string>;
  updateStage: (ids: Array<string>, stage: OrganizationStage) => void;
}) => {
  return (
    <>
      <CommandSubItem
        rightLabel='Lead'
        icon={<Columns03 />}
        leftLabel='Change org stage'
        rightAccessory={<Kbd className='px-1.5'>L</Kbd>}
        keywords={organizationKeywords.change_org_stage_to_lead}
        onSelectAction={() => {
          updateStage(selectedIds, OrganizationStage.Lead);
          closeMenu();
        }}
      />

      <CommandSubItem
        rightLabel='Target'
        icon={<Columns03 />}
        leftLabel='Change org stage'
        rightAccessory={<Kbd className='px-1.5'>T</Kbd>}
        keywords={organizationKeywords.change_org_stage_to_target}
        onSelectAction={() => {
          updateStage(selectedIds, OrganizationStage.Target);
          closeMenu();
        }}
      />

      <CommandSubItem
        rightLabel='Engaged'
        icon={<Columns03 />}
        leftLabel='Change org stage'
        keywords={organizationKeywords.change_org_stage_to_engaged}
        onSelectAction={() => {
          updateStage(selectedIds, OrganizationStage.Engaged);
          closeMenu();
        }}
      />
      <CommandSubItem
        rightLabel='Trial'
        icon={<Columns03 />}
        leftLabel='Change org stage'
        keywords={organizationKeywords.change_org_stage_to_trial}
        onSelectAction={() => {
          updateStage(selectedIds, OrganizationStage.Trial);
          closeMenu();
        }}
      />

      <CommandSubItem
        icon={<Columns03 />}
        rightLabel='Unqualified'
        leftLabel='Change org stage'
        keywords={organizationKeywords.change_org_stage_to_not_a_fit}
        onSelectAction={() => {
          updateStage(selectedIds, OrganizationStage.Unqualified);
          closeMenu();
        }}
      />
    </>
  );
};
