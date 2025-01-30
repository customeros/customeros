import { Organization } from '@store/Organizations/Organization.dto';
import { CountryCell } from '@finder/components/Columns/Cells/country';
import { OrganizationStageCell } from '@finder/components/Columns/Cells/stage';
import { DateCell } from '@finder/components/Columns/shared/Cells/DateCell/DateCell';
import {
  ColumnDef,
  ColumnDef as ColumnDefinition,
} from '@tanstack/react-table';
import { AvatarHeader } from '@finder/components/Columns/organizations/Headers/Avatar';
import { getColumnConfig } from '@finder/components/Columns/shared/util/getColumnConfig';

import { cn } from '@ui/utils/cn';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import { Social, TableViewDef, ColumnViewType } from '@graphql/types';

import {
  OwnerCell,
  AvatarCell,
  DomainsCell,
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

type ColumnDatum = Organization;

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
      const enrichedOrg = props.getValue()?.value;

      const icon = enrichedOrg?.iconUrl;
      const logo = enrichedOrg?.logoUrl;
      const isEnriching = props.getValue()?.isEnriching;

      return (
        <AvatarCell
          icon={icon}
          logo={logo}
          isEnriching={isEnriching}
          id={props.getValue()?.value?.id}
          name={props.getValue()?.value?.name}
        />
      );
    },
    header: AvatarHeader,
    skeleton: () => <Skeleton className='size-[24px]' />,
  }),
  [ColumnViewType.OrganizationsName]: columnHelper.accessor('value.id', {
    id: ColumnViewType.OrganizationsName,
    minSize: 160,
    size: 160,
    maxSize: 400,
    enableColumnFilter: false,
    enableResizing: true,
    cell: (props) => {
      return <OrganizationCell id={props.getValue()} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Company'
        id={ColumnViewType.OrganizationsName}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsPrimaryDomains]: columnHelper.accessor(
    'value.domainsDetails',
    {
      id: ColumnViewType.OrganizationsPrimaryDomains,
      minSize: 152,
      maxSize: 500,
      enableColumnFilter: false,
      enableResizing: true,
      enableSorting: false,
      cell: (props) => {
        const organizationId = props.row.original.value.id;

        return <DomainsCell organizationId={organizationId} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Primary Domains'
          id={ColumnViewType.OrganizationsPrimaryDomains}
          {...getTHeadProps<Organization>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[50%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsRelationship]: columnHelper.accessor(
    'value.relationship',
    {
      id: ColumnViewType.OrganizationsRelationship,
      minSize: 119,
      maxSize: 400,
      enableColumnFilter: false,
      enableResizing: true,
      header: (props) => (
        <THead
          title='Relationship'
          id={ColumnViewType.OrganizationsRelationship}
          {...getTHeadProps<Organization>(props)}
        />
      ),
      cell: (props) => {
        const id = props.row.original.value?.id;

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
    'value',
    {
      id: ColumnViewType.OrganizationsOnboardingStatus,
      minSize: 114,
      maxSize: 400,
      enableColumnFilter: false,
      enableResizing: true,
      cell: (props) => {
        const value = props.getValue();
        const status = value.onboardingStatus;
        const updatedAt = value.onboardingStatusUpdatedAt;

        return <OnboardingCell status={status} updatedAt={updatedAt} />;
      },
      header: (props) => (
        <THead
          title='Onboarding'
          id={ColumnViewType.OrganizationsOnboardingStatus}
          {...getTHeadProps<Organization>(props)}
        />
      ),
      skeleton: () => (
        <div className='flex flex-col gap-1'>
          <Skeleton className='w-[33%] h-[14px]' />
        </div>
      ),
    },
  ),
  [ColumnViewType.OrganizationsRenewalLikelihood]: columnHelper.accessor('id', {
    id: ColumnViewType.OrganizationsRenewalLikelihood,
    minSize: 82,
    size: 110,
    maxSize: 400,
    enableColumnFilter: false,
    enableResizing: true,
    enableSorting: true,
    cell: (props) => {
      return <RenewalLikelihoodCell id={props.getValue()} />;
    },
    header: (props) => (
      <THead
        title='Health'
        data-test='renewal-likelihood'
        id={ColumnViewType.OrganizationsRenewalLikelihood}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[25%] h-[14px]' />
      </div>
    ),
  }),
  [ColumnViewType.OrganizationsRenewalDate]: columnHelper.accessor('value', {
    id: ColumnViewType.OrganizationsRenewalDate,
    minSize: 156,
    size: 156,
    maxSize: 400,
    enableColumnFilter: false,
    enableResizing: true,
    enableSorting: true,
    cell: (props) => {
      const nextRenewalDate = props.getValue()?.renewalSummaryNextRenewalAt;

      return <TimeToRenewalCell nextRenewalDate={nextRenewalDate} />;
    },

    header: (props) => (
      <THead
        title='Renewal Date'
        id={ColumnViewType.OrganizationsRenewalDate}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[50%] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsForecastArr]: columnHelper.accessor('value', {
    id: ColumnViewType.OrganizationsForecastArr,
    minSize: 154,
    size: 154,
    maxSize: 400,
    enableColumnFilter: false,
    enableResizing: true,
    enableSorting: true,
    cell: (props) => {
      const value = props.getValue();
      const amount = value?.renewalSummaryArrForecast;
      const potentialAmount = value?.renewalSummaryMaxArrForecast;

      return (
        <RenewalForecastCell
          id={value?.id}
          amount={amount}
          potentialAmount={potentialAmount}
        />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='ARR Forecast'
        id={ColumnViewType.OrganizationsForecastArr}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[50%] h-[14px]' />
        <Skeleton className='w-[25%] h-[14px]' />
      </div>
    ),
  }),
  [ColumnViewType.OrganizationsOwner]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.OrganizationsOwner,
    minSize: 82,
    size: 154,
    maxSize: 400,
    enableColumnFilter: false,
    enableResizing: true,
    cell: (props) => {
      const row = props.getValue();

      const owner = row?.value?.owner;

      return <OwnerCell ownerId={owner?.id} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Owner'
        id={ColumnViewType.OrganizationsOwner}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsLeadSource]: columnHelper.accessor(
    'value.leadSource',
    {
      id: ColumnViewType.OrganizationsLeadSource,
      minSize: 84,
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
          {...getTHeadProps<Organization>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsCreatedDate]: columnHelper.accessor(
    'value.createdAt',
    {
      id: ColumnViewType.OrganizationsCreatedDate,
      size: 145,
      minSize: 122,
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
          {...getTHeadProps<Organization>(props)}
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
      minSize: 95,
      maxSize: 400,
      enableResizing: true,
      enableColumnFilter: false,
      cell: (props) => {
        const value = props.getValue();
        const isEnriching = props.row.original.isEnriching;

        if (!value) {
          return (
            <p className='text-gray-400'>
              {isEnriching ? 'Enriching...' : 'Not set'}
            </p>
          );
        }

        return <p className='text-gray-700 cursor-default truncate'>{value}</p>;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Founded'
          id={ColumnViewType.OrganizationsYearFounded}
          {...getTHeadProps<Organization>(props)}
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
      minSize: 110,
      maxSize: 400,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: false,
      cell: (props) => {
        const value = props.getValue();
        const isEnriching = props.row.original.isEnriching;

        if (!value) {
          return (
            <p className='text-gray-400'>
              {isEnriching ? 'Enriching...' : 'Not set'}
            </p>
          );
        }

        return <p className='text-gray-700 cursor-default truncate'>{value}</p>;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Employees'
          id={ColumnViewType.OrganizationsEmployeeCount}
          {...getTHeadProps<Organization>(props)}
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
      minSize: 94,
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
          {...getTHeadProps<Organization>(props)}
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
      minSize: 140,
      maxSize: 400,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: true,
      cell: (props) => (
        <LastTouchpointCell
          lastTouchPointAt={props.row.original?.value?.lastTouchPointAt}
          lastTouchPointType={props.row.original?.value?.lastTouchPointType}
        />
      ),
      header: (props) => (
        <THead<HTMLInputElement>
          title='Last Touchpoint'
          id={ColumnViewType.OrganizationsLastTouchpoint}
          {...getTHeadProps<Organization>(props)}
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
    'value.lastTouchPointAt',
    {
      id: ColumnViewType.OrganizationsLastTouchpointDate,
      size: 154,
      minSize: 136,
      maxSize: 400,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: true,
      cell: (props) => (
        <LastTouchpointDateCell lastTouchPointAt={props.getValue()} />
      ),
      header: (props) => (
        <THead<HTMLInputElement>
          title='Last Interacted'
          id={ColumnViewType.OrganizationsLastTouchpointDate}
          {...getTHeadProps<Organization>(props)}
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
    'value.churnedAt',
    {
      id: ColumnViewType.OrganizationsChurnDate,
      size: 115,
      minSize: 110,
      maxSize: 400,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: true,
      cell: (props) => {
        return <DateCell value={props.getValue()} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Churn Date'
          id={ColumnViewType.OrganizationsChurnDate}
          {...getTHeadProps<Organization>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsLtv]: columnHelper.accessor('value.ltv', {
    id: ColumnViewType.OrganizationsLtv,
    size: 110,
    minSize: 64,
    maxSize: 600,
    enableResizing: true,
    enableColumnFilter: false,
    cell: (props) => {
      const value = props.getValue();
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
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsIndustry]: columnHelper.accessor(
    'value.industryName',
    {
      id: ColumnViewType.OrganizationsIndustry,
      minSize: 95,
      maxSize: 600,
      enableResizing: true,
      cell: (props) => {
        const value = props.getValue();
        const isEnriching = props.row.original.isEnriching;
        const flaggedArWrong = props.row.original.value.wrongIndustry;

        return (
          <IndustryCell
            value={value}
            id={props.row.original.id}
            enrichingStatus={isEnriching}
            flaggedAsIncorrect={flaggedArWrong}
          />
        );
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Industry'
          id={ColumnViewType.OrganizationsIndustry}
          {...getTHeadProps<Organization>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsContactCount]: columnHelper.accessor('value', {
    id: ColumnViewType.OrganizationsContactCount,
    minSize: 94,
    maxSize: 400,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,

    cell: (props) => {
      const value = props.row.original.value.contactCount;

      return (
        <div data-test='organization-contacts-in-all-orgs-table'>{value}</div>
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Contacts'
        filterWidth='auto'
        id={ColumnViewType.OrganizationsContactCount}
        {...getTHeadProps<Organization>(props)}
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
          {...getTHeadProps<Organization>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsTags]: columnHelper.accessor('id', {
    id: ColumnViewType.OrganizationsTags,
    size: 154,
    minSize: 70,
    maxSize: 400,
    enableResizing: true,
    enableSorting: false,
    cell: (props) => {
      return <OrganizationsTagsCell id={props.getValue()} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Tags'
        filterWidth='auto'
        id={ColumnViewType.OrganizationsTags}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsIsPublic]: columnHelper.accessor(
    'value.public',
    {
      id: ColumnViewType.OrganizationsIsPublic,
      size: 154,
      minSize: 142,
      maxSize: 400,
      enableResizing: true,
      enableColumnFilter: false,
      cell: (props) => {
        const value = props.getValue();

        if (value === undefined) {
          return <div className='text-gray-400'>Unknown</div>;
        }

        return <div>{value ? 'Public' : 'Private'}</div>;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Ownership Type'
          id={ColumnViewType.OrganizationsIsPublic}
          {...getTHeadProps<Organization>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OrganizationsStage]: columnHelper.accessor('id', {
    id: ColumnViewType.OrganizationsStage,
    size: 154,
    minSize: 76,
    maxSize: 400,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      return <OrganizationStageCell id={props.getValue()} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Stage'
        filterWidth='auto'
        id={ColumnViewType.OrganizationsStage}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),

  [ColumnViewType.OrganizationsCountry]: columnHelper.accessor('id', {
    id: ColumnViewType.OrganizationsCountry,
    size: 210,
    minSize: 91,
    maxSize: 400,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      return <CountryCell type='organization' id={props.getValue()} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Country'
        filterWidth='auto'
        id={ColumnViewType.OrganizationsCountry}
        {...getTHeadProps<Organization>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[75%] h-[14px]' />,
  }),
  [ColumnViewType.OrganizationsUpdatedDate]: columnHelper.accessor(
    'value.updatedAt',
    {
      id: ColumnViewType.OrganizationsUpdatedDate,
      minSize: 125,
      maxSize: 600,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: true,
      cell: (props) => {
        const lastUpdatedAt = props.getValue();

        return <DateCell value={lastUpdatedAt} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Last Updated'
          id={ColumnViewType.OrganizationsUpdatedDate}
          {...getTHeadProps<Organization>(props)}
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
