import { observer } from 'mobx-react-lite';
import { SearchSortContact } from '@domain/usecases/people-contact-card/search-sort-contacts.usecase';

import { IconButton } from '@ui/form/IconButton';
import { SwitchVertical01 } from '@ui/media/icons/SwitchVertical01';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

interface SortOptionsMenuProps {
  searchSortContact: SearchSortContact;
}

export const SortOptionsMenu = observer(
  ({ searchSortContact }: SortOptionsMenuProps) => {
    const sortDirection = searchSortContact.getSortDirection();

    return (
      <div className='flex gap-2'>
        <IconButton
          size='xs'
          variant='ghost'
          aria-label='sort direction'
          icon={<SwitchVertical01 />}
          onClick={() => {
            searchSortContact.setSortDirection(
              sortDirection === 'asc' ? 'desc' : 'asc',
            );
          }}
        />
        <Menu>
          <MenuButton>{searchSortContact.getSort()}</MenuButton>
          <MenuList>
            <MenuItem onClick={() => searchSortContact.setSort('First name')}>
              First name
            </MenuItem>
            <MenuItem onClick={() => searchSortContact.setSort('Created')}>
              Created
            </MenuItem>
            <MenuItem onClick={() => searchSortContact.setSort('Updated')}>
              Updated
            </MenuItem>
            <MenuItem onClick={() => searchSortContact.setSort('Tenure')}>
              Tenure
            </MenuItem>
          </MenuList>
        </Menu>
      </div>
    );
  },
);
