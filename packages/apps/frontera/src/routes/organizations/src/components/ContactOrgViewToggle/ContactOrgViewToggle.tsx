import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

import { cn } from '@ui/utils/cn.ts';
import { useStore } from '@shared/hooks/useStore';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { Button } from '@ui/form/Button/Button.tsx';
import { TableIdType, TableViewType } from '@graphql/types';

export const ContactOrgViewToggle = observer(() => {
  const store = useStore();
  const [searchParams, setSearchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const [tabs, setLastActivePosition] = useLocalStorage<{
    [key: string]: string;
  }>('customeros-player-last-position', { root: 'organizations' });
  const search = searchParams.get('search');
  const [lastSearchForPreset, setLastSearchForPreset] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-last-search-for-preset`, { root: 'root' });

  const tableViewDefs = store.tableViewDefs.toArray();

  const tableViewDef = store.tableViewDefs.getById(preset || '')?.value;
  const tableViewId = tableViewDef?.tableId;
  const tableViewType = tableViewDef?.tableType;

  const tableDefsMap = tableViewDefs.reduce((map, def) => {
    map[def.value.tableId] = def.value.id;

    return map;
  }, {} as { [key: string]: string });

  const nurtureTableDef = tableDefsMap[TableIdType.Nurture];
  const contactsForTargetOrgsTableDef =
    tableDefsMap[TableIdType.ContactsForTargetOrganizations];
  const orgTableDef = tableDefsMap[TableIdType.Organizations];
  const contactsTableDef = tableDefsMap[TableIdType.Contacts];

  const getTargetTable = () => {
    switch (tableViewId) {
      case TableIdType.Nurture:
      case TableIdType.ContactsForTargetOrganizations:
        return nurtureTableDef;
      case TableIdType.Organizations:
      case TableIdType.Contacts:
        return orgTableDef;
      default:
        return null;
    }
  };

  const getContactTable = () => {
    switch (tableViewId) {
      case TableIdType.Nurture:
      case TableIdType.ContactsForTargetOrganizations:
        return contactsForTargetOrgsTableDef;
      case TableIdType.Organizations:
      case TableIdType.Contacts:
        return contactsTableDef;
      default:
        return null;
    }
  };

  const targetTableDef = getTargetTable();
  const contactTableDef = getContactTable();

  const handleNavigate = (newPreset: string) => {
    const newParams = new URLSearchParams(searchParams.toString());

    newParams.set('preset', newPreset);
    setSearchParams(newParams);
    setLastActivePosition({
      ...tabs,
      root: `finder?preset=${newPreset}`,
    });

    if (preset) {
      setLastSearchForPreset({
        ...lastSearchForPreset,
        [preset]: search ?? '',
      });
    }
  };

  const showToggle =
    (tableViewType && [TableViewType.Contacts].includes(tableViewType)) ||
    (tableViewId &&
      [TableIdType.Organizations, TableIdType.Nurture].includes(tableViewId));

  return (
    <>
      {showToggle && (
        <ButtonGroup className='flex items-center'>
          <Button
            size='xs'
            onClick={() => targetTableDef && handleNavigate(targetTableDef)}
            className={cn('bg-white !border-r px-4', {
              'bg-gray-50 text-gray-500 font-normal': preset !== targetTableDef,
            })}
          >
            Orgs
          </Button>
          <Button
            size='xs'
            onClick={() => contactTableDef && handleNavigate(contactTableDef)}
            className={cn('bg-white px-4', {
              'bg-gray-50 text-gray-500 font-normal':
                preset !== contactTableDef,
            })}
          >
            Contacts
          </Button>
        </ButtonGroup>
      )}
    </>
  );
});
