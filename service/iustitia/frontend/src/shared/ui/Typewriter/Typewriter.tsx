import type { HTMLAttributes, ReactNode } from 'react';
import styles from './Typewriter.module.css';

interface TypewriterProps extends HTMLAttributes<HTMLSpanElement> {
  block?: boolean;
  muted?: boolean;
  children: ReactNode;
}

export const Typewriter = ({
  block = false,
  muted = false,
  className,
  children,
  ...rest
}: TypewriterProps) => {
  const classes = [styles.typewriter, block && styles.block, muted && styles.muted, className]
    .filter(Boolean)
    .join(' ');
  if (block) {
    return (
      <div className={classes} {...(rest as HTMLAttributes<HTMLDivElement>)}>
        {children}
      </div>
    );
  }
  return (
    <span className={classes} {...rest}>
      {children}
    </span>
  );
};
