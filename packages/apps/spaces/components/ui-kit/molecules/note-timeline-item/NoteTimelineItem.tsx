import React, {
  MutableRefObject,
  Ref,
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
import { Controller, useForm } from 'react-hook-form';
import { Editor } from '../editor';
import { NoteEditorModes } from '../editor/Editor';
import { ContactAvatar } from '../organization-avatar';
import { NotedEntity } from '../../../../graphQL/__generated__/generated';
import { OrganizationAvatar } from '../organization-avatar/OrganizationAvatar';
import { getContactDisplayName } from '../../../../utils';

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
  const [editNote, setEditNote] = useState(false);
  const elementRef = useRef<MutableRefObject<Ref<HTMLDivElement>>>(null);

  const [note, setNote] = useState({
    id,
    html: noteContent,
    htmlEnhanced: noteContent,
  });

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
              .catch((reason: any) => {
                toast.error(
                  'There was a problem on our side and we are doing our best to solve it!',
                );
              });
          }
        },
      });
    } else {
      setNote({ id, html: noteContent, htmlEnhanced: noteContent });
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

      setNote({
        id,
        html: noteContent,
        htmlEnhanced: html,
      });
    }
  }, [id, images, noteContent]);

  const { handleSubmit, setValue, getValues, control, reset } = useForm({
    defaultValues: {
      id: note?.id || '',
      html: note?.html || '',
      htmlEnhanced: note.htmlEnhanced || '',
    },
  });

  const onSubmit = handleSubmit(({ htmlEnhanced, ...data }) => {
    const dataToSubmit = {
      ...data,
      html: htmlEnhanced?.replaceAll(/.src(\S*)/g, ''), //remove src attribute to not send the file bytes in here
    };
    onUpdateNote(dataToSubmit).then(() => {
      setEditNote(false);
    });
  });

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
  console.log('🏷️ ----- no: ', noted);
  return (
    <div
      className={styles.noteWrapper}
      //@ts-expect-error fixme
      ref={elementRef}
    >
      <div className={styles.noteContainer}>
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
        {editNote && (
          <div className={styles.noteContent}>
            <div
              style={{
                display: 'flex',
                flexDirection: 'column',
              }}
            >
              <Controller
                name='htmlEnhanced'
                control={control}
                render={({ field }) => (
                  <Editor
                    onCancel={() => null} // not used
                    mode={NoteEditorModes.EDIT}
                    onGetFieldValue={getValues}
                    value={field.value}
                    onSave={() => null} //not used
                    label='Save'
                    onTextChange={(e) => setValue('htmlEnhanced', e.htmlValue)}
                  />
                )}
              />
            </div>

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
        )}

        {!editNote && (
          <div
            className={styles.noteContent}
            dangerouslySetInnerHTML={{
              __html: sanitizeHtml(
                linkifyHtml(note.htmlEnhanced, {
                  defaultProtocol: 'https',
                  rel: 'noopener noreferrer',
                }),
              ),
            }}
          />
        )}
        <div className={styles.actions}>
          <Avatar
            name={createdBy?.firstName || ''}
            surname={createdBy?.lastName || ''}
            size={30}
          />
          {editNote ? (
            <IconButton
              size='xxxs'
              onClick={onSubmit}
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
