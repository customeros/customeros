import { Cake } from '@ui/media/icons/Cake';
import { Key01 } from '@ui/media/icons/Key01';
import { Tag01 } from '@ui/media/icons/Tag01';
import { Hash02 } from '@ui/media/icons/Hash02';
import { Globe01 } from '@ui/media/icons/Globe01';
import { Users03 } from '@ui/media/icons/Users03';
import { Trophy01 } from '@ui/media/icons/Trophy01';
import { Calendar } from '@ui/media/icons/Calendar';
import { Activity } from '@ui/media/icons/Activity';
import { Columns03 } from '@ui/media/icons/Columns03';
import { Building07 } from '@ui/media/icons/Building07';
import { Calculator } from '@ui/media/icons/Calculator';
import { Building05 } from '@ui/media/icons/Building05';
import { ArrowCircleDownRight } from '@ui/media/icons/ArrowCircleDownRight';
import { CurrencyDollarCircle } from '@ui/media/icons/CurrencyDollarCircle';
import { AlignHorizontalCentre02 } from '@ui/media/icons/AlignHorizontalCentre02';
import {
  ColumnViewType,
  OnboardingStatus,
  OrganizationStage,
  ComparisonOperator,
  LastTouchpointType,
  OrganizationRelationship,
  OpportunityRenewalLikelihood,
} from '@shared/types/__generated__/graphql.types';

export type FilterType = {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  options?: any[];
  icon: JSX.Element;
  filterName: string;
  filterAccesor: ColumnViewType;
  filterOperators: ComparisonOperator[];
  filterType: 'text' | 'date' | 'number' | 'list';
  groupOptions?: { label: string; options: { id: string; label: string }[] };
};

import { uniqBy } from 'lodash';
import { type RootStore } from '@store/root';

import { Globe04 } from '@ui/media/icons/Globe04';
import { LinkedinOutline } from '@ui/media/icons/LinkedinOutline';

