import { observer } from 'mobx-react-lite';

import { InternalStage } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';

interface StageCellProps {
  id: string;
  stage: string;
}

const internalStageOptions: Record<string, string> = {
  [InternalStage.ClosedLost]: 'Closed Lost',
  [InternalStage.ClosedWon]: 'Closed Won',
};

export const StageCell = observer(({ stage, id }: StageCellProps) => {
  const store = useStore();

  const stages = store.settings.tenant.value?.opportunityStages;
  const internalStage = store.opportunities.value.get(id)?.value?.internalStage;
  const isClosed = internalStage
    ? [InternalStage.ClosedWon, InternalStage.ClosedLost].includes(
        internalStage,
      )
    : false;

  if (internalStage && isClosed) {
    return <p>{internalStageOptions[internalStage]}</p>;
  }

  const stageObj = stages?.find((s) => s.value === stage);

  if (!stageObj) {
    return <p className='text-gray-400'>Unknown</p>;
  }

  return <p>{stageObj?.label}</p>;
});
