import React, { useCallback, useState } from 'react';

import styles from './meeting-timeline-item.module.scss';
import { extraAttributes } from '../editor/SocialEditor';
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
import classNames from 'classnames';
import { DebouncedEditor } from '../editor/DebouncedEditor';
import {
  DebouncedInput,
  FileUpload,
  GroupLight,
  IconButton,
  Pencil,
  PinAltLight,
  Upload,
  VoiceWave,
  CalendarPlus,
  CloudUpload,
  VoiceWaveRecord,
  ChevronDown,
} from '../../atoms';
import Autocomplete from 'react-google-autocomplete';
import { ContactAvatar } from '../contact-avatar/ContactAvatar';
import { ContactAutocomplete } from './components/ContactAutocomplete';
import { PreviewAttendees } from './components/EditAttendees';
import { DateTimeUtils } from '../../../../utils';
import { DateRangePicker } from '@wojtekmaj/react-daterange-picker';
import FileO from '../../atoms/icons/FileO';
import Timekeeper from 'react-timekeeper';
import { TimePicker } from './components/time-picker';
import { DatePicker } from 'react-date-picker';
import {
  useCreateMeetingFromContact,
  useUpdateMeeting,
} from '../../../../hooks/useMeeting';
import { useRecoilState } from 'recoil';
import { contactNewItemsToEdit } from '../../../../state';

interface MeetingTimelineItemProps {}

