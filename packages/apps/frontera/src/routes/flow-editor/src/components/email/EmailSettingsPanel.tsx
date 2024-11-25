import { useParams } from 'react-router-dom';
import { useRef, useMemo, useState, useEffect } from 'react';

import { LexicalEditor } from 'lexical';
import { observer } from 'mobx-react-lite';
import { useReactFlow } from '@xyflow/react';
import { render } from '@react-email/render';
import { FlowActionType } from '@store/Flows/types';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button';
import { Editor } from '@ui/form/Editor/Editor';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
// import { MailAdd } from '@ui/media/icons/MailAdd';
import { ZenCircle } from '@ui/media/icons/ZenCircle';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { EmailTemplate } from '@shared/components/EmailTemplate';
import { extractPlainText } from '@ui/form/Editor/utils/extractPlainText';
import { convertPlainTextToHtml } from '@ui/form/Editor/utils/convertPlainTextToHtml';

import { useUndoRedo } from '../../hooks';
import { EmailEditorModal } from './EmailEditorModal';

export const EmailSettingsPanel = observer(() => {
  const { ui, flowEmailVariables, flows } = useStore();
  const flowId = useParams()?.id as string;
  const [focusMode, setFocusMode] = useState(false);
  const sidePanelStore = ui.flowActionSidePanel;
  const { setNodes } = useReactFlow();

  const inputRef = useRef<LexicalEditor>(null);
  const editorRef = useRef<LexicalEditor>(null);
  const placeholder = useMemo(() => getRandomEmailPrompt(), []);
  const variables = flowEmailVariables?.value.get('CONTACT')?.variables;
  const data = ui.flowActionSidePanel.context.node?.data;

  const [subject, setSubject] = useState(data?.subject ?? '');
  const [bodyTemplate, setBodyTemplate] = useState(data?.bodyTemplate ?? '');

  useEffect(() => {
    setSubject(data?.subject ?? '');
    setBodyTemplate(data?.bodyTemplate ?? '');

    if (
      data?.action !== FlowActionType.EMAIL_REPLY &&
      data?.subject?.trim()?.length === 0
    ) {
      setTimeout(() => {
        inputRef.current?.focus();
      }, 0);
    }
  }, [data?.subject, data?.bodyTemplate, data?.action]);

  const { takeSnapshot } = useUndoRedo();

  const flow = flows.value.get(flowId)?.value?.name;

  useEffect(() => {
    const isSubjectChanged = subject !== data?.subject;
    const plainTextBodyBefore =
      data?.bodyTemplate && extractPlainText(data?.bodyTemplate)?.trim();
    const plainTextBody = extractPlainText(bodyTemplate)?.trim();

    if (isSubjectChanged || plainTextBodyBefore !== plainTextBody) {
      sidePanelStore.setContext({
        ...sidePanelStore.context,
        hasUnsavedChanges: true,
      });
    }

    if (data?.subject === subject && data?.bodyTemplate === bodyTemplate) {
      sidePanelStore.setContext({
        ...sidePanelStore.context,
        hasUnsavedChanges: false,
      });
    }
  }, [data?.subject, data?.bodyTemplate, subject, bodyTemplate]);

  const prepareEmailContent = async (bodyHtml: string) => {
    try {
      const emailHtml = await render(<EmailTemplate bodyHtml={bodyHtml} />, {
        pretty: true,
      });

      return {
        html: emailHtml,
      };
    } catch (error) {
      ui.toastError(
        'Unable to process email content',
        'email-content-parsing-error',
      );
    }
  };

  const handleEmailDataChange = ({
    subject,
    bodyTemplate,
  }: {
    subject: string;
    bodyTemplate: string;
  }) => {
    setNodes((nds) =>
      nds.map((node) => {
        if (node.id === sidePanelStore.context.node?.id) {
          return {
            ...node,
            data: {
              ...node.data,
              subject,
              bodyTemplate,
              isEditing: false,
            },
          };
        }

        if (
          !!node.data?.replyTo &&
          node.data?.replyTo === sidePanelStore.context.node?.id
        ) {
          return {
            ...node,
            data: {
              ...node.data,
              subject: `RE: ${subject}`,
            },
          };
        }

        return node;
      }),
    );
  };

  const handleSave = async () => {
    try {
      const emailContent = await prepareEmailContent(bodyTemplate);

      if (emailContent?.html) {
        handleEmailDataChange({
          subject,
          bodyTemplate: emailContent.html,
        });
      }

      setTimeout(() => {
        takeSnapshot();
        ui.commandMenu.setContext({
          entity: 'Flow',
          ids: [flowId],
        });
        sidePanelStore.clearContext();
        sidePanelStore.setOpen(false);
      }, 0);
    } catch (error) {
      console.error('Error saving email:', error);
      ui.toastError('Error saving email', 'email-save-error');
    }
  };

  useEffect(() => {
    ui.commandMenu.setCallback(handleSave);
  }, [handleSave]);

  const handleCancelChanges = () => {
    ui.commandMenu.clearContext();
    sidePanelStore.clearContext();
    sidePanelStore.setOpen(false);
  };

  return (
    <div
      onClick={(e) => e.stopPropagation()}
      onKeyDown={(e) => e.stopPropagation()}
      className='pl-4 -ml-4 fixed z-50 top-[0px] bottom-0 right-0 w-[400px] h-full  '
    >
      <article className=' h-full bg-white  border-l flex flex-col gap-1 animate-slideLeft'>
        <div className='flex justify-between items-center border-b border-gray-200 p-4 y-2 h-[41px]'>
          <div className='flex items-center gap-2'>
            <h1 className='font-medium'>Email {}</h1>
            <Tooltip label='Focus mode'>
              <IconButton
                size='xs'
                variant='ghost'
                icon={<ZenCircle />}
                aria-label={'Focus mode'}
                onClick={() => setFocusMode(true)}
              />
            </Tooltip>
          </div>

          <div className='flex gap-2'>
            <div>
              {ui.flowActionSidePanel.context.hasUnsavedChanges && (
                <Button size='xs' variant='ghost' onClick={handleCancelChanges}>
                  Discard changes
                </Button>
              )}
            </div>
            <div>
              <Button size='xs' variant='outline' onClick={handleSave}>
                Done
              </Button>
            </div>
          </div>
        </div>

        <div className='px-4 overflow-y-auto min-h-[50%] h-full'>
          <Tooltip
            align='start'
            label={
              data?.action === FlowActionType.EMAIL_REPLY
                ? `Reply to email subjects can't be edited`
                : ''
            }
          >
            <div
              className={cn({
                'cursor-not-allowed':
                  data?.action === FlowActionType.EMAIL_REPLY,
              })}
            >
              <Editor
                size='md'
                usePlainText
                ref={inputRef}
                placeholder='Subject'
                variableOptions={variables}
                namespace='flow-email-editor-subject'
                onChange={(html) => setSubject(extractPlainText(html))}
                defaultHtmlValue={convertPlainTextToHtml(subject ?? '')}
                onKeyDown={(e) => {
                  if (e.key === 'Tab') {
                    e.preventDefault();
                    editorRef.current?.focus();
                  }
                }}
                placeholderClassName={cn(
                  'text-sm font-medium h-auto cursor-text ',
                  {
                    'pointer-events-none text-gray-400':
                      data?.action === FlowActionType.EMAIL_REPLY,
                  },
                )}
                className={cn(
                  `text-sm font-medium h-auto cursor-text email-editor-subject`,
                  {
                    'pointer-events-none text-gray-400':
                      data?.action === FlowActionType.EMAIL_REPLY,
                  },
                )}
              />
            </div>
          </Tooltip>
          <div className='min-h-[60vh] mb-2'>
            <Editor
              ref={editorRef}
              placeholder={placeholder}
              variableOptions={variables}
              dataTest='flow-email-editor'
              namespace='flow-email-editor'
              defaultHtmlValue={bodyTemplate}
              placeholderClassName='text-sm '
              className='text-sm cursor-text email-editor h-full'
              onChange={(e) => {
                setBodyTemplate(e);
              }}
            />
          </div>
        </div>
        {/*<Button*/}
        {/*  size='lg'*/}
        {/*  isDisabled*/}
        {/*  variant='ghost'*/}
        {/*  colorScheme='primary'*/}
        {/*  leftIcon={<MailAdd className='text-inherit size-4' />}*/}
        {/*  className='text-primary-700 w-full text-sm absolute bottom-0 py-4 rounded-none border-t border-solid border-gray-200'*/}
        {/*>*/}
        {/*  Set up a test email...*/}
        {/*</Button>*/}
        <EmailEditorModal
          flowName={flow}
          subject={subject}
          action={data?.action}
          variables={variables}
          setSubject={setSubject}
          handleSave={handleSave}
          isEditorOpen={focusMode}
          placeholder={placeholder}
          bodyTemplate={bodyTemplate}
          setBodyTemplate={setBodyTemplate}
          handleCancel={() => setFocusMode(false)}
        />
      </article>
    </div>
  );
});

const emailPrompts = [
  'Write something {{contact_first_name}} will want to share with their boss',
  'Craft an email that makes {{contact_first_name}} say "Wow!"',
  'Compose an email {{contact_first_name}} will quote in their presentation',
  "Make {{contact_first_name}} feel like they've discovered a hidden treasure",
  'Write an email that makes {{contact_first_name}} rethink their strategy',
  "Write something {{contact_first_name}} can't get from a Google search",
  'Compose the email that ends {{contact_first_name}}’s decision paralysis',
  "Write an email {{contact_first_name}} can't ignore",
  'Write something that makes {{contact_first_name}} feel stupid for not replying',
  'Write something that makes {{contact_first_name}} say, “Yes, this is what we need!”',
  'Show {{contact_first_name}} what they’re missing—start typing...',
  'Type an email that helps {{contact_first_name}} win',
  'Write something {{contact_first_name}} remember',
];

function getRandomEmailPrompt(): string {
  const randomIndex = Math.floor(Math.random() * emailPrompts.length);

  return emailPrompts[randomIndex];
}
