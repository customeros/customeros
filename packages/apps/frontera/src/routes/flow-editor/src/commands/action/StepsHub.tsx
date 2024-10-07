import { observer } from 'mobx-react-lite';
import { FlowActionType } from '@store/Flows/types.ts';
import { Node, MarkerType, useReactFlow } from '@xyflow/react';

import { cn } from '@ui/utils/cn.ts';
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

import { keywords } from './keywords.ts';
import { useUndoRedo } from '../../hooks';

export const StepsHub = observer(() => {
  const { ui } = useStore();
  const { takeSnapshot } = useUndoRedo();

  const { setEdges, setNodes, getNodes, getEdges } = useReactFlow();

  const findPreviousEmailNode = (
    nodes: Node[],
    currentNodeId: string,
  ): Node | null => {
    const currentNodeIndex = nodes.findIndex(
      (node) => node.id === currentNodeId,
    );

    for (let i = currentNodeIndex; i >= 0; i--) {
      if (nodes[i].data.action === FlowActionType.EMAIL_NEW) {
        return nodes[i];
      }
    }

    return null;
  };

  const hasEmailNodeBeforeCurrent = (
    nodes: Node[],
    currentNodeId: string,
  ): boolean => {
    return findPreviousEmailNode(nodes, currentNodeId) !== null;
  };

  const handleAddNode = async (type: FlowActionType | 'WAIT') => {
    takeSnapshot();

    const nodes = getNodes();
    const edges = getEdges();

    const sourceNode = nodes.find(
      (node) => node.id === ui.flowCommandMenu.context.meta?.source,
    );
    const targetNode = nodes.find(
      (node) => node.id === ui.flowCommandMenu.context.meta?.target,
    );

    if (!sourceNode || !targetNode) return;

    let typeBasedContent: {
      subject?: string;
      replyTo?: string;
      bodyTemplate?: string;
      waitDuration?: number;
    } = {};

    if (type === 'WAIT') {
      typeBasedContent = { waitDuration: 1 };
    } else if (type === FlowActionType.EMAIL_REPLY) {
      const prevEmailNode = findPreviousEmailNode(nodes, sourceNode.id);
      const prevSubject = prevEmailNode?.data?.subject || '';

      typeBasedContent = {
        replyTo: prevEmailNode?.id,
        subject: `RE: ${prevSubject}`,
        bodyTemplate: '',
      };
    } else {
      typeBasedContent = { subject: '', bodyTemplate: '' };
    }

    const newNode = {
      id: `${type}-${nodes.length + 1}`,
      type: type === 'WAIT' ? 'wait' : 'action',
      position: { x: sourceNode.position.x, y: sourceNode.position.y + 10 },
      data: {
        action: type,
        ...typeBasedContent,
      },
    };

    const edgeToNewNode = {
      id: `e${ui.flowCommandMenu.context.meta?.source}-${newNode.id}`,
      source: ui.flowCommandMenu.context.meta?.source,
      target: newNode.id,
      type: 'baseEdge',
      markerEnd: {
        type: MarkerType.Arrow,
        width: 20,
        height: 20,
      },
    };

    const edgeFromNewNode = {
      id: `e${newNode.id}-${ui.flowCommandMenu.context.meta?.target}`,
      source: newNode.id,
      target: ui.flowCommandMenu.context.meta?.target,
      type: 'baseEdge',
      markerEnd: {
        type: MarkerType.Arrow,
        width: 20,
        height: 20,
      },
    };

    const updatedEdges = edges.filter(
      (e) =>
        !(
          e.source === ui.flowCommandMenu.context.meta?.source &&
          e.target === ui.flowCommandMenu.context.meta?.target
        ),
    );

    const updatedNodes = [...nodes, newNode];
    const newEdges = [...updatedEdges, edgeToNewNode, edgeFromNewNode];

    setNodes(updatedNodes);
    setEdges(newEdges);
    // getLayoutedElements(updatedNodes, newEdges, elkOptions).then(
    //   // Use getLayoutedElements to calculate new positions
    //   // @ts-expect-error fix type later
    //   ({ nodes: layoutedNodes, edges: layoutedEdges }) => {
    //
    //     window.requestAnimationFrame(() => {
    //       const fitViewOptions: FitViewOptions = {
    //         padding: 0.1,
    //         maxZoom: 1,
    //         minZoom: 1,
    //         duration: 0,
    //         nodes: layoutedNodes,
    //       };
    //
    //       fitView(fitViewOptions);
    //     });
    //   },
    // );
  };

  const updateSelectedNode = (type: FlowActionType | 'WAIT') => {
    handleAddNode(type);
    ui.flowCommandMenu.setOpen(false);
    ui.flowCommandMenu.reset();
  };

  const currentNodeId = ui.flowCommandMenu.context.meta?.source;

  const canReplyToEmail = currentNodeId
    ? hasEmailNodeBeforeCurrent(getNodes(), currentNodeId)
    : false;

  return (
    <>
      <CommandItem
        leftAccessory={<Mail01 />}
        keywords={keywords.send_email}
        onSelect={() => {
          updateSelectedNode(FlowActionType.EMAIL_NEW);
        }}
      >
        Send email
      </CommandItem>
      <CommandItem
        disabled={!canReplyToEmail}
        leftAccessory={<MailReply />}
        keywords={keywords.reply_to_previous_email}
        className={cn({
          hidden: !canReplyToEmail,
        })}
        onSelect={() => {
          updateSelectedNode(FlowActionType.EMAIL_REPLY);
        }}
      >
        Reply to previous email
      </CommandItem>
      <CommandItem
        keywords={keywords.wait}
        leftAccessory={<Hourglass02 />}
        onSelect={() => {
          updateSelectedNode('WAIT');
        }}
      >
        Wait
      </CommandItem>
      <CommandItem
        disabled
        leftAccessory={<LinkedinOutline />}
        keywords={keywords.send_linkedin_message}
      >
        <span className='text-gray-700'>Send LinkedIn message</span>
        <span className='text-gray-500'>(Coming soon)</span>
      </CommandItem>
      <CommandItem
        disabled
        leftAccessory={<PlusSquare />}
        keywords={keywords.create_record}
      >
        <span className='text-gray-700'>Create record</span>
        <span className='text-gray-500'>(Coming soon)</span>
      </CommandItem>
      <CommandItem
        disabled
        leftAccessory={<RefreshCw01 />}
        keywords={keywords.update_record}
      >
        <span className='text-gray-700'>Update record</span>
        <span className='text-gray-500'>(Coming soon)</span>
      </CommandItem>
      <CommandItem
        disabled
        leftAccessory={<Star06 />}
        keywords={keywords.enrich_record}
      >
        <span className='text-gray-700'>Enrich record</span>
        <span className='text-gray-500'>(Coming soon)</span>
      </CommandItem>
      <CommandItem
        disabled
        leftAccessory={<Star06 />}
        keywords={keywords.verify_record_property}
      >
        <span className='text-gray-700'>Verify record property</span>
        <span className='text-gray-500'>(Coming soon)</span>
      </CommandItem>
      <CommandItem
        disabled
        keywords={keywords.conditions}
        leftAccessory={<ArrowIfPath />}
      >
        <span className='text-gray-700'>Conditions</span>
        <span className='text-gray-500'>(Coming soon)</span>
      </CommandItem>
      <CommandItem
        disabled
        keywords={keywords.create_to_do}
        leftAccessory={<ClipboardCheck />}
      >
        <span className='text-gray-700'>Create to-do</span>
        <span className='text-gray-500'>(Coming soon)</span>
      </CommandItem>
    </>
  );
});
