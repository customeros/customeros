'use client';

import {Flex} from "@ui/layout/Flex";

export const SettingsMainSection = ({ children }: { children?: React.ReactNode }) => {
  return (
      <>
        <Flex>
          {children}
        </Flex>
      </>
  );
};
