import { useSearchParams } from 'react-router-dom';
import { useState, RefObject, startTransition } from 'react';

import { observer } from 'mobx-react-lite';
import difference from 'lodash/difference';
import { FilterItem } from '@store/types.ts';
import countries from '@assets/countries/countries.json';

import { Input } from '@ui/form/Input';
import { flags } from '@ui/media/flags';
import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { SearchSm } from '@ui/media/icons/SearchSm.tsx';
import { FilterHeader } from '@shared/components/Filters';
import { InputGroup, LeftElement } from '@ui/form/InputGroup';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

interface ContactFilterProps {
  placeholder?: string;
  property?: ColumnViewType;
  type: 'contacts' | 'organizations';
  initialFocusRef: RefObject<HTMLInputElement>;
  locationType: 'countryCodeA2' | 'locality' | 'region';
}

const defaultFilter: FilterItem = {
  property: ColumnViewType.ContactsPersona,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.Contains,
};

export const LocationFilter = observer(
  ({
    initialFocusRef,
    property,
    placeholder,
    locationType,
    type,
  }: ContactFilterProps) => {
    const [searchParams] = useSearchParams();
    const [searchValue, setSearchValue] = useState('');
    const preset = searchParams.get('preset');

    const store = useStore();
    const tableViewDef = store.tableViewDefs.getById(preset ?? '');

    const filter = tableViewDef?.getFilter(
      property || defaultFilter.property,
    ) ?? { ...defaultFilter, property: property || defaultFilter.property };
    const dataArray =
      type === 'contacts'
        ? store.contacts.toArray()
        : store.organizations.toArray();
    const allLocations = [
      ...new Set(
        dataArray
          .map((e) => e.value.locations.map((d) => d[locationType]))
          .flat()
          .filter((e) => !!e),
      ),
    ] as Array<string>;

    const filteredLocations = allLocations
      .filter((e) => {
        if (!searchValue) return true;
        if (!e) return false;
        if (locationType === 'countryCodeA2') {
          const country = countries
            .find((d) => d.alpha2 === e?.toLowerCase())
            ?.name?.toLowerCase();

          return (
            e?.toLowerCase().includes(searchValue.toLowerCase()) ||
            country?.includes(searchValue.toLowerCase())
          );
        }

        return e?.toLowerCase().includes(searchValue.toLowerCase());
      })
      .sort((a, b) => a.localeCompare(b));

    const toggle = () => {
      tableViewDef?.toggleFilter(filter);
    };

    const handleChange = (newValue: string) => {
      const filterValue = Array.isArray(filter.value) ? filter.value : [];
      const value = filterValue?.includes(newValue)
        ? filterValue.filter((e) => e !== newValue)
        : [...filterValue, newValue];
      startTransition(() => {
        tableViewDef?.setFilter({
          ...filter,
          value,
          active: filter.active || true,
        });
      });
      setSearchValue('');
    };
    const isAllChecked =
      filter.value.length === allLocations?.length && allLocations?.length > 0;
    const handleSelectAll = () => {
      let nextValue: string[] = [];

      if (isAllChecked) {
        tableViewDef?.setFilter({
          ...filter,
          value: difference(filter.value, allLocations),
          active: false,
        });

        return;
      }

      if (searchValue) {
        nextValue = [
          ...allLocations,
          ...difference(filter.value, allLocations),
        ];
      } else {
        nextValue = allLocations;
      }

      tableViewDef?.setFilter({
        ...filter,
        value: nextValue,
        active: nextValue.length > 0,
      });
    };

    const handleShowEmpty = (isChecked: boolean) => {
      tableViewDef?.setFilter({
        ...filter,
        includeEmpty: isChecked,
        active: filter.active || true,
      });
    };

    return (
      <div className='max-h-[500px]'>
        <FilterHeader
          onToggle={toggle}
          onDisplayChange={() => {}}
          isChecked={filter.active ?? false}
        />
        <InputGroup>
          <LeftElement>
            <SearchSm color='gray.500' />
          </LeftElement>
          <Input
            value={searchValue}
            size='sm'
            ref={initialFocusRef}
            onChange={(e) => {
              setSearchValue(e.target.value);
              tableViewDef?.setFilter({
                ...filter,
                active: e.target.value.length > 0,
              });
            }}
            placeholder={placeholder || 'e.g. United States'}
            className='border-none'
          />
        </InputGroup>

        <div className='flex flex-col pt-2 pb-2 gap-2 border-b border-gray-200'>
          <Checkbox
            isChecked={filter.includeEmpty ?? false}
            onChange={(value) => handleShowEmpty(value as boolean)}
          >
            <span className='text-sm'>Unknown</span>
          </Checkbox>
          <Checkbox isChecked={isAllChecked} onChange={handleSelectAll}>
            <span className='text-sm'>
              {isAllChecked ? 'Deselect all' : 'Select all'}
            </span>
          </Checkbox>
        </div>
        {!!filteredLocations.length && (
          <div className='mt-2 overflow-hidden overflow-y-auto max-h-[360px]  -mr-3 '>
            {filteredLocations?.map((e) =>
              e ? (
                <Checkbox
                  key={e}
                  className='mt-2 min-w-5 flex items-center'
                  size='md'
                  isChecked={filter.value.includes(e) ?? false}
                  labelProps={{ className: 'text-sm mt-2' }}
                  onChange={() => handleChange(e)}
                >
                  <div className='flex items-center overflow-ellipsis'>
                    {locationType === 'countryCodeA2' ? (
                      <>
                        {flags[e as keyof typeof flags]}

                        <span className='overflow-hidden overflow-ellipsis whitespace-nowrap ml-2 max-w-[110px]'>
                          {countries.find((d) => d.alpha2 === e.toLowerCase())
                            ?.name ?? e}
                        </span>
                      </>
                    ) : (
                      e ?? 'Unnamed'
                    )}
                  </div>
                </Checkbox>
              ) : null,
            )}{' '}
          </div>
        )}
      </div>
    );
  },
);
