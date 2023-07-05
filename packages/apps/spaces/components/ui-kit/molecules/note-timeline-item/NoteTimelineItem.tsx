import React, {
  MutableRefObject,
  Ref,
  useCallback,
  useEffect,
  useRef,
  useState,
} from 'react';
import styles from './note.module.scss';
import { toast } from 'react-toastify';
import parse from 'html-react-parser';
import ReactDOMServer from 'react-dom/server';
import axios from 'axios';
import Check from '@spaces/atoms/icons/Check';
import Trash from '@spaces/atoms/icons/Trash';
import Pencil from '@spaces/atoms/icons/Pencil';
import { Avatar } from '@spaces/atoms/avatar';
import { IconButton } from '@spaces/atoms/icon-button/IconButton';
import sanitizeHtml from 'sanitize-html';
import {
  useDeleteNote,
  useLinkNoteAttachment,
  useUnlinkNoteAttachment,
  useUpdateNote,
} from '@spaces/hooks/useNote';
import linkifyHtml from 'linkify-html';
import { getContactDisplayName } from '../../../../utils';
import classNames from 'classnames';
import { extraAttributes, SocialEditor } from '../editor/SocialEditor';
import { TableExtension } from '@remirror/extension-react-tables';
import {
  AnnotationExtension,
  BlockquoteExtension,
  BoldExtension,
  BulletListExtension,
  FontSizeExtension,
  HistoryExtension,
  ImageExtension,
  ItalicExtension,
  LinkExtension,
  MentionAtomExtension,
  OrderedListExtension,
  StrikeExtension,
  TextColorExtension,
  UnderlineExtension,
  wysiwygPreset,
} from 'remirror/extensions';
import { useRemirror } from '@remirror/react';
import { useRecoilState } from 'recoil';
import { prosemirrorNodeToHtml } from 'remirror';
import { contactNewItemsToEdit } from '../../../../state';
import { useFileUpload } from '@spaces/hooks/useFileUpload';
import { Note } from '../../../../hooks/useNote/types';
import Paperclip from '@spaces/atoms/icons/Paperclip';
import { DeleteConfirmationDialog } from '@spaces/atoms/delete-confirmation-dialog';
import { DataSource } from '@spaces/graphql';

interface Props {
  note: Note;
}

