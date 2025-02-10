import { useRef, useState, useEffect } from 'react';

import { MaskElement } from 'imask';
import { observer } from 'mobx-react-lite';
import { NodeProps, useNodesData, useReactFlow } from '@xyflow/react';

import { Button } from '@ui/form/Button/Button';
import { Hourglass02 } from '@ui/media/icons/Hourglass02';
import { MaskedResizableInput } from '@ui/form/Input/MaskedResizableInput';

import { Handle } from '../components';

type DurationUnit = 'minutes' | 'hours' | 'days';

const UNITS: Record<DurationUnit, { plural: string; singular: string }> = {
  minutes: { singular: 'min', plural: 'min' },
  hours: { singular: 'hour', plural: 'hours' },
  days: { singular: 'day', plural: 'days' },
};

const CONVERSION_RATES = {
  days: 1440, // minutes per day
  hours: 60, // minutes per hour
  minutes: 1,
};

export const WaitNode = observer(
  ({
    id,
    data,
  }: NodeProps & { data: Record<string, string | number | boolean> }) => {
    const { setNodes, getNode } = useReactFlow();
    const nodeData = useNodesData(id);
    const inputRef = useRef<MaskElement>();

    const [unit, setUnit] = useState<DurationUnit>(
      (data.fe_waitDurationUnit as DurationUnit) || 'days',
    );
    const [minutes, setMinutes] = useState<number>(
      (data.waitDuration as number) || 0,
    );
    const [editValue, setEditValue] = useState<string>('');

    const isEditing = nodeData?.data?.isEditing;
    const selected = getNode(id)?.selected;

    const convertDuration = (
      value: number,
      fromUnit: DurationUnit,
      toUnit: DurationUnit,
    ) => (value * CONVERSION_RATES[fromUnit]) / CONVERSION_RATES[toUnit];

    const formatNumber = (value: number) =>
      new Intl.NumberFormat('en-US', {
        minimumFractionDigits: 0,
        maximumFractionDigits: unit === 'minutes' ? 0 : 2,
      }).format(value);

    const cycleUnit = (direction: 'up' | 'down') => {
      if (!isEditing) return;
      const units = Object.keys(UNITS) as DurationUnit[];
      const currentIndex = units.indexOf(unit);
      const newIndex =
        direction === 'up'
          ? (currentIndex + 1) % units.length
          : (currentIndex - 1 + units.length) % units.length;

      setUnit(units[newIndex]);
    };

    const updateNodes = (newMinutes: number, newUnit: DurationUnit) => {
      setNodes((nodes) => {
        const nodeIndex = nodes.findIndex((node) => node.id === id);
        const updatedNodes = nodes.map((node) => {
          if (node.id === id) {
            return {
              ...node,
              data: {
                ...node.data,
                waitDuration: newMinutes,
                fe_waitDurationUnit: newUnit,
              },
            };
          }

          return node;
        });

        if (nodeIndex < nodes.length - 1) {
          updatedNodes[nodeIndex + 1] = {
            ...updatedNodes[nodeIndex + 1],
            data: {
              ...updatedNodes[nodeIndex + 1].data,
              waitBefore: newMinutes,
            },
          };
        }

        return updatedNodes;
      });
    };

    useEffect(() => {
      if (nodeData?.data?.isEditing) {
        const value = formatNumber(convertDuration(minutes, 'minutes', unit));

        setEditValue(value);
        inputRef.current?.select(0, value.length);
      }
    }, [data.waitDuration, nodeData?.data?.isEditing]);

    useEffect(() => {
      if (isEditing && !selected) {
        const newMinutes = Math.round(
          convertDuration(parseFloat(editValue) || 0, unit, 'minutes'),
        );

        setMinutes(newMinutes);
        updateNodes(newMinutes, unit);
        setNodes((nodes) =>
          nodes.map((n) =>
            n.id === id ? { ...n, data: { ...n.data, isEditing: false } } : n,
          ),
        );
      }
    }, [selected, isEditing]);

    const displayValue = isEditing
      ? editValue
      : minutes === 0
      ? '0'
      : formatNumber(convertDuration(minutes, 'minutes', unit));
    const unitDisplay =
      UNITS[unit][parseFloat(displayValue) === 1 ? 'singular' : 'plural'];

    return (
      <div className='relative w-[156px] h-[56px] bg-white border border-grayModern-300 p-4 rounded-lg group cursor-pointer flex items-center'>
        <div className='truncate text-sm flex items-center justify-between w-full'>
          <div className='flex items-center'>
            <div className='size-6 mr-2 bg-gray-50 border border-gray-100 rounded flex items-center justify-center'>
              <Hourglass02 className='text-gray-500' />
            </div>

            {isEditing ? (
              <div className='flex mr-1 items-baseline'>
                <MaskedResizableInput
                  unmask
                  size='xs'
                  autoFocus
                  mask='num'
                  placeholder='0'
                  variant='unstyled'
                  // @ts-expect-error - unmask is not in the types
                  inputRef={inputRef}
                  value={displayValue}
                  onFocus={(e) => e.target.select()}
                  className='min-w-2.5 min-h-0 max-h-4'
                  onAccept={(_val, maskRef) =>
                    setEditValue(maskRef._unmaskedValue)
                  }
                  blocks={{
                    num: {
                      mask: Number,
                      radix: '.',
                      scale: 3,
                      max: 9990,
                      min: 0,
                      mapToRadix: [','],
                      lazy: false,
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
        </div>
        <Handle type='target' />
        <Handle type='source' />
      </div>
    );
  },
);
