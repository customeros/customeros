import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

import { cn } from '@ui/utils/cn.ts';
import { useStore } from '@shared/hooks/useStore';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { Button } from '@ui/form/Button/Button.tsx';
import { TableIdType, TableViewType } from '@graphql/types';

export const TableViewsToggleNavigation = observer(() => {
  const store = useStore();
  const [searchParams, setSearchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const [tabs, setLastActivePosition] = useLocalStorage<{
    [key: string]: string;
  }>('customeros-player-last-position', { root: 'organizations' });

  const tableViewDefs = store.tableViewDefs.toArray();
  const tableViewDef = store.tableViewDefs.getById(preset || '')?.value;
  const tableViewId = tableViewDef?.tableId;
  const tableViewType = tableViewDef?.tableType;

  const findPresetTable = (tableIdTypes: TableIdType[]): string | null => {
    const presetTable = tableViewDefs.find((def) =>
      tableIdTypes.includes(def.value.tableId),
    );

    return presetTable ? presetTable.value.id : null;
  };

  const getTablePair = (): [string | null, string | null] => {
    switch (tableViewId) {
      case TableIdType.Targets:

      case TableIdType.ContactsForTargetOrganizations: {
        return [
          findPresetTable([TableIdType.Targets]),
          findPresetTable([TableIdType.ContactsForTargetOrganizations]),
        ];
      }
      case TableIdType.UpcomingInvoices:
      case TableIdType.PastInvoices:
        return [
          findPresetTable([TableIdType.UpcomingInvoices]),
          findPresetTable([TableIdType.PastInvoices]),
        ];
      default:
        return [null, null];
    }
  };

  const [firstTableDef, secondTableDef] = getTablePair();

  const handleNavigate = (newPreset: string) => {
    const newParams = new URLSearchParams(searchParams.toString());

    newParams.set('preset', newPreset);
    setSearchParams(newParams);
    setLastActivePosition({
      ...tabs,
      root: `finder?preset=${newPreset}`,
    });
  };

  const showToggle =
    (tableViewType &&
      [TableViewType.Contacts, TableViewType.Invoices].includes(
        tableViewType,
      ) &&
      tableViewDef?.isPreset) ||
    (tableViewId &&
      [
        TableIdType.Organizations,
        TableIdType.Targets,
        TableIdType.UpcomingInvoices,
        TableIdType.PastInvoices,
      ].includes(tableViewId) &&
      tableViewDef?.isPreset);

  const getButtonLabels = (): [string, string] => {
    if (
      [TableIdType.UpcomingInvoices, TableIdType.PastInvoices].includes(
        tableViewId as TableIdType,
      )
    ) {
      return ['Upcoming', 'Past'];
    }

    return ['Orgs', 'Contacts'];
  };

  const [firstButtonLabel, secondButtonLabel] = getButtonLabels();

  return (
    <>
      {showToggle && firstTableDef && secondTableDef && (
        <ButtonGroup className='flex items-center'>
          <Button
            size='xs'
            onClick={() => firstTableDef && handleNavigate(firstTableDef)}
            className={cn('bg-white !border-r px-4', {
              'bg-gray-50 text-gray-500 font-normal': preset !== firstTableDef,
            })}
          >
            {firstButtonLabel}
          </Button>
          <Button
            size='xs'
            onClick={() => secondTableDef && handleNavigate(secondTableDef)}
            className={cn('bg-white px-4', {
              'bg-gray-50 text-gray-500 font-normal': preset !== secondTableDef,
            })}
          >
            {secondButtonLabel}
          </Button>
        </ButtonGroup>
      )}
    </>
  );
});
