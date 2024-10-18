import { useEffect } from 'react';

import { NodeProps, useNodesData, useReactFlow } from '@xyflow/react';

import { cn } from '@ui/utils/cn';
import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { Hourglass02 } from '@ui/media/icons/Hourglass02';
import { MaskedResizableInput } from '@ui/form/Input/MaskedResizableInput';

import { Handle } from '../components';

const MINUTES_PER_DAY = 1440;

export const WaitNode = ({
  id,
  data,
}: NodeProps & { data: Record<string, string | number | boolean> }) => {
  const { setNodes, getNode } = useReactFlow();
  const nodeData = useNodesData(id);

  const handleDurationChange = (newValue: string) => {
    setNodes((nds) => {
      const durationInDays = parseFloat(newValue);
      const durationInMinutes = Math.round(durationInDays * MINUTES_PER_DAY);

      const updatedNodes = nds.map((node) => {
        if (node.id === id) {
          return {
            ...node,
            data: {
              ...node.data,
              waitDuration: durationInMinutes,
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
            waitBefore: durationInMinutes,
          },
        };
      }

      return updatedNodes;
    });
  };

  const selected = getNode(id)?.selected;
  const isEditing = nodeData?.data?.isEditing;

  useEffect(() => {
    if (isEditing && !selected) {
      setNodes((nds) =>
        nds.map((n) =>
          n.id === id ? { ...n, data: { ...n.data, isEditing: false } } : n,
        ),
      );
    }
  }, [selected, id, setNodes, isEditing]);

  const toggleEditing = () => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === id
          ? { ...n, selected: true, data: { ...n.data, isEditing: true } }
          : n,
      ),
    );
  };

  const durationInDays = data.waitDuration
    ? (data.waitDuration as number) / MINUTES_PER_DAY
    : 0;

  const displayDuration =
    durationInDays === 0
      ? '0'
      : new Intl.NumberFormat('en-US', {
          minimumFractionDigits: 0,
          maximumFractionDigits: 3,
        }).format(durationInDays);

  const isDaySingular = durationInDays === 1;

  return (
    <div className='w-[156px] h-[56px] bg-white border border-grayModern-300 p-4 rounded-lg group cursor-pointer flex items-center'>
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
                value={displayDuration ?? ''}
                onFocus={(e) => e.target.select()}
                className='min-w-2.5  min-h-0 max-h-4'
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
              <span className='ml-1'>{isDaySingular ? 'day' : 'days'}</span>
            </div>
          ) : (
            <span className='truncate'>
              {displayDuration} {isDaySingular ? 'day' : 'days'}
            </span>
          )}
        </div>

        <IconButton
          size='xxs'
          variant='ghost'
          aria-label='Edit'
          icon={<Edit03 />}
          onClick={toggleEditing}
          className={cn(
            'ml-2 opacity-0 group-hover:opacity-100 pointer-events-all',
            {
              'opacity-0 group-hover:opacity-0': isEditing,
            },
          )}
        />
      </div>
      <Handle type='target' />
      <Handle type='source' />
    </div>
  );
};