export const getFilterTypes = (store?: RootStore) => {
  const filterTypes: Partial<Record<ColumnViewType, FilterType>> = {
    [ColumnViewType.OrganizationsName]: {
      filterType: 'text',
      filterName: 'Organization name',
      filterAccesor: ColumnViewType.OrganizationsName,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Building07 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.OrganizationsWebsite]: {
      filterType: 'text',
      filterName: 'Website',
      filterAccesor: ColumnViewType.OrganizationsWebsite,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Globe01 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.OrganizationsRelationship]: {
      filterType: 'list',
      filterName: 'Relationship',
      filterAccesor: ColumnViewType.OrganizationsRelationship,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <AlignHorizontalCentre02 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: [
        {
          label: 'Customer',
          id: OrganizationRelationship.Customer,
        },
        {
          label: 'Prospect',
          id: OrganizationRelationship.Prospect,
        },
        {
          label: 'Not a Fit',
          id: OrganizationRelationship.NotAFit,
        },
        {
          label: 'Former Customer',
          id: OrganizationRelationship.FormerCustomer,
        },
      ],
    },
    [ColumnViewType.OrganizationsRenewalLikelihood]: {
      filterType: 'list',
      filterName: 'Health',
      filterAccesor: ColumnViewType.OrganizationsRenewalLikelihood,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Activity className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: [
        { id: OpportunityRenewalLikelihood.HighRenewal, label: 'High' },
        { id: OpportunityRenewalLikelihood.MediumRenewal, label: 'Medium' },
        { id: OpportunityRenewalLikelihood.LowRenewal, label: 'Low' },
        { id: OpportunityRenewalLikelihood.ZeroRenewal, label: 'Zero' },
      ],
    },
    [ColumnViewType.OrganizationsOnboardingStatus]: {
      filterType: 'list',
      filterName: 'Onboarding status',
      filterAccesor: ColumnViewType.OrganizationsOnboardingStatus,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Trophy01 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: [
        {
          id: OnboardingStatus.Done,
          label: 'Done',
        },
        {
          id: OnboardingStatus.Stuck,
          label: 'Stuck',
        },
        {
          id: OnboardingStatus.Late,
          label: 'Late',
        },
        {
          id: OnboardingStatus.OnTrack,
          label: 'On track',
        },
        {
          id: OnboardingStatus.Successful,
          label: 'Successful',
        },
        {
          id: OnboardingStatus.NotStarted,
          label: 'Not started',
        },
        {
          id: OnboardingStatus.NotApplicable,
          label: 'Not applicable',
        },
      ],
    },
    [ColumnViewType.OrganizationsRenewalDate]: {
      filterType: 'date',
      filterName: 'Renewal date',
      filterAccesor: ColumnViewType.OrganizationsRenewalDate,
      filterOperators: [ComparisonOperator.Gt, ComparisonOperator.Lt],
      icon: (
        <Calendar className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.OrganizationsForecastArr]: {
      filterType: 'number',
      filterName: 'ARR forecast',
      filterAccesor: ColumnViewType.OrganizationsForecastArr,
      filterOperators: [
        ComparisonOperator.Gt,
        ComparisonOperator.Lt,
        ComparisonOperator.Eq,
        ComparisonOperator.NotEqual,
      ],
      icon: (
        <Calculator className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.OrganizationsOwner]: {
      filterType: 'list',
      filterName: 'Owner',
      filterAccesor: ColumnViewType.OrganizationsOwner,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: <Key01 />,
      options: store?.users.toArray().map((user) => ({
        id: user?.id,
        label: user?.name,
        avatar: user?.value?.profilePhotoUrl,
      })),
    },
    [ColumnViewType.OrganizationsLeadSource]: {
      filterType: 'list',
      filterName: 'Source',
      filterAccesor: ColumnViewType.OrganizationsLeadSource,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <ArrowCircleDownRight className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: uniqBy(store?.organizations.toArray(), 'value.leadSource')
        .map((v) => v.value.leadSource)
        .filter(Boolean)
        .map((leadSource) => ({
          id: leadSource,
          label: leadSource,
        })),
    },
    [ColumnViewType.OrganizationsCreatedDate]: {
      filterType: 'date',
      filterName: 'Created date',
      filterAccesor: ColumnViewType.OrganizationsCreatedDate,
      filterOperators: [ComparisonOperator.Lt, ComparisonOperator.Gt],
      icon: (
        <Calendar className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.OrganizationsYearFounded]: {
      filterType: 'number',
      filterName: 'Founded',
      filterAccesor: ColumnViewType.OrganizationsYearFounded,
      filterOperators: [
        ComparisonOperator.Gt,
        ComparisonOperator.Lt,
        ComparisonOperator.Eq,
        ComparisonOperator.NotEqual,
      ],
      icon: <Cake className='group-hover:text-gray-700 text-gray-500 mb-0.5' />,
    },
    [ColumnViewType.OrganizationsEmployeeCount]: {
      filterType: 'number',
      filterName: 'Employees',
      filterAccesor: ColumnViewType.OrganizationsEmployeeCount,
      filterOperators: [
        ComparisonOperator.Gt,
        ComparisonOperator.Lt,
        ComparisonOperator.Eq,
        ComparisonOperator.NotEqual,
      ],
      icon: (
        <Users03 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.OrganizationsSocials]: {
      filterType: 'text',
      filterName: 'LinkedIn URL',
      filterAccesor: ColumnViewType.OrganizationsSocials,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <LinkedinOutline className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.OrganizationsLastTouchpoint]: {
      filterType: 'list',
      filterName: 'Last touchpoint',
      filterAccesor: ColumnViewType.OrganizationsLastTouchpoint,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Calendar className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: [
        {
          id: LastTouchpointType.InteractionEventEmailSent,
          label: 'Email sent',
        },
        { id: LastTouchpointType.IssueCreated, label: 'Issue created' },
        { id: LastTouchpointType.IssueUpdated, label: 'Issue updated' },
        { id: LastTouchpointType.LogEntry, label: 'Log entry' },
        { id: LastTouchpointType.Meeting, label: 'Meeting' },
        {
          id: LastTouchpointType.InteractionEventChat,
          label: 'Message received',
        },
        {
          id: LastTouchpointType.ActionCreated,
          label: 'Organization created',
        },
      ],
    },
    [ColumnViewType.OrganizationsChurnDate]: {
      filterType: 'date',
      filterName: 'Churn date',
      filterAccesor: ColumnViewType.OrganizationsChurnDate,
      filterOperators: [ComparisonOperator.Lt, ComparisonOperator.Gt],
      icon: (
        <Calendar className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.OrganizationsLastTouchpointDate]: {
      filterType: 'date',
      filterName: 'Last interacted',
      filterAccesor: ColumnViewType.OrganizationsLastTouchpointDate,
      filterOperators: [ComparisonOperator.Lt, ComparisonOperator.Gt],
      icon: (
        <Calendar className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.OrganizationsLtv]: {
      filterType: 'number',
      filterName: 'LTV',
      filterAccesor: ColumnViewType.OrganizationsLtv,
      filterOperators: [
        ComparisonOperator.Gt,
        ComparisonOperator.Lt,
        ComparisonOperator.Eq,
        ComparisonOperator.NotEqual,
      ],
      icon: (
        <CurrencyDollarCircle className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.OrganizationsIndustry]: {
      filterType: 'list',
      filterName: 'Industry',
      filterAccesor: ColumnViewType.OrganizationsIndustry,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: <Building05 />,
      options: uniqBy(store?.organizations.toArray(), 'value.industry')
        .map((v) => v.value.industry)
        .filter(Boolean)
        .sort((a, b) => (a && b ? a?.localeCompare(b) : -1))
        .map((industry) => ({
          id: industry,
          label: industry,
        })),
    },
    [ColumnViewType.OrganizationsContactCount]: {
      filterType: 'number',
      filterName: 'Contact count',
      filterAccesor: ColumnViewType.OrganizationsContactCount,
      filterOperators: [
        ComparisonOperator.Gt,
        ComparisonOperator.Lt,
        ComparisonOperator.Eq,
        ComparisonOperator.NotEqual,
      ],
      icon: (
        <Hash02 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.OrganizationsTags]: {
      filterType: 'list',
      filterName: 'Tags',
      filterAccesor: ColumnViewType.OrganizationsTags,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Tag01 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: store?.tags.toArray().map((tag) => ({
        id: tag.value.metadata.id,
        label: tag.value.name,
      })),
    },

    [ColumnViewType.OrganizationsHeadquarters]: {
      filterType: 'list',
      filterName: 'Country',
      filterAccesor: ColumnViewType.OrganizationsHeadquarters,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Globe04 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: uniqBy(
        store?.organizations.toArray().map((org) => ({
          id: org.value.locations?.[0]?.countryCodeA2,
          label: org.value.locations?.[0]?.country,
        })),
        'id',
      ),
    },

    [ColumnViewType.OrganizationsIsPublic]: {
      filterType: 'list',
      filterName: 'Ownership type',
      filterAccesor: ColumnViewType.OrganizationsIsPublic,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Key01 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: [
        { id: 'Public', label: 'Public' },
        { id: 'Private', label: 'Private' },
      ],
    },
    [ColumnViewType.OrganizationsStage]: {
      filterType: 'list',
      filterName: 'Stage',
      filterAccesor: ColumnViewType.OrganizationsStage,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Columns03 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: [
        {
          label: 'Lead',
          id: OrganizationStage.Lead,
        },
        {
          label: 'Target',
          id: OrganizationStage.Target,
        },
        {
          label: 'Engaged',
          id: OrganizationStage.Engaged,
        },
        {
          label: 'Trial',
          id: OrganizationStage.Trial,
        },
        {
          label: 'Unqualified',
          id: OrganizationStage.Unqualified,
        },
      ],
    },
    [ColumnViewType.OrganizationsParentOrganization]: {
      filterType: 'text',
      filterName: 'Parent org',
      filterAccesor: ColumnViewType.OrganizationsParentOrganization,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Globe01 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
  };

  return filterTypes;
};
