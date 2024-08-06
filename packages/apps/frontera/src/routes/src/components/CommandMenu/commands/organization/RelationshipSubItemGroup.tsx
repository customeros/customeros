import { CommandSubItem } from '@ui/overlay/CommandMenu';
import { OrganizationRelationship } from '@graphql/types';
import { AlignHorizontalCentre02 } from '@ui/media/icons/AlignHorizontalCentre02';

export const RelationshipSubItemGroup = ({
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
        leftLabel='Change relationship'
        icon={<AlignHorizontalCentre02 />}
        onSelectAction={() => {
          updateRelationship(selectedIds, OrganizationRelationship.Customer);
          closeMenu();
        }}
      />

      <CommandSubItem
        rightLabel='Former Customer'
        leftLabel='Change relationship'
        icon={<AlignHorizontalCentre02 />}
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
        leftLabel='Change relationship'
        icon={<AlignHorizontalCentre02 />}
        onSelectAction={() => {
          updateRelationship(selectedIds, OrganizationRelationship.NotAFit);
          closeMenu();
        }}
      />
      <CommandSubItem
        rightLabel='Prospect'
        leftLabel='Change relationship'
        icon={<AlignHorizontalCentre02 />}
        onSelectAction={() => {
          updateRelationship(selectedIds, OrganizationRelationship.Prospect);
          closeMenu();
        }}
      />
    </>
  );
};
