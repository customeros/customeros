import {
  BaseEdge,
  EdgeProps,
  MarkerType,
  getBezierPath,
  EdgeLabelRenderer,
} from '@xyflow/react';

import { Hourglass02 } from '@ui/media/icons/Hourglass02.tsx';

export const TimeTriggerEdge: React.FC<EdgeProps> = ({
  id,
  data,
  ...props
}) => {
  const [edgePath, labelX, labelY] = getBezierPath(props);

  return (
    <>
      <>
        <BaseEdge id={id} path={edgePath} markerEnd={MarkerType.Arrow} />
        <EdgeLabelRenderer>
          <div
            className='nodrag nopan bg-white flex items-center'
            style={{
              position: 'absolute',
              transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
            }}
          >
            <Hourglass02 className='mr-1' />
            {data?.timeValue} {data?.timeUnit}
          </div>
        </EdgeLabelRenderer>
      </>
    </>
  );
};
