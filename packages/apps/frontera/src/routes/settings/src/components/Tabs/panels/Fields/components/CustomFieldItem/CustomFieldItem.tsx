import React, { useState } from 'react';

import { type RootStore } from '@store/root';

import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { Archive } from '@ui/media/icons/Archive';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import {
  ColumnView,
  CustomField,
} from '@shared/types/__generated__/graphql.types';

import { CustomFieldModal } from '../CustomFieldModal';
import {
  getCustomFieldTypes,
  getDefaultFieldTypes,
} from '../../Organizations/filedTypes';

interface CustomFieldItemProps {
  store: RootStore;
  isEditable?: boolean;
  field: ColumnView | CustomField;
}

export const CustomFieldItem = ({
  field,
  store,
  isEditable = false,
}: CustomFieldItemProps) => {
  const { onOpen, onToggle, open } = useDisclosure();
  const [isEdit, setIsEdit] = useState<boolean>(false);

  const defaultCustomField = getDefaultFieldTypes(store);
  const customField = getCustomFieldTypes();

  const isColumnView = (
    field: ColumnView | CustomField,
  ): field is ColumnView => {
    return (field as ColumnView).columnType !== undefined;
  };

  const customFieldStore = store.customFields.value.get(
    (field as CustomField).id || '',
  );

  const fieldData = isColumnView(field)
    ? defaultCustomField[field.columnType]
    : customField[field.template?.type as keyof typeof customField];

  const fieldName = isColumnView(field)
    ? (fieldData as { fieldName: string })?.fieldName
    : customFieldStore?.value.name;

  const fieldIcon = fieldData?.icon;
  const fieldType = fieldData?.fieldTypeName;

  return (
    <>
      <div className='flex justify-between items-center py-2 w-full px-2'>
        <div className='flex justify-between w-full'>
          <div className='flex items-center gap-2 flex-2'>
            {fieldIcon}
            {fieldName}
          </div>
          <div className='flex items-center justify-between flex-1'>
            {fieldType}
            <Menu>
              <MenuButton asChild disabled={!isEditable}>
                <IconButton
                  size='xs'
                  variant='ghost'
                  aria-label='Edit field'
                  icon={<DotsVertical />}
                  isDisabled={!isEditable}
                />
              </MenuButton>
              <MenuList>
                <MenuItem
                  className='group/edit'
                  onClick={() => {
                    setIsEdit(true);
                    onOpen();
                  }}
                >
                  <div className='flex items-center'>
                    <Edit03 className='mr-2 text-gray-500 group-hover/edit:text-gray-700' />
                    Edit field
                  </div>
                </MenuItem>
                <MenuItem
                  className='group/archive'
                  onClick={() => {
                    store.customFields.deleteCustomField(
                      (field as CustomField)?.id,
                    );
                  }}
                >
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
      {isEdit && (
        <CustomFieldModal
          isOpen={open}
          isEdit={isEdit}
          onOpenChange={onToggle}
          fieldId={(field as CustomField)?.id}
        />
      )}
    </>
  );
};
