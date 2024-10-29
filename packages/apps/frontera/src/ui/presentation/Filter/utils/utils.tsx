import { match } from 'ts-pattern';
import { EmailVerificationStatus } from '@finder/components/Columns/contacts/Filters/Email/utils';

import { Equal } from '@ui/media/icons/Equal';
import { Cube01 } from '@ui/media/icons/Cube01';
import { EqualNot } from '@ui/media/icons/EqualNot';
import { CheckCircle } from '@ui/media/icons/CheckCircle';
import { CubeOutline } from '@ui/media/icons/CubeOutline';
import { ChevronLeft } from '@ui/media/icons/ChevronLeft';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { CalendarAfter } from '@ui/media/icons/CalendarAfter';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { SpacingWidth01 } from '@ui/media/icons/SpacingWidth01';
import { CalendarBefore } from '@ui/media/icons/CalendarBefore';
import {
  EmailDeliverable,
  ComparisonOperator,
} from '@shared/types/__generated__/graphql.types';

export const handleOperatorName = (
  operator: ComparisonOperator,
  type?: string,
  plural?: boolean,
) => {
  return match(operator)
    .with(ComparisonOperator.Between, () => 'between')
    .with(ComparisonOperator.Contains, () =>
      type === 'list' && !plural
        ? 'is'
        : plural && type === 'list'
        ? 'is any of'
        : 'contains',
    )
    .with(ComparisonOperator.Eq, () => 'equals')
    .with(ComparisonOperator.Gt, () =>
      type === 'date' ? 'after' : 'more than',
    )
    .with(ComparisonOperator.Gte, () => 'greater than or equal to')
    .with(ComparisonOperator.In, () => 'in')
    .with(ComparisonOperator.IsEmpty, () => 'is empty')
    .with(ComparisonOperator.IsNull, () => 'is null')
    .with(ComparisonOperator.Lt, () =>
      type === 'date' ? 'before' : 'less than',
    )
    .with(ComparisonOperator.Lte, () => 'less than or equal to')
    .with(ComparisonOperator.StartsWith, () => 'starts with')
    .with(ComparisonOperator.IsNotEmpty, () => 'is not empty')
    .with(ComparisonOperator.NotContains, () =>
      type === 'list' && !plural
        ? 'is not'
        : plural && type === 'list'
        ? 'is none of'
        : 'does not contain',
    )
    .with(ComparisonOperator.NotEqual, () => 'not equal to')
    .otherwise(() => 'unknown');
};

export const handleOperatorIcon = (
  operator: ComparisonOperator,
  type?: string,
) => {
  return match(operator)
    .with(ComparisonOperator.Between, () => (
      <SpacingWidth01 className='text-gray-500 group-hover:text-gray-700' />
    ))
    .with(ComparisonOperator.Contains, () => (
      <CheckCircle className='text-gray-500 group-hover:text-gray-700' />
    ))
    .with(ComparisonOperator.Eq, () => (
      <Equal className='text-gray-500 group-hover:text-gray-700' />
    ))
    .with(ComparisonOperator.Gt, () =>
      type === 'date' ? (
        <CalendarAfter className='text-gray-500 group-hover:text-gray-700' />
      ) : (
        <ChevronRight className='text-gray-500 group-hover:text-gray-700' />
      ),
    )
    .with(ComparisonOperator.Gte, () => 'greater than or equal to')
    .with(ComparisonOperator.In, () => 'in')
    .with(ComparisonOperator.IsEmpty, () => (
      <CubeOutline className='text-gray-500 group-hover:text-gray-700' />
    ))
    .with(ComparisonOperator.IsNull, () => 'is null')
    .with(ComparisonOperator.Lt, () =>
      type === 'date' ? (
        <CalendarBefore className='text-gray-500 group-hover:text-gray-700' />
      ) : (
        <ChevronLeft className='text-gray-500 group-hover:text-gray-700' />
      ),
    )
    .with(ComparisonOperator.Lte, () => 'less than or equal to')
    .with(ComparisonOperator.StartsWith, () => 'starts with')
    .with(ComparisonOperator.IsNotEmpty, () => (
      <Cube01 className='text-gray-500 group-hover:text-gray-700' />
    ))
    .with(ComparisonOperator.NotContains, () => (
      <SlashCircle01 className='text-gray-500 group-hover:text-gray-700' />
    ))
    .with(ComparisonOperator.NotEqual, () => (
      <EqualNot className='text-gray-500 group-hover:text-gray-700' />
    ))
    .otherwise(() => 'unknown');
};

