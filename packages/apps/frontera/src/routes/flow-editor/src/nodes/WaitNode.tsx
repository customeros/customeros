import { useEffect } from 'react';

import { NodeProps, useNodesData, useReactFlow } from '@xyflow/react';

import { cn } from '@ui/utils/cn';
import { ResizableInput } from '@ui/form/Input';
import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { Hourglass02 } from '@ui/media/icons/Hourglass02';

import { Handle } from '../components';

export const WaitNode = ({
  id,
  data,
}: NodeProps & { data: Record<string, string | number | boolean> }) => {
  const { setNodes, getNode } = useReactFlow();
  const nodeData = useNodesData(id);

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
        n.id === id ? { ...n, data: { ...n.data, isEditing: true } } : n,
      ),
    );
  };

  return (
    <div className='w-[150px] bg-white border border-grayModern-300 p-3 rounded-lg group cursor-pointer'>
      <div className='truncate text-sm flex items-center justify-between'>
        <div className='flex items-center'>
          <div className='size-6 mr-2 bg-gray-50 border border-gray-100 rounded flex items-center justify-center'>
            <Hourglass02 className='text-gray-500' />
          </div>

          {isEditing ? (
            <div className='flex mr-1 items-baseline'>
              <ResizableInput
                min={0}
                size='xs'
                autoFocus
                step={0.5}
                type='number'
                placeholder={'0'}
                variant='unstyled'
                onFocus={(e) => e.target.select()}
                className='min-w-2.5 min-h-0 max-h-4'
                value={data?.waitDuration?.toString() || ''}
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
