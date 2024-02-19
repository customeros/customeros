import { DateTimeUtils } from '@spaces/utils/date';
import {
  OnboardingPlanStatus,
  OnboardingPlanMilestoneStatus,
} from '@graphql/types';

import { PlanDatum, MilestoneDatum } from '../types';

export function getMilestoneDueDate(
  value: string,
  status: OnboardingPlanMilestoneStatus,
  isDone?: boolean,
) {
  if (!value) return '';
  const isLate = checkMilestoneLate(status);

  const dueDate = DateTimeUtils.toISOMidnight(value);
  const now = DateTimeUtils.toISOMidnight(new Date());

  const days = DateTimeUtils.differenceInDays(dueDate, now);
  const displayText = (() => {
    if (isDone) {
      return days !== 0
        ? `Done ${isLate ? 'late on' : 'on'} ${DateTimeUtils.format(
            value,
            DateTimeUtils.date,
          )}`
        : 'Done today';
    }

    const prefix =
      days < 0 && isLate ? 'Late by' : days === 0 ? 'Due' : 'Due in';
    const absDays = Math.abs(days);
    const suffix = absDays === 1 ? 'day' : 'days';

    if (absDays === 0) return 'Due today';

    return `${prefix} ${absDays} ${suffix}`;
  })();

  return displayText;
}

export function getMilestoneDoneDate(
  value: string,
  status: OnboardingPlanMilestoneStatus,
) {
  if (!value) return '';
  const isLate = checkMilestoneLate(status);

  return `Done ${isLate ? 'late on' : 'on'} ${DateTimeUtils.format(
    value,
    DateTimeUtils.date,
  )}`;
}

function checkMilestoneLate(status?: OnboardingPlanMilestoneStatus) {
  if (!status) return false;

  return [
    OnboardingPlanMilestoneStatus.StartedLate,
    OnboardingPlanMilestoneStatus.NotStartedLate,
    OnboardingPlanMilestoneStatus.DoneLate,
  ].includes(status);
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