export const handlePropertyPlural = (property: string, selection: string[]) => {
  return match(property)
    .with('Health', () => (selection.length === 1 ? 'health' : 'health'))
    .with('Relationship', () =>
      selection.length === 1 ? 'relationship' : 'relationships',
    )
    .with('Owner', () => (selection.length === 1 ? 'owner' : 'owners'))
    .with('Industry', () =>
      selection.length === 1 ? 'industry' : 'industries',
    )
    .with('Tags', () => (selection.length === 1 ? 'tag' : 'tags'))
    .with('Stage', () => (selection.length === 1 ? 'stage' : 'stages'))
    .with('Source', () => (selection.length === 1 ? 'source' : 'sources'))
    .with('Onboarding status', () =>
      selection.length === 1 ? 'onboarding status' : 'onboarding statuses',
    )
    .with('Ownership type', () =>
      selection.length === 1 ? 'ownership type' : 'ownership types',
    )
    .with('Last touchpoint', () =>
      selection.length === 1 ? 'last touchpoint' : 'last touchpoints',
    )
    .with('Country', () => (selection.length === 1 ? 'country' : 'countries'))

    .with('Primary email status', () =>
      selection.length === 1 ? 'email status' : 'email statuses',
    )
    .with('Flow status', () =>
      selection.length === 1 ? 'flow status' : 'flow statuses',
    )
    .with('Region', () => (selection.length === 1 ? 'region' : 'regions'))
    .with('City', () => (selection.length === 1 ? 'city' : 'cities'))
    .with('LinkedIn Connections', () =>
      selection.length === 1 ? 'connection' : 'connections',
    )
    .with('Persona', () => (selection.length === 1 ? 'persona' : 'personas'))
    .with('Current flow', () => (selection.length === 1 ? 'flow' : 'flows'))
    .with('Period', () => (selection.length === 1 ? 'period' : 'period'))
    .with('Currency', () =>
      selection.length === 1 ? 'currency' : 'currencies',
    )
    .with('Status', () => (selection.length === 1 ? 'status' : 'statuses'))
    .with('Renewal', () =>
      selection.length === 1 ? 'renewal' : 'renewal types',
    )
    .with('Biling cycle', () =>
      selection.length === 1 ? 'billing cycle' : 'billing cycles',
    )
    .with('Invoice status', () =>
      selection.length === 1 ? 'invoice status' : 'invoice statuses',
    )
    .otherwise(() => 'unknown');
};

export const categorySelected = (id: string) => {
  return match(id)
    .with(
      EmailVerificationStatus.FirewallProtected,
      () => EmailDeliverable.Deliverable,
    )
    .with(
      EmailVerificationStatus.FreeAccount,
      () => EmailDeliverable.Deliverable,
    )
    .with(EmailVerificationStatus.NoRisk, () => EmailDeliverable.Deliverable)
    .with(
      EmailVerificationStatus.IncorrectFormat,
      () => EmailDeliverable.Undeliverable,
    )
    .with(
      EmailVerificationStatus.InvalidMailbox,
      () => EmailDeliverable.Undeliverable,
    )
    .with(
      EmailVerificationStatus.MailboxFull,
      () => EmailDeliverable.Undeliverable,
    )
    .with(EmailVerificationStatus.CatchAll, () => EmailDeliverable.Unknown)
    .with(EmailVerificationStatus.NotVerified, () => EmailDeliverable.Unknown)
    .with(
      EmailVerificationStatus.VerificationInProgress,
      () => EmailDeliverable.Unknown,
    )
    .otherwise(() => EmailDeliverable.Unknown);
};

export const emailVerificationLabels = (id: string) => {
  match(id)
    .with('no_risk', () => 'No risk')
    .with('firewall_protected', () => 'Firewall protected')
    .with('free_account', () => 'Free acount')
    .with('invalid_mailbox', () => "Mailbox dosen't exist")
    .with('mailbox_full', () => 'Mailbox full')
    .with('incorrect_format', () => 'Incorrect email format')
    .with('catch_all', () => 'Catach all')
    .with('not_verified', () => 'Not verified yet')
    .with('verifrication_in_progress', () => 'Verification in progress')
    .otherwise(() => '');
};
