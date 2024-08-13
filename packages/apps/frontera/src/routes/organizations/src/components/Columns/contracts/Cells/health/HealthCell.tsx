import { useMemo, useState } from 'react';

import { set } from 'lodash';
import { observer } from 'mobx-react-lite';
import {
  getLikelihoodColor,
  getRenewalLikelihoodLabel,
} from '@renewals/components/Columns/Cells/renewal/utils.ts';

import { cn } from '@ui/utils/cn';
import { useStore } from '@shared/hooks/useStore';
import { Opportunity, OpportunityRenewalLikelihood } from '@graphql/types';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';

interface RenewalLikelihoodCellProps {
  id: string;
  value?: OpportunityRenewalLikelihood | null;
}

export const HealthCell = observer(({ id }: RenewalLikelihoodCellProps) => {
  const store = useStore();
  const contract = store.contracts.value.get(id);

  const opportunityId = useMemo(() => {
    return (
      contract?.value?.opportunities?.find((e) => e.internalStage === 'OPEN')
        ?.id || contract?.value?.opportunities?.[0]?.id
    );
  }, []);

  const opportunity = store.opportunities.value.get(opportunityId ?? '');
  const value = opportunity?.value?.renewalLikelihood;

  const [isEditing, setIsEditing] = useState(false);
  const colors = value ? getLikelihoodColor(value) : 'text-gray-400';

  const handleClick = (value: OpportunityRenewalLikelihood) => {
    opportunity?.update((opp: Opportunity) => {
      const potentialAmount = opp?.maxAmount ?? 0;

      set(opp, 'renewalLikelihood', value);
      set(
        opp,
        'arrForecast',
        (() => {
          switch (value) {
            case OpportunityRenewalLikelihood.HighRenewal:
              return potentialAmount;
            case OpportunityRenewalLikelihood.MediumRenewal:
              return (50 / 100) * potentialAmount;
            case OpportunityRenewalLikelihood.LowRenewal:
              return (25 / 100) * potentialAmount;
            default:
              return (50 / 100) * potentialAmount;
          }
        })(),
      );

      return opp;
    });
  };

  return (
    <div className='flex gap-1 items-center group/likelihood'>
      <Menu open={isEditing} onOpenChange={setIsEditing}>
        <MenuButton asChild disabled>
          <span
            className={cn('cursor-default', colors)}
            onDoubleClick={() => setIsEditing(true)}
            data-test='oppanization-health-in-all-opps-table'
          >
            {value ? getRenewalLikelihoodLabel(value) : 'No set'}
          </span>
        </MenuButton>
        <MenuList align='center'>
          <MenuItem
            onClick={() =>
              handleClick(OpportunityRenewalLikelihood.HighRenewal)
            }
          >
            High
          </MenuItem>
          <MenuItem
            onClick={() =>
              handleClick(OpportunityRenewalLikelihood.MediumRenewal)
            }
          >
            Medium
          </MenuItem>
          <MenuItem
            onClick={() => handleClick(OpportunityRenewalLikelihood.LowRenewal)}
          >
            Low
          </MenuItem>
        </MenuList>
      </Menu>
      {/*todo*/}
      {/*<IconButton*/}
      {/*  size='xxs'*/}
      {/*  variant='ghost'*/}
      {/*  aria-label='edit renewal likelihood'*/}
      {/*  icon={<Edit03 className='text-gray-500' />}*/}
      {/*  onClick={() => {*/}
      {/*    setIsEditing(true);*/}
      {/*  }}*/}
      {/*  className={cn(*/}
      {/*    'rounded-md opacity-0 group-hover/likelihood:opacity-100',*/}
      {/*    isEditing && 'opacity-100',*/}
      {/*  )}*/}
      {/*/>*/}
    </div>
  );
});
