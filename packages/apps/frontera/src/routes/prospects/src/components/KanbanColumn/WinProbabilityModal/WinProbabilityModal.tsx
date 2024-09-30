import { useRef, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { InternalStage } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
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
  ModalHeader,
  ModalOverlay,
  ModalContent,
  ModalCloseButton,
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
        <ModalContent data-test='win-rate-modal'>
          <ModalHeader>
            <p className='text-md font-semibold mb-1'>
              Set win probability for stage {columnName}
            </p>
            <p className='text-sm'>
              Based on historical data or projections, what percentage of
              opportunities in the <b>{columnName}</b> stage is typically won?
            </p>
          </ModalHeader>
          <ModalBody className='flex flex-col gap-2'>
            <ModalCloseButton />
            <div className='flex justify-between w-full'>
              <label htmlFor='range-slider' className='font-medium text-sm'>
                Probability to win
              </label>
              <span className='font-medium text-sm'>
                {stageLikelihoodRate}%
              </span>
            </div>
            <RangeSlider
              min={0}
              step={5}
              max={100}
              id='range-slider'
              className='w-full'
              value={[stageLikelihoodRate]}
              onValueChange={handleSetProbability}
            >
              <RangeSliderTrack
                dataTest='slider-bar'
                className='bg-gray-400 h-0.5'
              >
                <RangeSliderFilledTrack className='bg-gray-500' />
              </RangeSliderTrack>
              <RangeSliderThumb className='ring-1 shadow-md cursor-pointer ring-gray-400' />
            </RangeSlider>
          </ModalBody>
          <ModalFooter className='flex gap-3'>
            <ModalClose className='w-full'>
              <Button className='w-full' onClick={handleReset}>
                Cancel
              </Button>
            </ModalClose>
            <Button
              className='w-full'
              onClick={handleSave}
              colorScheme='primary'
              dataTest='win-rate-confirm'
            >
              Confirm
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    );
  },
);
