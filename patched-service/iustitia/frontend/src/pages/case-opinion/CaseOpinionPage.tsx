import { ROUTES } from '@/app/router/routes';
import { useCase } from '@features/case-view';
import { OpinionForm } from '@features/opinion-file';
import { Paper, useToast } from '@shared/ui';
import { useNavigate, useParams } from 'react-router-dom';
import styles from './CaseOpinionPage.module.css';

export const CaseOpinionPage = (): JSX.Element => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const toast = useToast();
  const { data: caseData, isLoading, isError } = useCase(id);

  if (!id) {
    return <div className={styles.status}>Идентификатор дела не указан.</div>;
  }

  return (
    <div className={styles.wrap}>
      <div className={styles.intro}>
        <h1 className={styles.title}>Заключение прокурора</h1>
        <span className={styles.subtitle}>Предварительная квалификация · Трибунал</span>
      </div>

      {isLoading && <div className={styles.status}>Загрузка дела…</div>}
      {isError && <div className={styles.status}>Не удалось получить дело.</div>}

      {caseData && (
        <div className={styles.meta}>
          Дело №{caseData.seq_num}/Δ · {caseData.defendant} · {caseData.crime}
        </div>
      )}

      <Paper variant="paper">
        <OpinionForm
          caseId={id}
          onSuccess={() => {
            toast.push({
              kind: 'success',
              title: 'Заключение передано',
              message: 'Дело направлено судье.',
            });
            navigate(ROUTES.caseView(id), { replace: true });
          }}
        />
      </Paper>
    </div>
  );
};
