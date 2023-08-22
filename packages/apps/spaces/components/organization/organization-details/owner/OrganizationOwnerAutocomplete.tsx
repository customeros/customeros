import React, { useCallback, useMemo, useState } from 'react';
import { useLinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useLinkOrganizationOwner';
import { useUnlinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useUnlinkOrganizationOwner';
import { useRecoilValue } from 'recoil';
import { ownerListData } from '../../../../state/userData';
import { User } from '@spaces/graphql';
import { useUsers } from '@spaces/hooks/useUser';

import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { IconButton } from '@ui/form/IconButton';

import { Select } from '@ui/form/SyncSelect/Select';
import { SelectOption } from '@shared/types/SelectOptions';

type Owner = Pick<User, 'id' | 'firstName' | 'lastName'> | null;
interface OrganizationOwnerProps {
  id: string;
  owner?: Owner;
}

export const OrganizationOwnerAutocomplete: React.FC<
  OrganizationOwnerProps
> = ({ id, owner }) => {
  useUsers();
  const [isEditing, setIsEditing] = useState(false);
  const [prevSelection, setPrevSelection] = useState(owner);
  const { onLinkOrganizationOwner, saving } = useLinkOrganizationOwner({
    organizationId: id,
    onCompleted: () => {
      setIsEditing(false);
    },
  });
  const { onUnlinkOrganizationOwner } = useUnlinkOrganizationOwner({
    organizationId: id,
    onCompleted: () => {
      setIsEditing(false);
    },
  });
  const { ownerList } = useRecoilValue(ownerListData);
  const value = owner ? ownerList.find((o) => o.value === owner.id) : null;

  const ownerOptionsList = useMemo(() => {
    return ownerList;
  }, [ownerList]);

  const handleSelect = useCallback(
    (option: SelectOption) => {
      if (!option && !!prevSelection) {
        onUnlinkOrganizationOwner();
        setPrevSelection(null);
      }

      if (option) {
        onLinkOrganizationOwner({
          userId: option.value,
          name: ownerList?.find((e) => e.value === option.value).label || '',
        });
        setPrevSelection(owner);
      }
    },
    [
      prevSelection,
      onUnlinkOrganizationOwner,
      onLinkOrganizationOwner,
      ownerList,
      owner,
    ],
  );

  if (!isEditing) {
    return (
      <Flex
        w='full'
        gap='1'
        align='center'
        _hover={{
          '& #edit-button': {
            opacity: 1,
          },
        }}
      >
        <Text
          cursor='default'
          color={value ? 'gray.700' : 'gray.500'}
          onDoubleClick={() => setIsEditing(true)}
        >
          {value?.label ?? 'Owner'}
        </Text>
        <IconButton
          aria-label='erc'
          size='xs'
          borderRadius='md'
          minW='4'
          w='4'
          minH='4'
          h='4'
          opacity='0'
          variant='ghost'
          id='edit-button'
          onClick={() => setIsEditing(true)}
          icon={<Icons.Edit3 color='gray.500' boxSize='3' />}
        />
      </Flex>
    );
  }

  return (
    <Select
      size='sm'
      isClearable
      value={value}
      isLoading={saving}
      variant='unstyled'
      placeholder='Owner'
      autoFocus
      onKeyDown={(e) => {
        if (e.key === 'Escape') {
          setIsEditing(false);
        }
      }}
      defaultMenuIsOpen
      onBlur={() => setIsEditing(false)}
      backspaceRemovesValue
      openMenuOnClick={false}
      onChange={handleSelect}
      options={ownerOptionsList}
      chakraStyles={{
        valueContainer: (props) => ({
          ...props,
          p: 0,
        }),
        singleValue: (props) => ({
          ...props,
          paddingBottom: 0,
          ml: 0,
        }),
        control: (props) => ({
          ...props,
          minH: '0',
        }),
        clearIndicator: (props) => ({
          ...props,
          boxSize: '3',
        }),
        placeholder: (props) => ({
          ...props,
          ml: 0,
          color: 'gray.500',
        }),
        inputContainer: (props) => ({
          ...props,
          py: 0,
          ml: 0,
        }),
      }}
    />
  );
};
