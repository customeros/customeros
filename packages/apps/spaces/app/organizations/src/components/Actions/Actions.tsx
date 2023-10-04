import { useState } from 'react';

import { TableInstance, RowSelectionState } from '@ui/presentation/Table';
import { Icons } from '@ui/media/Icon';
import { Organization } from '@graphql/types';
import { IconButton } from '@ui/form/IconButton';
import { Archive } from '@ui/media/icons/Archive';
import { Menu, MenuButton, MenuList, MenuItem } from '@ui/overlay/Menu';

interface OrganizationListActionsProps {
  table: TableInstance<Organization>;
  isSelectionEnabled: boolean;
  selection: RowSelectionState;
  toggleSelection: (value: boolean) => void;
  onMergeOrganizations: (table: TableInstance<Organization>) => void;
  onArchiveOrganizations: (table: TableInstance<Organization>) => void;
  onCreateOrganization: () => void;
}

export const OrganizationListActions = ({
  table,
  selection,
  toggleSelection,
  isSelectionEnabled,
  onMergeOrganizations,
  onArchiveOrganizations,
  onCreateOrganization,
}: OrganizationListActionsProps) => {
  const [mode, setMode] = useState<'MERGE' | 'ARCHIVE' | null>(null);

  if (isSelectionEnabled) {
    if (Object.keys(selection).length > 1 && mode === 'MERGE') {
      return (
        <IconButton
          size='xs'
          variant='ghost'
          colorScheme='green'
          onClick={() => onMergeOrganizations(table)}
          aria-label='Merge Organizations'
          icon={<Icons.Check boxSize='4' />}
        />
      );
    }
    if (Object.keys(selection).length >= 1 && mode === 'ARCHIVE') {
      return (
        <IconButton
          size='xs'
          variant='ghost'
          colorScheme='red'
          onClick={() => onArchiveOrganizations(table)}
          aria-label='Archive Organizations'
          icon={<Archive />}
        />
      );
    }

    return (
      <IconButton
        size='xs'
        aria-label='Discard'
        variant='ghost'
        onClick={() => {
          toggleSelection(false);
          table.resetRowSelection();
        }}
        icon={<Icons.XClose boxSize='4' color='gray.400' />}
      />
    );
  }

  return (
    <Menu>
      <MenuButton
        size='xs'
        variant='ghost'
        as={IconButton}
        aria-label='Table Actions'
        icon={<Icons.DotsVertical color='gray.400' boxSize='4' />}
      />
      <MenuList boxShadow='xl'>
        <MenuItem onClick={onCreateOrganization}>Add organization</MenuItem>
        <MenuItem
          onClick={() => {
            toggleSelection(true);
            setMode('MERGE');
          }}
        >
          Merge organizations
        </MenuItem>
        <MenuItem
          onClick={() => {
            toggleSelection(true);
            setMode('ARCHIVE');
          }}
        >
          Archive organizations
        </MenuItem>
      </MenuList>
    </Menu>
  );
};
