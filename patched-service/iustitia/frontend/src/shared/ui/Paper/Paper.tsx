import type { HTMLAttributes, ReactNode } from 'react';
import styles from './Paper.module.css';

type PaperVariant = 'paper' | 'tight' | 'panel';

interface PaperProps extends HTMLAttributes<HTMLDivElement> {
  variant?: PaperVariant;
  children: ReactNode;
}

export const Paper = ({ variant = 'paper', className, children, ...rest }: PaperProps) => {
  const base = variant === 'panel' ? styles.panel : styles.paper;
  const modifier = variant === 'tight' ? styles.tight : '';
  const classes = [base, modifier, className].filter(Boolean).join(' ');
  return (
    <div className={classes} {...rest}>
      {children}
    </div>
  );
};
