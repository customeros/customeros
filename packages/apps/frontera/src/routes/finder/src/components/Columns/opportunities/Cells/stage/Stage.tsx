import { RootStore } from '@store/root';
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

export const makeStageLabels = (store: RootStore, preset?: string) => {
  const cacheKey = 'cache' as keyof typeof makeStageLabels;

  if (makeStageLabels?.[cacheKey]) {
    return makeStageLabels[cacheKey] as Record<string, string>;
  } else {
    const tableViewDefOpportunity = store.tableViewDefs.getById(preset ?? '1');

    const columns = tableViewDefOpportunity?.value?.columns;
    const opportunityStages = columns?.reduce((acc, curr) => {
      if (!curr.filter?.includes('STAGE')) return acc;

      return {
        ...acc,
        ['STAGE' + curr.filter.split('STAGE')[1][0]]: curr.name,
      };
    }, {} as Record<string, string>);

    // @ts-expect-error ignore
    makeStageLabels[cacheKey] = opportunityStages as Record<string, string>;

    return opportunityStages;
  }
};

export const StageCell = observer(({ stage, id }: StageCellProps) => {
  const store = useStore();

  const opportunityDef = store.tableViewDefs.opportunitiesPreset;

  const internalStage = store.opportunities.value.get(id)?.value?.internalStage;

  const isClosed = internalStage
    ? [InternalStage.ClosedWon, InternalStage.ClosedLost].includes(
        internalStage,
      )
    : false;

  if (internalStage && isClosed) {
    return <p>{internalStageOptions[internalStage]}</p>;
  }

  const stageLabel = makeStageLabels(store, opportunityDef);
  const stageObj = stageLabel?.[stage];

  if (!stageObj) {
    return <p className='text-gray-400'>Unknown</p>;
  }

  return <p>{stageObj}</p>;
});
