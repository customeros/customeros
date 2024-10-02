import { match } from 'ts-pattern';

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
import { ComparisonOperator } from '@shared/types/__generated__/graphql.types';

export const handleOperatorName = (
  operator: ComparisonOperator,
  type?: string,
) => {
  return match(operator)
    .with(ComparisonOperator.Between, () => 'between')
    .with(ComparisonOperator.Contains, () => 'contains')
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
    .with(ComparisonOperator.NotContains, () => 'does not contain')
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
