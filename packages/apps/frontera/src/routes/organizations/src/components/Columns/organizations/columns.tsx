import React from 'react';

import { Store } from '@store/store.ts';
import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { cn } from '@ui/utils/cn.ts';
import { DateTimeUtils } from '@utils/date.ts';
import { createColumnHelper } from '@ui/presentation/Table';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton.tsx';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber.ts';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead.tsx';
import { CountryCell } from '@organizations/components/Columns/Cells/country';
import { OrganizationStageCell } from '@organizations/components/Columns/Cells/stage';
import { SocialsFilter } from '@organizations/components/Columns/shared/Filters/Socials';
import { StageFilter } from '@organizations/components/Columns/organizations/Filters/Stage';
import {
  Social,
  Organization,
  TableViewDef,
  ColumnViewType,
} from '@graphql/types';
import { AvatarHeader } from '@organizations/components/Columns/organizations/Headers/Avatar';
import { getColumnConfig } from '@organizations/components/Columns/shared/util/getColumnConfig.ts';
import { NumericValueFilter } from '@organizations/components/Columns/shared/Filters/NumericValueFilter';

import { OwnershipTypeFilter } from '../shared/Filters/OwnershipTypeFilter';
import { LocationFilter } from '../shared/Filters/LocationFilter/LocationFilter';
import {
  OwnerCell,
  AvatarCell,
  WebsiteCell,
  IndustryCell,
  OnboardingCell,
  OrganizationCell,
  TimeToRenewalCell,
  LastTouchpointCell,
  RenewalForecastCell,
  OrganizationsTagsCell,
  RenewalLikelihoodCell,
  LastTouchpointDateCell,
  OrganizationLinkedInCell,
  OrganizationRelationshipCell,
} from './Cells';
import {
  LtvFilter,
  OwnerFilter,
  SourceFilter,
  ChurnedFilter,
  WebsiteFilter,
  ForecastFilter,
  IndustryFilter,
  OnboardingFilter,
  CreatedDateFilter,
  OrganizationFilter,
  RelationshipFilter,
  TimeToRenewalFilter,
  LastInteractedFilter,
  LastTouchpointFilter,
  RenewalLikelihoodFilter,
} from './Filters';

type ColumnDatum = Store<Organization>;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

