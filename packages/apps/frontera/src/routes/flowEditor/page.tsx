import React, { useState, useEffect, useCallback } from 'react';

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
    <Handle
      type='source'
      position={Position.Right}
      className='h-2 w-2 bg-green-500'
    />
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

const SidePanel = ({
  open,
  setOpen,
  nodeId,
  nodes,
  setNodes,
  edges,
  setEdges,
}) => {
  const node = nodes.find((n) => n.id === nodeId);
  const outgoingEdges = edges.filter((e) => e.source === nodeId);

  const handleSubjectChange = (e) => {
    setNodes((nds) =>
      nds.map((node) =>
        node.id === nodeId
          ? { ...node, data: { ...node.data, subject: e.target.value } }
          : node,
      ),
    );
  };

  const handleContentChange = (contentD) => {
    const content = extractPlainText(contentD);

    setNodes((nds) =>
      nds.map((node) =>
        node.id === nodeId
          ? { ...node, data: { ...node.data, content } }
          : node,
      ),
    );
  };

  const handleTriggerTypeChange = (edgeId, newTriggerType) => {
    setEdges((eds) =>
      eds.map((e) => {
        if (e.id === edgeId) {
          const newData = { ...e.data, triggerType: newTriggerType };
          const newLabel =
            newTriggerType === 'completion'
              ? 'On completion'
              : `${e.data.timeValue} ${e.data.timeUnit}`;

          return { ...e, data: newData, label: newLabel };
        }

        return e;
      }),
    );
  };

  const handleTimeValueChange = (edgeId, newTimeValue) => {
    setEdges((eds) =>
      eds.map((e) => {
        if (e.id === edgeId) {
          const newData = { ...e.data, timeValue: newTimeValue };

          return {
            ...e,
            data: newData,
            label: `${newTimeValue} ${e.data.timeUnit}`,
          };
        }

        return e;
      }),
    );
  };

  const handleTimeUnitChange = (edgeId, newTimeUnit) => {
    setEdges((eds) =>
      eds.map((e) => {
        if (e.id === edgeId) {
          const newData = { ...e.data, timeUnit: newTimeUnit };

          return {
            ...e,
            data: newData,
            label: `${e.data.timeValue} ${newTimeUnit}`,
          };
        }

        return e;
      }),
    );
  };

  console.log('🏷️ ----- open: ', open);
  console.log('🏷️ ----- node: ', outgoingEdges);

  if (!open) return null;

  if (node) {
    return (
      <div className='min-w-[400px] w-[450px] bg-white absolute top-0 right-0  py-4 px-6 flex flex-col h-[100vh] border-t border-l animate-slideLeft shadow-xl'>
        <div className='flex mt-3 mb-9'>
          <div className='flex-1'>
            <div className='font-bold text-2xl'>Edit {node.type}</div>
          </div>
          <IconButton
            variant='ghost'
            icon={<Check />}
            aria-label='Done'
            onClick={() => setOpen(false)}
          />
        </div>
        <Input
          autoFocus
          placeholder='Subject'
          onChange={handleSubjectChange}
          value={node.data.subject || ''}
        />
        <br />
        <Editor
          usePlainText
          className='mb-10'
          mentionsOptions={[]}
          hashtagsOptions={[]}
          namespace='LogEntryCreator'
          onHashtagSearch={() => null}
          onMentionsSearch={() => null}
          onHashtagsChange={() => null}
          onChange={handleContentChange}
          dataTest='timeline-log-editor'
          value={node.data.content || ''}
          defaultHtmlValue={convertPlainTextToHtml(node.data.content || '')}
          placeholder={`We're excited to invite you to join the early access version of ...`}
        />
      </div>
    );
  }

  return (
    <div className='min-w-[400px] w-[450px] bg-white absolute top-0 right-0  py-4 px-6 flex flex-col h-[100vh] border-t border-l animate-slideLeft shadow-xl'>
      <div className='flex items-center mt-3 mb-9 text-gray-700'>
        <Link05 className='size-8 mr-2' />

        <div className='flex-1'>
          <div className='font-bold text-2xl'>Edit trigger</div>
        </div>
        <IconButton
          variant='ghost'
          icon={<Check />}
          aria-label='Done'
          onClick={() => setOpen(false)}
        />
      </div>

      <div className='mt-4'>
        <div className='flex items-center gap-1.5 text-base'>
          <Hourglass02 />
          <span>Wait </span>
          <ResizableInput
            className='text-gray-700 underline min-h-3'
            value={outgoingEdges?.data?.timeValue ?? '1'}
            onChange={(e) => handleTimeValueChange(nodeId, e.target.value)}
          />{' '}
          day
        </div>

        {outgoingEdges.map((edge) => (
          <div key={edge.id} className='mb-4 p-2 border rounded'>
            <Select
              value={edge.data.triggerType}
              onChange={(e) => handleTriggerTypeChange(edge.id, e.target.value)}
            >
              <option value='time'>Time Delay</option>
              <option value='completion'>On Completion</option>
            </Select>
            {edge.data.triggerType === 'time' && (
              <div className='mt-2 flex items-center'>
                <Input
                  type='number'
                  className='w-20 mr-2'
                  value={edge.data.timeValue}
                  onChange={(e) =>
                    handleTimeValueChange(edge.id, e.target.value)
                  }
                />
                <Select
                  value={edge.data.timeUnit}
                  onChange={(e) =>
                    handleTimeUnitChange(edge.id, e.target.value)
                  }
                >
                  <option value='minutes'>Minutes</option>
                  <option value='hours'>Hours</option>
                  <option value='days'>Days</option>
                </Select>
              </div>
            )}
          </div>
        ))}

        <Button
          variant='ghost'
          className='px-0 mt-4'
          leftIcon={<Plus className='size-4 text-inherit' />}
        >
          Add condition
        </Button>
      </div>
    </div>
  );
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
