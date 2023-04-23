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
  CalendarPlus,
  CloudUpload,
  VoiceWaveRecord,
  ChevronDown,
} from '../../atoms';
import Autocomplete from 'react-google-autocomplete';
import { ContactAvatar } from '../contact-avatar/ContactAvatar';
import { ContactAutocomplete } from './components/ContactAutocomplete';
import { PreviewAttendees } from './components/PreviewAttendees';
import FileO from '../../atoms/icons/FileO';
import { TimePicker } from './components/time-picker';
import { DatePicker } from 'react-date-picker';
import {
  useUpdateMeeting,
  Meeting,
  useLinkMeetingAttachement,
  useUnlinkMeetingAttachement,
} from '../../../../hooks/useMeeting';
import { useRecoilState } from 'recoil';
import { contactNewItemsToEdit } from '../../../../state';
import { getAttendeeDataFromParticipant } from './utils';
import { MeetingParticipant } from '../../../../graphQL/__generated__/generated';
import { className } from 'jsx-dom-cjs';

interface MeetingTimelineItemProps {
  meeting: Meeting;
}

export const MeetingTimelineItem = ({
  meeting,
}: MeetingTimelineItemProps): JSX.Element => {
  const { onUpdateMeeting } = useUpdateMeeting({ meetingId: meeting.id });
  const { onLinkMeetingAttachement } = useLinkMeetingAttachement({
    meetingId: 'meeting.id',
  });
  const { onUnlinkMeetingAttachement } = useUnlinkMeetingAttachement({
    meetingId: 'meeting.id',
  });

  const [itemsInEditMode, setItemToEditMode] = useRecoilState(
    contactNewItemsToEdit,
  );

  const [editNote, setEditNote] = useState(false);
  const [editAgenda, setEditAgenda] = useState(false);

  const [files, setFiles] = useState([] as any);

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
      <section>
        <div className={styles.rangePicker}>
          <DatePicker
            onChange={(e) => {
              console.log('🏷️ ----- : ', e);
            }}
            value={meeting.start}
            calendarIcon={<CalendarPlus />}
            required={false}
          />
        </div>
      </section>

      <div className={classNames(styles.folder)}>
        <section className={styles.dateAndAvatars}>
          <TimePicker
            alignment='left'
            dateTime={meeting.start}
            label={'from'}
          />
          <section className={styles.folderTab}>
            <div className={styles.leftShape}>
              <div
                className={styles.avatars}
                // style={{ width: meeting?.attendees.length * 25 }}
              >
                {meeting?.attendedBy.map(
                  (attendeeData: MeetingParticipant, index: number) => {
                    const attendee =
                      getAttendeeDataFromParticipant(attendeeData);
                    if (meeting?.attendedBy.length > 3 && index === 3) {
                      return (
                        <PreviewAttendees
                          hiddenAttendeesNumber={meeting.attendedBy.length - 3}
                          selectedAttendees={meeting?.attendedBy.slice(index)}
                        />
                      );
                    }
                    if (index > 3) {
                      return null;
                    }

                    return (
                      <div
                        key={`${index}-${attendee.id}`}
                        className={styles.avatar}
                        style={{
                          zIndex: index,
                          left: index === 0 ? 0 : 20 * index,
                        }}
                      >
                        <ContactAvatar contactId={attendee.id} />
                      </div>
                    );
                  },
                )}

                <div className={styles.addUserButton}>
                  <ContactAutocomplete
                    selectedAttendees={meeting.attendedBy}
                    onRemoveAttendee={(attendeeId) => {
                      const newAttendeeList = meeting.attendedBy.filter(
                        (attendeeData) => {
                          const attendee =
                            getAttendeeDataFromParticipant(attendeeData);

                          return attendee.id !== attendeeId;
                        },
                      );

                      return onUpdateMeeting({
                        attendedBy: newAttendeeList,
                      });
                    }}
                    onAddAttendee={(newParticipant) => {
                      onUpdateMeeting({
                        attendedBy: [newParticipant],
                      });
                    }}
                  />
                </div>
              </div>
            </div>
          </section>
          <TimePicker alignment='right' dateTime={meeting.end} label={'to'} />
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

                return onLinkMeetingAttachement(newFile.id);
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

                return onUnlinkMeetingAttachement(fileId);
              }}
            />
          </div>
        </div>

        {/* RECORDING SECTION*/}
        <article
          className={classNames(styles.recordingSection, {
            [styles.recordingUploaded]: meeting.recoding,
          })}
        >
          <div
            className={classNames(styles.recordingCta, {
              [styles.recordingUploaded]: meeting.recoding,
            })}
          >
            <div className={styles.recordingIcon}>
              {meeting.recoding ? (
                <CloudUpload height={24} width={24} />
              ) : (
                <FileO height={24} width={24} aria-label='Meeting recording' />
              )}
            </div>
            {meeting.recoding ? (
              <span> 1h 27min 32s </span>
            ) : (
              <h3>Upload the recording</h3>
            )}
            <VoiceWaveRecord />
          </div>
        </article>

        <div className={styles.collapsibleSection}>
          <div
            className={classNames(styles.transcriptionSection, {
              [styles.recordingUploaded]: meeting.recoding,
            })}
          />
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
