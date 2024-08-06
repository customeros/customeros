import { CommandSubItem } from '@ui/overlay/CommandMenu';
import { OrganizationRelationship } from '@graphql/types';
import { AlignHorizontalCentre02 } from '@ui/media/icons/AlignHorizontalCentre02';

export const RelationshipSubMenu = ({
  selectedIds,
  updateRelationship,
  closeMenu,
}: {
  closeMenu: () => void;
  selectedIds: Array<string>;
  updateRelationship: (
    ids: Array<string>,
    relationship: OrganizationRelationship,
  ) => void;
}) => {
  return (
    <>
      <CommandSubItem
        rightLabel='Customer'
        icon={<AlignHorizontalCentre02 />}
        leftLabel='Change org relationship'
        onSelectAction={() => {
          updateRelationship(selectedIds, OrganizationRelationship.Customer);
          closeMenu();
        }}
      />

      <CommandSubItem
        rightLabel='Former Customer'
        icon={<AlignHorizontalCentre02 />}
        leftLabel='Change org relationship'
        onSelectAction={() => {
          updateRelationship(
            selectedIds,
            OrganizationRelationship.FormerCustomer,
          );
          closeMenu();
        }}
      />

      <CommandSubItem
        rightLabel='Not A Fit'
        icon={<AlignHorizontalCentre02 />}
        leftLabel='Change org relationship'
        onSelectAction={() => {
          updateRelationship(selectedIds, OrganizationRelationship.NotAFit);
          closeMenu();
        }}
      />
      <CommandSubItem
        rightLabel='Prospect'
        icon={<AlignHorizontalCentre02 />}
        leftLabel='Change org relationship'
        onSelectAction={() => {
          updateRelationship(selectedIds, OrganizationRelationship.Prospect);
          closeMenu();
        }}
      />
    </>
  );
};
