import { useRef, useState, useEffect } from 'react';

import { useKey } from 'rooks';
import { NodeProps, useNodesData, useReactFlow } from '@xyflow/react';

import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { Button } from '@ui/form/Button/Button.tsx';
import { Hourglass02 } from '@ui/media/icons/Hourglass02';
import { MaskedResizableInput } from '@ui/form/Input/MaskedResizableInput';

import { Handle } from '../components';

const MINUTES_PER_DAY = 1440;
const MINUTES_PER_HOUR = 60;

type DurationUnit = 'minutes' | 'hours' | 'days';
const waitDurationOptions: DurationUnit[] = ['minutes', 'hours', 'days'];

const unitDisplayText: Record<
  DurationUnit,
  { plural: string; singular: string }
> = {
  minutes: { singular: 'min', plural: 'min' },
  hours: { singular: 'hour', plural: 'hours' },
  days: { singular: 'day', plural: 'days' },
};

const convertFromMinutes = (minutes: number, unit: DurationUnit): number => {
  switch (unit) {
    case 'days':
      return minutes / MINUTES_PER_DAY;
    case 'hours':
      return minutes / MINUTES_PER_HOUR;
    case 'minutes':
    default:
      return minutes;
  }
};

const convertToMinutes = (value: number, unit: DurationUnit): number => {
  switch (unit) {
    case 'days':
      return value * MINUTES_PER_DAY;
    case 'hours':
      return value * MINUTES_PER_HOUR;
    case 'minutes':
    default:
      return value;
  }
};

const getUnitDisplay = (value: number, unit: DurationUnit): string => {
  return value === 1
    ? unitDisplayText[unit].singular
    : unitDisplayText[unit].plural;
};

