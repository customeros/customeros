import { useState, useCallback } from 'react';

import { toJS } from 'mobx';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Cake } from '@ui/media/icons/Cake';
import { Play } from '@ui/media/icons/Play';
import { Key01 } from '@ui/media/icons/Key01';
import { Tag01 } from '@ui/media/icons/Tag01';
import { Button } from '@ui/form/Button/Button';
import { Star06 } from '@ui/media/icons/Star06';
import { Users03 } from '@ui/media/icons/Users03';
import { useStore } from '@shared/hooks/useStore';
import { Globe05 } from '@ui/media/icons/Globe05';
import { Linkedin } from '@ui/media/icons/Linkedin';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { TableViewType } from '@shared/types/tableDef';
import { Building05 } from '@ui/media/icons/Building05';
import { StopCircle } from '@ui/media/icons/StopCircle';
import { SelectOption } from '@shared/types/SelectOptions';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { ConfirmDialog } from '@ui/overlay/AlertDialog/ConfirmDialog';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';
import {
  WorkflowType,
  ColumnViewType,
  OrganizationStage,
  ComparisonOperator,
} from '@shared/types/__generated__/graphql.types';

import { industryOptions, locationsOptions } from '../../utils';
import { ICPRangeSelector, MultiSelectFilter } from '../components';
import { getOrganizationFilterFns } from '../../Columns/organizations';
import { getFlowFilterFns } from '../../Columns/organizations/flowFilters';

const options = ['between', 'less than', 'more than'];

