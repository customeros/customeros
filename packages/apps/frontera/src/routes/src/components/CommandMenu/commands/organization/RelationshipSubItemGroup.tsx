import { CommandSubItem } from '@ui/overlay/CommandMenu';
import { OrganizationRelationship } from '@graphql/types';
import { AlignHorizontalCentre02 } from '@ui/media/icons/AlignHorizontalCentre02';

import { organizationKeywords } from './keywords';

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
        keywords={organizationKeywords.change_relationship_to_customer}
        onSelectAction={() => {
          updateRelationship(selectedIds, OrganizationRelationship.Customer);
          closeMenu();
        }}
      />

      <CommandSubItem
        rightLabel='Former Customer'
        leftLabel='Change relationship'
        icon={<AlignHorizontalCentre02 />}
        keywords={organizationKeywords.change_relationship_to_former_customer}
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
        keywords={organizationKeywords.change_relationship_to_unqualified}
        onSelectAction={() => {
          updateRelationship(selectedIds, OrganizationRelationship.NotAFit);
          closeMenu();
        }}
      />
      <CommandSubItem
        rightLabel='Prospect'
        leftLabel='Change relationship'
        icon={<AlignHorizontalCentre02 />}
        keywords={organizationKeywords.change_relationship_to_prospect}
        onSelectAction={() => {
          updateRelationship(selectedIds, OrganizationRelationship.Prospect);
          closeMenu();
        }}
      />
    </>
  );
};
