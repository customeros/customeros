import { Hash02 } from '@ui/media/icons/Hash02';
import { Calendar } from '@ui/media/icons/Calendar';
import {
  ColumnViewType,
  OnboardingStatus,
  OrganizationStage,
  LastTouchpointType,
  CustomFieldTemplateType,
  OrganizationRelationship,
  OpportunityRenewalLikelihood,
} from '@shared/types/__generated__/graphql.types';

export type fieldType = {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  options?: any[];
  icon: JSX.Element;
  fieldName: string;
  columnAccesor: ColumnViewType;
  fieldType: CustomFieldTemplateType | 'multi-select' | 'date';
  groupOptions?: { label: string; options: { id: string; label: string }[] };
  fieldTypeName: 'Text' | 'Number' | 'Date' | 'Single select' | 'Multi select';
};

import { uniqBy } from 'lodash';
import { type RootStore } from '@store/root';

import { Type01 } from '@ui/media/icons/Type01';
import { RadioButton } from '@ui/media/icons/RadioButton';
import { ListBulleted } from '@ui/media/icons/ListBulleted';

export const getDefaultFieldTypes = (store?: RootStore) => {
  const filterTypes: Partial<Record<ColumnViewType, fieldType>> = {
    [ColumnViewType.OrganizationsName]: {
      fieldType: CustomFieldTemplateType.FreeText,
      fieldTypeName: 'Text',
      fieldName: 'Organization name',
      columnAccesor: ColumnViewType.OrganizationsName,
      icon: <Type01 className='mb-0.5' />,
    },
    [ColumnViewType.OrganizationsWebsite]: {
      fieldType: CustomFieldTemplateType.FreeText,
      fieldTypeName: 'Text',
      fieldName: 'Website',
      columnAccesor: ColumnViewType.OrganizationsWebsite,
      icon: <Type01 className='mb-0.5' />,
    },
    [ColumnViewType.OrganizationsRelationship]: {
      fieldType: CustomFieldTemplateType.SingleSelect,
      fieldTypeName: 'Single select',
      fieldName: 'Relationship',
      columnAccesor: ColumnViewType.OrganizationsRelationship,
      icon: <RadioButton />,
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
      fieldType: CustomFieldTemplateType.SingleSelect,
      fieldTypeName: 'Single select',
      fieldName: 'Health',
      columnAccesor: ColumnViewType.OrganizationsRenewalLikelihood,
      icon: <RadioButton />,
      options: [
        { id: OpportunityRenewalLikelihood.HighRenewal, label: 'High' },
        { id: OpportunityRenewalLikelihood.MediumRenewal, label: 'Medium' },
        { id: OpportunityRenewalLikelihood.LowRenewal, label: 'Low' },
        { id: OpportunityRenewalLikelihood.ZeroRenewal, label: 'Zero' },
      ],
    },
    [ColumnViewType.OrganizationsOnboardingStatus]: {
      fieldType: CustomFieldTemplateType.SingleSelect,
      fieldTypeName: 'Single select',
      fieldName: 'Onboarding status',
      columnAccesor: ColumnViewType.OrganizationsOnboardingStatus,
      icon: <RadioButton />,
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
      fieldType: 'date',
      fieldTypeName: 'Date',
      fieldName: 'Renewal date',
      columnAccesor: ColumnViewType.OrganizationsRenewalDate,
      icon: <Calendar />,
    },
    [ColumnViewType.OrganizationsForecastArr]: {
      fieldType: CustomFieldTemplateType.Number,
      fieldTypeName: 'Number',
      fieldName: 'ARR forecast',
      columnAccesor: ColumnViewType.OrganizationsForecastArr,

      icon: <Hash02 />,
    },
    [ColumnViewType.OrganizationsOwner]: {
      fieldType: CustomFieldTemplateType.SingleSelect,
      fieldTypeName: 'Single select',
      fieldName: 'Owner',
      columnAccesor: ColumnViewType.OrganizationsOwner,
      icon: <RadioButton />,
      options: store?.users.toArray().map((user) => ({
        id: user?.id,
        label: user?.name,
        avatar: user?.value?.profilePhotoUrl,
      })),
    },
    [ColumnViewType.OrganizationsLeadSource]: {
      fieldType: CustomFieldTemplateType.FreeText,
      fieldTypeName: 'Text',
      fieldName: 'Source',
      columnAccesor: ColumnViewType.OrganizationsLeadSource,

      icon: <Type01 />,
    },
    [ColumnViewType.OrganizationsCreatedDate]: {
      fieldType: 'date',
      fieldTypeName: 'Date',
      fieldName: 'Created date',
      columnAccesor: ColumnViewType.OrganizationsCreatedDate,
      icon: <Calendar />,
    },
    [ColumnViewType.OrganizationsYearFounded]: {
      fieldType: CustomFieldTemplateType.Number,
      fieldTypeName: 'Number',
      fieldName: 'Founded',
      columnAccesor: ColumnViewType.OrganizationsYearFounded,
      icon: <Hash02 />,
    },
    [ColumnViewType.OrganizationsEmployeeCount]: {
      fieldType: CustomFieldTemplateType.Number,
      fieldTypeName: 'Number',
      fieldName: 'Employees',
      columnAccesor: ColumnViewType.OrganizationsEmployeeCount,
      icon: <Hash02 />,
    },
    [ColumnViewType.OrganizationsSocials]: {
      fieldType: CustomFieldTemplateType.FreeText,
      fieldTypeName: 'Text',
      fieldName: 'LinkedIn URL',
      columnAccesor: ColumnViewType.OrganizationsSocials,
      icon: <Type01 />,
    },
    [ColumnViewType.OrganizationsLastTouchpoint]: {
      fieldType: CustomFieldTemplateType.SingleSelect,
      fieldTypeName: 'Single select',
      fieldName: 'Last touchpoint',
      columnAccesor: ColumnViewType.OrganizationsLastTouchpoint,
      icon: <RadioButton />,
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
      fieldType: 'date',
      fieldTypeName: 'Date',
      fieldName: 'Churn date',
      columnAccesor: ColumnViewType.OrganizationsChurnDate,
      icon: <Calendar className='mb-0.5' />,
    },
    [ColumnViewType.OrganizationsLastTouchpointDate]: {
      fieldType: 'date',
      fieldTypeName: 'Date',
      fieldName: 'Last interacted',
      columnAccesor: ColumnViewType.OrganizationsLastTouchpointDate,
      icon: <Calendar className='mb-0.5' />,
    },
    [ColumnViewType.OrganizationsLtv]: {
      fieldType: CustomFieldTemplateType.Number,
      fieldTypeName: 'Number',
      fieldName: 'LTV',
      columnAccesor: ColumnViewType.OrganizationsLtv,
      icon: <Hash02 />,
    },
    [ColumnViewType.OrganizationsIndustry]: {
      fieldType: CustomFieldTemplateType.SingleSelect,
      fieldTypeName: 'Single select',
      fieldName: 'Industry',
      columnAccesor: ColumnViewType.OrganizationsIndustry,
      icon: <RadioButton />,
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
      fieldType: CustomFieldTemplateType.Number,
      fieldTypeName: 'Number',
      fieldName: 'Contact count',
      columnAccesor: ColumnViewType.OrganizationsContactCount,
      icon: <Hash02 />,
    },
    [ColumnViewType.OrganizationsTags]: {
      fieldType: 'multi-select',
      fieldTypeName: 'Multi select',
      fieldName: 'Tags',
      columnAccesor: ColumnViewType.OrganizationsTags,
      icon: <ListBulleted />,
      options: store?.tags.toArray().map((tag) => ({
        id: tag.value.metadata.id,
        label: tag.value.name,
      })),
    },
    [ColumnViewType.OrganizationsHeadquarters]: {
      fieldType: CustomFieldTemplateType.SingleSelect,
      fieldTypeName: 'Single select',
      fieldName: 'Country',
      columnAccesor: ColumnViewType.OrganizationsHeadquarters,
      icon: <RadioButton />,
      options: uniqBy(
        store?.organizations.toArray().map((org) => ({
          id: org.value.locations?.[0]?.countryCodeA2,
          label: org.value.locations?.[0]?.country,
        })),
        'id',
      ),
    },
    [ColumnViewType.OrganizationsIsPublic]: {
      fieldType: CustomFieldTemplateType.SingleSelect,
      fieldTypeName: 'Single select',
      fieldName: 'Ownership type',
      columnAccesor: ColumnViewType.OrganizationsIsPublic,
      icon: <RadioButton />,
      options: [
        { id: 'Public', label: 'Public' },
        { id: 'Private', label: 'Private' },
      ],
    },
    [ColumnViewType.OrganizationsStage]: {
      fieldType: CustomFieldTemplateType.SingleSelect,
      fieldTypeName: 'Single select',
      fieldName: 'Stage',
      columnAccesor: ColumnViewType.OrganizationsStage,
      icon: <RadioButton />,
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
      fieldType: CustomFieldTemplateType.FreeText,
      fieldTypeName: 'Text',
      fieldName: 'Parent org',
      columnAccesor: ColumnViewType.OrganizationsParentOrganization,
      icon: <Type01 className='mb-0.5' />,
    },
  };

  return filterTypes;
};

export interface customFieldType {
  [key: string]: {
    icon: JSX.Element;
    fieldTypeName: string;
    fieldType: CustomFieldTemplateType;
  };
}

export const getCustomFieldTypes = () => {
  const customFieldTypes: customFieldType = {
    [CustomFieldTemplateType.FreeText]: {
      fieldType: CustomFieldTemplateType.FreeText,
      fieldTypeName: 'Text',
      icon: <Type01 />,
    },
    [CustomFieldTemplateType.Number]: {
      fieldType: CustomFieldTemplateType.Number,
      fieldTypeName: 'Number',
      icon: <Hash02 />,
    },
    [CustomFieldTemplateType.SingleSelect]: {
      fieldType: CustomFieldTemplateType.SingleSelect,
      fieldTypeName: 'Single select',
      icon: <RadioButton />,
    },
  };

  return customFieldTypes;
};
