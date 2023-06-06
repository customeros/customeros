import { useContext } from 'react';
import { SelectContext } from './context';

export const useSelect = () => useContext(SelectContext);
