import React, { useCallback, useState } from 'react';

import styles from './meeting-timeline-item.module.scss';
import DateTimePicker from 'react-datetime-picker';
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
  Input,
  Paperclip,
  Pencil,
  PinAltLight,
  UserEdit,
} from '../../atoms';
import Autocomplete from 'react-google-autocomplete';
import { ContactAvatar } from '../contact-avatar/ContactAvatar';
import { ContactAutocomplete } from './components/ContactAutocomplete';

interface MeetingTimelineItemProps {}

export const MeetingTimelineItem = ({ meeting }: any): JSX.Element => {
  const [value, onChange] = useState(new Date());
  const [editNote, setEditNote] = useState(false);
  const [editAgenda, setEditAgenda] = useState(false);
  const [attendeesDropdownOpen, setAttendeesDropdownOpen] = useState(false);

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
    <div className={styles.folder}>
      <section className={styles.dateAndAvatars}>
        <div className={styles.date}>
          {/*<DateTimePicker onChange={onChange} value={value} />*/}
        </div>
        <section className={styles.folderTab}>
          <div className={styles.leftShape}>
            <div
              className={styles.avatars}
              // style={{ width: meeting?.attendees.length * 25 }}
            >
              {meeting?.attendees.map(({ id }, index) => {
                if (meeting?.attendees.length > 3 && index === 3) {
                  return (
                    <div className={styles.numberAvatar}>
                      {' '}
                      +{meeting.attendees.length - 3}
                    </div>
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
                <ContactAutocomplete />
              </div>
            </div>
          </div>
        </section>

        <div className={styles.date}>
          {/*<DateTimePicker onChange={onChange} value={value} />*/}
        </div>
      </section>

      <div className={styles.bottomPart}>
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
                  // style={{ width: '12ch' }}
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
                  onChange={() => null}
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
              onDebouncedSave={(data) => console.log('data', data)}
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
              onDebouncedSave={(data) => console.log('data', data)}
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

          <FileUpload onFileUpload={() => null} />
        </div>
      </div>
    </div>
  );
};
