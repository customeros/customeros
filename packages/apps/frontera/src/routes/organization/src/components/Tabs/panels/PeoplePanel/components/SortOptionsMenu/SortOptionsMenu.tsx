import { observer } from 'mobx-react-lite';
import { SearchSortContact } from '@domain/usecases/contact-details/search-sort-contacts.usecase';

import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { SwitchVertical01 } from '@ui/media/icons/SwitchVertical01';
import { SwitchVertical02 } from '@ui/media/icons/SwitchVertical02';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

interface SortOptionsMenuProps {
  searchSortContact: SearchSortContact;
}

export const SortOptionsMenu = observer(
  ({ searchSortContact }: SortOptionsMenuProps) => {
    const sortDirection = searchSortContact.getSortDirection();

    return (
      <div className='flex'>
        <IconButton
          size='xs'
          variant='ghost'
          aria-label='sort direction'
          onClick={() => {
            searchSortContact.setSortDirection(
              sortDirection === 'asc' ? 'desc' : 'asc',
            );
          }}
          icon={
            searchSortContact.getSortDirection() === 'asc' ? (
              <SwitchVertical01 />
            ) : (
              <SwitchVertical02 />
            )
          }
        />
        <Menu>
          <MenuButton asChild>
            <Button size='xs' variant='ghost'>
              {searchSortContact.getSort()}
            </Button>
          </MenuButton>
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
