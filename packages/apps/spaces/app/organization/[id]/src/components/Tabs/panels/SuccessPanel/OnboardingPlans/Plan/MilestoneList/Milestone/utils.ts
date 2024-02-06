import { DateTimeUtils } from '@spaces/utils/date';
import {
  OnboardingPlanMilestoneStatus,
  OnboardingPlanMilestoneItemStatus,
} from '@graphql/types';

import { MilestoneForm } from './types';

export function checkTaskDone(status: OnboardingPlanMilestoneItemStatus) {
  return [
    OnboardingPlanMilestoneItemStatus.Done,
    OnboardingPlanMilestoneItemStatus.DoneLate,
    OnboardingPlanMilestoneItemStatus.Skipped,
    OnboardingPlanMilestoneItemStatus.SkippedLate,
  ].includes(status);
}

export function checkTaskLate(status: OnboardingPlanMilestoneItemStatus) {
  return [
    OnboardingPlanMilestoneItemStatus.DoneLate,
    OnboardingPlanMilestoneItemStatus.NotDoneLate,
    OnboardingPlanMilestoneItemStatus.SkippedLate,
  ].includes(status);
}

export function checkMilestoneDone(milestone: MilestoneForm) {
  if (
    [
      OnboardingPlanMilestoneStatus.Done,
      OnboardingPlanMilestoneStatus.DoneLate,
    ].includes(milestone?.statusDetails?.status)
  )
    return true;

  if (!milestone?.items?.length) return false;

  return milestone?.items?.every((item) => checkTaskDone(item?.status));
}

export function checkMilestoneLate(milestone: MilestoneForm) {
  if (
    [
      OnboardingPlanMilestoneStatus.DoneLate,
      OnboardingPlanMilestoneStatus.StartedLate,
      OnboardingPlanMilestoneStatus.NotStartedLate,
    ].includes(milestone?.statusDetails?.status)
  )
    return true;

  if (!milestone?.items?.length) return false;

  return milestone?.items?.some((item) => checkTaskLate(item?.status));
}

export function computeMilestoneStatus(milestone: MilestoneForm) {
  const isPastDueDate =
    DateTimeUtils.differenceInDays(
      milestone?.dueDate,
      new Date().toISOString(),
    ) < 0;
  const allTasksDone = milestone?.items?.every((i) => checkTaskDone(i.status));
  const someTasksDone = milestone?.items?.some((i) => checkTaskDone(i.status));

  if (allTasksDone) {
    return isPastDueDate
      ? OnboardingPlanMilestoneStatus.DoneLate
      : OnboardingPlanMilestoneStatus.Done;
  }
  if (!allTasksDone && someTasksDone) {
    return isPastDueDate
      ? OnboardingPlanMilestoneStatus.StartedLate
      : OnboardingPlanMilestoneStatus.Started;
  }
  if (!allTasksDone && !someTasksDone) {
    return isPastDueDate
      ? OnboardingPlanMilestoneStatus.NotStartedLate
      : OnboardingPlanMilestoneStatus.NotStarted;
  }

  return milestone?.statusDetails?.status;
}
