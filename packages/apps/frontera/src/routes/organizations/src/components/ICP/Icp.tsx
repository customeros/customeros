import { useSearchParams } from 'react-router-dom';
import { useMemo, useState, useCallback } from 'react';

import { observer } from 'mobx-react-lite';

import { Cake } from '@ui/media/icons/Cake';
import { Play } from '@ui/media/icons/Play';
import { Key01 } from '@ui/media/icons/Key01';
import { Tag01 } from '@ui/media/icons/Tag01';
import { Button } from '@ui/form/Button/Button';
import { Star06 } from '@ui/media/icons/Star06';
import { Globe03 } from '@ui/media/icons/Globe03';
import { Users03 } from '@ui/media/icons/Users03';
import { useStore } from '@shared/hooks/useStore';
import { Linkedin } from '@ui/media/icons/Linkedin';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { Building05 } from '@ui/media/icons/Building05';
import { getContainerClassNames } from '@ui/form/Select';
import { SelectOption } from '@shared/types/SelectOptions';
import { ComparisonOperator } from '@shared/types/__generated__/graphql.types';

import { RangeSelector, MultiSelectFilter } from '../shared';
import { industryOptions, locationsOptions } from '../utils';
import { getAllFilterFns } from '../Columns/Dictionaries/columnsDictionary';
import { getOrganizationFilterFn } from '../Columns/Dictionaries/SortAndFilterDictionary';
import { getFlowFilters } from '../Columns/Dictionaries/SortAndFilterDictionary/flowFilters';

const options = ['between', 'less than', 'more than'];
const ownershipOptions = ['Private', 'Public'];

