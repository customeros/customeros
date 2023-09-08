import { Organization } from '@spaces/graphql';

import { Icons } from '@ui/media/Icon';
import { IconButton } from '@ui/form/IconButton';
import {
  TableInstance,
  RowSelectionState,
} from '@spaces/ui/presentation/Table';
import { Menu, MenuButton, MenuList, MenuItem } from '@ui/overlay/Menu';
import { useRecoilState } from 'recoil';
import { tableMode } from '@spaces/finder/state';
import { SvgIcon } from '@ui/media/icons/SvgIcon';
import Archive from '@ui/media/icons/Archive';

interface OrganizationListActionsProps {
  table: TableInstance<Organization>;
  isSelectionEnabled: boolean;
  selection: RowSelectionState;
  toggleSelection: (value: boolean) => void;
  onMergeOrganizations: (table: TableInstance<Organization>) => void;
  onArchiveOrganizations: (table: TableInstance<Organization>) => void;
  onCreateOrganization: () => void;
}

const OrganizationListActions = ({
  table,
  selection,
  toggleSelection,
  isSelectionEnabled,
  onMergeOrganizations,
  onArchiveOrganizations,
  onCreateOrganization,
}: OrganizationListActionsProps) => {
  const [mode, setMode] = useRecoilState(tableMode);

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
          icon={
            <SvgIcon boxSize='4'>
              <Archive />
            </SvgIcon>
          }
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
export default OrganizationListActions;
