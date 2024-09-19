import { Handle, Position } from '@xyflow/react';

import { Play } from '@ui/media/icons/Play.tsx';

import '@xyflow/react/dist/style.css';

export const StartNode = ({ data }) => (
  <div className='aspect-[9/1] w-[200px] bg-white border-2 border-green-200 p-3 rounded-lg shadow-md'>
    <div className='flex items-center font-bold'>
      <div className='flex items-center'>
        <Play className='mr-2 w-4 h-4' />
      </div>
      <span className='truncate'>Start Flow</span>
    </div>
    <div className='truncate text-sm mt-1'>Double click to edit trigger</div>
    <Handle type='source' className='h-2 w-2' position={Position.Right} />
  </div>
);

export const BasicNode = ({ type, data }) => (
  <div
    className={`aspect-[9/1] bg-white border-2 border-${data.color}-200 p-3 rounded-lg shadow-md`}
  >
    <div className='font-bold'>Send {type}</div>
    <div>{data.content || 'No content'}</div>

    <Handle
      type='target'
      position={Position.Top}
      className={`h-2 w-2 bg-${data.color}-500`}
    />
    <Handle
      type='source'
      position={Position.Bottom}
      className={`h-2 w-2 bg-${data.color}-500`}
    />
  </div>
);
export const TriggerNode = ({ data }) => (
  <div
    className={`aspect-[9/1] w-[180px] bg-white border-2 border-${data.color}-200 p-3 rounded-lg shadow-md`}
  >
    <div className='flex items-center text-gray-400 uppercase text-xs'>
      Trigger
    </div>
    <div className='truncate  text-sm '>
      {data.content || <span className='text-gray-500'>Webhook</span>}
    </div>
    <Handle
      type='target'
      position={Position.Top}
      className={`h-2 w-2 bg-gray-400`}
    />
    <Handle
      type='source'
      position={Position.Bottom}
      className={`h-2 w-2 bg-gray-400`}
    />
  </div>
);

export const nodeTypes = {
  startNode: StartNode,
  step: BasicNode,
  trigger: TriggerNode,
};