export const MeetingTimelineItem = ({ meeting }: any): JSX.Element => {
  const { onUpdateMeeting } = useUpdateMeeting({ meetingId: 'meeting.id' });

  const [itemsInEditMode, setItemToEditMode] = useRecoilState(
    contactNewItemsToEdit,
  );

  const [value, onChange] = useState([new Date(), new Date()]);

  const [editNote, setEditNote] = useState(false);
  const [editAgenda, setEditAgenda] = useState(false);
  const [attendeesDropdownOpen, setAttendeesDropdownOpen] = useState(false);

  const [files, setFiles] = useState([] as any);
  const [fileIdsToAdd, setFileIdsToAdd] = useState([] as any); //HERE ARE THE attachments ID to save
  const [fileIdsToRemove, setFileIdsToRemove] = useState([] as any); //HERE ARE THE attachments ID to remove from the meeting

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
    new ImageExtension({}),
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
  const extensions = useCallback(() => [...remirrorExtentions], []);

  const { manager, state, setState, getContext } = useRemirror({
    extensions,
    extraAttributes,
    // state can created from a html string.
    stringHandler: 'html',

    // This content is used to create the initial value. It is never referred to again after the first render.
    content: '',
  });
  return (
    <div>
      <section className={styles.timeTickSection}>
        <div className={styles.rangePicker}>
          <DatePicker
            onChange={onChange}
            value={value}
            calendarIcon={<CalendarPlus />}
            required={false}
          />
        </div>
      </section>

      <div className={classNames(styles.folder)}>
        <section className={styles.dateAndAvatars}>
          <TimePicker alignment='left' dateTime={new Date()} label={'from'} />
          <section className={styles.folderTab}>
            <div className={styles.leftShape}>
              <div
                className={styles.avatars}
                // style={{ width: meeting?.attendees.length * 25 }}
              >
                {meeting?.attendees.map(({ id }: any, index: number) => {
                  if (meeting?.attendees.length > 3 && index === 3) {
                    return (
                      <PreviewAttendees
                        hiddenContactsNumber={meeting.attendees.length - 3}
                        selectedContacts={meeting?.attendees.slice(index)}
                      />
                    );
                  }
                  if (index > 3) {
                    return null;
                  }

                  return (
                    <div
                      key={`${index}-${id}`}
                      className={styles.avatar}
                      style={{
                        zIndex: index,
                        left: index === 0 ? 0 : 20 * index,
                      }}
                    >
                      <ContactAvatar contactId={id} />
                    </div>
                  );
                })}

                <div className={styles.addUserButton}>
                  <ContactAutocomplete
                    selectedContacts={meeting.attendees}
                    onSelectContact={(data) =>
                      onUpdateMeeting({
                        attendedBy: [...meeting.attendedBy, data],
                      })
                    }
                  />
                </div>
              </div>
            </div>
          </section>
          <TimePicker alignment='right' dateTime={new Date()} label={'to'} />
        </section>

        <div
          className={classNames(styles.editableMeetingProperties, {
            [styles.draftMode]: true,
            [styles.pastMode]: false,
          })}
        >
          <div className={styles.contentWithBorderWrapper}>
            <section className={styles.meetingLocationSection}>
              <div
                className={classNames(styles.meetingLocation, {
                  [styles.selectedMeetingLocation]: false,
                })}
              >
                <span className={styles.meetingLocationLabel}>meeting at</span>
                <div className={styles.meetingLocationInputWrapper}>
                  <PinAltLight />
                  <Autocomplete
                    className={styles.meetingInput}
                    placeholder={'Add meeting location'}
                    apiKey={process.env.GOOGLE_MAPS_API_KEY}
                    onPlaceSelected={(place, inputRef, autocomplete) => {
                      console.log(autocomplete);
                    }}
                  />
                </div>
              </div>

              <div
                className={classNames(styles.meetingLocation, {
                  [styles.selectedMeetingLocation]: false,
                })}
              >
                <span className={styles.meetingLocationLabel}>meeting at</span>
                <div className={styles.meetingLocationInputWrapper}>
                  <GroupLight />
                  <DebouncedInput
                    className={styles.meetingInput}
                    inlineMode
                    onChange={(event) =>
                      onUpdateMeeting({ location: event.target.value })
                    }
                    placeholder='Add conference link'
                  />
                </div>
              </div>
            </section>

            <section className={styles.agenda}>
              <div className={styles.agendaTitleSection}>
                <h3 className={styles.agendaTitle}>Agenda</h3>
              </div>
              <DebouncedEditor
                value={`
           <p>INTRODUCTION</p>
           <p>DISCUSSION</p>
           <p>NEXT STEPS</p>
          `}
                className={classNames({
                  [styles.readMode]: !editAgenda,
                  [styles.editorEditMode]: editAgenda,
                })}
                isEditMode={editAgenda}
                manager={manager}
                state={state}
                setState={setState}
                context={getContext()}
                onDebouncedSave={(data: string) =>
                  onUpdateMeeting({ agenda: data })
                }
              />

              {!editAgenda && (
                <IconButton
                  className={styles.editButton}
                  mode={'text'}
                  onClick={() => setEditAgenda(true)}
                  icon={<Pencil style={{ transform: 'scale(0.8)' }} />}
                />
              )}
            </section>

            <section className={styles.meetingNoteSection}>
              <DebouncedEditor
                isEditMode={editNote}
                className={classNames(styles.meetingNoteWrapper, {
                  [styles.readMode]: !editNote,
                })}
                value={'NOTES:'}
                manager={manager}
                state={state}
                setState={setState}
                context={getContext()}
                onDebouncedSave={(data: string) =>
                  onUpdateMeeting({ note: { html: data } })
                }
              />
              {!editNote && (
                <IconButton
                  className={styles.editButton}
                  mode='text'
                  onClick={() => setEditNote(true)}
                  icon={<Pencil style={{ transform: 'scale(0.8)' }} />}
                />
              )}
            </section>

            <FileUpload
              files={files}
              onBeginFileUpload={(fileKey: string) => {
                setFiles((prevFiles: any) => [
                  ...prevFiles,
                  {
                    key: fileKey,
                    uploaded: false,
                  },
                ]);
              }}
              onFileUpload={(newFile: any) => {
                setFiles((prevFiles: any) => {
                  return prevFiles.map((file: any) => {
                    if (file.key === newFile.key) {
                      file = {
                        id: newFile.id,
                        key: newFile.key,
                        name: newFile.name,
                        extension: newFile.extension,
                        uploaded: true,
                      };
                    }
                    return file;
                  });
                });
                setFileIdsToAdd((prevFileIdsToAdd: any) => [
                  ...prevFileIdsToAdd,
                  newFile.id,
                ]);
              }}
              onFileUploadError={(fileKey: any) => {
                setFiles((prevFiles: any) => {
                  // TODO do not remove the file from the list
                  // show the error instead for that particular file
                  return prevFiles.filter((file: any) => file.key !== fileKey);
                });
              }}
              onFileRemove={(fileId: any) => {
                setFiles((prevFiles: any) => {
                  return prevFiles.filter((file: any) => file.id !== fileId);
                });
                setFileIdsToRemove((prevFileIdsToRemove: any) => [
                  ...prevFileIdsToRemove,
                  fileId,
                ]);
              }}
            />
          </div>
        </div>

        {/* RECORDING SECTION*/}
        <article
          className={classNames(styles.recordingSection, {
            [styles.recordingUploaded]: false,
          })}
        >
          <div
            className={classNames(styles.recordingCta, {
              [styles.fileUploaded]: false,
            })}
          >
            <div className={styles.recordingIcon}>
              {meeting.recording ? (
                <CloudUpload height={24} width={24} />
              ) : (
                <FileO height={24} width={24} aria-label='Meeting recording' />
              )}
            </div>
            {false ? (
              <span> 1h 27min 32s </span>
            ) : (
              <h3>Upload the recording</h3>
            )}
            <VoiceWaveRecord />
          </div>
        </article>

        <div className={styles.collapsibleSection}>
          <div className={styles.transcriptionSection} />
          <div className={styles.collapseExpandButtonWrapper}>
            <IconButton
              className={styles.collapseExpandButton}
              isSquare
              disabled
              mode='secondary'
              size='xxxxs'
              icon={<ChevronDown width={24} height={24} />}
              onClick={() => console.log('collapse / expand button click')}
            />
          </div>
        </div>
      </div>
    </div>
  );
};
