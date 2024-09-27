import { observer } from 'mobx-react-lite';
import { FlowActionType } from '@store/Flows/types.ts';
import { MarkerType, useReactFlow, FitViewOptions } from '@xyflow/react';

import { useStore } from '@shared/hooks/useStore';
import { Mail01 } from '@ui/media/icons/Mail01.tsx';
import { Star06 } from '@ui/media/icons/Star06.tsx';
import { CommandItem } from '@ui/overlay/CommandMenu';
import { MailReply } from '@ui/media/icons/MailReply.tsx';
import { PlusSquare } from '@ui/media/icons/PlusSquare.tsx';
import { Hourglass02 } from '@ui/media/icons/Hourglass02.tsx';
import { RefreshCw01 } from '@ui/media/icons/RefreshCw01.tsx';
import { ArrowIfPath } from '@ui/media/icons/ArrowIfPath.tsx';
import { ClipboardCheck } from '@ui/media/icons/ClipboardCheck.tsx';
import { LinkedinOutline } from '@ui/media/icons/LinkedinOutline.tsx';

import { useLayoutedElements } from '../../hooks';

const elkOptions = {
  'elk.algorithm': 'layered',
  'elk.layered.spacing.nodeNodeBetweenLayers': '100',
  'elk.spacing.nodeNode': '80',
  'elk.direction': 'DOWN',
};
export const StepsHub = observer(() => {
  const { ui } = useStore();
  const { getLayoutedElements } = useLayoutedElements();

  const { setEdges, setNodes, getNodes, getEdges, fitView } = useReactFlow();

  const handleAddNode = async (type: FlowActionType | 'WAIT') => {
    const nodes = getNodes();
    const edges = getEdges();

    const sourceNode = nodes.find(
      (node) => node.id === ui.flowCommandMenu.context.meta?.source,
    );
    const targetNode = nodes.find(
      (node) => node.id === ui.flowCommandMenu.context.meta?.target,
    );

    if (!sourceNode || !targetNode) return;

    const typeBasedContent =
      type === 'WAIT' ? { subject: '', bodyTemplate: '' } : { waitDuration: 1 };

    // Create the new node
    const newNode = {
      id: `${type}-${nodes.length + 1}`,
      type: type === 'WAIT' ? 'wait' : 'action',
      position: { x: sourceNode.position.x, y: 0 }, // Initial position will be adjusted by ELK
      data: {
        action: type,
        ...typeBasedContent,
      },
    };

    // Create two new edges
    const edgeToNewNode = {
      id: `e${ui.flowCommandMenu.context.meta?.source}-${newNode.id}`,
      source: ui.flowCommandMenu.context.meta?.source,
      target: newNode.id,
      type: 'baseEdge',
      markerEnd: { type: MarkerType.Arrow },
    };

    const edgeFromNewNode = {
      id: `e${newNode.id}-${ui.flowCommandMenu.context.meta?.target}`,
      source: newNode.id,
      target: ui.flowCommandMenu.context.meta?.target,
      type: 'baseEdge',
      markerEnd: { type: MarkerType.Arrow },
    };

    // Remove the old edge
    const updatedEdges = edges.filter(
      (e) =>
        !(
          e.source === ui.flowCommandMenu.context.meta?.source &&
          e.target === ui.flowCommandMenu.context.meta?.target
        ),
    );

    // Add the new node and edges
    const updatedNodes = [...nodes, newNode];
    const newEdges = [...updatedEdges, edgeToNewNode, edgeFromNewNode];

    getLayoutedElements(updatedNodes, newEdges, elkOptions).then(
      // @ts-expect-error not for poc
      ({ nodes: layoutedNodes, edges: layoutedEdges }) => {
        setNodes(layoutedNodes);
        setEdges(layoutedEdges);
        window.requestAnimationFrame(() => {
          const fitViewOptions: FitViewOptions = {
            padding: 0.1,
            maxZoom: 1,
            minZoom: 1,
            duration: 0,
            nodes: layoutedNodes,
          };

          fitView(fitViewOptions);
        });
      },
    );
  };

  const updateSelectedNode = (type: FlowActionType | 'WAIT') => {
    handleAddNode(type);
    ui.flowCommandMenu.setOpen(false);
    ui.flowCommandMenu.reset();
  };

  return (
    <>
      <CommandItem
        leftAccessory={<Mail01 />}
        keywords={['send', 'email']}
        onSelect={() => {
          updateSelectedNode(FlowActionType.EMAIL_NEW);
        }}
      >
        Send email
      </CommandItem>
      <CommandItem
        leftAccessory={<MailReply />}
        keywords={['reply', 'to', 'previous', 'email']}
        onSelect={() => {
          updateSelectedNode(FlowActionType.EMAIL_REPLY);
        }}
      >
        Reply to previous email
      </CommandItem>
      <CommandItem
        keywords={['wait', 'delay']}
        leftAccessory={<Hourglass02 />}
        onSelect={() => {
          updateSelectedNode('WAIT');
        }}
      >
        Wait
      </CommandItem>
      <CommandItem disabled leftAccessory={<LinkedinOutline />}>
        <span className='text-gray-700'>Send LinkedIn message</span>
        <span className='text-gray-500'>(Coming soon)</span>
      </CommandItem>
      <CommandItem disabled leftAccessory={<PlusSquare />}>
        <span className='text-gray-700'>Create record</span>
        <span className='text-gray-500'>(Coming soon)</span>
      </CommandItem>
      <CommandItem disabled leftAccessory={<RefreshCw01 />}>
        <span className='text-gray-700'>Update record</span>
        <span className='text-gray-500'>(Coming soon)</span>
      </CommandItem>
      <CommandItem disabled leftAccessory={<Star06 />}>
        <span className='text-gray-700'>Enrich record</span>
        <span className='text-gray-500'>(Coming soon)</span>
      </CommandItem>
      <CommandItem disabled leftAccessory={<Star06 />}>
        <span className='text-gray-700'>Verify record property</span>
        <span className='text-gray-500'>(Coming soon)</span>
      </CommandItem>
      <CommandItem disabled leftAccessory={<ArrowIfPath />}>
        <span className='text-gray-700'>Conditions</span>
        <span className='text-gray-500'>(Coming soon)</span>
      </CommandItem>
      <CommandItem disabled leftAccessory={<ClipboardCheck />}>
        <span className='text-gray-700'>Create to-do</span>
        <span className='text-gray-500'>(Coming soon)</span>
      </CommandItem>
    </>
  );
});
