import { useSearchParams } from 'react-router-dom';
import { useState, RefObject, startTransition } from 'react';

import { observer } from 'mobx-react-lite';
import { FilterItem } from '@store/types.ts';
import flags from '@assets/countries/flags.json';
import countries from '@assets/countries/countries.json';

import { Input } from '@ui/form/Input';
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
  locationType: 'countryCodeA2' | 'locality';
  initialFocusRef: RefObject<HTMLInputElement>;
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
    ].filter((e) => {
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
    });

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

    return (
      <>
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
            ref={initialFocusRef}
            onChange={(e) => setSearchValue(e.target.value)}
            placeholder={placeholder || 'e.g. United States'}
            className='border-none'
          />
        </InputGroup>
        {!!allLocations.length && (
          <div className='mt-2 overflow-y-auto  -mr-3 h-[13rem] max-w-[12rem]'>
            {allLocations?.map((e) =>
              e ? (
                <Checkbox
                  key={e}
                  className='mt-2 flex items-center'
                  size='md'
                  isChecked={filter.value.includes(e) ?? false}
                  labelProps={{ className: 'text-sm mt-2' }}
                  onChange={() => handleChange(e)}
                >
                  <div className='flex items-center'>
                    {locationType === 'countryCodeA2' ? (
                      <>
                        <img
                          src={flags[e.toLowerCase() as keyof typeof flags]}
                          alt={e}
                          className='rounded-full mr-2'
                          style={{ clipPath: 'circle(35%)' }}
                        />
                        <span className='overflow-hidden overflow-ellipsis whitespace-nowrap'>
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
            )}
          </div>
        )}
      </>
    );
  },
);