// TODO - FE should not be responsible for handling duration unit related logic COS-5474
export const WaitNode = ({
  id,
  data,
}: NodeProps & { data: Record<string, string | number | boolean> }) => {
  const { setNodes, getNode } = useReactFlow();
  const nodeData = useNodesData(id);
  const containerRef = useRef<HTMLDivElement | null>(null);

  const [displayDurationUnit, setDisplayDurationUnit] = useState<DurationUnit>(
    (data.fe_waitDurationUnit as DurationUnit) || 'days',
  );

  const initialMinutes = (data.waitDuration as number) || 0;
  const [durationInMinutes, setDurationInMinutes] =
    useState<number>(initialMinutes);

  const [editingValue, setEditingValue] = useState<string>('');

  const isEditing = nodeData?.data?.isEditing;
  const selected = getNode(id)?.selected;

  const updateNodes = (minutes: number, unit: DurationUnit) => {
    setNodes((nds) => {
      const updatedNodes = nds.map((node) => {
        if (node.id === id) {
          return {
            ...node,
            data: {
              ...node.data,
              waitDuration: minutes,
              fe_waitDurationUnit: unit,
            },
          };
        }

        return node;
      });

      const currentNodeIndex = updatedNodes.findIndex((node) => node.id === id);

      if (currentNodeIndex < updatedNodes.length - 1) {
        const nextNode = updatedNodes[currentNodeIndex + 1];

        updatedNodes[currentNodeIndex + 1] = {
          ...nextNode,
          data: {
            ...nextNode.data,
            waitBefore: minutes,
          },
        };
      }

      return updatedNodes;
    });
  };

  const handleDurationChange = (newValue: string) => {
    setEditingValue(newValue);
  };

  const cycleUnit = (direction: 'up' | 'down') => {
    if (!isEditing) return;

    const currentIndex = waitDurationOptions.indexOf(displayDurationUnit);
    const newUnit =
      direction === 'up'
        ? waitDurationOptions[(currentIndex + 1) % waitDurationOptions.length]
        : waitDurationOptions[
            (currentIndex - 1 + waitDurationOptions.length) %
              waitDurationOptions.length
          ];

    setDisplayDurationUnit(newUnit);
  };

  const toggleEditing = () => {
    const displayValue = convertFromMinutes(
      durationInMinutes,
      displayDurationUnit,
    );

    setEditingValue(displayValue.toString());

    setNodes((nds) =>
      nds.map((n) =>
        n.id === id
          ? { ...n, selected: true, data: { ...n.data, isEditing: true } }
          : n,
      ),
    );
  };

  // Handle exiting edit mode
  useEffect(() => {
    if (isEditing && !selected) {
      // Convert the current editing value to minutes based on the final unit
      const parsedValue = parseFloat(editingValue) || 0;
      const finalMinutes = Math.round(
        convertToMinutes(parsedValue, displayDurationUnit),
      );

      setDurationInMinutes(finalMinutes);
      updateNodes(finalMinutes, displayDurationUnit);

      setNodes((nds) =>
        nds.map((n) =>
          n.id === id ? { ...n, data: { ...n.data, isEditing: false } } : n,
        ),
      );
    }
  }, [selected, id, setNodes, isEditing, displayDurationUnit, editingValue]);

  useKey(
    ['ArrowUp'],
    (e) => {
      e.preventDefault();
      e.stopPropagation();
      cycleUnit('up');
    },
    { target: containerRef },
  );

  useKey(
    ['ArrowDown'],
    (e) => {
      e.preventDefault();
      cycleUnit('down');
    },
    { target: containerRef },
  );

  // Format display value
  const displayValue = isEditing
    ? editingValue
    : durationInMinutes === 0
    ? '0'
    : new Intl.NumberFormat('en-US', {
        minimumFractionDigits: 0,
        maximumFractionDigits: displayDurationUnit === 'minutes' ? 0 : 2,
      }).format(convertFromMinutes(durationInMinutes, displayDurationUnit));

  const unitDisplay = getUnitDisplay(
    parseFloat(displayValue) || 0,
    displayDurationUnit,
  );

  return (
    <div
      ref={containerRef}
      className='relative w-[156px] h-[56px] bg-white border border-grayModern-300 p-4 rounded-lg group cursor-pointer flex items-center'
    >
      <div className='truncate text-sm flex items-center justify-between w-full'>
        <div className='flex items-center'>
          <div className='size-6 mr-2 bg-gray-50 border border-gray-100 rounded flex items-center justify-center'>
            <Hourglass02 className='text-gray-500' />
          </div>

          {isEditing ? (
            <div className='flex mr-1 items-baseline'>
              <MaskedResizableInput
                size='xs'
                autoFocus
                mask={`num`}
                unmask={true}
                placeholder={'0'}
                variant='unstyled'
                value={displayValue}
                onFocus={(e) => e.target.select()}
                className='min-w-2.5 min-h-0 max-h-4'
                onAccept={(_val, maskRef) => {
                  const unmaskedValue = maskRef._unmaskedValue;

                  handleDurationChange(unmaskedValue);
                }}
                blocks={{
                  num: {
                    mask: Number,
                    radix: '.',
                    scale: 3,
                    max: 9990,
                    mapToRadix: [','],
                    lazy: false,
                    min: 0,
                    placeholderChar: '#',
                    thousandsSeparator: ',',
                    normalizeZeros: true,
                    padFractionalZeros: false,
                    autofix: true,
                  },
                }}
              />
              <Button
                variant='link'
                onClick={() => cycleUnit('up')}
                className='p-0 ml-1 shadow-none'
              >
                {unitDisplay}
              </Button>
            </div>
          ) : (
            <span className='truncate'>
              {displayValue} {unitDisplay}
            </span>
          )}
        </div>

        <IconButton
          size='xxs'
          variant='ghost'
          aria-label='Edit'
          icon={<Edit03 />}
          onClick={toggleEditing}
          className={`ml-2 opacity-0 group-hover:opacity-100 pointer-events-all ${
            isEditing ? 'opacity-0 group-hover:opacity-0' : ''
          }`}
        />
      </div>
      <Handle type='target' />
      <Handle type='source' />
    </div>
  );
};
