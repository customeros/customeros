import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Reminders } from './reminders';

export const FutureZone = () => {
  const isRemindersEnabled = useFeatureIsOn('reminders');

  if (!isRemindersEnabled) return null;

  return (
    <>
      <Reminders />
    </>
  );
};
