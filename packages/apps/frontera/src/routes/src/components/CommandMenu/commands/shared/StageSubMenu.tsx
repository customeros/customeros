import { OrganizationStage } from '@graphql/types';
import { Columns03 } from '@ui/media/icons/Columns03';
import { CommandSubItem } from '@ui/overlay/CommandMenu';

export const StageSubMenu = ({
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
        onSelectAction={() => {
          updateStage(selectedIds, OrganizationStage.Lead);
          closeMenu();
        }}
      />

      <CommandSubItem
        rightLabel='Target'
        icon={<Columns03 />}
        leftLabel='Change org stage'
        onSelectAction={() => {
          updateStage(selectedIds, OrganizationStage.Target);
          closeMenu();
        }}
      />

      <CommandSubItem
        rightLabel='Engaged'
        icon={<Columns03 />}
        leftLabel='Change org stage'
        onSelectAction={() => {
          updateStage(selectedIds, OrganizationStage.Engaged);
          closeMenu();
        }}
      />
      <CommandSubItem
        rightLabel='Trial'
        icon={<Columns03 />}
        leftLabel='Change org stage'
        onSelectAction={() => {
          updateStage(selectedIds, OrganizationStage.Trial);
          closeMenu();
        }}
      />
    </>
  );
};