export const Icp = observer(() => {
  const store = useStore();
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const [employeesFilter, setEmployeesFilter] = useState(options[1]);
  const [followersFilter, setFollowersFilter] = useState(options[1]);
  const [organizationFilter, setOrganizationFilter] = useState(options[1]);
  const [ownership, setOwnership] = useState(ownershipOptions[0]);

  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const tableType = tableViewDef?.value?.tableType;

  const dataSet = useMemo(() => {
    return store.organizations;
  }, [tableType]);

  const filterFunction = useMemo(() => {
    return getOrganizationFilterFn;
  }, [tableType]);

  const handleEmployeesFilter = () => {
    const currentIndex = options.indexOf(employeesFilter);
    const nextIndex = (currentIndex + 1) % options.length;
    setEmployeesFilter(options[nextIndex]);
  };

  const handleTagsFilter = () => {
    const currentIndex = options.indexOf(followersFilter);
    const nextIndex = (currentIndex + 1) % options.length;
    setFollowersFilter(options[nextIndex]);
  };

  const handleOrganizationFilter = () => {
    const currentIndex = options.indexOf(organizationFilter);
    const nextIndex = (currentIndex + 1) % options.length;
    setOrganizationFilter(options[nextIndex]);
  };

  const handleOwnershipFilter = useCallback(() => {
    const currentIndex = ownershipOptions.indexOf(ownership);
    const nextIndex = (currentIndex + 1) % ownershipOptions.length;
    setOwnership(ownershipOptions[nextIndex]);
    store.workFlows.setFilter({
      property: 'ownership',
      value: ownership,
    });
  }, [ownership]);

  const tagsOptions = store.tags
    .toArray()
    .map((tag) => ({ value: tag.value.id, label: tag.value.name }));

  const handleChange = (selectedOptions: SelectOption[], property: string) => {
    if (!selectedOptions) {
      store.workFlows.setFilter({
        property: property,
        value: [],
      });

      return;
    }

    const newValues = selectedOptions.map(
      (option: SelectOption) => option.value,
    );

    store.workFlows.setFilter({
      property: property,
      value: newValues,
      operation: ComparisonOperator.Contains,
    });
  };

  const handleFilterSelected = (property: string) => {
    const filter = store.workFlows.getFilter(property);

    return filter ? filter.value : [];
  };

  const data = dataSet?.toComputedArray((arr) => {
    const filters = getAllFilterFns(tableViewDef?.getFilters(), filterFunction);
    if (filters) {
      // @ts-expect-error fixme
      arr = arr.filter((v) => filters.every((fn) => fn(v)));
    }

    return arr;
  });

  const filteredData = dataSet?.toComputedArray((arr) => {
    const filters = getAllFilterFns(tableViewDef?.getFilters(), filterFunction);

    const flowFilters = getAllFilterFns(
      store.workFlows.getFilters(),
      getFlowFilters,
    );
    if (flowFilters.length && true) {
      // @ts-expect-error fixme
      arr = arr.filter((v) => flowFilters.every((fn) => fn(v)));
    }
    if (filters) {
      // @ts-expect-error fixme
      arr = arr.filter((v) => filters.every((fn) => fn(v)));
    }

    return arr;
  });

  const toatalResults = data?.length;
  const filteredResults = filteredData?.length;

  return (
    <>
      <div className='flex items-center justify-between'>
        <p className='font-semibold'>Auto-qualify leads</p>
        <Button size='xxs' leftIcon={<Play />}>
          Start flow
        </Button>
      </div>
      <p className='mt-1'>
        Create your <span className='font-medium'>Ideal Company Profile </span>{' '}
        and automatically qualify
        <span className='font-medium'> Leads </span>
        into <span className='font-medium'>Targets</span>
      </p>
      <p className='font-medium leading-5 text-gray-500 mt-4 mb-2'>WHEN</p>

      <MultiSelectFilter
        icon={<Building05 className='mr-2 text-gray-500' />}
        label='Industry'
        description='is any of'
        placeholder='Industries'
        classNames={{
          container: () => getContainerClassNames(undefined, 'unstyled', {}),
        }}
        options={industryOptions}
        onChange={(value) => handleChange(value, 'industry')}
        value={handleFilterSelected('industry').map((value: string) => ({
          value: value,
          label: industryOptions.find(
            (option: SelectOption) => option.value === value,
          )?.label,
        }))}
      />

      <div className='flex items-center w-full'>
        <div className='flex-1 items-center flex'>
          <Users03 className='mr-2 text-gray-500' />
          <p className='font-medium'>
            Employees <span className='font-normal'>are </span>
            <span
              className='cursor-pointer underline font-normal text-gray-500 hover:text-gray-700'
              onClick={handleEmployeesFilter}
            >
              {employeesFilter}
            </span>
          </p>
        </div>

        <RangeSelector
          filter={employeesFilter}
          placeholder='Number of employees'
          onChange={(values) => {
            if (values[0] !== undefined) {
              store.workFlows.setFilter({
                property: 'employees',
                value: values,
                operation:
                  employeesFilter === 'between'
                    ? ComparisonOperator.Between
                    : employeesFilter === 'less than'
                    ? ComparisonOperator.Lte
                    : ComparisonOperator.Gte,
              });
            }
          }}
        />
      </div>
      <MultiSelectFilter
        icon={<Globe03 className='mr-2 text-gray-500' />}
        label='Headquarters'
        description='is any of'
        placeholder='Headquarter countries'
        value={handleFilterSelected('headquarters').map((value: string) => ({
          value: value,
          label: value,
        }))}
        onChange={(value) => handleChange(value, 'headquarters')}
        options={locationsOptions}
      />

      <MultiSelectFilter
        icon={<Tag01 className='mr-2 text-gray-500' />}
        label='Tag'
        description='is any of'
        placeholder='Organization tags'
        value={handleFilterSelected('tags').map((value: string) => ({
          value: value,
          label: value,
        }))}
        onChange={(value) => handleChange(value, 'tags')}
        options={tagsOptions}
      />

      <div className='flex items-center w-full '>
        <div className='flex flex-1 items-center'>
          <Linkedin className='mr-2 text-gray-500 ' />
          <p className='font-medium'>
            Followers <span className='font-normal'>is </span>
            <span
              className='cursor-pointer underline font-normal text-gray-500 hover:text-gray-700'
              onClick={handleTagsFilter}
            >
              {followersFilter}
            </span>
          </p>
        </div>

        <RangeSelector
          filter={followersFilter}
          placeholder='Number of followers'
          onChange={(values) => {
            if (values[0] !== undefined) {
              store.workFlows.setFilter({
                property: 'followers',
                value: values,
                operation:
                  employeesFilter === 'between'
                    ? ComparisonOperator.Between
                    : employeesFilter === 'less than'
                    ? ComparisonOperator.Lte
                    : ComparisonOperator.Gte,
              });
            }
          }}
        />
      </div>
      <div className='flex items-center w-full '>
        <div className='flex flex-1 items-center'>
          <Cake className='mr-2 text-gray-500' />
          <p className='font-medium'>
            Organization age <span className='font-normal'>is </span>
            <span
              className='cursor-pointer underline font-normal text-gray-500 hover:text-gray-700'
              onClick={handleOrganizationFilter}
            >
              {organizationFilter}
            </span>
          </p>
        </div>

        <RangeSelector
          filter={organizationFilter}
          placeholder='Age'
          onChange={(values) => {
            if (values[0] !== undefined) {
              store.workFlows.setFilter({
                property: 'age',
                value: values,
                operation:
                  employeesFilter === 'between'
                    ? ComparisonOperator.Between
                    : employeesFilter === 'less than'
                    ? ComparisonOperator.Lte
                    : ComparisonOperator.Gte,
              });
            }
          }}
          years
        />
      </div>

      <div className='flex items-center w-full mt-2 '>
        <div className='flex flex-1 items-center'>
          <Key01 className='mr-2 text-gray-500' />
          <p className='font-medium'>
            Ownership <span className='font-normal'>is </span>
          </p>
        </div>
        <div className='flex-1 flex items-center'>
          <span
            onClick={() => {
              handleOwnershipFilter();
            }}
            className='cursor-pointer underline'
          >
            {ownership}
          </span>
        </div>
      </div>

      <div className='mt-4 border rounded-md flex items-start gap-2 p-3 bg-grayModern-50'>
        <div className='flex flex-col w-fit'>
          <Star06 className='mt-1 text-grayModern-500' />
        </div>
        <div className='flex flex-col'>
          <p>
            This flow will qualify{' '}
            <span className='font-medium'>
              {' '}
              {`${filteredResults}/${toatalResults} Leads`}
            </span>{' '}
            into <span className='font-medium'>Targets</span>
          </p>
          <Checkbox
            onChange={(v) => store.workFlows.setFiltersStatus(v as boolean)}
            isChecked={store.workFlows.value.filterStatus as boolean}
          >
            See filtered leads before starting the flow
          </Checkbox>
        </div>
      </div>
    </>
  );
});
