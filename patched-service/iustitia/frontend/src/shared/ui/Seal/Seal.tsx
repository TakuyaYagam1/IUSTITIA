import type { HTMLAttributes, ReactNode } from 'react';
import styles from './Seal.module.css';

type SealVariant = 'topsecret' | 'classified' | 'approved' | 'rejected';

interface SealProps extends HTMLAttributes<HTMLSpanElement> {
  variant?: SealVariant;
  large?: boolean;
  children?: ReactNode;
}

const labels: Record<SealVariant, string> = {
  topsecret: 'TOP SECRET',
  classified: 'СЕКРЕТНО',
  approved: 'УТВЕРЖДЕНО',
  rejected: 'ОТКЛОНЕНО',
};

export const Seal = ({
  variant = 'topsecret',
  large = false,
  className,
  children,
  ...rest
}: SealProps) => {
  const classes = [styles.seal, styles[variant], large && styles.large, className]
    .filter(Boolean)
    .join(' ');
  return (
    <span className={classes} {...rest}>
      {children ?? labels[variant]}
    </span>
  );
};
