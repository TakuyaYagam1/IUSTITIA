import type { HTMLAttributes } from 'react';
import styles from './Nameplate.module.css';

interface NameplateProps extends HTMLAttributes<HTMLDivElement> {
  name: string;
  role?: string;
}

export const Nameplate = ({ name, role, className, ...rest }: NameplateProps) => {
  const classes = [styles.nameplate, className].filter(Boolean).join(' ');
  return (
    <div className={classes} {...rest}>
      <span className={styles.name}>{name}</span>
      {role && <span className={styles.role}>{role}</span>}
    </div>
  );
};
