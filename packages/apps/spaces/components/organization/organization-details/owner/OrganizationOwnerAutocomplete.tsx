import React, {
  FC,
  PropsWithChildren,
  useCallback,
  useMemo,
  useState,
} from 'react';
import { useLinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useLinkOrganizationOwner';
import { useUnlinkOrganizationOwner } from '@spaces/hooks/useOrganizationOwner/useUnlinkOrganizationOwner';
import { useRecoilValue } from 'recoil';
import { ownerListData } from '../../../../state/userData';
import { Select, useSelect } from '@spaces/atoms/select';
import styles from '@spaces/organization/organization-details/relationship/organization-relationship.module.scss';
import classNames from 'classnames';
import { User } from '@spaces/graphql';
import { InlineLoader } from '@spaces/atoms/inline-loader';
import { useUsers } from '@spaces/hooks/useUser';

type Owner = Pick<User, 'id' | 'firstName' | 'lastName'> | null;
interface OrganizationOwnerProps {
  id: string;
  owner?: Owner;
}
interface SelectMenuProps {
  noOfVisibleItems?: number;
  itemSize?: number;
}
const SelectMenu = ({
  noOfVisibleItems = 7,
  itemSize = 38,
}: SelectMenuProps) => {
  const { state, getMenuProps, getMenuItemProps } = useSelect();
  const maxMenuHeight = itemSize * noOfVisibleItems;
  return (
    <ul
      className={styles.dropdownMenu}
      {...getMenuProps({ maxHeight: maxMenuHeight })}
    >
      {state.items.length ? (
        state.items.map(({ value, label }, index) => (
          <li
            key={value}
            className={classNames(styles.dropdownMenuItem, {
              [styles.isFocused]: state.currentIndex === index,
              [styles.isSelected]: state.selection === value,
            })}
            {...getMenuItemProps({ value, index })}
          >
            {label}
          </li>
        ))
      ) : (
        <li className={styles.dropdownMenuItem} data-dropdown='menuitem'>
          No options available
        </li>
      )}
    </ul>
  );
};

const SelectInput: FC<{ saving: boolean }> = ({ saving }) => {
  const { state, getInputProps, autofillValue } = useSelect();

  return (
    <>
      <span
        role='textbox'
        placeholder='Owner'
        contentEditable={state.isEditing}
        className={classNames(styles.dropdownInput)}
        {...getInputProps()}
      />
      <span className={styles.autofill}>{autofillValue}</span>
      {saving && <InlineLoader />}
    </>
  );
};

const SelectWrapper = ({ children }: PropsWithChildren) => {
  const { getWrapperProps } = useSelect();

  return (
    <div {...getWrapperProps()} className={styles.dropdownWrapper}>
      {children}
    </div>
  );
};
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
        <SelectInput saving={saving} />
        <SelectMenu />
      </SelectWrapper>
    </Select>
  );
};
