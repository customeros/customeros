import { useRef, useState } from 'react';

import {
  BaseEdge,
  EdgeProps,
  MarkerType,
  useReactFlow,
  getBezierPath,
  ViewportPortal,
  EdgeLabelRenderer,
} from '@xyflow/react';

import { X } from '@ui/media/icons/X.tsx';
import { Plus } from '@ui/media/icons/Plus.tsx';
import { IconButton } from '@ui/form/IconButton';
import { Button } from '@ui/form/Button/Button.tsx';

export const BasicEdge: React.FC<EdgeProps> = ({ id, data, ...props }) => {
  const [isOpen, setIsOpen] = useState(false);
  const { setEdges, setNodes, getNodes } = useReactFlow();
  const [edgePath, labelX, labelY] = getBezierPath({
    ...props,
  });

  const nodes = getNodes();

  const toggleOpen = () => {
    setIsOpen(!isOpen);
  };

  const handleDisconnectSteps = () => {
    setEdges((edges) => edges.filter((edge) => edge.id !== id));
  };

  const handleAddNode = (type: 'trigger' | 'step') => {
    const sourceNode = nodes.find((node) => node.id === props.source);
    const targetNode = nodes.find((node) => node.id === props.target);

    if (!sourceNode || !targetNode) return;

    // Calculate the midpoint between source and target nodes
    const midX = (sourceNode.position.x + targetNode.position.x) / 2;
    const midY = (sourceNode.position.y + targetNode.position.y) / 2;

    // Create the new node
    const newNode = {
      id: `step-${nodes.length + 1}`,
      type: 'default', // or any other type you want
      position: { x: midX, y: midY },
      data: { label: `New Step ${nodes.length + 1}` },
    };

    // Create two new edges
    const edgeToNewNode = {
      id: `e${props.source}-${newNode.id}`,
      source: props.source,
      target: newNode.id,
      type: 'baseEdge',
    };

    const edgeFromNewNode = {
      id: `e${newNode.id}-${props.target}`,
      source: newNode.id,
      target: props.target,
      type: 'baseEdge',
    };

    // Update nodes and edges
    setNodes((nds) => nds.concat(newNode));
    setEdges(
      (eds) =>
        eds
          .filter((e) => e.id !== id) // Remove the old edge
          .concat([edgeToNewNode, edgeFromNewNode]), // Add the new edges
    );
  };

  return (
    <>
      <BaseEdge path={edgePath} markerEnd={MarkerType.Arrow} />
      <EdgeLabelRenderer>
        <div
          className='nodrag nopan'
          style={{
            position: 'absolute',
            transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
            fontSize: 12,
            pointerEvents: 'all',
          }}
        >
          <IconButton
            size='xxs'
            onClick={toggleOpen}
            aria-label='Add step or trigger'
            className='text-white bg-gray-700 hover:bg-gray-600 hover:text-white focus:bg-gray-600 focus:text-white rounded-full'
            icon={
              isOpen ? (
                <X className='text-inherit' />
              ) : (
                <Plus className='text-inherit' />
              )
            }
          />
        </div>
      </EdgeLabelRenderer>

      {isOpen && (
        <ViewportPortal>
          <div
            className={`bg-white shadow-md border rounded-lg p-2 flex flex-col justify-start pointer-events-auto`}
            style={{
              transform: `translate(-50%, -50%) translate(${labelX + 60}px,${
                labelY + 30
              }px)`,
              position: 'absolute',
            }}
          >
            <Button
              variant='ghost'
              onClick={() => handleAddNode('step')}
              className='px-1 py-0.5 cursor-pointer justify-start'
            >
              Add step
            </Button>
            <Button
              variant='ghost'
              className='px-1 py-0.5 cursor-pointer justify-start'
            >
              Add trigger
            </Button>

            <Button
              variant='ghost'
              onClick={handleDisconnectSteps}
              className='px-1 py-0.5 cursor-pointer justify-start'
            >
              Disconnect those steps
            </Button>
          </div>
        </ViewportPortal>
      )}
    </>
  );
};
