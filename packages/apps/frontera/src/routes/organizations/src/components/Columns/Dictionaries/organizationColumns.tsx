import { Store } from '@store/store.ts';
import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { cn } from '@ui/utils/cn.ts';
import { DateTimeUtils } from '@utils/date.ts';
import { createColumnHelper } from '@ui/presentation/Table';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton.tsx';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber.ts';
import { Contact, Organization, ColumnViewType } from '@graphql/types';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead.tsx';
import { OrganizationLinkedInCell } from '@organizations/components/Columns/Cells/socials/OrganizationLinkedInCell.tsx';

import { AvatarHeader } from '../Headers/Avatar';
import { LastTouchpointDateCell } from '../Cells/touchpointDate';
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
  RenewalLikelihoodCell,
  OrganizationRelationshipCell,
} from '../Cells';
import {
  LtvFilter,
  OwnerFilter,
  SourceFilter,
  WebsiteFilter,
  ChurnedFilter,
  SocialsFilter,
  ForecastFilter,
  IndustryFilter,
  EmployeesFilter,
  OnboardingFilter,
  CreatedDateFilter,
  OrganizationFilter,
  RelationshipFilter,
  TimeToRenewalFilter,
  LastTouchpointFilter,
  LastInteractedFilter,
  RenewalLikelihoodFilter,
} from '../Filters';

type ColumnDatum = Store<Organization>;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

export const organizationColumns: Record<string, Column> = {
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
      size: 125,
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
          title='Next Renewal'
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
      enableColumnFilter: false,
      cell: (props) => {
        const value = props.getValue();

        if (!value) {
          return <p className='text-gray-400'>Unknown</p>;
        }

        return <p className='text-gray-700 cursor-default truncate'>{value}</p>;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          id={ColumnViewType.OrganizationsForecastArr}
          title='Founded'
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
          renderFilter={() => <EmployeesFilter />}
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
            {/* {value || 'Unknown'} */}
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
          filterWidth='auto'
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
    size: 200,
    enableColumnFilter: false,
    cell: (props) => {
      const value = props
        .getValue()
        ?.contacts?.content?.filter((e: Contact) => e.tags)?.length;

      return <div>{value}</div>;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        id={ColumnViewType.OrganizationsContactCount}
        title='Tagged Contacts'
        filterWidth='auto'
        {...getTHeadProps<Store<Organization>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
};
