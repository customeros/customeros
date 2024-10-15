import { useSearchParams } from 'react-router-dom';
import { useRef, Fragment, useState } from 'react';

import set from 'lodash/set';
import { useKeyBindings } from 'rooks';
import cityTimezone from 'city-timezones';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { flags } from '@ui/media/flags';
import { Avatar } from '@ui/media/Avatar';
import { Plus } from '@ui/media/icons/Plus';
import { Tag01 } from '@ui/media/icons/Tag01';
import { Spinner } from '@ui/feedback/Spinner';
import { Mail02 } from '@ui/media/icons/Mail02';
import { Star06 } from '@ui/media/icons/Star06';
import { Button } from '@ui/form/Button/Button';
import { getTimezone } from '@utils/getTimezone';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive';
import { Tags } from '@organization/components/Tabs';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { TextInput } from '@ui/media/icons/TextInput';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { getFormattedLink } from '@utils/getExternalLink';
import { Tag, Social, TableViewType } from '@graphql/types';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { LinkedInSolid02 } from '@ui/media/icons/LinkedInSolid02';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { EmailValidationMessage } from '@organization/components/Tabs/panels/PeoplePanel/ContactCard/EmailValidationMessage';
import {
  Modal,
  ModalBody,
  ModalClose,
  ModalPortal,
  ModalFooter,
  ModalOverlay,
  ModalCloseButton,
  ModalFeaturedHeader,
  ModalFeaturedContent,
} from '@ui/overlay/Modal';

