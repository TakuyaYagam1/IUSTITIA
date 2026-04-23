import { ROUTES } from '@/app/router/routes';
import { CaseCreateForm } from '@features/case-create';
import { Paper, useToast } from '@shared/ui';
import { useNavigate } from 'react-router-dom';
import styles from './CaseNewPage.module.css';

export const CaseNewPage = (): JSX.Element => {
  const navigate = useNavigate();
  const toast = useToast();

  return (
    <div className={styles.wrap}>
      <div className={styles.intro}>
        <h1 className={styles.title}>Новое дело</h1>
        <span className={styles.subtitle}>Гражданский долг · Купол №7</span>
      </div>

      <div className={styles.warning}>
        Заведомо ложное заявление карается изоляцией. Все обращения фиксируются и передаются
        регистратору для канцелярской сверки.
      </div>

      <Paper variant="paper">
        <CaseCreateForm
          onSuccess={(c) => {
            toast.push({
              kind: 'success',
              title: 'Принято',
              message: `Дело №${c.seq_num}/Δ зарегистрировано.`,
            });
            navigate(ROUTES.caseView(c.id), { replace: true });
          }}
        />
      </Paper>
    </div>
  );
};
