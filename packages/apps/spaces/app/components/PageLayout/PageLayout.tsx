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
      h='100vh'
      backgroundColor='grey.100'
      templateAreas={`"sidebar content"`}
      templateColumns='80px 1fr'
      transition='all ease 0.25s'
    >
      <Sidebar isOwner={isOwner} />
      <GridItem
        p='4'
        h='100%'
        area='content'
        overflowX='hidden'
        overflowY='auto'
        bg='gray.50'
      >
        {children}
      </GridItem>
    </Grid>
  );
};
