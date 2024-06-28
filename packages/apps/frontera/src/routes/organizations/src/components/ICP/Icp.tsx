import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { Cake } from '@ui/media/icons/Cake';
import { Play } from '@ui/media/icons/Play';
import { Key01 } from '@ui/media/icons/Key01';
import { Tag01 } from '@ui/media/icons/Tag01';
import { Button } from '@ui/form/Button/Button';
import { Star06 } from '@ui/media/icons/Star06';
import { Globe03 } from '@ui/media/icons/Globe03';
import { Users03 } from '@ui/media/icons/Users03';
import { Linkedin } from '@ui/media/icons/Linkedin';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { Building05 } from '@ui/media/icons/Building05';
import { Select, getContainerClassNames } from '@ui/form/Select';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

export const Icp = observer(() => {
  const options = ['between', 'less than', 'more than'];
  const [employeesFilter, setEmployeesFilter] = useState(options[1]);
  const [followersFilter, setFollowersFilter] = useState(options[1]);
  const [organizationFilter, setOrganizationFilter] = useState(options[1]);

  const handleEmployeesFilter = () => {
    const currentIndex = options.indexOf(employeesFilter);
    const nextIndex = (currentIndex + 1) % options.length;
    setEmployeesFilter(options[nextIndex]);
  };

  const handleTagsFilter = () => {
    const currentIndex = options.indexOf(followersFilter);
    const nextIndex = (currentIndex + 1) % options.length;
    setFollowersFilter(options[nextIndex]);
  };

  const handleOrganizationFilter = () => {
    const currentIndex = options.indexOf(organizationFilter);
    const nextIndex = (currentIndex + 1) % options.length;
    setOrganizationFilter(options[nextIndex]);
  };

  return (
    <>
      <div className='flex items-center justify-between'>
        <p className='font-semibold'>Auto-qualify leads</p>
        <Button size='xxs' leftIcon={<Play />}>
          Start flow
        </Button>
      </div>
      <p className='mt-1'>
        Create your <span className='font-medium'>Ideal Company Profile </span>{' '}
        and automatically qualify
        <span className='font-medium'> Leads </span>
        into <span className='font-medium'>Targets</span>
      </p>
      <p className='font-medium leading-5 text-gray-500 mt-4 mb-2'>WHEN</p>

      <div className='flex items-center w-full'>
        <div className='flex items-center flex-1'>
          <Building05 className='mr-2 text-gray-500' />
          <p className='font-medium'>
            Industry <span className='font-normal'> is any of </span>
          </p>
        </div>
        <div className='flex-1'>
          <Select
            isMulti
            options={[
              { label: 'Tech', value: 'tech' },
              { label: 'Health', value: 'health' },
            ]}
            classNames={{
              container: () =>
                getContainerClassNames(undefined, 'unstyled', {}),
            }}
            placeholder='Industries'
          />
        </div>
      </div>

      <div className='flex items-center w-full'>
        <div className='flex-1 items-center flex'>
          <Users03 className='mr-2 text-gray-500' />
          <p className='font-medium'>
            Employees <span className='font-normal'>are </span>
            <span
              className='cursor-pointer underline font-normal text-gray-500 hover:text-gray-700'
              onClick={handleEmployeesFilter}
            >
              {employeesFilter}
            </span>
          </p>
        </div>
        <div className='flex-1 flex items-center'>
          <Input
            variant='unstyled'
            type='number'
            placeholder={
              employeesFilter === 'between' ? 'Min' : 'Number of employees'
            }
            style={{
              width: employeesFilter !== 'between' ? '100%' : '50px',
            }}
          />
          <span
            className='mr-[30px] '
            style={{
              display: employeesFilter === 'between' ? 'block' : 'none',
            }}
          >
            -{' '}
          </span>
          <Input
            style={{
              display: employeesFilter === 'between' ? 'block' : 'none',
            }}
            variant='unstyled'
            type='number'
            placeholder='Max'
            className='w-[50px]'
          />
        </div>
      </div>

      <div className='flex items-center w-full '>
        <div className='flex flex-1 items-center'>
          <Globe03 className='mr-2 text-gray-500' />
          <p className='font-medium'>
            Headquarters <span className='font-normal'>is any of</span>{' '}
          </p>
        </div>
        <div className='flex-1'>
          <Select isMulti options={[]} placeholder='Headquarter countries' />
        </div>
      </div>

      <div className='flex items-center w-full '>
        <div className='flex flex-1 items-center'>
          <Tag01 className='mr-2 text-gray-500' />
          <p className='font-medium'>
            Tag <span className='font-normal'>is any of</span>{' '}
          </p>
        </div>
        <div className='flex-1'>
          <Select isMulti options={[]} placeholder='Organization tags' />
        </div>
      </div>

      <div className='flex items-center w-full '>
        <div className='flex flex-1 items-center'>
          <Linkedin className='mr-2 text-gray-500 ' />
          <p className='font-medium'>
            Followers <span className='font-normal'>is </span>
            <span
              className='cursor-pointer underline font-normal text-gray-500 hover:text-gray-700'
              onClick={handleTagsFilter}
            >
              {followersFilter}
            </span>
          </p>
        </div>
        <div className='flex-1 flex items-center'>
          <Input
            variant='unstyled'
            type='number'
            placeholder={
              followersFilter === 'between' ? 'Min' : 'Number of followers'
            }
            style={{
              width: followersFilter !== 'between' ? '100%' : '50px',
            }}
          />

          <span
            className='mr-[30px] '
            style={{
              display: followersFilter === 'between' ? 'block' : 'none',
            }}
          >
            -{' '}
          </span>
          <Input
            style={{
              display: followersFilter === 'between' ? 'block' : 'none',
            }}
            variant='unstyled'
            type='number'
            placeholder='Max'
            className='w-[50px]'
          />
        </div>
      </div>
      <div className='flex items-center w-full '>
        <div className='flex flex-1 items-center'>
          <Cake className='mr-2 text-gray-500' />
          <p className='font-medium'>
            Organization age <span className='font-normal'>is </span>
            <span
              className='cursor-pointer underline font-normal text-gray-500 hover:text-gray-700'
              onClick={handleOrganizationFilter}
            >
              {organizationFilter}
            </span>
          </p>
        </div>
        <div className='flex-1 flex items-center'>
          <Input
            variant='unstyled'
            type='number'
            placeholder={organizationFilter === 'between' ? 'Min' : 'Age'}
            className={cn(
              organizationFilter === 'between' ? 'w-[32px] ' : 'w-[32px]',
            )}
          />
          <span>/yrs</span>

          <span
            className='mx-4 '
            style={{
              display: organizationFilter === 'between' ? 'block' : 'none',
            }}
          >
            -{' '}
          </span>
          <Input
            style={{
              display: organizationFilter === 'between' ? 'block' : 'none',
            }}
            variant='unstyled'
            type='number'
            placeholder='Max'
            className={cn(
              organizationFilter === 'between' ? 'w-[32px] ' : 'w-[24px]',
            )}
          />
          {organizationFilter === 'between' && <span>/yrs</span>}
        </div>
      </div>

      <div className='flex items-center w-full mt-2 '>
        <div className='flex flex-1 items-center'>
          <Key01 className='mr-2 text-gray-500' />
          <p className='font-medium'>
            Ownership <span className='font-normal'>is </span>
          </p>
        </div>
        <div className='flex-1 flex items-center'>
          <Menu>
            <MenuButton>test</MenuButton>
            <MenuList>
              <MenuItem>Private</MenuItem>
              <MenuItem>Public</MenuItem>
            </MenuList>
          </Menu>
        </div>
      </div>

      <div className='mt-4 border rounded-md flex items-start gap-2 p-3 bg-grayModern-50'>
        <div className='flex flex-col w-fit'>
          <Star06 className='mt-1 text-grayModern-500' />
        </div>
        <div className='flex flex-col'>
          <p>This flow will qualify 34/89 Leads into Targets</p>
          <Checkbox>See filtered leads before starting the flow</Checkbox>
        </div>
      </div>
    </>
  );
});
