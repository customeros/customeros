import { type RootStore } from '@store/root';

import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { Archive } from '@ui/media/icons/Archive';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { ColumnView } from '@shared/types/__generated__/graphql.types';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

import { getFieldTypes } from '../../Organizations/filedTypes';
interface CustomFieldItemProps {
  store: RootStore;
  field: ColumnView;
}

export const CustomFieldItem = ({ field, store }: CustomFieldItemProps) => {
  const customField = getFieldTypes(store);
  const fieldName = customField[field.columnType]?.fieldName;
  const fieldIcon = customField[field.columnType]?.icon;
  const fieldType = customField[field.columnType]?.fieldTypeName;

  return (
    <div className='flex justify-between items-center py-2 w-full px-2'>
      <div className='flex justify-between w-full'>
        <div className='flex items-center gap-2 flex-2'>
          {fieldIcon}
          {fieldName}
        </div>
        <div className='flex items-center justify-between flex-1'>
          {fieldType}
          <Menu>
            <MenuButton asChild>
              <IconButton
                size='xs'
                variant='ghost'
                aria-label='Edit field'
                icon={<DotsVertical />}
              />
            </MenuButton>
            <MenuList>
              <MenuItem className='group/edit'>
                <div className='flex items-center'>
                  <Edit03 className='mr-2 text-gray-500 group-hover/edit:text-gray-700' />
                  Edit field
                </div>
              </MenuItem>
              <MenuItem className='group/archive'>
                <div className='flex items-center'>
                  <Archive className='mr-2 group-hover/archive:text-gray-700 text-gray-500' />
                  Archive field
                </div>
              </MenuItem>
            </MenuList>
          </Menu>
        </div>
      </div>
    </div>
  );
};
