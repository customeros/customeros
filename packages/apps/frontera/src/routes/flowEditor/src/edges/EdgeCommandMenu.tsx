import { observer } from 'mobx-react-lite';
import { MarkerType, useReactFlow, ViewportPortal } from '@xyflow/react';

import { useStore } from '@shared/hooks/useStore';
import { Mail01 } from '@ui/media/icons/Mail01.tsx';
import { Star06 } from '@ui/media/icons/Star06.tsx';
import { MailReply } from '@ui/media/icons/MailReply.tsx';
import { PlusSquare } from '@ui/media/icons/PlusSquare.tsx';
import { Hourglass02 } from '@ui/media/icons/Hourglass02.tsx';
import { ArrowIfPath } from '@ui/media/icons/ArrowIfPath.tsx';
import { RefreshCw01 } from '@ui/media/icons/RefreshCw01.tsx';
import { ClipboardCheck } from '@ui/media/icons/ClipboardCheck.tsx';
import { LinkedinOutline } from '@ui/media/icons/LinkedinOutline.tsx';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

import { getLayoutedElements } from '../controls/LayoutButton.tsx';

export const StepViewportPortal = observer(
  ({
    id,
    positionAbsoluteX,
    positionAbsoluteY,
  }: {
    id: string;
    positionAbsoluteX: number;
    positionAbsoluteY: number;
  }) => {
    const { ui } = useStore();

    return (
      <>
        {ui.flowCommandMenu?.isOpen && id === ui.flowCommandMenu.context.id && (
          <ViewportPortal>
            <div
              style={{
                transform: `translate(${positionAbsoluteX / 2}px,${
                  positionAbsoluteY + 15
                }px)`,
                position: 'absolute',
                pointerEvents: 'all',
              }}
            >
              <StepMenu />
            </div>
          </ViewportPortal>
        )}
      </>
    );
  },
);
const elkOptions = {
  'elk.algorithm': 'layered',
  'elk.layered.spacing.nodeNodeBetweenLayers': '100',
  'elk.spacing.nodeNode': '80',
  'elk.direction': 'DOWN',
};
export const StepsHub = observer(() => {
  const { ui } = useStore();

  const { setEdges, setNodes, getNodes, getEdges } = useReactFlow();

  const handleAddNode = async (type: 'SendEmail' | 'ReplyToEmail' | 'Wait') => {
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
      type === 'ReplyToEmail'
        ? { subject: '', content: '' }
        : type === 'SendEmail'
        ? { subject: '', content: '' }
        : { waitDuration: 1 };

    // Create the new node
    const newNode = {
      id: `${type}-${nodes.length + 1}`,
      type: 'action',
      position: { x: 0, y: 0 }, // Initial position will be adjusted by ELK
      data: {
        stepType: type,
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

    // Use getLayoutedElements to calculate new positions
    // @ts-expect-error fix type later
    const { nodes: layoutedNodes, edges: layoutedEdges } =
      await getLayoutedElements(updatedNodes, newEdges, elkOptions);

    // Update the React Flow state
    setNodes(layoutedNodes);
    setEdges(layoutedEdges);
  };

  const updateSelectedNode = (type: 'SendEmail' | 'ReplyToEmail' | 'Wait') => {
    handleAddNode(type);
    ui.flowCommandMenu.setOpen(false);
    ui.flowCommandMenu.reset();
  };

  return (
    <Command>
      <CommandInput
        autoFocus
        className='p-1 text-sm'
        placeholder='Search trigger'
        inputWrapperClassName='min-h-4'
        wrapperClassName='py-2 px-4 mt-2'
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
      />

      <Command.List>
        <CommandItem
          leftAccessory={<Mail01 />}
          onSelect={() => {
            updateSelectedNode('SendEmail');
          }}
        >
          Send email
        </CommandItem>{' '}
        <CommandItem
          leftAccessory={<MailReply />}
          onSelect={() => {
            updateSelectedNode('ReplyToEmail');
          }}
        >
          Reply to previous email
        </CommandItem>{' '}
        <CommandItem
          leftAccessory={<Hourglass02 />}
          onSelect={() => {
            updateSelectedNode('Wait');
          }}
        >
          Wait
        </CommandItem>
        <CommandItem disabled leftAccessory={<LinkedinOutline />}>
          <span className='text-gray-700'>Send connection request</span>{' '}
          <span className='text-gray-500'>(Coming soon)</span>
        </CommandItem>{' '}
        <CommandItem disabled leftAccessory={<LinkedinOutline />}>
          <span className='text-gray-700'>Send LinkedIn message</span>{' '}
          <span className='text-gray-500'>(Coming soon)</span>
        </CommandItem>{' '}
        <CommandItem disabled leftAccessory={<PlusSquare />}>
          <span className='text-gray-700'>Create record</span>{' '}
          <span className='text-gray-500'>(Coming soon)</span>
        </CommandItem>{' '}
        <CommandItem disabled leftAccessory={<RefreshCw01 />}>
          <span className='text-gray-700'>Update record</span>{' '}
          <span className='text-gray-500'>(Coming soon)</span>
        </CommandItem>{' '}
        <CommandItem disabled leftAccessory={<Star06 />}>
          <span className='text-gray-700'>Enrich record</span>{' '}
          <span className='text-gray-500'>(Coming soon)</span>
        </CommandItem>{' '}
        <CommandItem disabled leftAccessory={<Star06 />}>
          <span className='text-gray-700'>Verify record property</span>{' '}
          <span className='text-gray-500'>(Coming soon)</span>
        </CommandItem>{' '}
        <CommandItem disabled leftAccessory={<ArrowIfPath />}>
          <span className='text-gray-700'>Conditions</span>{' '}
          <span className='text-gray-500'>(Coming soon)</span>
        </CommandItem>{' '}
        <CommandItem disabled leftAccessory={<ClipboardCheck />}>
          <span className='text-gray-700'>Create to-do</span>{' '}
          <span className='text-gray-500'>(Coming soon)</span>
        </CommandItem>
      </Command.List>
    </Command>
  );
});

const StepMenu = observer(() => {
  return (
    <>
      <StepsHub />
    </>
  );
});
