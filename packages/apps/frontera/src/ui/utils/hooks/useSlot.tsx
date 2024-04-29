/* eslint-disable @typescript-eslint/no-explicit-any */
import { FC, Children, ReactNode, isValidElement } from 'react';

const getSlot = (children: ReactNode, component: FC<any>) => {
  let slot: ReactNode | null = null;

  Children.forEach(children, (child) => {
    if (!isValidElement(child) || child.type !== component) return null;
    slot = child;
  });

  return slot;
};

const omitSlots = (children: ReactNode, ...components: FC<any>[]) =>
  Children.toArray(children).filter(
    (child) =>
      !isValidElement(child) || !components.some((c) => c === child.type),
  );

export const useSlot = (children: ReactNode, component: FC<any>) => {
  return getSlot(children, component);
};

export const useSlots = (children: ReactNode, ...components: FC<any>[]) => {
  return components.map((component) => getSlot(children, component));
};

export const useOmitSlots = (children: ReactNode, ...components: FC<any>[]) => {
  return omitSlots(children, ...components);
};
