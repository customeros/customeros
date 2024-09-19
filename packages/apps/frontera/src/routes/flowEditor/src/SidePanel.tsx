import React from 'react';

import { Input } from '@ui/form/Input';
import { Select } from '@ui/form/Select';
import { Plus } from '@ui/media/icons/Plus.tsx';
import { Editor } from '@ui/form/Editor/Editor';
import { IconButton } from '@ui/form/IconButton';
import { Check } from '@ui/media/icons/Check.tsx';
import { Button } from '@ui/form/Button/Button.tsx';
import { Hourglass02 } from '@ui/media/icons/Hourglass02.tsx';
import { extractPlainText } from '@ui/form/Editor/utils/extractPlainText.ts';
import { convertPlainTextToHtml } from '@ui/form/Editor/utils/convertPlainTextToHtml.ts';

const EmailNodeEdit = ({ node, handleSubjectChange, handleContentChange }) => (
  <>
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
  </>
);

const LinkedInNodeEdit = ({ node, handleContentChange }) => (
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
    placeholder={`Enter your LinkedIn message here...`}
    defaultHtmlValue={convertPlainTextToHtml(node.data.content || '')}
  />
);

const TriggerEdit = ({
  outgoingEdges,
  handleTriggerTypeChange,
  handleTimeValueChange,
  handleTimeUnitChange,
}) => (
  <div className='mt-4'>
    <div className='flex items-center gap-1.5 text-base'>
      <Hourglass02 />
      <span>Wait </span>
      <input
        type='number'
        className='text-gray-700 underline min-h-3 w-12'
        value={outgoingEdges[0]?.data?.timeValue ?? '1'}
        onChange={(e) =>
          handleTimeValueChange(outgoingEdges[0].id, e.target.value)
        }
      />{' '}
      <Select
        value={outgoingEdges[0]?.data?.timeUnit ?? 'days'}
        onChange={(e) =>
          handleTimeUnitChange(outgoingEdges[0].id, e.target.value)
        }
      >
        <option value='minutes'>Minutes</option>
        <option value='hours'>Hours</option>
        <option value='days'>Days</option>
      </Select>
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
              onChange={(e) => handleTimeValueChange(edge.id, e.target.value)}
            />
            <Select
              value={edge.data.timeUnit}
              onChange={(e) => handleTimeUnitChange(edge.id, e.target.value)}
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
);

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

  if (!open) return null;

  const renderNodeContent = () => {
    switch (node.type) {
      case 'emailNode':
        return (
          <EmailNodeEdit
            node={node}
            handleSubjectChange={handleSubjectChange}
            handleContentChange={handleContentChange}
          />
        );
      case 'linkedInMessageNode':
      case 'linkedInInviteNode':
        return (
          <LinkedInNodeEdit
            node={node}
            handleContentChange={handleContentChange}
          />
        );
      default:
        return null;
    }
  };

  const labels = {
    emailNode: 'email',
    startNode: 'start flow ',
    linkedInMessageNode: 'LinkedIn message',
    linkedInInviteNode: 'LinkedIn invite',
  };

  return (
    <div className='min-w-[400px] w-[400px] bg-white absolute top-0 right-0 py-4 px-6 flex flex-col h-[100vh] border-t border-l animate-slideLeft shadow-xl'>
      <div className='flex mt-3 mb-9'>
        <div className='flex-1'>
          <div className='font-bold text-2xl'>
            {node ? `Edit ${labels?.[node?.type]}` : 'Edit trigger'}
          </div>
        </div>
        <IconButton
          variant='ghost'
          icon={<Check />}
          aria-label='Done'
          onClick={() => setOpen(false)}
        />
      </div>

      {node ? (
        renderNodeContent()
      ) : (
        <TriggerEdit
          outgoingEdges={outgoingEdges}
          handleTimeUnitChange={handleTimeUnitChange}
          handleTimeValueChange={handleTimeValueChange}
          handleTriggerTypeChange={handleTriggerTypeChange}
        />
      )}
    </div>
  );
};

export default SidePanel;
