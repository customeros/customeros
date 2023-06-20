import React, { useCallback, useMemo, useState } from 'react';
import { useLinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useLinkOrganizationOwner';
import { useUnlinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useUnlinkOrganizationOwner';
import { useRecoilValue } from 'recoil';
import { ownerListData } from '../../../../state/userData';
import { User } from '@spaces/graphql';
import { useUsers } from '@spaces/hooks/useUser';
import {
  Select,
  SelectMenu,
  SelectInput,
  SelectWrapper,
} from '@spaces/ui/form/select';

type Owner = Pick<User, 'id' | 'firstName' | 'lastName'> | null;
interface OrganizationOwnerProps {
  id: string;
  owner?: Owner;
}

export const OrganizationOwnerAutocomplete: React.FC<
  OrganizationOwnerProps
> = ({ id, owner }) => {
  useUsers();
  const [prevSelection, setPrevSelection] = useState<any>(owner);
  const { onLinkOrganizationOwner, saving } = useLinkOrganizationOwner({
    organizationId: id,
  });
  const { onUnlinkOrganizationOwner } = useUnlinkOrganizationOwner({
    organizationId: id,
  });
  const { ownerList } = useRecoilValue(ownerListData);

  const ownerOptionsList = useMemo(() => {
    return ownerList;
  }, [ownerList]);

  const handleSelect = useCallback(
    (ownerId: string) => {
      if (!ownerId && prevSelection) {
        onUnlinkOrganizationOwner();
      } else {
        onLinkOrganizationOwner({
          userId: ownerId,
          name: ownerList?.find((e) => e.value === ownerId).label || '',
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
    <Select<string>
      onSelect={handleSelect}
      value={owner?.id}
      options={ownerOptionsList}
    >
      <SelectWrapper>
        <SelectInput saving={saving} placeholder='Owner' />
        <SelectMenu />
      </SelectWrapper>
    </Select>
  );
};
