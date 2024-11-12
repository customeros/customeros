import { CountryCell } from '@finder/components/Columns/Cells/country';
import { OrganizationStore } from '@store/Organizations/Organization.store';
import { OrganizationStageCell } from '@finder/components/Columns/Cells/stage';
import {
  ColumnDef,
  ColumnDef as ColumnDefinition,
} from '@tanstack/react-table';
import { AvatarHeader } from '@finder/components/Columns/organizations/Headers/Avatar';
import { DateCell } from '@finder/components/Columns/shared/Cells/DateCell/DateCell.tsx';
import { getColumnConfig } from '@finder/components/Columns/shared/util/getColumnConfig.ts';

import { cn } from '@ui/utils/cn.ts';
import { createColumnHelper } from '@ui/presentation/Table';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton.tsx';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber.ts';
import { Social, TableViewDef, ColumnViewType } from '@graphql/types';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead.tsx';

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

type ColumnDatum = OrganizationStore;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

export const columns: Record<string, Column> = {
  [ColumnViewType.OrganizationsAvatar]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.OrganizationsAvatar,
    size: 29,
    minSize: 29,
    maxSize: 29,
    enableColumnFilter: false,
    enableResizing: false,
    cell: (props) => {
      const icon = props.getValue()?.value?.icon;
      const logo = props.getValue()?.value?.logo;
      const description = props.getValue()?.value?.valueProposition;

      return (
        <AvatarCell
          icon={icon}
          logo={logo}
          description={description}
          name={props.getValue()?.value?.name}
          id={props.getValue()?.value?.metadata?.id}
        />
      );
    },
    header: AvatarHeader,
    skeleton: () => <Skeleton className='size-[24px]' />,
  }),
  [ColumnViewType.OrganizationsName]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.OrganizationsName,
    minSize: 160,
    size: 160,
    maxSize: 400,
    enableColumnFilter: false,
    enableResizing: true,
    cell: (props) => {
      return (
        <OrganizationCell
          name={props.getValue().value.name}
          id={props.getValue().value.metadata?.id}
          isSubsidiary={!!props.getValue()?.value?.subsidiaryOf?.length}
          parentOrganizationName={
            props.getValue()?.value?.subsidiaryOf?.[0]?.organization.name
          }
        />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Organization'
        id={ColumnViewType.OrganizationsName}
        {...getTHeadProps<OrganizationStore>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsWebsite]: columnHelper.accessor(
    'value.website',
    {
      id: ColumnViewType.OrganizationsWebsite,
      minSize: 125,
      maxSize: 400,
      enableColumnFilter: false,
      enableResizing: true,
      enableSorting: false,
      cell: (props) => {
        const organizationId = props.row.original.value.metadata.id;

        return <WebsiteCell organizationId={organizationId} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Website'
          id={ColumnViewType.OrganizationsWebsite}
          {...getTHeadProps<OrganizationStore>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[50%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsRelationship]: columnHelper.accessor(
    'value.relationship',
    {
      id: ColumnViewType.OrganizationsRelationship,
      minSize: 160,
      maxSize: 400,
      enableColumnFilter: false,
      enableResizing: true,
      header: (props) => (
        <THead
          title='Relationship'
          id={ColumnViewType.OrganizationsRelationship}
          {...getTHeadProps<OrganizationStore>(props)}
        />
      ),
      cell: (props) => {
        const id = props.row.original.value.metadata?.id;

        return (
          <OrganizationRelationshipCell
            id={id}
            dataTest='organization-relationship-button-in-all-orgs-table'
          />
        );
      },
      skeleton: () => <Skeleton className='w-[100%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsOnboardingStatus]: columnHelper.accessor(
    'value.accountDetails',
    {
      id: ColumnViewType.OrganizationsOnboardingStatus,
      minSize: 125,
      maxSize: 400,
      enableColumnFilter: false,
      enableResizing: true,
      cell: (props) => {
        const status = props.getValue()?.onboarding?.status;
        const updatedAt = props.getValue()?.onboarding?.updatedAt;

        return <OnboardingCell status={status} updatedAt={updatedAt} />;
      },
      header: (props) => (
        <THead
          title='Onboarding'
          id={ColumnViewType.OrganizationsOnboardingStatus}
          {...getTHeadProps<OrganizationStore>(props)}
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
      minSize: 110,
      size: 110,
      maxSize: 400,
      enableColumnFilter: false,
      enableResizing: true,
      enableSorting: true,
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
          title='Health'
          data-test='renewal-likelihood'
          id={ColumnViewType.OrganizationsRenewalLikelihood}
          {...getTHeadProps<OrganizationStore>(props)}
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
      minSize: 156,
      size: 156,
      maxSize: 400,
      enableColumnFilter: false,
      enableResizing: true,
      enableSorting: true,
      cell: (props) => {
        const nextRenewalDate =
          props.getValue()?.renewalSummary?.nextRenewalDate;

        return <TimeToRenewalCell nextRenewalDate={nextRenewalDate} />;
      },

      header: (props) => (
        <THead
          title='Renewal Date'
          id={ColumnViewType.OrganizationsRenewalDate}
          {...getTHeadProps<OrganizationStore>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[50%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsForecastArr]: columnHelper.accessor(
    'value.accountDetails',
    {
      id: ColumnViewType.OrganizationsForecastArr,
      minSize: 154,
      size: 154,
      maxSize: 400,
      enableColumnFilter: false,
      enableResizing: true,
      enableSorting: true,
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
          title='ARR Forecast'
          id={ColumnViewType.OrganizationsForecastArr}
          {...getTHeadProps<OrganizationStore>(props)}
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
    minSize: 154,
    size: 154,
    maxSize: 400,
    enableColumnFilter: false,
    enableResizing: true,
    cell: (props) => {
      return (
        <OwnerCell
          owner={props.getValue()}
          id={props.row.original.value.metadata?.id}
        />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Owner'
        id={ColumnViewType.OrganizationsOwner}
        {...getTHeadProps<OrganizationStore>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsLeadSource]: columnHelper.accessor(
    'value.leadSource',
    {
      id: ColumnViewType.OrganizationsLeadSource,
      minSize: 100,
      size: 100,
      maxSize: 400,
      enableColumnFilter: false,
      enableResizing: true,
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
          title='Source'
          id={ColumnViewType.OrganizationsLeadSource}
          {...getTHeadProps<OrganizationStore>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsCreatedDate]: columnHelper.accessor(
    'value.metadata.created',
    {
      id: ColumnViewType.OrganizationsCreatedDate,
      size: 145,
      minSize: 145,
      maxSize: 400,
      enableResizing: true,
      enableColumnFilter: false,
      cell: (props) => {
        const value = props.getValue();

        return <DateCell value={value} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Created Date'
          id={ColumnViewType.OrganizationsCreatedDate}
          {...getTHeadProps<OrganizationStore>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsYearFounded]: columnHelper.accessor(
    'value.yearFounded',
    {
      id: ColumnViewType.OrganizationsYearFounded,
      size: 120,
      minSize: 120,
      maxSize: 400,
      enableResizing: true,
      enableColumnFilter: false,
      cell: (props) => {
        const value = props.getValue();
        const enrichedOrg = props.row.original.value.enrichDetails;
        const enrichingStatus =
          !enrichedOrg?.enrichedAt &&
          enrichedOrg?.requestedAt &&
          !enrichedOrg?.failedAt;

        if (!value) {
          return (
            <p className='text-gray-400'>
              {enrichingStatus ? 'Enriching...' : 'Not set'}
            </p>
          );
        }

        return <p className='text-gray-700 cursor-default truncate'>{value}</p>;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Founded'
          id={ColumnViewType.OrganizationsYearFounded}
          {...getTHeadProps<OrganizationStore>(props)}
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
      minSize: 125,
      maxSize: 400,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: false,
      cell: (props) => {
        const value = props.getValue();
        const enrichedOrg = props.row.original.value.enrichDetails;
        const enrichingStatus =
          !enrichedOrg?.enrichedAt &&
          enrichedOrg?.requestedAt &&
          !enrichedOrg?.failedAt;

        if (!value) {
          return (
            <p className='text-gray-400'>
              {enrichingStatus ? 'Enriching...' : 'Not set'}
            </p>
          );
        }

        return <p className='text-gray-700 cursor-default truncate'>{value}</p>;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Employees'
          id={ColumnViewType.OrganizationsEmployeeCount}
          {...getTHeadProps<OrganizationStore>(props)}
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
      minSize: 125,
      maxSize: 400,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: false,
      cell: (props) => (
        <OrganizationLinkedInCell organizationId={props.row.original.id} />
      ),
      header: (props) => (
        <THead<HTMLInputElement>
          title='LinkedIn'
          id={ColumnViewType.OrganizationsSocials}
          {...getTHeadProps<OrganizationStore>(props)}
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
      minSize: 200,
      maxSize: 400,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: true,
      cell: (props) => (
        <LastTouchpointCell
          lastTouchPointAt={
            props.row.original?.value?.lastTouchpoint?.lastTouchPointAt
          }
          lastTouchPointType={
            props.row.original?.value?.lastTouchpoint?.lastTouchPointType
          }
          lastTouchPointTimelineEvent={
            props.row.original?.value?.lastTouchpoint
              ?.lastTouchPointTimelineEvent
          }
        />
      ),
      header: (props) => (
        <THead<HTMLInputElement>
          title='Last Touchpoint'
          id={ColumnViewType.OrganizationsLastTouchpoint}
          {...getTHeadProps<OrganizationStore>(props)}
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
      size: 154,
      minSize: 154,
      maxSize: 400,
      enableResizing: true,
      enableColumnFilter: false,
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
          title='Last Interacted'
          id={ColumnViewType.OrganizationsLastTouchpointDate}
          {...getTHeadProps<OrganizationStore>(props)}
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
      minSize: 115,
      maxSize: 400,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: true,
      cell: (props) => {
        const value = props.row.original.value.accountDetails?.churned;

        return <DateCell value={value} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Churn Date'
          id={ColumnViewType.OrganizationsChurnDate}
          {...getTHeadProps<OrganizationStore>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsLtv]: columnHelper.accessor(
    'value.accountDetails',
    {
      id: ColumnViewType.OrganizationsLtv,
      size: 110,
      minSize: 100,
      maxSize: 600,
      enableResizing: true,
      enableColumnFilter: false,
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
          title='LTV'
          id={ColumnViewType.OrganizationsLtv}
          {...getTHeadProps<OrganizationStore>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsIndustry]: columnHelper.accessor(
    'value.industry',
    {
      id: ColumnViewType.OrganizationsIndustry,
      minSize: 175,
      maxSize: 600,
      enableResizing: true,
      cell: (props) => {
        const value = props.getValue();
        const enrichedOrg = props.row.original.value.enrichDetails;
        const enrichingStatus =
          !enrichedOrg?.enrichedAt &&
          enrichedOrg?.requestedAt &&
          !enrichedOrg?.failedAt;

        return <IndustryCell value={value} enrichingStatus={enrichingStatus} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Industry'
          id={ColumnViewType.OrganizationsIndustry}
          {...getTHeadProps<OrganizationStore>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsContactCount]: columnHelper.accessor('value', {
    id: ColumnViewType.OrganizationsContactCount,
    minSize: 90,
    maxSize: 400,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,

    cell: (props) => {
      const value = props.getValue()?.contacts?.content?.length;

      return (
        <div data-test='organization-contacts-in-all-orgs-table'>{value}</div>
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Contacts'
        filterWidth='auto'
        id={ColumnViewType.OrganizationsContactCount}
        {...getTHeadProps<OrganizationStore>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsLinkedinFollowerCount]: columnHelper.accessor(
    'value',
    {
      id: ColumnViewType.OrganizationsLinkedinFollowerCount,
      minSize: 175,
      maxSize: 400,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: true,

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
          title='LinkedIn Followers'
          id={ColumnViewType.OrganizationsLinkedinFollowerCount}
          {...getTHeadProps<OrganizationStore>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsTags]: columnHelper.accessor('value', {
    id: ColumnViewType.OrganizationsTags,
    size: 154,
    minSize: 154,
    maxSize: 400,
    enableResizing: true,
    enableSorting: false,
    cell: (props) => {
      const value = props.getValue()?.metadata?.id;

      return <OrganizationsTagsCell id={value} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Tags'
        filterWidth='auto'
        id={ColumnViewType.OrganizationsTags}
        {...getTHeadProps<OrganizationStore>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsIsPublic]: columnHelper.accessor('value', {
    id: ColumnViewType.OrganizationsIsPublic,
    size: 154,
    minSize: 154,
    maxSize: 400,
    enableResizing: true,
    enableColumnFilter: false,
    cell: (props) => {
      const value = props.getValue()?.public;

      if (value === undefined) {
        return <div className='text-gray-400'>Unknown</div>;
      }

      return <div>{value ? 'Public' : 'Private'}</div>;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Ownership Type'
        id={ColumnViewType.OrganizationsIsPublic}
        {...getTHeadProps<OrganizationStore>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsStage]: columnHelper.accessor('value', {
    id: ColumnViewType.OrganizationsStage,
    size: 154,
    minSize: 154,
    maxSize: 400,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      return (
        <OrganizationStageCell id={props.row.original.value.metadata?.id} />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Stage'
        filterWidth='auto'
        id={ColumnViewType.OrganizationsStage}
        {...getTHeadProps<OrganizationStore>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),

  [ColumnViewType.OrganizationsHeadquarters]: columnHelper.accessor(
    'value.metadata',
    {
      id: ColumnViewType.OrganizationsHeadquarters,
      size: 210,
      minSize: 210,
      maxSize: 400,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: true,
      cell: (props) => {
        const value = props.getValue()?.id;

        return <CountryCell id={value} type='organization' />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Country'
          filterWidth='auto'
          id={ColumnViewType.OrganizationsHeadquarters}
          {...getTHeadProps<OrganizationStore>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),

  [ColumnViewType.OrganizationsParentOrganization]: columnHelper.accessor(
    (row) => row,
    {
      id: ColumnViewType.OrganizationsParentOrganization,
      size: 210,
      minSize: 210,
      maxSize: 400,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: true,
      cell: (props) => {
        const parentOrg =
          props.getValue()?.value?.parentCompanies?.[0]?.organization;

        if (!parentOrg) return null;

        return (
          <OrganizationCell
            isSubsidiary={false}
            name={parentOrg?.name}
            parentOrganizationName={''}
            id={parentOrg?.metadata?.id}
          />
        );
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Parent Org'
          id={ColumnViewType.OrganizationsParentOrganization}
          {...getTHeadProps<OrganizationStore>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
};

export const getOrganizationColumnsConfig = (
  tableViewDef?: Array<TableViewDef>[0],
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
): ColumnDef<ColumnDatum, any>[] =>
  getColumnConfig<ColumnDatum>(columns, tableViewDef);
