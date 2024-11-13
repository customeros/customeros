import type { HeaderContext } from '@tanstack/react-table';

import { memo, RefObject } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { ArrowUp } from '@ui/media/icons/ArrowUp';
import { useStore } from '@shared/hooks/useStore';
import { ArrowDown } from '@ui/media/icons/ArrowDown';

interface THeadProps<
  InitialRefType extends { focus(): void } = HTMLButtonElement,
> {
  id: string;
  py?: string;
  title: string;
  padding?: string;
  subTitle?: string;
  canSort?: boolean;
  canFilter?: boolean;
  isFiltered?: boolean;
  borderTopColor?: string;
  isSorted?: string | boolean;
  filterWidth?: string | number;
  onToggleSort?: (e: unknown) => void;
  renderFilter?: (
    initialFocusRef: RefObject<InitialRefType>,
  ) => React.ReactNode;
}

const THead = observer(
  <InitialRefType extends { focus(): void } = HTMLButtonElement>({
    id,
    title,
    canSort,
    isSorted,
    subTitle,
    onToggleSort,
    py,
  }: THeadProps<InitialRefType>) => {
    const store = useStore();
    const [searchParams] = useSearchParams();
    const preset = searchParams.get('preset');

    const isActive = store.tableViewDefs
      .getById(preset ?? '')
      ?.getFilter(id)?.active;

    return (
      <div className='flex w-full ml-[-22px] flex-col justify-start items-start group'>
        <div
          style={{ paddingTop: py ?? '0', paddingBottom: py ?? '0' }}
          className={cn(
            isActive ? 'border-gray-300 shadow-sm' : 'border-transparent',
            (canSort && isActive) || isSorted ? 'ml-0' : 'ml-3',
            !canSort ? '' : 'group-hover:ml-0',
            'flex items-center border rounded-[4px] transition-opacity duration-200 ease-in-out pr-2',
          )}
        >
          {canSort ? (
            isSorted === 'asc' ? (
              <ArrowUp
                role='button'
                id='sort-icon'
                onClick={onToggleSort}
                className={cn(
                  isSorted || isActive ? 'w-3 ' : 'w-0 ',
                  !isSorted ? 'text-gray-400' : 'text-gray-700',
                  'mx-1 w-3 h-3 cursor-pointer group-hover:transition-opacity group-hover:opacity-100 group-hover:w-3 group-hover:duration-200 group-hover:ease-in-out',
                )}
              />
            ) : (
              <ArrowDown
                role='button'
                id='sort-icon'
                onClick={onToggleSort}
                className={cn(
                  isSorted || isActive ? 'w-3 opacity-100' : 'w-0 opacity-0',
                  !isSorted ? 'text-gray-400' : 'text-gray-700',
                  'mx-1 h-3 cursor-pointer group-hover:transition-opacity group-hover:opacity-100 group-hover:w-3 group-hover:duration-200 group-hover:ease-in-out',
                )}
              />
            )
          ) : (
            <div className={cn(canSort ? 'w-3' : 'w-0', 'flex mx-1')} />
          )}
          <p
            onClick={onToggleSort}
            data-test={`org-header-${id}`}
            className={cn(
              isSorted ? ' tracking-[-0.3px] ' : 'mt-0',
              canSort ? 'cursor-pointer' : 'cursor-default',
              !isSorted ? 'font-base' : 'font-medium',
              'text-sm text-gray-700',
            )}
          >
            {title}
          </p>
        </div>
        {subTitle && <p className='text-xs text-gray-500'>{subTitle}</p>}
      </div>
    );
  },
);

export function getTHeadProps<T extends object>(
  context: HeaderContext<T, unknown>,
) {
  const header = context.header;

  const canSort = header.column.getCanSort();
  const isSorted = header.column.getIsSorted();
  const canFilter = header.column.getCanFilter();
  const isFiltered = header.column.getIsFiltered();
  const onToggleSort = header.column.getToggleSortingHandler();

  return {
    canSort,
    isSorted,
    canFilter,
    isFiltered,
    onToggleSort,
  };
}

export default memo(THead) as typeof THead;
