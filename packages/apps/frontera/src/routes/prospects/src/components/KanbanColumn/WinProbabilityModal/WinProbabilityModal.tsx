import { useRef, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { InternalStage } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Percent03 } from '@ui/media/icons/Percent03';
import {
  RangeSlider,
  RangeSliderThumb,
  RangeSliderTrack,
  RangeSliderFilledTrack,
} from '@ui/form/RangeSlider';
import {
  Modal,
  ModalBody,
  ModalClose,
  ModalFooter,
  ModalOverlay,
  ModalCloseButton,
  ModalFeaturedHeader,
  ModalFeaturedContent,
} from '@ui/overlay/Modal/Modal';

interface WinProbabilityModalProps {
  open: boolean;
  columnName: string;
  onToggle: () => void;
  onUpdateProbability: () => void;
  stage: string | InternalStage.ClosedLost | InternalStage.ClosedWon;
}

export const WinProbabilityModal = observer(
  ({
    open,
    stage,
    onToggle,
    columnName,
    onUpdateProbability,
  }: WinProbabilityModalProps) => {
    const [initialValue, setInitialValue] = useState<number>(0);
    const store = useStore();
    const hasSetInitialValue = useRef(false);

    const stageLikelihoodRate =
      store.settings.tenant.value?.opportunityStages.find(
        (s) => s.value === stage,
      )?.likelihoodRate ?? 0;

    const handleSetProbability = (values: number[]) => {
      store.settings.tenant.update(
        (value) => {
          const stageIndex = value.opportunityStages.findIndex(
            (s) => s.value === stage,
          );
          value.opportunityStages[stageIndex].likelihoodRate = values[0];

          return value;
        },
        { mutate: false },
      );
    };

    const handleReset = () => {
      handleSetProbability([initialValue]);
    };

    const handleSave = () => {
      store.settings.tenant.saveOpportunityStage(stage);
      onUpdateProbability();

      onToggle();
    };

    useEffect(() => {
      if (!open) {
        setInitialValue(0);

        hasSetInitialValue.current = false;

        return;
      }

      if (hasSetInitialValue.current) return;
      setInitialValue(stageLikelihoodRate);
      hasSetInitialValue.current = true;
    }, [open]);

    return (
      <Modal open={open} onOpenChange={onToggle}>
        <ModalOverlay />
        <ModalFeaturedContent>
          <ModalFeaturedHeader featuredIcon={<Percent03 />}>
            <p className='text-lg font-semibold mb-1'>
              Set win probability for stage, {columnName}
            </p>
            <p className='text-sm'>
              Based on historical data or projections, what percentage of
              opportunities in the <b>{columnName}</b> stage is typically Won?
            </p>
          </ModalFeaturedHeader>
          <ModalBody className='flex flex-col gap-4'>
            <ModalCloseButton />
            <div className='flex justify-between w-full'>
              <label htmlFor='range-slider' className='font-semibold'>
                Probability to win
              </label>
              <span className='font-semibold'>{stageLikelihoodRate}%</span>
            </div>
            <RangeSlider
              min={0}
              step={1}
              max={100}
              id='range-slider'
              className='w-full'
              value={[stageLikelihoodRate]}
              onValueChange={handleSetProbability}
            >
              <RangeSliderTrack className='bg-gray-400 h-0.5'>
                <RangeSliderFilledTrack className='bg-gray-500' />
              </RangeSliderTrack>
              <RangeSliderThumb className='ring-1 shadow-md cursor-pointer ring-gray-400' />
            </RangeSlider>
          </ModalBody>
          <ModalFooter className='flex gap-3'>
            <ModalClose className='w-full'>
              <Button className='w-full' onClick={handleReset}>
                Close
              </Button>
            </ModalClose>
            <Button
              className='w-full'
              colorScheme='primary'
              onClick={handleSave}
            >
              Set probability
            </Button>
          </ModalFooter>
        </ModalFeaturedContent>
      </Modal>
    );
  },
);
