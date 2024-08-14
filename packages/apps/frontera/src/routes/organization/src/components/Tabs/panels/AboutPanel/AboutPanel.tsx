import { useRef } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Input } from '@ui/form/Input';
import { Select } from '@ui/form/Select';
import { UrlInput } from '@ui/form/UrlInput';
import { Tag01 } from '@ui/media/icons/Tag01.tsx';
import { Users03 } from '@ui/media/icons/Users03';
import { Share07 } from '@ui/media/icons/Share07';
import { useStore } from '@shared/hooks/useStore';
import { Seeding } from '@ui/media/icons/Seeding';
import { Target05 } from '@ui/media/icons/Target05';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Textarea } from '@ui/form/Textarea/Textarea';
import { Building07 } from '@ui/media/icons/Building07';
import { Tag, TagLabel } from '@ui/presentation/Tag/Tag';
import { BrokenHeart } from '@ui/media/icons/BrokenHeart';
import { Briefcase02 } from '@ui/media/icons/Briefcase02';
import { ActivityHeart } from '@ui/media/icons/ActivityHeart';
import { MessageXCircle } from '@ui/media/icons/MessageXCircle';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { HorizontalBarChart03 } from '@ui/media/icons/HorizontalBarChart03';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import {
  Social,
  DataSource,
  Tag as TagType,
  OrganizationRelationship,
} from '@graphql/types';

import { OwnerInput } from './components/owner';
import { Tags, SocialIconInput } from '../../shared';
import { Branches, ParentOrgInput } from './components/branches';
import {
  stageOptions,
  industryOptions,
  getStageOptions,
  employeesOptions,
  businessTypeOptions,
  relationshipOptions,
  lastFundingRoundOptions,
} from './util';

const placeholders = {
  valueProposition: `What is this organization about - what they do, who they serve, and what makes them unique?`,
};

const iconMap = {
  Customer: <ActivityHeart className='text-gray-500' />,
  Prospect: <Seeding className='text-gray-500' />,
  'Not a fit': <MessageXCircle className='text-gray-500' />,
  'Former Customer': <BrokenHeart className='text-gray-500' />,
};