export const Icp = observer(() => {
  const store = useStore();
  const getWorkFlow = store.workFlows
    .toArray()
    .filter((wf) => toJS(wf.value.type === WorkflowType.IdealCustomerProfile));

  const getWorkFlowId = getWorkFlow.map((wf) => wf.value.id);

  const workFlow = store.workFlows.getByType(getWorkFlowId[0]);

  const { onOpen, onClose, open } = useDisclosure();

  const [employeesFilter, setEmployeesFilter] = useState(
    workFlow?.getFilter(`${ColumnViewType.OrganizationsEmployeeCount}`)
      ?.operation === ComparisonOperator.Between
      ? 'between'
      : ComparisonOperator.Lt ===
        workFlow?.getFilter(`${ColumnViewType.OrganizationsEmployeeCount}`)
          ?.operation
      ? 'less than'
      : 'more than' ?? options[1],
  );
  const [followersFilter, setFollowersFilter] = useState(
    workFlow?.getFilter(`${ColumnViewType.OrganizationsLinkedinFollowerCount}`)
      ?.operation === ComparisonOperator.Between
      ? 'between'
      : ComparisonOperator.Lt ===
        workFlow?.getFilter(
          `${ColumnViewType.OrganizationsLinkedinFollowerCount}`,
        )?.operation
      ? 'less than'
      : 'more than' ?? options[1],
  );
  const [yearsFilter, setYearsFilter] = useState(
    workFlow?.getFilter(`${ColumnViewType.OrganizationsYearFounded}`)
      ?.operation === ComparisonOperator.Between
      ? 'between'
      : ComparisonOperator.Lt ===
        workFlow?.getFilter(`${ColumnViewType.OrganizationsYearFounded}`)
          ?.operation
      ? 'less than'
      : 'more than' ?? options[1],
  );

  const leadsPreset = store.tableViewDefs.leadsPreset;

  const handleEmployeesFilter = () => {
    const currentIndex = options.indexOf(employeesFilter);
    const nextIndex = (currentIndex + 1) % options.length;

    setEmployeesFilter(options[nextIndex]);

    if (workFlow?.value.live) {
      workFlow?.update((workflow) => {
        workflow.live = false;

        return workflow;
      });
    }
    workFlow?.removeFilter(ColumnViewType.OrganizationsEmployeeCount);
  };

  const handleTagsFilter = () => {
    const currentIndex = options.indexOf(followersFilter);
    const nextIndex = (currentIndex + 1) % options.length;

    setFollowersFilter(options[nextIndex]);

    if (workFlow?.value.live) {
      workFlow?.update((workflow) => {
        workflow.live = false;

        return workflow;
      });
    }
    workFlow?.removeFilter(ColumnViewType.OrganizationsLinkedinFollowerCount);
  };

  const handleYearsFilter = () => {
    const currentIndex = options.indexOf(yearsFilter);
    const nextIndex = (currentIndex + 1) % options.length;

    setYearsFilter(options[nextIndex]);

    if (workFlow?.value.live) {
      workFlow?.update((workflow) => {
        workflow.live = false;

        return workflow;
      });
    }
    workFlow?.removeFilter(ColumnViewType.OrganizationsYearFounded);
  };

  const tagsOptions = store.tags
    .toArray()
    .map((tag) => ({ value: tag.value.id, label: tag.value.name }));

  const handleChange = useCallback(
    (selectedOptions: SelectOption[], property: string) => {
      workFlow?.update((workflow) => {
        workflow.live = false;

        return workflow;
      });

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
    },
    [workFlow?.value.condition],
  );

  const handleFilterSelected = (property: string) => {
    const filter = workFlow?.getFilter(property);

    return filter ? filter.value : [];
  };

  const tableViewDef = store.tableViewDefs.getById(leadsPreset ?? '1');

  const tableType = tableViewDef?.value.tableType;

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

  const organizationsChangeStage = () => {
    const selectedIds = organizationsData.map((org) => org.value.metadata.id);

    store.ui.setMovedIcpOrganization(selectedIds.length);
    store.organizations.updateStage(
      selectedIds,
      OrganizationStage.Target,
      false,
    );
  };

  const handleFiltersEmpty = useCallback(() => {
    const filters = workFlow?.getFilters();

    if (filters) {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const filterValues = Object.values(filters).flatMap((v: any) =>
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        v.map((f: any) => f.filter.value.length),
      );

      if (filterValues.length === 0) {
        return true;
      }
    }

    return false;
  }, [workFlow?.value.condition]);

  return (
    <>
      <div className='flex items-center justify-between'>
        <p className='font-semibold'>Auto-qualify leads</p>

        {workFlow?.value.live === false ? (
          <Tooltip
            hasArrow
            side='bottom'
            align='center'
            label='First add at least 1 filter'
            open={handleFiltersEmpty() ? undefined : false}
            className={cn(handleFiltersEmpty() ? 'flex' : 'hidden')}
          >
            <span>
              <Button
                size='xxs'
                leftIcon={<Play />}
                isDisabled={handleFiltersEmpty()}
                onClick={() => {
                  onOpen();
                }}
              >
                Start automation
              </Button>
            </span>
          </Tooltip>
        ) : (
          <Button
            size='xxs'
            colorScheme='warning'
            leftIcon={<StopCircle />}
            onClick={() => {
              workFlow?.update((workflow) => {
                workflow.live = false;

                return workflow;
              });
            }}
          >
            Stop automation
          </Button>
        )}
      </div>
      <p className='mt-1'>
        Create your <span className='font-medium'>Ideal Company Profile </span>{' '}
        and automatically qualify
        <span className='font-medium'> Leads </span>
        from the Organizations view into
        <span className='font-medium'> Targets</span>
      </p>
      <p className='font-medium leading-5 text-gray-500 mt-4 mb-2'>WHEN</p>

      <MultiSelectFilter
        label='Industry'
        description='is any of'
        placeholder='Industries'
        options={industryOptions}
        icon={<Building05 className='mr-2 text-gray-500' />}
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
              onClick={handleEmployeesFilter}
              className='cursor-pointer underline font-normal text-gray-500 hover:text-gray-700'
            >
              {employeesFilter}
            </span>
          </p>
        </div>

        <ICPRangeSelector
          filter={employeesFilter}
          placeholder='Number of employees'
          property={ColumnViewType.OrganizationsEmployeeCount}
        />
      </div>
      <MultiSelectFilter
        label='Headquarters'
        description='is any of'
        options={locationsOptions}
        placeholder='Headquarter countries'
        icon={<Globe05 className='mr-2 text-gray-500' />}
        onChange={(value) =>
          handleChange(value, ColumnViewType.OrganizationsHeadquarters)
        }
        value={handleFilterSelected(
          ColumnViewType.OrganizationsHeadquarters,
        ).map((value: string) => ({
          value: value,
          label: value,
        }))}
      />

      <MultiSelectFilter
        label='Tag'
        options={tagsOptions}
        description='is any of'
        placeholder='Organization tags'
        icon={<Tag01 className='mr-2 text-gray-500' />}
        onChange={(value) =>
          handleChange(value, ColumnViewType.OrganizationsTags)
        }
        value={handleFilterSelected(ColumnViewType.OrganizationsTags).map(
          (value: string) => ({
            value: value,
            label: tagsOptions
              .filter((option) => option.value === value)
              .map((option) => option.label),
          }),
        )}
      />

      <div className='flex items-center w-full '>
        <div className='flex flex-1 items-center'>
          <Linkedin className='mr-2 text-gray-500 ' />
          <p className='font-medium'>
            Followers <span className='font-normal'>is </span>
            <span
              onClick={handleTagsFilter}
              className='cursor-pointer underline font-normal text-gray-500 hover:text-gray-700'
            >
              {followersFilter}
            </span>
          </p>
        </div>

        <ICPRangeSelector
          filter={followersFilter}
          placeholder='Number of followers'
          property={ColumnViewType.OrganizationsLinkedinFollowerCount}
        />
      </div>
      <div className='flex items-center w-full '>
        <div className='flex flex-1 items-center'>
          <Cake className='mr-2 text-gray-500' />
          <p className='font-medium'>
            Organization age <span className='font-normal'>is </span>
            <span
              onClick={handleYearsFilter}
              className='cursor-pointer underline font-normal text-gray-500 hover:text-gray-700'
            >
              {yearsFilter}
            </span>
          </p>
        </div>

        <ICPRangeSelector
          years
          placeholder='Age'
          filter={yearsFilter}
          property={ColumnViewType.OrganizationsYearFounded}
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
      {!handleFiltersEmpty() && (
        <div className='mt-4 border rounded-md flex items-start gap-2 p-3 bg-grayModern-50'>
          <div className='flex flex-col w-fit'>
            <Star06 className='mt-1 text-grayModern-500' />
          </div>
          <div className='flex flex-col'>
            <p className='font-medium'>
              This automation will qualify{' '}
              <span className='font-medium'>
                {' '}
                {`${filteredResults}/${toatalResults} Leads`}
              </span>{' '}
              into <span className='font-medium'>Targets</span>
            </p>
            {/* <Checkbox
              isChecked={store.ui.isFilteringICP as boolean}
              onChange={(v) => store.ui.setIsFilteringICP(v as boolean)}
            >
              Show leads that will not be qualified
            </Checkbox> */}
          </div>
        </div>
      )}
      <ConfirmDialog
        isOpen={open}
        onClose={onClose}
        title='Start auto-qualifying leads'
        confirmButtonLabel='Start automation'
        body={'You can manually change this stage at any time again.'}
        description={`Starting this automation will immediately qualify ${filteredResults} Leads into Targets and continue to qualify matching leads in the background until stopped. 
         `}
        onConfirm={() => {
          organizationsChangeStage();

          workFlow?.update((workflow) => {
            workflow.live = true;

            return workflow;
          });
        }}
      />
    </>
  );
});
