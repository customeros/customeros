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
import {
  DeleteConfirmationDialog,
  Trash,
  Pencil,
  IconButton,
  Avatar,
  Check,
} from '../../atoms';
import sanitizeHtml from 'sanitize-html';
import { useDeleteNote, useUpdateNote } from '../../../../hooks/useNote';
import linkifyHtml from 'linkify-html';
import { NotedEntity } from '../../../../graphQL/__generated__/generated';
import { getContactDisplayName } from '../../../../utils';
import classNames from 'classnames';
import { extraAttributes, SocialEditor } from '../editor/SocialEditor';
import { TableExtension } from '@remirror/extension-react-tables';
import {
  AnnotationExtension,
  BlockquoteExtension,
  BoldExtension,
  BulletListExtension,
  EmojiExtension,
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
import data from 'svgmoji/emoji.json';
import { useRemirror } from '@remirror/react';
import { prosemirrorNodeToHtml } from 'remirror';
import { useRecoilState } from 'recoil';
import { contactNewItemsToEdit } from '../../../../state';

interface Props {
  noteContent: string;
  createdAt: string;
  id: string;
  createdBy?: {
    firstName?: string;
    lastName?: string;
  };
  source?: string;
  noted?: Array<NotedEntity>;
}

export const NoteTimelineItem: React.FC<Props> = ({
  noteContent,
  id,
  createdBy,
  noted,
}) => {
  const [images, setImages] = useState({});
  const [deleteConfirmationModalVisible, setDeleteConfirmationModalVisible] =
    useState(false);
  const { onUpdateNote } = useUpdateNote();
  const { onRemoveNote } = useDeleteNote();
  const [itemsInEditMode, setItemToEditMode] = useRecoilState(
    contactNewItemsToEdit,
  );

  const [editNote, setEditNote] = useState(false);
  const elementRef = useRef<MutableRefObject<Ref<HTMLDivElement>>>(null);

  const remirrorExtentions = [
    new TableExtension(),
    new MentionAtomExtension({
      matchers: [
        { name: 'at', char: '@' },
        { name: 'tag', char: '#' },
      ],
    }),

    new EmojiExtension({ plainText: true, data, moji: 'noto' }),
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
  const extensions = useCallback(() => [...remirrorExtentions], [id]);

  const { manager, state, setState, getContext } = useRemirror({
    extensions,
    extraAttributes,
    // state can created from a html string.
    stringHandler: 'html',

    // This content is used to create the initial value. It is never referred to again after the first render.
    content: sanitizeHtml(
      linkifyHtml(noteContent, {
        defaultProtocol: 'https',
        rel: 'noopener noreferrer',
      }),
    ),
  });

  useEffect(() => {
    if (
      itemsInEditMode.timelineEvents.findIndex(
        (data: { id: string }) => data.id === id,
      ) !== -1
    ) {
      setEditNote(true);
    }
  }, []);

  useEffect(() => {
    if ((noteContent.match(/<img/g) || []).length > 0) {
      parse(noteContent, {
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
  }, [id, noteContent]);

  useEffect(() => {
    const imagesToLoad = (noteContent.match(/<img/g) || []).length;
    if (imagesToLoad > 0 && Object.keys(images).length === imagesToLoad) {
      const htmlParsed = parse(noteContent, {
        replace: (domNode: any) => {
          if (
            domNode.name === 'img' &&
            domNode.attribs &&
            domNode.attribs.alt
          ) {
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-ignore
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
  }, [id, images, noteContent, editNote]);

  const handleUpdateNote = () => {
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
        //@ts-expect-error fixme
        elementRef.current.scrollIntoView({
          behavior: 'smooth',
          inline: 'start',
        });
      }
    }, 0);
  };

  return (
    <div
      className={styles.noteWrapper}
      //@ts-expect-error fixme
      ref={elementRef}
    >
      <div
        className={classNames(styles.noteContainer, {
          [styles.withToolbar]: editNote,
        })}
      >
        <div className={styles.actions}>
          {noted?.map((data, index) => {
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
                  // @ts-expect-error this is correct, alias was added and ts does not recognize it
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
              icon={<Trash style={{ transform: 'scale(0.9)', color: 'red' }} />}
              mode='text'
              title='Delete'
              style={{ marginBottom: 0 }}
            />
          )}
        </div>
        <div
          className={classNames(styles.noteContent, {
            [styles.editNoteContent]: editNote,
          })}
        >
          <SocialEditor
            mode={editNote ? 'EDIT' : ''}
            editable={editNote}
            manager={manager}
            state={state}
            setState={setState}
            items={[]}
          />

          <DeleteConfirmationDialog
            deleteConfirmationModalVisible={deleteConfirmationModalVisible}
            setDeleteConfirmationModalVisible={
              setDeleteConfirmationModalVisible
            }
            deleteAction={() =>
              onRemoveNote(id).then(() =>
                setDeleteConfirmationModalVisible(false),
              )
            }
            confirmationButtonLabel='Delete note'
          />
        </div>

        <div className={styles.actions}>
          <Avatar
            name={createdBy?.firstName || ''}
            surname={createdBy?.lastName || ''}
            size={30}
          />
          {editNote ? (
            <IconButton
              size='xxxs'
              onClick={handleUpdateNote}
              icon={<Check style={{ transform: 'scale(0.9)' }} />}
              mode='text'
              title='Edit'
              style={{ marginBottom: 0, color: 'green' }}
            />
          ) : (
            <IconButton
              size='xxxs'
              onClick={() => handleToggleEditMode(true)}
              icon={<Pencil style={{ transform: 'scale(0.9)' }} />}
              mode='text'
              title='Edit'
              style={{ marginBottom: 0 }}
            />
          )}
        </div>
      </div>
    </div>
  );
};
