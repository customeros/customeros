import React, { useCallback, useState } from 'react';
import styles from './meeting-timeline-item.module.scss';
import { extraAttributes } from '../editor/SocialEditor';
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
import classNames from 'classnames';
import { DebouncedEditor } from '../editor/DebouncedEditor';
import {
  FileUpload,
  GroupLight,
  CalendarPlus,
  EditableContentInput,
} from '../../atoms';
import { ContactAvatar } from '../contact-avatar/ContactAvatar';
import { AttendeeAutocomplete } from './components/AttendeeAutocomplete';
import { PreviewAttendees } from './components/PreviewAttendees';
import { TimePicker } from './components/time-picker';
import { DatePicker } from 'react-date-picker';
import {
  useUpdateMeeting,
  Meeting,
  useLinkMeetingAttachement,
  useUnlinkMeetingAttachement,
} from '../../../../hooks/useMeeting';
import { getAttendeeDataFromParticipant } from './utils';
import { MeetingParticipant } from '../../../../graphQL/__generated__/generated';
import { MeetingRecording } from './components/meeting-recording';
import { toast } from 'react-toastify';
import { getDate } from 'date-fns';
import { DateTimeUtils } from '../../../../utils';

interface MeetingTimelineItemProps {
  meeting: Meeting;
}

export const MeetingTimelineItem = ({
  meeting,
}: MeetingTimelineItemProps): JSX.Element => {
  const { onUpdateMeeting } = useUpdateMeeting({
    meetingId: meeting.id,
    appSource: meeting.appSource,
  });
  const { onLinkMeetingAttachement } = useLinkMeetingAttachement({
    meetingId: meeting.id,
  });
  const { onUnlinkMeetingAttachement } = useUnlinkMeetingAttachement({
    meetingId: meeting.id,
  });

  const [editNote, setEditNote] = useState(false);
  const [editAgenda, setEditAgenda] = useState(false);

  const [files, setFiles] = useState(meeting.includes || []);

  const remirrorExtentions = [
    new TableExtension(),
    new MentionAtomExtension({
      matchers: [
        { name: 'at', char: '@' },
        { name: 'tag', char: '#' },
      ],
    }),

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
    <div style={{ width: '100%' }}>
      <section>
        <div className={styles.rangePicker}>
          <DatePicker
            onChange={(e: any) => {
              const day = getDate(new Date(e));
              const month = new Date(e).getMonth();
              const year = new Date(e).getFullYear();

              const startDate = new Date(meeting.startedAt).setFullYear(
                year,
                month,
                day,
              );
              const endDate = new Date(meeting.endedAt).setFullYear(
                year,
                month,
                day,
              );

              try {
                onUpdateMeeting({
                  startedAt: new Date(startDate),
                  endedAt: new Date(endDate),
                });
                // const newDateTime = new Date(date);
              } catch (e) {
                toast.error('Invalid date selected');
              }
            }}
            value={meeting.startedAt}
            calendarIcon={<CalendarPlus />}
            required={false}
          />
        </div>
      </section>

      <div className={classNames(styles.folder)}>
        <section className={styles.dateAndAvatars}>
          <TimePicker
            alignment='left'
            dateTime={meeting.startedAt}
            label={'from'}
            onUpdateTime={(startDate) =>
              onUpdateMeeting({ startedAt: startDate })
            }
          />
          <section className={styles.folderTab}>
            <div className={styles.leftShape}>
              <div className={styles.avatars}>
                {(meeting?.attendedBy || []).map(
                  (attendeeData: MeetingParticipant, index: number) => {
                    const attendee =
                      getAttendeeDataFromParticipant(attendeeData);
                    if (meeting.attendedBy.length > 3 && index === 3) {
                      return (
                        <PreviewAttendees
                          key={`attendee-preview-hidden-item-${meeting.id}-${attendee.id}`}
                          hiddenAttendeesNumber={meeting.attendedBy.length - 3}
                          selectedAttendees={meeting.attendedBy.slice(index)}
                        />
                      );
                    }
                    if (index > 3) {
                      return null;
                    }

                    return (
                      <div
                        key={`attendee-preview-item-${meeting.id}-${attendee.id}-${index}`}
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
                  <AttendeeAutocomplete
                    meetingId={meeting.id}
                    selectedAttendees={meeting.attendedBy || []}
                  />
                </div>
              </div>
            </div>
          </section>
          <TimePicker
            alignment='right'
            dateTime={meeting.endedAt}
            label={'to'}
            onUpdateTime={(endDate) => onUpdateMeeting({ endedAt: endDate })}
          />
        </section>

        <div
          className={classNames(styles.editableMeetingProperties, {
            [styles.draftMode]: DateTimeUtils.isBeforeNow(meeting.startedAt),
            [styles.pastMode]: !DateTimeUtils.isBeforeNow(meeting.startedAt),
          })}
        >
          <div className={styles.contentWithBorderWrapper}>
            <section className={styles.meetingLocationSection}>
              {/*<div*/}
              {/*  className={classNames(styles.meetingLocation, {*/}
              {/*    [styles.selectedMeetingLocation]: meeting?.conferenceUrl,*/}
              {/*  })}*/}
              {/*>*/}
              {/*  <span className={styles.meetingLocationLabel}>meeting at</span>*/}
              {/*  <div className={styles.meetingLocationInputWrapper}>*/}
              {/*    <PinAltLight />*/}
              {/*    <Autocomplete*/}
              {/*      className={styles.meetingInput}*/}
              {/*      placeholder={'Add meeting location'}*/}
              {/*      apiKey={process.env.GOOGLE_MAPS_API_KEY}*/}
              {/*      onPlaceSelected={(place, inputRef, autocomplete) => {*/}
              {/*        console.log(autocomplete);*/}
              {/*      }}*/}
              {/*    />*/}
              {/*  </div>*/}
              {/*</div>*/}

              <div
                className={classNames(styles.meetingLocation, {
                  [styles.selectedMeetingLocation]: meeting?.conferenceUrl,
                })}
              >
                <span className={styles.meetingLocationLabel}>meeting at</span>
                <div className={styles.meetingLocationInputWrapper}>
                  <GroupLight
                    color={meeting?.conferenceUrl ? '#fff' : '#878787'}
                  />
                  <EditableContentInput
                    isEditMode
                    className={styles.meetingInput}
                    value={meeting?.conferenceUrl || ''}
                    onChange={(data: string) =>
                      onUpdateMeeting({
                        conferenceUrl: data,
                      })
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
                onToggleEditMode={(newMode: boolean) => {
                  setEditAgenda(newMode);
                }}
                manager={manager}
                state={state}
                setState={setState}
                context={getContext()}
                onDebouncedSave={(data: string) =>
                  onUpdateMeeting({
                    agenda: data,
                  })
                }
              />
            </section>

            <section className={styles.meetingNoteSection}>
              <DebouncedEditor
                isEditMode={editNote}
                className={classNames(styles.meetingNoteWrapper, {
                  [styles.readMode]: !editNote,
                })}
                value={meeting?.note?.html || 'NOTES:'}
                manager={manager}
                state={state}
                onToggleEditMode={(newMode: boolean) => {
                  setEditNote(newMode);
                }}
                setState={setState}
                context={getContext()}
                onDebouncedSave={(data: string) => {

                  return onUpdateMeeting({
                    note: { html: data, appSource: 'OPENLINE' },
                  });
                }}
              />
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

        <MeetingRecording
          meeting={meeting}
          onUpdateMeetingRecording={(id) => onUpdateMeeting({ recording: id })}
        />
      </div>
    </div>
  );
};
