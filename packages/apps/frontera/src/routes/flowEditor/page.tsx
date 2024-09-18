import React, { useState, useEffect, useCallback } from 'react';

import { useKey } from 'rooks';
import {
  Panel,
  Handle,
  addEdge,
  MiniMap,
  Controls,
  Position,
  ReactFlow,
  Background,
  MarkerType,
  useNodesState,
  useEdgesState,
} from '@xyflow/react';

import { Select } from '@ui/form/Select';
import { Play } from '@ui/media/icons/Play.tsx';
import { Plus } from '@ui/media/icons/Plus.tsx';
import { Link05 } from '@ui/media/icons/Link05';
import { IconButton } from '@ui/form/IconButton';
import { Check } from '@ui/media/icons/Check.tsx';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { Button } from '@ui/form/Button/Button.tsx';
import { Editor } from '@ui/form/Editor/Editor.tsx';
import { Mail01 } from '@ui/media/icons/Mail01.tsx';
import { Users01 } from '@ui/media/icons/Users01.tsx';
import { Input, ResizableInput } from '@ui/form/Input';
import { UserPlus01 } from '@ui/media/icons/UserPlus01.tsx';
import { Hourglass02 } from '@ui/media/icons/Hourglass02.tsx';
import { ChevronRight } from '@ui/media/icons/ChevronRight.tsx';
import { Tag, TagLabel, TagLeftIcon } from '@ui/presentation/Tag';
import { extractPlainText } from '@ui/form/Editor/utils/extractPlainText.ts';
import { MessageTextSquare01 } from '@ui/media/icons/MessageTextSquare01.tsx';
import { convertPlainTextToHtml } from '@ui/form/Editor/utils/convertPlainTextToHtml.ts';

import SidePanel from './src/SidePanel';
import { TimeTriggerEdge } from './src/edges';

import '@xyflow/react/dist/style.css';

const StartNode = ({ data }) => (
  <div className='aspect-[9/1] max-w-[230px] bg-white border-2 border-green-200 p-3 rounded-lg shadow-md'>
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

const BasicNode = ({ type, data }) => (
  <div
    className={`aspect-[9/1] bg-white border-2 border-${data.color}-200 p-3 rounded-lg shadow-md`}
  >
    <Handle
      type='target'
      position={Position.Left}
      className={`h-2 w-2 bg-${data.color}-500`}
    />
    <div className='font-bold'>Send {type}</div>
    <div>{data.content || 'No content'}</div>

    <Handle
      type='source'
      position={Position.Right}
      className={`h-2 w-2 bg-${data.color}-500`}
    />
  </div>
);
const EmailNode = ({ data }) => (
  <div
    className={`aspect-[9/1] max-w-[230px] bg-white border-2 border-${data.color}-200 p-3 rounded-lg shadow-md`}
  >
    <Handle
      type='target'
      position={Position.Left}
      className={`h-2 w-2 bg-${data.color}-500`}
    />
    <div className='flex items-center font-bold'>
      <div className='flex items-center'>
        <Mail01 className='mr-2 w-4 h-4' />
      </div>
      <span className='truncate'>{data.subject || `Send Email`} </span>
    </div>
    <div className='truncate  text-sm mt-1'>
      {data.content || (
        <span className='text-gray-400'>Double click to edit content </span>
      )}
    </div>

    <Handle
      type='source'
      position={Position.Right}
      className={`h-2 w-2 bg-${data.color}-500`}
    />
  </div>
);

const nodeTypes = {
  startNode: StartNode,

  emailNode: (props) => <EmailNode {...props} type='Email' />,
  linkedInMessageNode: (props) => (
    <BasicNode {...props} type='LinkedIn Message' />
  ),
  linkedInInviteNode: (props) => (
    <BasicNode {...props} type='LinkedIn Invite' />
  ),
};

const edgeTypes = {
  triggerEdge: TimeTriggerEdge,
};

export const MarketingFlowBuilder = () => {
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [open, setOpen] = useState(false);

  useEffect(() => {
    // Add the start node if it doesn't exist
    if (!nodes.some((node) => node.type === 'startNode')) {
      setNodes([
        {
          id: 'start-node',
          type: 'startNode',
          position: { x: 0, y: 0 },
          data: {
            triggerType: 'manual',
            interval: 1,
            unit: 'days',
          },
        },
      ]);
    }
  }, []);

  const onConnect = useCallback(
    (params) => {
      const edgeData = {
        triggerType: 'time',
        timeValue: 1,
        timeUnit: 'days',
      };

      setEdges((eds) =>
        addEdge(
          {
            ...params,
            type: 'triggerEdge',
            data: edgeData,
            style: { strokeWidth: 2 },
            markerEnd: { type: MarkerType.Arrow },
            label: `1 day`, // Default label
          },
          eds,
        ),
      );
    },
    [setEdges],
  );

  const editNode = useCallback((nodeId) => {
    setOpen(nodeId);
  }, []);

  const editEdge = useCallback(
    (edge) => {
      console.log('🏷️ ----- edge: ', edge);
      setOpen(edge.id);
    },
    [setEdges],
  );

  const addNode = useCallback(
    (type) => {
      const newNode = {
        id: `${type}-${nodes.length + 1}`,
        type,
        position: { x: Math.random() * 500, y: Math.random() * 500 },
        data: {
          color:
            type === 'emailNode'
              ? 'blue'
              : type === 'linkedInMessageNode'
              ? 'green'
              : 'yellow',
          subject: '',
        },
      };

      setNodes((nds) => nds.concat(newNode));
    },
    [nodes.length, setNodes],
  );

  // Keyboard shortcuts
  useKey(['E'], () => addNode('emailNode'), {
    when: !open,
  });
  useKey(['L'], () => addNode('linkedInMessageNode'), {
    when: !open,
  });
  useKey(['I'], () => addNode('linkedInInviteNode'), {
    when: !open,
  });

  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <ReactFlow
        fitView
        nodes={nodes}
        edges={edges}
        onConnect={onConnect}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onEdgeDoubleClick={(_, edge) => editEdge(edge)}
        onNodeDoubleClick={(_, node) => editNode(node?.id)}
      >
        <div className='bg-white px-10 h-14 border-b flex items-center text-base font-bold'>
          <span className='font-medium text-gray-500'>Flows</span>
          <ChevronRight className='size-4 mx-1.5 text-gray-500' />
          <span className='mr-2'>Waitlist</span>
          <Tag
            size='md'
            variant='outline'
            colorScheme='gray'
            onClick={() => {}}
            className='bg-transparent py-0.5 cursor-pointer text-gray-500'
          >
            <TagLeftIcon>
              <Users01 />
            </TagLeftIcon>
            <TagLabel>4 contacts</TagLabel>
          </Tag>
        </div>
        <Background />
        <Controls />
        <MiniMap />
        <Panel position='bottom-center'>
          <ButtonGroup className='bg-white'>
            <IconButton
              icon={<Mail01 />}
              aria-label='Add Email'
              onClick={() => addNode('emailNode')}
            />
            <IconButton
              icon={<MessageTextSquare01 />}
              aria-label='Send LinkedIn Message'
              onClick={() => addNode('linkedInMessageNode')}
            />
            <IconButton
              icon={<UserPlus01 />}
              aria-label='Send LinkedIn Invite'
              onClick={() => addNode('linkedInInviteNode')}
            />
          </ButtonGroup>
        </Panel>
      </ReactFlow>
      <SidePanel
        open={open}
        nodeId={open}
        nodes={nodes}
        edges={edges}
        setOpen={setOpen}
        setNodes={setNodes}
        setEdges={setEdges}
      />
    </div>
  );
};
