import React, { useEffect, useState } from 'react';
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
} from '../../atoms';
import sanitizeHtml from 'sanitize-html';

interface Props {
  noteContent: string;
  createdAt: string;
  contactId?: string;
  id: string;
  refreshNoteData: (id: string) => void;
  createdBy?: {
    firstName?: string;
    lastName?: string;
  };
  source?: string;
}

export const NoteTimelineItem: React.FC<Props> = ({
  noteContent,
  id,
  createdBy,
  contactId,
  refreshNoteData,
  source,
}) => {
  // const client =  useGraphQLClient();
  const [images, setImages] = useState({});
  const [editNote, setEditNote] = useState(false);

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

  const [deleteConfirmationModalVisible, setDeleteConfirmationModalVisible] =
    useState(false);
  const deleteNote = () => {
    // DeleteNote(client, id)
    //   .then((result: boolean) => {
    //     if (result) {
    //       refreshNoteData(id);
    //       setDeleteConfirmationModalVisible(false);
    //       toast.success('Note removed successfully!');
    //     } else {
    //       toast.error(
    //         'There was a problem on our side and we are doing our best to solve it!',
    //       );
    //     }
    //   })
    //   .catch((reason: any) => {
    //     toast.error(
    //       'There was a problem on our side and we are doing our best to solve it!',
    //     );
    //   });
  };

  return (
    <>
      {/*{editNote && (*/}
      {/*  <ContactNoteModalTemplate*/}
      {/*    isEdit*/}
      {/*    note={note}*/}
      {/*    contactId={contactId}*/}
      {/*    notifyChanged={(data) => {*/}
      {/*      setEditNote(false);*/}
      {/*      refreshNoteData(data);*/}
      {/*    }}*/}
      {/*    notifyCancel={() => setEditNote(false)}*/}
      {/*  />*/}
      {/*)}*/}
      <DeleteConfirmationDialog
        deleteConfirmationModalVisible={deleteConfirmationModalVisible}
        setDeleteConfirmationModalVisible={setDeleteConfirmationModalVisible}
        deleteAction={deleteNote}
        confirmationButtonLabel='Delete note'
      />

      {!editNote && (
        <div className='flex justify-content-between'>
          <div className={styles.noteContainer}>
            <div
              className={`${styles.noteContent}`}
              dangerouslySetInnerHTML={{
                __html: sanitizeHtml(note.htmlEnhanced),
              }}
            ></div>
          </div>
          <div className={styles.actionContainer}>
            <div className={styles.actions}>
              <IconButton
                onClick={() => setDeleteConfirmationModalVisible(true)}
                icon={<Trash />}
                mode='secondary'
                title='Delete'
                style={{ marginRight: 0, marginBottom: '8px', height: '1rem' }}
              />

              <IconButton
                onClick={() => setEditNote(true)}
                icon={<Pencil />}
                mode='secondary'
                title='Edit'
                style={{ marginRight: 0, height: '1rem' }}
              />
            </div>
            <div className={styles.noteData}>
              <div className='text-sm text-gray-600'>
                {(createdBy?.firstName || createdBy?.lastName) && '- '}
                {createdBy?.firstName} {createdBy?.lastName}
              </div>

              {source && (
                <div className='flex text-sm text-gray-600 '>
                  <div className='mr-2'>Source:</div>
                  <div className='capitaliseFirstLetter'>{source}</div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </>
  );
};
