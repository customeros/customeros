import { useRef } from 'react';
import { useParams } from 'react-router-dom';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Input } from '@ui/form/Input';
import { Select } from '@ui/form/Select';
import { UrlInput } from '@ui/form/UrlInput';
import { Spinner } from '@ui/feedback/Spinner';
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
  OrganizationStage,
  OrganizationRelationship,
} from '@graphql/types';

import { Tags } from './components/tags';
import { SocialIconInput } from '../../shared';
import { OwnerInput } from './components/owner';
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

  const handleSocialChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const id = (e.target as HTMLInputElement).id;
    const value = e.target.value;

    if (organization) {
      const idx = organization?.value.socialMedia.findIndex((s) => s.id === id);

      if (idx < 0) return;

      organization.value.socialMedia[idx].url = value;
    }
  };

  const handleSocialBlur = (
    e: React.ChangeEvent<HTMLInputElement>,
    newInputRef: React.RefObject<HTMLInputElement>,
  ) => {
    const id = (e.target as HTMLInputElement).id;

    const idx = organization?.value.socialMedia.findIndex((s) => s.id === id);

    if (organization.value.socialMedia[idx].url === '') {
      organization.value.socialMedia.splice(idx, 1);
      newInputRef.current?.focus();
    }

    organization.commit();
  };

  const handleSocialKeyDown = (
    e: React.KeyboardEvent<HTMLInputElement>,
    newInputRef: React.RefObject<HTMLInputElement>,
  ) => {
    const id = (e.target as HTMLInputElement).id;

    const idx = organization.value.socialMedia.findIndex((s) => s.id === id);
    const social = organization.value.socialMedia[idx];

    if (!social) return;

    if (social.url === '') {
      organization.value.socialMedia.splice(idx, 1);
      newInputRef.current?.focus();
    }

    organization.commit();
  };

  const handleCreateSocial = (value: string) => {
    organization.value.socialMedia.push({
      id: crypto.randomUUID(),
      url: value,
    } as Social);
    organization.commit();
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

  const enrichedOrg = organization?.value.enrichDetails;
  const enrichingStatus =
    !enrichedOrg?.enrichedAt &&
    enrichedOrg?.requestedAt &&
    !enrichedOrg?.failedAt;

  return (
    <div className='flex pt-[6px] px-6 w-full h-full overflow-y-auto flex-1 bg-gray-25 rounded-2xl'>
      <div className='flex h-full flex-col  overflow-visible w-full'>
        {enrichingStatus && (
          <div className='flex items-center justify-start gap-2 border-[1px] text-sm border-grayModern-100 bg-grayModern-50 rounded-[4px] py-1 px-2 '>
            <Spinner
              label='enriching org'
              className='text-grayModern-300 fill-grayModern-500 size-4'
            />
            <span className='font-medium'>
              We're enriching this organizations's details...
            </span>
          </div>
        )}

        <div className='flex items-center justify-between'>
          <Input
            size='xs'
            name='name'
            ref={nameRef}
            autoComplete='off'
            variant='unstyled'
            dataTest='org-about-name'
            placeholder='Company name'
            disabled={orgNameReadOnly}
            onFocus={(e) => e.target.select()}
            value={organization?.value.name || ''}
            onBlur={() => {
              organization.commit();
            }}
            onChange={(e) => {
              organization.value.name = e.target.value;
            }}
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
          dataTest='org-about-www'
          value={organization?.value?.website || ''}
          onBlur={() => {
            organization.commit();
          }}
          onChange={(e) => {
            organization.value.website = e.target.value;
          }}
        />
        <Textarea
          size='md'
          spellCheck={false}
          className='mb-6 mt-2'
          name='valueProposition'
          data-test='org-about-description'
          placeholder={placeholders.valueProposition}
          value={organization?.value?.valueProposition || ''}
          onBlur={() => {
            organization.commit();
          }}
          onChange={(e) => {
            organization.value.valueProposition = e.target.value;
          }}
        />
        <Tags
          dataTest='org-about-tags'
          className='min-h-10 py-2'
          inputPlaceholder='Search...'
          onCreate={handleCreateOption}
          placeholder='Organization tags'
          leftAccessory={<Tag01 className='mr-3 text-gray-500' />}
          options={store.tags.toArray().map((t) => ({
            value: t.id,
            label: t.value?.name,
          }))}
          value={
            organization.value.tags?.map((t) => ({
              value: t.id,
              label: t.name,
            })) ?? []
          }
          onChange={(selection) => {
            const tags = selection
              .map((o) => store.tags.getById(o.value)?.value)
              .filter(Boolean) as TagType[];

            organization.value.tags = tags;
            organization.commit();
          }}
        />

        {showParentRelationshipSelector && (
          <ParentOrgInput id={id} isReadOnly={parentRelationshipReadOnly} />
        )}
        <div className='flex items-center justify-center w-full'>
          <div className='flex-2' data-test='org-about-relationship'>
            <Menu>
              <MenuButton
                data-test='org-about-relationship'
                className='min-h-[40px] outline-none focus:outline-none'
              >
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
                      disabled={
                        (selectedRelationshipOption?.label === 'Customer' ||
                          selectedRelationshipOption?.label === 'Not a fit') &&
                        option.label === 'Prospect'
                      }
                      onClick={() => {
                        organization.value.relationship = option.value;
                        organization.value.stage = match(option.value)
                          .with(
                            OrganizationRelationship.Prospect,
                            () => OrganizationStage.Lead,
                          )
                          .with(
                            OrganizationRelationship.Customer,
                            () => OrganizationStage.InitialValue,
                          )
                          .with(
                            OrganizationRelationship.NotAFit,
                            () => OrganizationStage.Unqualified,
                          )
                          .with(
                            OrganizationRelationship.FormerCustomer,
                            () => OrganizationStage.Target,
                          )
                          .otherwise(() => undefined);

                        organization.commit();
                      }}
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
            <div className='flex-1' data-test='org-about-stage'>
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
                        organization.value.stage = option.value;
                        organization.commit();
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
            dataTest='org-about-industry'
            leftElement={<Building07 className='text-gray-500 mr-3' />}
            onChange={(option) => {
              organization.value.industry = option?.value;
              organization.commit();
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
            dataTest='org-about-business-type'
            leftElement={<Briefcase02 className='text-gray-500 mr-3' />}
            onChange={(option) => {
              organization.value.market = option?.value;
              organization.commit();
            }}
            value={businessTypeOptions.map((option) =>
              option.value === organization?.value.market ? option : null,
            )}
          />

          <div className='flex items-center justify-center w-full'>
            <div className='flex-1'>
              <Select
                isClearable
                name='lastFundingRound'
                placeholder='Last funding round'
                options={lastFundingRoundOptions}
                dataTest='org-about-last-funding-round'
                leftElement={
                  <HorizontalBarChart03 className='text-gray-500 mr-3' />
                }
                onChange={(option) => {
                  organization.value.lastFundingRound = option?.value;
                  organization.commit();
                }}
                value={lastFundingRoundOptions.map((option) =>
                  option.value === organization?.value.lastFundingRound
                    ? option
                    : null,
                )}
              />
            </div>
          </div>

          <Select
            isClearable
            name='employees'
            options={employeesOptions}
            placeholder='Number of employees'
            dataTest='org-about-number-of-employees'
            leftElement={<Users03 className='text-gray-500 mr-3' />}
            value={employeesOptions.map((option) =>
              option.value === organization?.value.employees ? option : null,
            )}
            onChange={(option) => {
              organization.value.employees = option.value;
              organization.commit();
            }}
          />

          <OwnerInput
            id={id}
            dataTest='org-about-org-owner'
            owner={organization?.value.owner}
          />
          <SocialIconInput
            name='socials'
            placeholder='Social link'
            onBlur={handleSocialBlur}
            onChange={handleSocialChange}
            onCreate={handleCreateSocial}
            onKeyDown={handleSocialKeyDown}
            dataTest='org-about-social-link'
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
