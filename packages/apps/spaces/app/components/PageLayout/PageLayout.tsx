'use client';
import { Sidebar } from '../Sidebar/Sidebar';

import { Grid, GridItem } from '@ui/layout/Grid';

interface PageLayoutProps {
  isOwner: boolean;
}

export const PageLayout = ({
  isOwner,
  children,
}: React.PropsWithChildren<PageLayoutProps>) => {
  return (
    <Grid
      gap='6'
      h='100vh'
      backgroundColor='grey.100'
      templateAreas={`"sidebar content"`}
      templateColumns='80px 1fr'
      transition='all ease 0.25s'
    >
      <Sidebar isOwner={isOwner} />
      <GridItem
        p='1.2rem'
        h='100%'
        area='content'
        overflowX='hidden'
        overflowY='auto'
      >
        {children}
      </GridItem>
    </Grid>
  );
};
