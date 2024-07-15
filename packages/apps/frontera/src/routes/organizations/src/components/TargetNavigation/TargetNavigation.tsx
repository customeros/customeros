import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

import { cn } from '@ui/utils/cn.ts';
import { useStore } from '@shared/hooks/useStore';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { Button } from '@ui/form/Button/Button.tsx';
import { TableIdType, TableViewType } from '@graphql/types';

export const TargetNavigation = observer(() => {
  const store = useStore();
  const [searchParams, setSearchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const [tabs, setLastActivePosition] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-player-last-position`, { root: 'organizations' });
  const tableViewName = store.tableViewDefs.getById(preset || '')?.value.name;
  const tableViewType = store.tableViewDefs.getById(preset || '')?.value
    .tableType;

  const contactTableDef = store.tableViewDefs
    .toArray()
    .find((e) => e.value.tableType === TableViewType.Contacts)?.value.id;
  const targetTableDef = store.tableViewDefs
    .toArray()
    .find((e) => e.value.tableId === TableIdType.Nurture)?.value?.id;
  const handleNavigate = (newPreset: string) => {
    setSearchParams((prev) => {
      prev.set('preset', newPreset as string);

      return prev;
    });
    setLastActivePosition({
      ...tabs,
      root: `finder?preset=${newPreset}`,
    });
  };

  return (
    <>
      {(tableViewType === TableViewType.Contacts ||
        tableViewName === 'Targets') && (
        <ButtonGroup className='flex items-center '>
          <Button
            size='xs'
            className={cn('bg-white !border-r px-4', {
              'bg-gray-50 text-gray-500 font-normal': preset !== targetTableDef,
            })}
            onClick={() => {
              handleNavigate(targetTableDef as string);
            }}
          >
            Orgs
          </Button>
          <Button
            size='xs'
            className={cn('bg-white px-4', {
              'bg-gray-50 text-gray-500 font-normal':
                preset !== contactTableDef,
            })}
            onClick={() => {
              handleNavigate(contactTableDef as string);
            }}
          >
            Contacts
          </Button>
        </ButtonGroup>
      )}
    </>
  );
});