export const NoteTimelineItem: React.FC<Props> = ({ note }) => {
  const [images, setImages] = useState({});
  const [deleteConfirmationModalVisible, setDeleteConfirmationModalVisible] =
    useState(false);
  const { onUpdateNote } = useUpdateNote();
  const { onRemoveNote } = useDeleteNote();
  const [itemsInEditMode, setItemToEditMode] = useRecoilState(
    contactNewItemsToEdit,
  );

  const [editNote, setEditNote] = useState(false);
  const elementRef = useRef<HTMLDivElement>(null);
  const { onLinkNoteAttachment } = useLinkNoteAttachment({
    noteId: note.id,
  });
  const { onUnlinkNoteAttachment } = useUnlinkNoteAttachment({
    noteId: note.id,
  });
  const uploadInputRef = React.useRef<HTMLInputElement>(null);

  const { handleInputFileChange } = useFileUpload({
    prevFiles: [],
    onBeginFileUpload: (data) => console.log(''),
    onFileUpload: (newFile) => {
      return onLinkNoteAttachment(newFile.id);
    },
    onFileUploadError: () =>
      toast.error('Something went wrong while uploading attachment'),
    onFileRemove: (fileId: string) => {
      return onUnlinkNoteAttachment(fileId);
    },
    uploadInputRef,
  });

  const remirrorExtentions = [
    new TableExtension(),
    new MentionAtomExtension({
      matchers: [
        { name: 'at', char: '@' },
        { name: 'tag', char: '#' },
      ],
    }),
    ...wysiwygPreset(),

    ...wysiwygPreset(),
    new BoldExtension(),
    new ItalicExtension(),
    new BlockquoteExtension(),
    new ImageExtension(),
    new LinkExtension({ autoLink: true }),
    new TextColorExtension(),
    new UnderlineExtension(),
    new FontSizeExtension(),
    new HistoryExtension(),
    new AnnotationExtension(),
    new BulletListExtension(),
    new OrderedListExtension(),
    new StrikeExtension(),
  ];
  const extensions = useCallback(() => [...remirrorExtentions], [note.id]);

  const { manager, state, setState, getContext } = useRemirror({
    extensions,
    extraAttributes,
    // state can created from a html string.
    stringHandler: 'html',

    // This content is used to create the initial value. It is never referred to again after the first render.
    content: sanitizeHtml(
      linkifyHtml(note.html, {
        defaultProtocol: 'https',
        rel: 'noopener noreferrer',
      }),
    ),
  });

  useEffect(() => {
    if (
      itemsInEditMode.timelineEvents.findIndex(
        (data: { id: string }) => data.id === note.id,
      ) !== -1 &&
      elementRef.current
    ) {
      setEditNote(true);
      elementRef.current?.scrollIntoView();
    }
  }, [elementRef]);

  useEffect(() => {
    if ((note.html.match(/<img/g) || []).length > 0) {
      parse(note.html, {
        replace: (domNode: any) => {
          if (
            domNode.name === 'img' &&
            domNode.attribs &&
            domNode.attribs.alt
          ) {
            const alt = domNode.attribs.alt;

            axios
              .get(`/fs/file/${alt}/base64`)
              .then(async (response: any) => {
                const dataUrl = response.data;

                setImages((prevImages: any) => {
                  const t = {} as any;
                  t[alt] = dataUrl as string;
                  return {
                    ...prevImages,
                    ...t,
                  };
                });
              })
              .catch(() => {
                toast.error(
                  'There was a problem on our side and we are doing our best to solve it!',
                );
              });
          }
        },
      });
    } else {
      // reset({ id, html: noteContent, htmlEnhanced: noteContent });
    }
  }, [note.id, note.html]);

  useEffect(() => {
    const imagesToLoad = (note.html.match(/<img/g) || []).length;
    if (imagesToLoad > 0 && Object.keys(images).length === imagesToLoad) {
      const htmlParsed = parse(note.html, {
        replace: (domNode: any) => {
          if (
            domNode.name === 'img' &&
            domNode.attribs &&
            domNode.attribs.alt
          ) {
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-expect-error
            const imageSrc = images[domNode.attribs.alt] as string;
            return (
              <img
                src={imageSrc}
                alt={domNode.attribs.alt}
                style={{ width: '200px' }}
              />
            );
          }
        },
      });

      const html = ReactDOMServer.renderToString(htmlParsed as any);

      getContext()?.setContent(html);
    }
  }, [note.id, images, note.html, editNote]);

  const handleUpdateNote = (id: string) => {
    const data = prosemirrorNodeToHtml(state.doc);

    const dataToSubmit = {
      id,
      html: data?.replaceAll(/.src(\S*)/g, '') || '',
    };
    onUpdateNote(dataToSubmit).then(() => {
      setEditNote(false);
    });
  };
  const handleToggleEditMode = (state: boolean) => {
    setEditNote(state);
    setTimeout(() => {
      if (elementRef?.current) {
        elementRef.current.scrollIntoView({
          behavior: 'smooth',
          inline: 'start',
        });
      }
    }, 0);
  };

  return (
    <div className={styles.noteWrapper} ref={elementRef}>
      <div
        className={classNames(styles.noteContainer, {
          [styles.withToolbar]: editNote,
        })}
      >
        <div className={styles.actions}>
          {note?.noted?.map((data: any, index: any) => {
            const isContact = data.__typename === 'Contact';
            const isOrg = data.__typename === 'Organization';

            if (isContact) {
              const name = getContactDisplayName(data).split(' ');
              const surname = name?.length === 2 ? name[1] : name[2];

              return (
                <Avatar
                  key={`${data.id}-${index}`}
                  name={name?.[0]}
                  surname={surname}
                  size={30}
                />
              );
            }

            if (isOrg) {
              return (
                <Avatar
                  key={`${data.id}-${index}`}
                  name={data.organizationName}
                  surname={''}
                  isSquare={data.__typename === 'Organization'}
                  size={30}
                />
              );
            }

            return <div key={`avatar-error-${data.id}-${index}`} />;
          })}

          {editNote && (
            <IconButton
              size='xxxs'
              onClick={() => setDeleteConfirmationModalVisible(true)}
              icon={<Trash color='red' height={20} />}
              mode='text'
              label='Delete'
              style={{ marginBottom: 0 }}
            />
          )}
        </div>
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            width: '100%',
          }}
        >
          <div
            className={classNames(styles.noteContent, {
              [styles.editNoteContent]: editNote,
              [styles.withFiles]: note.includes?.length > 0,
            })}
          >
            {note.source === DataSource.ZendeskSupport && (
              <p style={{ marginLeft: '1rem' }}>
                Comment on Zendesk issue:{' '}
                <span style={{ fontWeight: 600 }}>
                  {note?.mentioned?.[0]?.subject}
                </span>
              </p>
            )}
            <SocialEditor
              mode={editNote ? 'EDIT' : ''}
              editable={editNote}
              manager={manager}
              state={state}
              setState={setState}
              items={[]}
            >
              <input
                style={{ display: 'none' }}
                ref={uploadInputRef}
                type='file'
                onChange={handleInputFileChange}
              />
              <IconButton
                label='Attach file'
                isSquare
                mode='text'
                onClick={() => uploadInputRef?.current?.click()}
                icon={<Paperclip />}
              />
            </SocialEditor>
          </div>
          <DeleteConfirmationDialog
            deleteConfirmationModalVisible={deleteConfirmationModalVisible}
            setDeleteConfirmationModalVisible={
              setDeleteConfirmationModalVisible
            }
            deleteAction={() =>
              onRemoveNote(note.id).then(() =>
                setDeleteConfirmationModalVisible(false),
              )
            }
            confirmationButtonLabel='Delete note'
          />
        </div>

        <div className={styles.actions}>
          <Avatar
            name={note.createdBy?.firstName || ''}
            surname={note.createdBy?.lastName || ''}
            size={30}
          />
          {editNote ? (
            <IconButton
              size='xxxs'
              onClick={() => handleUpdateNote(note.id)}
              icon={<Check height={20} />}
              mode='text'
              label='Done'
              style={{ marginBottom: 0, color: 'green' }}
            />
          ) : note.source !== DataSource.ZendeskSupport ? (
            <IconButton
              size='xxxs'
              onClick={() => handleToggleEditMode(true)}
              icon={<Pencil height={20} />}
              mode='text'
              label='Edit'
              style={{ marginBottom: 0 }}
            />
          ) : null}
        </div>
      </div>
    </div>
  );
};