export const ContactPreviewCard = observer(() => {
  const store = useStore();
  const [searchParams] = useSearchParams();
  const [isEditName, setIsEditName] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const contactId = store.ui.focusRow;
  const preset = searchParams?.get('preset');
  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const tableType = tableViewDef?.value?.tableType;
  const [isOpen, setIsOpen] = useState(false);

  if (tableType !== TableViewType.Contacts && !contactId) {
    store.ui.setContactPreviewCardOpen(false);

    return null;
  }

  if (!contactId) return;

  const contact = store.contacts
    .toArray()
    .find((c) => c.id === String(contactId));

  const fullName = contact?.name || 'Unnamed';
  const src = contact?.value?.profilePhotoUrl;
  const activeCompany =
    (contact?.value?.organizations?.content?.length ?? 1) - 1;
  const company = contact?.value.organizations?.content?.[activeCompany]?.name;

  const role = contact?.value.jobRoles?.[0]?.jobTitle;
  const countryA3 = contact?.value.locations?.[0]?.countryCodeA3;
  const countryA2 = contact?.value.locations?.[0]?.countryCodeA2;
  const flag = flags[countryA2 || ''];
  const city = contact?.value.locations?.[0]?.locality;
  const timezone = city
    ? cityTimezone.lookupViaCity(city)?.[0]?.timezone
    : null;

  const validationDetails = contact?.value.emails?.[0]?.emailValidationDetails;

  const fromatedUrl = contact?.value?.socials?.[0]?.url.replace(
    'https://www.',
    '',
  );

  const formatedFollowersCount = contact?.value?.socials?.[0]?.followersCount
    ?.toLocaleString()
    .replace(/\B(?=(\d{3})+(?!\d))/g, ',');

  const handleUpdateSocial = (url: string) => {
    const linkedinId = contact?.value.socials.find((social) =>
      social.url.includes('linkedin'),
    )?.id;

    if (fromatedUrl === undefined && url.trim() !== '') {
      contact?.update((contactData) => {
        const formattedValue =
          url.includes('https://www') || url.includes('linkedin.com')
            ? getFormattedLink(url).replace(/^linkedin\.com\//, '')
            : url;

        contactData.socials.push({
          id: crypto.randomUUID(),
          url: `linkedin.com/${formattedValue}`,
        } as Social);

        return contactData;
      });
    }

    contact?.update((org) => {
      const idx = org.socials.findIndex((s) => s.id === linkedinId);
      const formattedValue =
        url.includes('https://www') || url.includes('linkedin.com')
          ? getFormattedLink(url).replace(/^linkedin\.com\//, '')
          : `in/${url}`;

      if (idx !== -1) {
        org.socials[idx].url = `linkedin.com/${formattedValue}`;
      }

      if (url === '') {
        org.socials.splice(idx, 1);
      }

      return org;
    });
  };

  useKeyBindings(
    {
      Escape: () => {
        store.ui.setContactPreviewCardOpen(false);
      },
      Space: (e) => {
        e.preventDefault();
        store.ui.setContactPreviewCardOpen(false);
      },
    },
    {
      when: store.ui.contactPreviewCardOpen,
    },
  );

  const userHasOrg = contact?.organization?.name;

  const userBeenEnriched =
    contact?.value?.enrichDetails.enrichedAt ||
    contact?.value?.enrichDetails.failedAt;

  const requestedEnrichment = contact?.value?.enrichDetails.requestedAt;

  return (
    <>
      {store.ui.contactPreviewCardOpen && (
        <div
          data-state={store.ui.contactPreviewCardOpen ? 'open' : 'closed'}
          className='data-[state=open]:animate-slideLeftAndFade data-[state=closed]:animate-slideRightAndFade flex flex-col absolute right-[12px] -top-[-53px] p-4 max-w-[390px] min-w-[350px] border border-gray-200 rounded-lg z-50 bg-white'
        >
          <div className='flex justify-between items-start'>
            <Avatar
              size='sm'
              textSize='xs'
              name={fullName}
              variant='circle'
              src={src || undefined}
            />
            {!userBeenEnriched && (
              <Tooltip label='Enrich this contact'>
                {!requestedEnrichment ? (
                  <IconButton
                    size='xxs'
                    icon={<Star06 />}
                    onClick={() => setIsOpen(true)}
                    aria-label='enrich this contact'
                  />
                ) : (
                  <Spinner
                    size='sm'
                    label='enriching'
                    className='text-gray-400 fill-gray-700'
                  />
                )}
              </Tooltip>
            )}
          </div>
          <div className='flex items-center gap-1'>
            {isEditName ? (
              <Input
                value={fullName}
                variant='unstyled'
                placeholder='Unknown'
                className='mb-[-8px]'
                onFocus={(e) => e.target.select()}
                onBlur={() => setIsEditName(false)}
                onChange={(e) => {
                  contact?.update((value) => {
                    set(value, 'name', e.target.value);

                    return value;
                  });
                }}
              />
            ) : (
              <>
                <span
                  onClick={() => setIsEditName(true)}
                  className='font-medium mt-2 overflow-hidden text-ellipsis whitespace-nowrap'
                >
                  {fullName}
                </span>
                {company ? (
                  <>
                    <span className='mt-2 text-gray-500'> at </span>
                    <span className='font-medium overflow-hidden text-ellipsis whitespace-nowrap mt-2 '>
                      {company || 'No org yet'}
                    </span>
                  </>
                ) : (
                  <span className='mt-2'>(No org yet)</span>
                )}
              </>
            )}
          </div>
          <Input
            size='xs'
            variant='unstyled'
            value={role || ''}
            placeholder='Enter title'
            onFocus={(e) => e.target.select()}
            className='w-[290px] overflow-hidden text-ellipsis whitespace-nowrap'
            onChange={(e) => {
              contact?.update((value) => {
                set(value, 'jobRoles[0].jobTitle', e.target.value);

                return value;
              });
            }}
          />
          <div className={cn('flex items-center mb-4', countryA3 && 'gap-1')}>
            <span className='mb-1'>{flag}</span>
            {countryA3 && <span className='ml-2 text-sm'>{countryA3}</span>}
            {countryA3 && city && timezone && <span>•</span>}
            {city && (
              <span className='overflow-hidden text-ellipsis whitespace-nowrap text-sm'>
                {city}
              </span>
            )}
            {city && timezone && <span>•</span>}
            {timezone && (
              <span className='w-[150px] text-sm'>
                {getTimezone(timezone || '')} local time
              </span>
            )}
          </div>
          <div className='flex justify-between gap-1 w-full mb-4 flex-col'>
            <div className='flex items-center justify-between w-full text-sm group/menu'>
              <div className='flex items-center gap-2'>
                <Mail02 className='mt-[1px] text-gray-500' />
                <span className='text-gray-500'>Emails</span>
              </div>
              {company && (
                <div className='flex items-center gap-2'>
                  <Menu>
                    <MenuButton>
                      <Tooltip
                        align='end'
                        side='bottom'
                        label={'Add new email'}
                      >
                        <IconButton
                          size='xxs'
                          variant='ghost'
                          icon={<Plus />}
                          aria-label='add new email'
                          className='group-hover/menu:opacity-100 opacity-0'
                        />
                      </Tooltip>
                    </MenuButton>
                    <MenuList>
                      <MenuItem
                        className='group/find-email'
                        onClick={() => {
                          setIsLoading(true);
                          contact
                            ?.findEmail()
                            .finally(() => setIsLoading(false));
                        }}
                      >
                        <div className='flex items-center gap-1'>
                          <Star06 className='group-hover/find-email:text-gray-700 text-gray-500' />
                          <span>{`Find email at ${company}`}</span>
                        </div>
                      </MenuItem>
                      <MenuItem
                        className='group/add-email'
                        onClick={() => {
                          store.ui.setSelectionId(contact.value.emails.length);
                          contact.update(
                            (c) => {
                              c.emails.push({
                                id: crypto.randomUUID(),
                                email: '',
                                appSource: '',
                                contacts: [],
                                createdAt: new Date().toISOString(),
                                updatedAt: new Date().toISOString(),
                                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                              } as any);

                              return c;
                            },
                            { mutate: false },
                          );
                          store.ui.commandMenu.setContext({
                            ids: [contact.value.id],
                            entity: 'Contact',
                            property: 'email',
                          });
                          store.ui.commandMenu.setType('EditEmail');
                          store.ui.commandMenu.setOpen(true);
                        }}
                      >
                        <PlusCircle className='text-gray-500 group-hover/add-email:text-gray-700' />
                        Add new email
                      </MenuItem>
                    </MenuList>
                  </Menu>
                  {isLoading && (
                    <Tooltip label={`Finding email at ${userHasOrg} `}>
                      <Spinner
                        size='sm'
                        label='finding email'
                        className='text-gray-400 fill-gray-700'
                      />
                    </Tooltip>
                  )}
                </div>
              )}
            </div>
            <div className='ml-6'>
              {contact?.value.emails.map((email, idx) => (
                <Fragment key={email.id}>
                  <div className=' flex items-center justify-between group/menu-email'>
                    <div key={email.id} className='flex items-center '>
                      <span>{email.email}</span>
                    </div>
                    <div className='flex items-center'>
                      {email && (
                        <EmailValidationMessage
                          email={email?.email || ''}
                          validationDetails={validationDetails}
                        />
                      )}
                      <Menu>
                        <MenuButton>
                          <IconButton
                            size='xxs'
                            variant='ghost'
                            icon={<DotsVertical />}
                            aria-label='add new email'
                            className='group-hover/menu-email:opacity-100 opacity-0'
                          />
                        </MenuButton>
                        <MenuList>
                          <MenuItem
                            className='group/edit-email'
                            onClick={() => {
                              store.ui.setSelectionId(idx);
                              store.ui.commandMenu.setType('EditEmail');
                              store.ui.commandMenu.setContext({
                                ids: [contact.value.id],
                                entity: 'Contact',
                                property: 'email',
                              });
                              store.ui.commandMenu.setOpen(true);
                            }}
                          >
                            <div className='flex items-center gap-2'>
                              <TextInput className='group-hover/edit-email:text-gray-700 text-gray-500' />
                              <span>Edit email</span>
                            </div>
                          </MenuItem>
                          <MenuItem
                            className='group/archive-email'
                            onClick={() => {
                              contact.update(
                                (c) => {
                                  c.emails.splice(idx, 1);

                                  return c;
                                },
                                { mutate: false },
                              );
                              contact.updateEmail(email?.email ?? '');
                            }}
                          >
                            <Archive className='text-gray-500 group-hover/archive-email:text-gray-700' />
                            Archive email
                          </MenuItem>
                        </MenuList>
                      </Menu>
                    </div>
                  </div>
                </Fragment>
              ))}
            </div>
          </div>

          <div className='flex justify-between gap-1 w-full mb-4'>
            <div className='flex items-center gap-2 mr-[52px] text-sm text-gray-500'>
              <Tag01 className='mt-[1px] text-gray-500' />
              Persona
            </div>
            <Tags
              hideBorder
              placeholder='No tags yet'
              value={
                contact?.value?.tags?.map((tag: Tag) => ({
                  value: tag.id,
                  label: tag.name,
                })) || []
              }
              onChange={(e) => {
                contact?.update((c) => {
                  c.tags =
                    (e
                      .map((tag) => store.tags?.value.get(tag.value)?.value)
                      .filter(Boolean) as Array<Tag>) ?? [];

                  return c;
                });
              }}
            />
          </div>
          <div className='flex flex-col gap-4'>
            <div className='flex gap- items-center w-full '>
              <div className='flex items-center gap-2 mr-[52px] text-sm text-gray-500'>
                <LinkedInSolid02 className='mt-[1px] text-gray-500 ' />
                LinkedIn
              </div>

              <Input
                size='xs'
                variant='unstyled'
                value={fromatedUrl}
                className='text-ellipsis'
                onFocus={(e) => e.target.select()}
                placeholder='LinkedIn profile link'
                onChange={(e) => {
                  handleUpdateSocial(e.target.value);
                }}
              />
            </div>
            <div className='flex gap-1 w-full'>
              <div className='flex items-center gap-2 mr-[42px] text-sm text-gray-500 '>
                <LinkedInSolid02 className='mt-[1px] text-gray-500' />
                Followers
              </div>
              <span
                className={cn(
                  'overflow-hidden text-ellipsis whitespace-nowrap cursor-not-allowed text-sm',
                  formatedFollowersCount ? 'text-gray-700' : 'text-gray-400',
                )}
              >
                {formatedFollowersCount || 'Unknown'}
              </span>
            </div>
            <div className='flex gap-1 w-full mt-[2px]'>
              <div className='flex items-center gap-2 mr-[19px] text-sm text-gray-500'>
                <LinkedInSolid02 className='mt-[1px] text-gray-500' />
                Connected to
              </div>
              <span
                className={cn(
                  'overflow-hidden text-ellipsis whitespace-nowrap cursor-not-allowed text-sm',
                  contact?.value?.connectedUsers?.[0]?.name
                    ? 'text-gray-700'
                    : 'text-gray-400',
                )}
              >
                {contact?.value?.connectedUsers?.[0]?.name || 'No one yet'}
              </span>
            </div>
          </div>
        </div>
      )}
      <EnrichContactModal
        isModalOpen={isOpen}
        contactId={contactId}
        onClose={() => {
          setIsOpen(false);
        }}
      />
    </>
  );
});

const EnrichContactModal = observer(
  ({
    isModalOpen = false,
    onClose,
    contactId,
  }: {
    onClose: () => void;
    isModalOpen: boolean;
    contactId: string | number;
  }) => {
    const store = useStore();
    const hasSubmitedRef = useRef(false);
    const [linkedin, setLinkedin] = useState(
      () => store.contacts.value.get(String(contactId))?.value.socials[0]?.url,
    );
    const [validation, setValidation] = useState<Record<'linkedin', boolean>>({
      linkedin: false,
    });

    const contactStore = store.contacts.value.get(String(contactId));

    const validate = () => {
      setValidation(() => ({
        linkedin: !linkedin,
      }));

      return linkedin;
    };

    const reset = () => {
      setLinkedin('');
      setValidation({
        linkedin: false,
      });
      hasSubmitedRef.current = false;
    };

    const handleSubmit = () => {
      hasSubmitedRef.current = true;

      if (!validate()) return;

      contactStore?.addSocial(linkedin || '', {
        onSuccess: () => {
          onClose();
          reset();
        },
      });
    };

    return (
      <Modal
        open={isModalOpen}
        onOpenChange={(open) => {
          if (!open) {
            reset();
            onClose();
          }
        }}
      >
        <ModalPortal>
          <ModalOverlay className='z-[999]' />
          <ModalFeaturedContent className='z-[9999]'>
            <ModalFeaturedHeader>
              <p className='text-lg font-semibold mb-1'>
                What’s this contact’s LinkedIn?
              </p>
              <p className='text-sm'>
                To enrich this contact, we need their LinkedIn URL
              </p>
            </ModalFeaturedHeader>
            <ModalCloseButton />
            <ModalBody className='flex flex-col gap-4'>
              <div className='flex flex-col'>
                <Input
                  id='linkedin'
                  value={linkedin}
                  placeholder='LinkedIn profile link'
                  className={cn(validation.linkedin && 'border-error-500')}
                  onChange={(e) => {
                    setLinkedin(e.target.value);
                  }}
                  onKeyDown={(e) => {
                    if (e.key === 'Escape') {
                      onClose(); // Close on escape key
                    }
                    e.stopPropagation();
                  }}
                />
                {validation.linkedin && (
                  <p className='text-sm text-error-500 mt-1'>
                    One does not simply skip LinkedIn
                  </p>
                )}
              </div>
            </ModalBody>
            <ModalFooter className='flex gap-3'>
              <ModalClose className='w-full'>
                <Button className='w-full'>Cancel</Button>
              </ModalClose>

              <Button
                className='w-full'
                colorScheme='primary'
                onClick={handleSubmit}
                loadingText='Creating contact'
                isLoading={store.contacts.isLoading}
                rightSpinner={
                  <Spinner
                    size='sm'
                    label='loading'
                    className='text-primary-500 fill-primary-200'
                  />
                }
              >
                Enrich contact
              </Button>
            </ModalFooter>
          </ModalFeaturedContent>
        </ModalPortal>
      </Modal>
    );
  },
);
