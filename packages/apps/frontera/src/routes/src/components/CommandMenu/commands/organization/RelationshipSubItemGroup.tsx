import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
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
  const store = useStore();

  const isSelected = () => {
    if (selectedIds.length > 1) {
      return;
    } else {
      const organization = store.organizations.value.get(selectedIds[0]);

      return organization?.value.relationship;
    }
  };

  return (
    <>
      <CommandSubItem
        rightLabel='Customer'
        leftLabel='Change relationship'
        icon={<AlignHorizontalCentre02 />}
        keywords={organizationKeywords.change_relationship_to_customer}
        rightAccessory={
          isSelected() === OrganizationRelationship.Customer ? <Check /> : null
        }
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
        rightAccessory={
          isSelected() === OrganizationRelationship.FormerCustomer ? (
            <Check />
          ) : null
        }
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
        rightAccessory={
          isSelected() === OrganizationRelationship.NotAFit ? <Check /> : null
        }
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
        rightAccessory={
          isSelected() === OrganizationRelationship.Prospect ? <Check /> : null
        }
        onSelectAction={() => {
          updateRelationship(selectedIds, OrganizationRelationship.Prospect);
          closeMenu();
        }}
      />
    </>
  );
};
