import { Organization } from '@spaces/graphql';

import { Icons } from '@ui/media/Icon';
import { Button } from '@ui/form/Button';
import { IconButton } from '@ui/form/IconButton';
import { TableInstance, RowSelectionState } from '@ui/presentation/_Table';
import { Menu, MenuButton, MenuList, MenuItem } from '@ui/overlay/Menu';

interface OrganizationListActionsProps {
  table: TableInstance<Organization>;
  isSelectionEnabled: boolean;
  selection: RowSelectionState;
  toggleSelection: (value: boolean) => void;
  onMergeOrganizations: (table: TableInstance<Organization>) => void;
  onCreateOrganization: () => void;
}

const OrganizationListActions = ({
  table,
  selection,
  toggleSelection,
  isSelectionEnabled,
  onMergeOrganizations,
  onCreateOrganization,
}: OrganizationListActionsProps) => {
  if (isSelectionEnabled) {
    if (Object.keys(selection).length > 1) {
      return (
        <Button
          variant='outline'
          onClick={() => onMergeOrganizations(table)}
          aria-label='Merge Organizations'
        >
          Merge
        </Button>
      );
    }
    return (
      <Button
        aria-label='Merge organizations'
        variant='outline'
        onClick={() => {
          toggleSelection(false);
          table.resetRowSelection();
        }}
      >
        Done
      </Button>
    );
  }

  return (
    <Menu>
      <MenuButton
        size='sm'
        variant='ghost'
        as={IconButton}
        aria-label='Table Actions'
        icon={<Icons.DotsVertical color='gray.400' boxSize='6' />}
      />
      <MenuList boxShadow='xl'>
        <MenuItem onClick={onCreateOrganization}>Add organization</MenuItem>
        <MenuItem
          onClick={() => {
            toggleSelection(true);
          }}
        >
          Merge organizations
        </MenuItem>
      </MenuList>
    </Menu>
  );
};
export default OrganizationListActions;
