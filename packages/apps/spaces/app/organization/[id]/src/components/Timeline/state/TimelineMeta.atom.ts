import { atom, useRecoilState } from 'recoil';

export const TimelineMetaState = atom({
  key: 'TimelineMetaState',
  default: {
    getTimelineVariables: {
      from: '',
      organizationId: '',
      size: 50,
    },
    reminders: {
      recentlyCreatedId: '',
      recentlyUpdatedId: '',
    },
  },
});

export const useTimelineMeta = () => {
  return useRecoilState(TimelineMetaState);
};
