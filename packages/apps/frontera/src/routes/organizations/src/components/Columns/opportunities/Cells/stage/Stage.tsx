import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';

interface StageCellProps {
  stage: string;
}

export const StageCell = observer(({ stage }: StageCellProps) => {
  const store = useStore();

  const stages = store.settings.tenant.value?.opportunityStages;

  const stageObj = stages?.find((s) => s.value === stage);

  if (!stageObj) {
    return <p className='text-gray-400'>Unknown</p>;
  }

  return <p>{stageObj?.label}</p>;
});