export const columns: Record<string, Column> = {
  [ColumnViewType.OrganizationsAvatar]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.OrganizationsAvatar,
    size: 26,
    enableColumnFilter: false,
    cell: (props) => {
      const icon = props.getValue()?.value?.icon;
      const logo = props.getValue()?.value?.logo;
      const description = props.getValue()?.value?.valueProposition;

      return (
        <AvatarCell
          icon={icon}
          logo={logo}
          description={description}
          id={props.getValue()?.value?.metadata?.id}
          name={props.getValue()?.value?.name}
        />
      );
    },
    header: AvatarHeader,
    skeleton: () => <Skeleton className='size-[24px]' />,
  }),
  [ColumnViewType.OrganizationsName]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.OrganizationsName,

    size: 150,
    cell: (props) => {
      return (
        <OrganizationCell
          id={props.getValue().value.metadata?.id}
          name={props.getValue().value.name}
          isSubsidiary={!!props.getValue()?.value?.subsidiaryOf?.length}
          parentOrganizationName={
            props.getValue()?.value?.subsidiaryOf?.[0]?.organization.name
          }
        />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.OrganizationsName}
        title='Organization'
        filterWidth='14rem'
        renderFilter={(initialFocusRef) => (
          <OrganizationFilter initialFocusRef={initialFocusRef} />
        )}
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsWebsite]: columnHelper.accessor(
    'value.website',
    {
      id: ColumnViewType.OrganizationsWebsite,
      size: 125,
      enableSorting: false,
      cell: (props) => {
        const organizationId = props.row.original.value.metadata.id;

        return <WebsiteCell organizationId={organizationId} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsWebsite}
          title='Website'
          filterWidth='14rem'
          renderFilter={(initialFocusRef) => (
            <WebsiteFilter initialFocusRef={initialFocusRef} />
          )}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[50%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsRelationship]: columnHelper.accessor(
    'value.relationship',
    {
      id: ColumnViewType.OrganizationsRelationship,
      size: 150,
      header: (props) => (
        <THead
          id={ColumnViewType.OrganizationsRelationship}
          title='Relationship'
          renderFilter={() => <RelationshipFilter />}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      cell: (props) => {
        const id = props.row.original.value.metadata?.id;

        return <OrganizationRelationshipCell id={id} />;
      },
      skeleton: () => <Skeleton className='w-[100%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsOnboardingStatus]: columnHelper.accessor(
    'value.accountDetails',
    {
      id: ColumnViewType.OrganizationsOnboardingStatus,
      size: 125,
      cell: (props) => {
        const status = props.getValue()?.onboarding?.status;
        const updatedAt = props.getValue()?.onboarding?.updatedAt;

        return <OnboardingCell status={status} updatedAt={updatedAt} />;
      },
      header: (props) => (
        <THead
          id={ColumnViewType.OrganizationsOnboardingStatus}
          title='Onboarding'
          renderFilter={() => <OnboardingFilter />}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => (
        <div className='flex flex-col gap-1'>
          <Skeleton className='w-[33%] h-[14px]' />
        </div>
      ),
    },
  ),
  [ColumnViewType.OrganizationsRenewalLikelihood]: columnHelper.accessor(
    'value.accountDetails',
    {
      id: ColumnViewType.OrganizationsRenewalLikelihood,
      size: 100,
      cell: (props) => {
        const value = props.getValue()?.renewalSummary?.renewalLikelihood;

        return (
          <RenewalLikelihoodCell
            value={value}
            id={props.row.original.value.metadata?.id}
          />
        );
      },
      header: (props) => (
        <THead
          id={ColumnViewType.OrganizationsRenewalLikelihood}
          title='Health'
          data-testid='renewal-likelihood'
          renderFilter={() => <RenewalLikelihoodFilter />}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => (
        <div className='flex flex-col gap-1'>
          <Skeleton className='w-[25%] h-[14px]' />
        </div>
      ),
    },
  ),
  [ColumnViewType.OrganizationsRenewalDate]: columnHelper.accessor(
    'value.accountDetails',
    {
      id: ColumnViewType.OrganizationsRenewalDate,
      size: 125,
      cell: (props) => {
        const nextRenewalDate =
          props.getValue()?.renewalSummary?.nextRenewalDate;

        return <TimeToRenewalCell nextRenewalDate={nextRenewalDate} />;
      },

      header: (props) => (
        <THead
          id={ColumnViewType.OrganizationsRenewalDate}
          title='Renewal Date'
          filterWidth='15rem'
          renderFilter={() => <TimeToRenewalFilter />}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[50%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsForecastArr]: columnHelper.accessor(
    'value.accountDetails',
    {
      id: ColumnViewType.OrganizationsForecastArr,
      size: 150,
      cell: (props) => {
        const value = props.getValue()?.renewalSummary;
        const amount = value?.arrForecast;
        const potentialAmount = value?.maxArrForecast;

        return (
          <RenewalForecastCell
            amount={amount}
            potentialAmount={potentialAmount}
            id={props.row.original.value.metadata?.id}
          />
        );
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsForecastArr}
          title='ARR Forecast'
          filterWidth='17rem'
          renderFilter={(initialFocusRef) => (
            <ForecastFilter initialFocusRef={initialFocusRef} />
          )}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => (
        <div className='flex flex-col gap-1'>
          <Skeleton className='w-[50%] h-[14px]' />
          <Skeleton className='w-[25%] h-[14px]' />
        </div>
      ),
    },
  ),
  [ColumnViewType.OrganizationsOwner]: columnHelper.accessor('value.owner', {
    id: ColumnViewType.OrganizationsOwner,
    size: 150,
    cell: (props) => {
      return (
        <OwnerCell
          id={props.row.original.value.metadata?.id}
          owner={props.getValue()}
        />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.OrganizationsOwner}
        title='Owner'
        filterWidth='14rem'
        renderFilter={(initialFocusRef) => (
          <OwnerFilter initialFocusRef={initialFocusRef} />
        )}
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsLeadSource]: columnHelper.accessor(
    'value.leadSource',
    {
      id: ColumnViewType.OrganizationsLeadSource,
      size: 100,
      cell: (props) => {
        if (!props.getValue()) {
          return <p className='text-gray-400'>Unknown</p>;
        }

        return (
          <p className='text-gray-700 cursor-default truncate'>
            {props.getValue()}
          </p>
        );
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsLeadSource}
          title='Source'
          renderFilter={() => <SourceFilter />}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsCreatedDate]: columnHelper.accessor(
    'value.metadata.created',
    {
      id: ColumnViewType.OrganizationsCreatedDate,
      size: 125,
      cell: (props) => {
        const value = props.getValue();

        if (!value) {
          return <p className='text-gray-400'>Unknown</p>;
        }

        return (
          <p className='text-gray-700 cursor-default truncate'>
            {DateTimeUtils.format(
              value,
              DateTimeUtils.defaultFormatShortString,
            )}
          </p>
        );
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsCreatedDate}
          title='Created Date'
          filterWidth='14rem'
          renderFilter={() => <CreatedDateFilter />}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsYearFounded]: columnHelper.accessor(
    'value.yearFounded',
    {
      id: ColumnViewType.OrganizationsYearFounded,
      size: 100,
      cell: (props) => {
        const value = props.getValue();

        if (!value) {
          return <p className='text-gray-400'>Unknown</p>;
        }

        return <p className='text-gray-700 cursor-default truncate'>{value}</p>;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsYearFounded}
          title='Founded'
          filterWidth='17.5rem'
          renderFilter={() => (
            <NumericValueFilter
              property={ColumnViewType.OrganizationsYearFounded}
              label='age'
              suffix={'year'}
            />
          )}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsEmployeeCount]: columnHelper.accessor(
    'value.employees',
    {
      id: ColumnViewType.OrganizationsEmployeeCount,
      size: 125,
      cell: (props) => {
        const value = props.getValue();

        if (!value) {
          return <p className='text-gray-400'>Unknown</p>;
        }

        return <p className='text-gray-700 cursor-default truncate'>{value}</p>;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsEmployeeCount}
          title='Employees'
          filterWidth='17.5rem'
          renderFilter={() => (
            <NumericValueFilter
              property={ColumnViewType.OrganizationsEmployeeCount}
              label='employees'
            />
          )}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsSocials]: columnHelper.accessor(
    'value.socialMedia',
    {
      id: ColumnViewType.OrganizationsSocials,
      size: 125,
      enableSorting: false,
      cell: (props) => (
        <OrganizationLinkedInCell organizationId={props.row.original.id} />
      ),
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsSocials}
          title='LinkedIn'
          filterWidth='14rem'
          renderFilter={(initialFocusRef) => (
            <SocialsFilter initialFocusRef={initialFocusRef} />
          )}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsLastTouchpoint]: columnHelper.accessor(
    (row) => row,
    {
      id: ColumnViewType.OrganizationsLastTouchpoint,
      size: 200,
      cell: (props) => (
        <LastTouchpointCell
          lastTouchPointAt={
            props.row.original?.value?.lastTouchpoint?.lastTouchPointAt
          }
          lastTouchPointTimelineEvent={
            props.row.original?.value?.lastTouchpoint
              ?.lastTouchPointTimelineEvent
          }
          lastTouchPointType={
            props.row.original?.value?.lastTouchpoint?.lastTouchPointType
          }
        />
      ),
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsLastTouchpoint}
          title='Last Touchpoint'
          renderFilter={() => <LastTouchpointFilter />}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => (
        <div className='flex flex-col gap-1'>
          <Skeleton className='w-[75%] h-[14px]' />
          <Skeleton className='w-[100%] h-[14px]' />
        </div>
      ),
    },
  ),
  [ColumnViewType.OrganizationsLastTouchpointDate]: columnHelper.accessor(
    (row) => row,
    {
      id: ColumnViewType.OrganizationsLastTouchpointDate,
      size: 150,
      enableSorting: true,
      cell: (props) => (
        <LastTouchpointDateCell
          lastTouchPointAt={
            props.row.original?.value?.lastTouchpoint?.lastTouchPointAt
          }
        />
      ),
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsLastTouchpointDate}
          title='Last Interacted'
          renderFilter={() => <LastInteractedFilter />}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => (
        <div className='flex flex-col gap-1'>
          <Skeleton className='w-[75%] h-[14px]' />
          <Skeleton className='w-[100%] h-[14px]' />
        </div>
      ),
    },
  ),
  [ColumnViewType.OrganizationsChurnDate]: columnHelper.accessor(
    'value.accountDetails',
    {
      id: ColumnViewType.OrganizationsChurnDate,
      size: 115,
      cell: (props) => {
        const value = props.row.original.value.accountDetails?.churned;

        return (
          <p
            className={cn(
              'text-gray-700 cursor-default',
              !value && 'text-gray-400',
            )}
          >
            {DateTimeUtils.format(
              value,
              DateTimeUtils.defaultFormatShortString,
            ) || 'Unknown'}
          </p>
        );
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsChurnDate}
          title='Churn Date'
          renderFilter={() => {
            return <ChurnedFilter />;
          }}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsLtv]: columnHelper.accessor(
    'value.accountDetails',
    {
      id: ColumnViewType.OrganizationsLtv,
      size: 100,
      cell: (props) => {
        const value = props.row.original.value.accountDetails?.ltv;

        const formatedValue = formatCurrency(value || 0, 0);

        return (
          <p
            className={cn(
              'text-gray-700 cursor-default',
              !value && 'text-gray-400',
            )}
          >
            {value ? `${formatedValue}` : 'Unknown'}
          </p>
        );
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsLtv}
          title='LTV'
          filterWidth='14rem'
          renderFilter={(initialFocusRef) => {
            return <LtvFilter initialFocusRef={initialFocusRef} />;
          }}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsIndustry]: columnHelper.accessor(
    'value.industry',
    {
      id: ColumnViewType.OrganizationsIndustry,
      size: 200,
      cell: (props) => {
        const value = props.getValue();

        return <IndustryCell value={value} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsIndustry}
          title='Industry'
          filterWidth='17rem'
          renderFilter={(initialFocusRef) => (
            <IndustryFilter initialFocusRef={initialFocusRef} />
          )}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsContactCount]: columnHelper.accessor('value', {
    id: ColumnViewType.OrganizationsContactCount,
    size: 80,
    enableColumnFilter: false,

    cell: (props) => {
      const value = props.getValue()?.contacts?.content?.length;

      return <div>{value}</div>;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.OrganizationsContactCount}
        title='Contacts'
        filterWidth='auto'
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsLinkedinFollowerCount]: columnHelper.accessor(
    'value',
    {
      id: ColumnViewType.OrganizationsLinkedinFollowerCount,
      size: 150,
      cell: (props) => {
        const value = props
          .getValue()
          ?.socialMedia.find((e: Social) =>
            e?.url?.includes('linkedin'),
          )?.followersCount;
        if (typeof value !== 'number')
          return <div className='text-gray-400'>Unknown</div>;

        return <div>{Number(value).toLocaleString()}</div>;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsLinkedinFollowerCount}
          title='LinkedIn Followers'
          filterWidth='17.5rem'
          renderFilter={() => (
            <NumericValueFilter
              property={ColumnViewType.OrganizationsLinkedinFollowerCount}
              label='followers'
            />
          )}
          {...getTHeadProps<Store<Organization>>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsTags]: columnHelper.accessor('value', {
    id: ColumnViewType.OrganizationsTags,
    size: 150,
    enableColumnFilter: false,
    cell: (props) => {
      const value = props.getValue()?.metadata?.id;

      return <OrganizationsTagsCell id={value} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.OrganizationsTags}
        title='Tags'
        filterWidth='auto'
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsIsPublic]: columnHelper.accessor('value', {
    id: ColumnViewType.OrganizationsIsPublic,
    size: 150,
    cell: (props) => {
      const value = props.getValue()?.public;
      if (value === undefined) {
        return <div className='text-gray-400'>Unknown</div>;
      }

      return <div>{value ? 'Public' : 'Private'}</div>;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.OrganizationsIsPublic}
        title='Ownership Type'
        renderFilter={(initialFocusRef) => (
          <OwnershipTypeFilter
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.OrganizationsIsPublic}
          />
        )}
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsStage]: columnHelper.accessor('value', {
    id: ColumnViewType.OrganizationsStage,
    size: 120,
    enableColumnFilter: true,
    enableSorting: true,
    cell: (props) => {
      return (
        <OrganizationStageCell id={props.row.original.value.metadata?.id} />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.OrganizationsStage}
        title='Stage'
        filterWidth='auto'
        renderFilter={() => <StageFilter />}
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsCity]: columnHelper.accessor('value.metadata', {
    id: ColumnViewType.OrganizationsCity,
    size: 210,
    cell: (props) => {
      const value = props.getValue()?.id;

      return <CountryCell id={value} type='organization' />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.OrganizationsCity}
        title='Headquarters'
        filterWidth='auto'
        renderFilter={(initialFocusRef) => (
          <LocationFilter
            type='organizations'
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.OrganizationsCity}
            locationType='countryCodeA2'
          />
        )}
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
};

export const getOrganizationColumnsConfig = (
  tableViewDef?: Array<TableViewDef>[0],
) => getColumnConfig<ColumnDatum>(columns, tableViewDef);
