import { useState, useEffect } from 'react';

import { NodeProps, useReactFlow } from '@xyflow/react';

import { cn } from '@ui/utils/cn.ts';
import { ResizableInput } from '@ui/form/Input';
import { IconButton } from '@ui/form/IconButton';
import { Edit03 } from '@ui/media/icons/Edit03.tsx';
import { Hourglass02 } from '@ui/media/icons/Hourglass02.tsx';

import { Handle } from '../components';

export const WaitNode = ({
  id,
  data,
}: NodeProps & { data: Record<string, string | number> }) => {
  const [isFocused, setFocused] = useState(false);
  const { setNodes, getNode } = useReactFlow();

  const handleDurationChange = (newValue: string) => {
    setNodes((nds) => {
      const duration = parseFloat(newValue);
      const updatedNodes = nds.map((node) => {
        if (node.id === id) {
          return {
            ...node,
            data: {
              ...node.data,
              waitDuration: duration,
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
            waitBefore: duration,
          },
        };
      }

      return updatedNodes;
    });
  };

  const selected = getNode(id)?.selected;

  useEffect(() => {
    if (isFocused && !selected) {
      setFocused(false);
    }
  }, [selected]);

  return (
    <div
      className={`w-[150px] bg-white border border-grayModern-300 p-3 rounded-lg group cursor-pointer`}
    >
      <div className='truncate text-sm flex items-center justify-between'>
        <div className='flex items-center'>
          <div
            className={`size-6 min-w-6 mr-2 bg-gray-50 text-gray-500 border-gray-100  rounded flex items-center justify-center`}
          >
            <Hourglass02 />
          </div>

          {isFocused ? (
            <div className='flex mr-1 items-baseline'>
              <ResizableInput
                min={0}
                size='xs'
                autoFocus
                step={0.5}
                type='number'
                placeholder={'0'}
                variant='unstyled'
                value={data.waitDuration || ''}
                onFocus={(e) => e.target.select()}
                className=' min-w-2.5 min-h-0 max-h-4'
                onChange={(e) => handleDurationChange(e.target.value)}
              />
              <span className='ml-1'>
                {data.waitDuration === 1 ? 'day' : 'days'}
              </span>
            </div>
          ) : (
            <span className='truncate'>
              {data.waitDuration || 0}{' '}
              {data.waitDuration === 1 ? 'day' : 'days'}
            </span>
          )}
        </div>

        <IconButton
          size='xxs'
          variant='ghost'
          aria-label='Edit'
          icon={<Edit03 />}
          onClick={() => setFocused(true)}
          className={cn(
            'ml-2  opacity-0 group-hover:opacity-100 pointer-events-all',
            {
              'opacity-0 group-hover:opacity-0': isFocused,
            },
          )}
        />
      </div>
      <Handle type='target' />
      <Handle type='source' />
    </div>
  );
};