export const AboutPanel = observer(() => {
  const store = useStore();
  const id = useParams()?.id as string;
  const [_, copyToClipboard] = useCopyToClipboard();
  const nameRef = useRef<HTMLInputElement | null>(null);

  const showParentRelationshipSelector = useFeatureIsOn(
    'show-parent-relationship-selector',
  );
  const parentRelationshipReadOnly = useFeatureIsOn(
    'parent-relationship-selector-read-only',
  );
  const orgNameReadOnly = useFeatureIsOn('org-name-readonly');

  const organization = store.organizations.value.get(id);

  if (!organization) return null;

  const selectedRelationshipOption = relationshipOptions.find(
    (option) => option.value === organization?.value.relationship,
  );

  const selectedStageOption = stageOptions.find(
    (option) => option.value === organization?.value.stage,
  );

  const applicableStageOptions = getStageOptions(
    organization?.value?.relationship,
  );

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    const { name, value } = e.target;

    organization?.update((org) => {
      // @ts-expect-error fixme
      org[name] = value;

      return org;
    });
  };

  const menuHandleChange = (name: string, value: string) => {
    organization?.update((org) => {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (org as any)[name] = value;

      return org;
    });
  };

  const handleSocialChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const id = (e.target as HTMLInputElement).id;
    const value = e.target.value;

    if (organization) {
      organization.update((org) => {
        const idx = organization?.value.socialMedia.findIndex(
          (s) => s.id === id,
        );

        if (idx !== -1) {
          org.socialMedia[idx].url = value;
        }

        return org;
      });
    }
  };

  const handleSocialBlur = (
    e: React.ChangeEvent<HTMLInputElement>,
    newInputRef: React.RefObject<HTMLInputElement>,
  ) => {
    const id = (e.target as HTMLInputElement).id;

    organization?.update((org) => {
      const idx = organization?.value.socialMedia.findIndex((s) => s.id === id);

      if (org.socialMedia[idx].url === '') {
        org.socialMedia.splice(idx, 1);
        newInputRef.current?.focus();
      }

      return org;
    });
  };

  const handleSocialKeyDown = (
    e: React.KeyboardEvent<HTMLInputElement>,
    newInputRef: React.RefObject<HTMLInputElement>,
  ) => {
    const id = (e.target as HTMLInputElement).id;

    organization?.update((org) => {
      const idx = org.socialMedia.findIndex((s) => s.id === id);
      const social = org.socialMedia[idx];

      if (!social) return org;

      if (social.url === '') {
        org.socialMedia.splice(idx, 1);
        newInputRef.current?.focus();
      }

      return org;
    });
  };

  const handleCreateSocial = (value: string) => {
    organization?.update((org) => {
      org.socialMedia.push({
        id: crypto.randomUUID(),
        url: value,
      } as Social);

      return org;
    });
  };

  const handleCreateOption = (value: string) => {
    store.tags?.create({ name: value });

    organization?.update((org) => {
      org.tags = [
        ...(org.tags || []),
        {
          id: value,
          name: value,
          metadata: {
            id: value,
            source: DataSource.Openline,
            sourceOfTruth: DataSource.Openline,
            appSource: 'organization',
            created: new Date().toISOString(),
            lastUpdated: new Date().toISOString(),
          },
          appSource: 'organization',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          source: DataSource.Openline,
        },
      ];

      return org;
    });
  };
  const filteredTags = organization?.value.tags
    ?.filter((e) => e.id)
    ?.map((tag) => ({
      label: tag.name,
      value: tag.id,
    }));

  return (
    <div className=' flex pt-[6px] px-6 w-full h-full overflow-y-auto flex-1 bg-gray-25 rounded-2xl'>
      <div className='flex h-full flex-col  overflow-visible w-full'>
        <div className='flex items-center justify-between'>
          <Input
            size='xs'
            name='name'
            ref={nameRef}
            autoComplete='off'
            variant='unstyled'
            onChange={handleChange}
            placeholder='Company name'
            disabled={orgNameReadOnly}
            onFocus={(e) => e.target.select()}
            value={organization?.value.name || ''}
            className='font-semibold text-[16px] mt-0.5 border-none overflow-hidden overflow-ellipsis'
          />
          {organization?.value.referenceId && (
            <div className='h-full ml-4'>
              <Tooltip asChild={false} label={'Copy ID'}>
                <Tag
                  variant='outline'
                  colorScheme='gray'
                  className='rounded-full cursor-pointer'
                  onClick={() => {
                    copyToClipboard(
                      organization?.value.referenceId ?? '',
                      'Reference ID copied ',
                    );
                  }}
                >
                  <TagLabel>{organization?.value.referenceId}</TagLabel>
                </Tag>
              </Tooltip>
            </div>
          )}
        </div>
        <UrlInput
          name='website'
          autoComplete='off'
          placeholder='www.'
          onChange={handleChange}
          value={organization?.value?.website || ''}
        />
        <Textarea
          size='md'
          spellCheck={false}
          className='mb-6 mt-2'
          name='valueProposition'
          onChange={handleChange}
          placeholder={placeholders.valueProposition}
          value={organization?.value?.valueProposition || ''}
        />
        <Tags
          value={filteredTags || []}
          placeholder='Organization tags'
          onCreateOption={handleCreateOption}
          icon={
            <Tag01 className='text-gray-500 min-w-[18px] min-h-4 mr-[10px] mt-[6px]' />
          }
          onChange={(e) => {
            organization?.update((org) => {
              org.tags =
                (e
                  .map((tag) => store.tags?.value.get(tag.value)?.value)
                  .filter(Boolean) as Array<TagType>) ?? [];

              return org;
            });
          }}
        />
        {showParentRelationshipSelector && (
          <ParentOrgInput id={id} isReadOnly={parentRelationshipReadOnly} />
        )}
        <div className='flex items-center justify-center w-full'>
          <div className='flex-2'>
            <Menu>
              <MenuButton className='min-h-[40px] outline-none focus:outline-none'>
                {
                  iconMap[
                    selectedRelationshipOption?.label as keyof typeof iconMap
                  ]
                }{' '}
                <span className='ml-2'>
                  {selectedRelationshipOption?.label}
                </span>
              </MenuButton>
              <MenuList side='bottom' align='start' className='min-w-[280px]'>
                {relationshipOptions
                  .filter(
                    (option) =>
                      !(
                        selectedRelationshipOption?.label === 'Customer' &&
                        option.label === 'Prospect'
                      ) &&
                      !(
                        selectedRelationshipOption?.label === 'Not a fit' &&
                        option.label === 'Prospect'
                      ),
                  )
                  .map((option) => (
                    <MenuItem
                      key={option.value}
                      onClick={() => {
                        menuHandleChange('relationship', option.value);
                      }}
                      disabled={
                        (selectedRelationshipOption?.label === 'Customer' ||
                          selectedRelationshipOption?.label === 'Not a fit') &&
                        option.label === 'Prospect'
                      }
                    >
                      {iconMap[option.label as keyof typeof iconMap]}
                      {option.label}
                    </MenuItem>
                  ))}
              </MenuList>
            </Menu>
          </div>

          {organization?.value?.relationship ===
            OrganizationRelationship.Prospect && (
            <div className='flex-1'>
              <Menu>
                <MenuButton className='min-h-[40px] outline-none focus:outline-none'>
                  <Target05 className='text-gray-500 mb-0.5' />
                  <span className='ml-2'>
                    {selectedStageOption?.label || 'Stage'}
                  </span>
                </MenuButton>
                <MenuList
                  side='bottom'
                  align='center'
                  className='min-w-[280px]'
                >
                  {applicableStageOptions.map((option) => (
                    <MenuItem
                      key={option.value}
                      onClick={() => {
                        menuHandleChange('stage', option.value);
                      }}
                    >
                      {iconMap[option.label as keyof typeof iconMap]}
                      {option.label}
                    </MenuItem>
                  ))}
                </MenuList>
              </Menu>
            </div>
          )}
        </div>
        <div className='flex flex-col w-full flex-1 items-start justify-start gap-0'>
          <Select
            isClearable
            name='industry'
            placeholder='Industry'
            options={industryOptions}
            leftElement={<Building07 className='text-gray-500 mr-3' />}
            onChange={(value) => {
              organization?.update((org) => {
                org.industry = value.value;

                return org;
              });
            }}
            value={
              industryOptions
                ? industryOptions.map((option) =>
                    option.options.find(
                      (v) => v.value === organization?.value?.industry,
                    ),
                  )
                : null
            }
          />

          <Select
            isClearable
            name='businessType'
            placeholder='Business Type'
            options={businessTypeOptions}
            leftElement={<Briefcase02 className='text-gray-500 mr-3' />}
            value={businessTypeOptions.map((option) =>
              option.value === organization?.value.market ? option : null,
            )}
            onChange={(value) => {
              organization?.update((org) => {
                if (value === null) org.market = null;
                else org.market = value.value;

                return org;
              });
            }}
          />

          <div className='flex items-center justify-center w-full'>
            <div className='flex-1'>
              <Select
                isClearable
                name='lastFundingRound'
                placeholder='Last funding round'
                options={lastFundingRoundOptions}
                leftElement={
                  <HorizontalBarChart03 className='text-gray-500 mr-3' />
                }
                value={lastFundingRoundOptions.map((option) =>
                  option.value === organization?.value.lastFundingRound
                    ? option
                    : null,
                )}
                onChange={(value) => {
                  organization?.update((org) => {
                    if (value === null) org.lastFundingRound = null;
                    else org.lastFundingRound = value.value;

                    return org;
                  });
                }}
              />
            </div>
          </div>

          <Select
            isClearable
            name='employees'
            options={employeesOptions}
            placeholder='Number of employees'
            leftElement={<Users03 className='text-gray-500 mr-3' />}
            value={employeesOptions.map((option) =>
              option.value === organization?.value.employees ? option : null,
            )}
            onChange={(value) => {
              organization?.update((org) => {
                org.employees = value.value;

                return org;
              });
            }}
          />

          <OwnerInput id={id} owner={organization?.value.owner} />
          <SocialIconInput
            name='socials'
            placeholder='Social link'
            onBlur={handleSocialBlur}
            onChange={handleSocialChange}
            onCreate={handleCreateSocial}
            onKeyDown={handleSocialKeyDown}
            leftElement={<Share07 className='text-gray-500' />}
            value={organization?.value.socialMedia.map((s) => ({
              value: s.id,
              label: s.url,
            }))}
          />

          {showParentRelationshipSelector &&
            organization?.subsidiaries?.length > 0 && (
              <Branches id={id} isReadOnly={parentRelationshipReadOnly} />
            )}
        </div>
        {organization?.value.customerOsId && (
          <Tooltip label='Copy ID'>
            <span
              className='py-3 w-fit text-gray-400 cursor-pointer'
              onClick={() =>
                copyToClipboard(
                  organization?.value.customerOsId ?? '',
                  'CustomerOS ID copied',
                )
              }
            >
              CustomerOS ID: {organization?.value.customerOsId}
            </span>
          </Tooltip>
        )}
      </div>
    </div>
  );
});
