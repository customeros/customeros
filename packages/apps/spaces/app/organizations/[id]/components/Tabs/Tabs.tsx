'use client';

import type { TabProps } from '@ui/disclosure/Tabs';
import {
  TabList,
  TabPanels,
  Tab as ChakraTab,
  Tabs as ChakraTabs,
} from '@ui/disclosure/Tabs';

const Tab = (props: TabProps) => (
  <ChakraTab
    mr='1'
    w='90px'
    fontSize='14px'
    bg='gray.100'
    _selected={{ bg: 'white', fontWeight: 'bold' }}
    {...props}
  />
);

export const Tabs = ({ children }: { children?: React.ReactNode }) => {
  return (
    <ChakraTabs
      mt='-38px'
      zIndex='1'
      variant='enclosed'
      defaultIndex={4}
      h='full'
      display='flex'
      flexDir='column'
    >
      <TabList>
        <Tab>Up Next</Tab>
        <Tab>Account</Tab>
        <Tab>Success</Tab>
        <Tab>People</Tab>
        <Tab>About</Tab>
      </TabList>

      <TabPanels h='full' position='relative' overflowY='auto' flex='1'>
        {children}
      </TabPanels>
    </ChakraTabs>
  );
};
