import { ROUTES } from '@/app/router/routes';
import { ComplaintForm } from '@features/complaint-submit';
import { Paper, useToast } from '@shared/ui';
import { useNavigate } from 'react-router-dom';
import styles from './ComplaintPage.module.css';

export const ComplaintPage = (): JSX.Element => {
  const navigate = useNavigate();
  const toast = useToast();

  return (
    <div className={styles.wrap}>
      <div className={styles.intro}>
        <h1 className={styles.title}>Подать донос</h1>
        <span className={styles.subtitle}>Гражданский долг · Купол №7</span>
      </div>

      <div className={styles.warning}>
        Заведомо ложный донос карается изоляцией. Все обращения фиксируются и передаются
        регистратору для канцелярской сверки.
      </div>

      <Paper variant="paper">
        <ComplaintForm
          onSuccess={() => {
            toast.push({
              kind: 'success',
              title: 'Принято',
              message: 'Донос зарегистрирован. Ожидайте повестки.',
            });
            navigate(ROUTES.cases, { replace: true });
          }}
        />
      </Paper>
    </div>
  );
};
