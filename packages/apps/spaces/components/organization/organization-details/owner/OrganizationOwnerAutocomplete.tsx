import React, { useCallback, useMemo, useState } from 'react';
import { useLinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useLinkOrganizationOwner';
import { useUnlinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useUnlinkOrganizationOwner';
import { useRecoilValue } from 'recoil';
import { ownerListData } from '../../../../state/userData';
import { User } from '@spaces/graphql';
import { useUsers } from '@spaces/hooks/useUser';

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
  const [prevSelection, setPrevSelection] = useState(owner);
  const { onLinkOrganizationOwner, saving } = useLinkOrganizationOwner({
    organizationId: id,
  });
  const { onUnlinkOrganizationOwner } = useUnlinkOrganizationOwner({
    organizationId: id,
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

  return (
    <Select
      size='sm'
      isClearable
      value={value}
      isLoading={saving}
      variant='unstyled'
      placeholder='Owner'
      backspaceRemovesValue
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
          display: 'none',
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
