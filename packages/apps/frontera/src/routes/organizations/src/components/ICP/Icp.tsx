import { useState } from 'react';
import { useSearchParams } from 'react-router-dom';

import { toJS } from 'mobx';
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
import { TableViewType } from '@shared/types/tableDef';
import { Building05 } from '@ui/media/icons/Building05';
import { StopCircle } from '@ui/media/icons/StopCircle';
import { SelectOption } from '@shared/types/SelectOptions';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';
import {
  WorkflowType,
  ColumnViewType,
  OrganizationStage,
  ComparisonOperator,
} from '@shared/types/__generated__/graphql.types';

import { RangeSelector, MultiSelectFilter } from '../shared';
import { industryOptions, locationsOptions } from '../utils';
import { getOrganizationFilterFns } from '../Columns/organizations';
import { getFlowFilterFns } from '../Columns/organizations/flowFilters';

const options = ['between', 'less than', 'more than'];
export const Icp = observer(() => {
  const store = useStore();
  const [searchParams] = useSearchParams();

  const [employeesFilter, setEmployeesFilter] = useState(options[1]);
  const [followersFilter, setFollowersFilter] = useState(options[1]);
  const [organizationFilter, setOrganizationFilter] = useState(options[1]);
  const getWorkFlow = store.workFlows
    .toArray()
    .filter((wf) => toJS(wf.value.type === WorkflowType.IdealCustomerProfile));

  const getWorkFlowId = getWorkFlow.map((wf) => wf.value.id);

  const workFlow = store.workFlows.getByType(getWorkFlowId[0]);
  const preset = searchParams?.get('preset');
  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const tableType = tableViewDef?.value?.tableType;

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

  const tagsOptions = store.tags
    .toArray()
    .map((tag) => ({ value: tag.value.id, label: tag.value.name }));

  const handleChange = (selectedOptions: SelectOption[], property: string) => {
    if (selectedOptions.length === 0) {
      workFlow?.removeFilter(property);

      return;
    }

    const newValues = selectedOptions.map(
      (option: SelectOption) => option.value,
    );

    workFlow?.setFilter({
      property: property,
      value: newValues,
      operation: ComparisonOperator.In,
    });
  };

  const handleFilterSelected = (property: string) => {
    const filter = workFlow?.getFilter(property);

    return filter ? filter.value : [];
  };

  const toatalResults = store.organizations
    ?.toArray()
    .filter((v) => v.value.stage === OrganizationStage.Lead).length;

  const organizationsData = store.organizations?.toComputedArray((arr) => {
    if (tableType !== TableViewType.Organizations) return arr;
    const filters = getOrganizationFilterFns(tableViewDef?.getFilters());

    const flowFilters = getFlowFilterFns(workFlow?.getFilters());

    if (flowFilters.length) {
      arr = arr.filter((v) => flowFilters.every((fn) => fn(v)));
    }
    if (filters) {
      arr = arr.filter((v) => filters.every((fn) => fn(v)));
    }

    return arr;
  });

  const filteredResults = organizationsData.length;

  return (
    <>
      <div className='flex items-center justify-between'>
        <p className='font-semibold'>Auto-qualify leads</p>

        {workFlow?.value.live === false ? (
          <Button
            size='xxs'
            leftIcon={<Play />}
            onClick={() => {
              workFlow?.update((workflow) => {
                workflow.live = true;

                return workflow;
              });
            }}
          >
            Start flow
          </Button>
        ) : (
          <Button
            size='xxs'
            leftIcon={<StopCircle />}
            onClick={() => {
              workFlow?.update((workflow) => {
                workflow.live = false;

                return workflow;
              });
            }}
          >
            Stop flow
          </Button>
        )}
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
        options={industryOptions}
        onChange={(value) =>
          handleChange(value, ColumnViewType.OrganizationsIndustry)
        }
        value={handleFilterSelected(ColumnViewType.OrganizationsIndustry).map(
          (value: string) => ({
            value: value,
            label: industryOptions.find(
              (option: SelectOption) => option.value === value,
            )?.label,
          }),
        )}
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
              workFlow?.setFilter({
                property: ColumnViewType.OrganizationsEmployeeCount,
                value: values,
                operation:
                  employeesFilter === 'between'
                    ? ComparisonOperator.Between
                    : employeesFilter === 'less than'
                    ? ComparisonOperator.Lte
                    : ComparisonOperator.Gte,
              });
            }
            if (values[0] === '') {
              workFlow?.removeFilter(ColumnViewType.OrganizationsEmployeeCount);
            }
          }}
        />
      </div>
      <MultiSelectFilter
        icon={<Globe03 className='mr-2 text-gray-500' />}
        label='Headquarters'
        description='is any of'
        placeholder='Headquarter countries'
        value={handleFilterSelected(
          ColumnViewType.OrganizationsHeadquarters,
        ).map((value: string) => ({
          value: value,
          label: value,
        }))}
        onChange={(value) =>
          handleChange(value, ColumnViewType.OrganizationsHeadquarters)
        }
        options={locationsOptions}
      />

      <MultiSelectFilter
        icon={<Tag01 className='mr-2 text-gray-500' />}
        label='Tag'
        description='is any of'
        placeholder='Organization tags'
        value={handleFilterSelected(ColumnViewType.OrganizationsTags).map(
          (value: string) => ({
            value: value,
            label: tagsOptions
              .filter((option) => option.value === value)
              .map((option) => option.label),
          }),
        )}
        onChange={(value) =>
          handleChange(value, ColumnViewType.OrganizationsTags)
        }
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
              workFlow?.setFilter({
                property: ColumnViewType.OrganizationsLinkedinFollowerCount,
                value: values,
                operation:
                  followersFilter === 'between'
                    ? ComparisonOperator.Between
                    : followersFilter === 'less than'
                    ? ComparisonOperator.Lte
                    : ComparisonOperator.Gte,
              });
            }
            if (values[0] === '') {
              workFlow?.removeFilter(
                ColumnViewType.OrganizationsLinkedinFollowerCount,
              );
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
              workFlow?.setFilter({
                property: ColumnViewType.OrganizationsYearFounded,
                value: values[1]
                  ? [
                      new Date().getFullYear() - (values[0] as number),
                      new Date().getFullYear() - (values[1] as number),
                    ]
                  : [new Date().getFullYear() - (values[0] as number)],
                operation:
                  organizationFilter === 'between'
                    ? ComparisonOperator.Between
                    : organizationFilter === 'less than'
                    ? ComparisonOperator.Lte
                    : ComparisonOperator.Gte,
              });
            }
            if (values[0] === '') {
              workFlow?.removeFilter(ColumnViewType.OrganizationsYearFounded);
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
          <Menu>
            <MenuButton>
              {workFlow?.getFilter(ColumnViewType.OrganizationsIsPublic)
                ?.value === true
                ? 'Public'
                : workFlow?.getFilter(ColumnViewType.OrganizationsIsPublic)
                    ?.value === undefined
                ? 'Not applicable'
                : 'Private'}
            </MenuButton>
            <MenuList>
              <MenuItem
                onClick={() => {
                  workFlow?.setFilter({
                    property: ColumnViewType.OrganizationsIsPublic,
                    value: false,
                    operation: ComparisonOperator.Eq,
                  });
                }}
              >
                Private
              </MenuItem>
              <MenuItem
                onClick={() => {
                  workFlow?.setFilter({
                    property: ColumnViewType.OrganizationsIsPublic,
                    value: true,
                    operation: ComparisonOperator.Eq,
                  });
                }}
              >
                Public
              </MenuItem>
              <MenuItem
                onClick={() => {
                  workFlow?.removeFilter(ColumnViewType.OrganizationsIsPublic);
                }}
              >
                Not applicable
              </MenuItem>
            </MenuList>
          </Menu>
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
            onChange={(v) => store.ui.setIsFilteringICP(v as boolean)}
            isChecked={store.ui.isFilteringICP as boolean}
          >
            See unqualified leads before starting the flow
          </Checkbox>
        </div>
      </div>
    </>
  );
});