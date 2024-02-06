import setHours from 'date-fns/setHours';

import { DateTimeUtils } from '@spaces/utils/date';
import {
  OnboardingPlanStatus,
  OnboardingPlanMilestoneStatus,
} from '@graphql/types';

import { PlanDatum, MilestoneDatum } from '../types';

export function getMilestoneDueDate(value: string, isDone?: boolean) {
  if (!value) return '';

  const dueDate = setHours(new Date(value), 0).toISOString();
  const now = setHours(new Date(), 0).toISOString();

  const days = DateTimeUtils.differenceInDays(dueDate, now);
  const displayText = (() => {
    if (isDone) {
      return days !== 0
        ? `Done on ${DateTimeUtils.format(value, DateTimeUtils.date)}`
        : 'Done today';
    }

    const prefix = days < 0 ? 'Late by' : days === 0 ? 'Due' : 'Due in';
    const absDays = Math.abs(days);
    const suffix = absDays === 1 ? 'day' : 'days';

    if (absDays === 0) return 'Due today';

    return `${prefix} ${absDays} ${suffix}`;
  })();

  return displayText;
}

export function checkMilestoneDue(milestone: MilestoneDatum) {
  return (
    [
      OnboardingPlanMilestoneStatus.Started,
      OnboardingPlanMilestoneStatus.StartedLate,
      OnboardingPlanMilestoneStatus.NotStarted,
      OnboardingPlanMilestoneStatus.NotStartedLate,
    ].includes(milestone.statusDetails?.status) && milestone.retired === false
  );
}

export function checkPlanDone(plan: PlanDatum) {
  return [OnboardingPlanStatus.Done, OnboardingPlanStatus.DoneLate].includes(
    plan.statusDetails.status,
  );
}
