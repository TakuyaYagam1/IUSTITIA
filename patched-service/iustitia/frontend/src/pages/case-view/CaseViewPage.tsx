import { ROUTES } from '@/app/router/routes';
import { CaseView } from '@entities/case';
import { useAuthStore } from '@entities/user';
import { useCase, useCaseComplaints, useCaseOpinion } from '@features/case-view';
import { Button } from '@shared/ui';
import { ChevronLeft, FileText, Gavel } from 'lucide-react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import styles from './CaseViewPage.module.css';

export const CaseViewPage = (): JSX.Element => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const session = useAuthStore((s) => s.session);
  const role = session?.role;
  const userId = session?.userId;
  const canSeeComplaints = role === 'judge' || role === 'prosecutor';
  const canSeeOpinion = role === 'judge' || role === 'prosecutor';

  const caseQuery = useCase(id);
  const caseData = caseQuery.data;

  const complaintsQuery = useCaseComplaints(canSeeComplaints ? id : undefined);
  const opinionEnabled =
    canSeeOpinion && (caseData?.status === 'hearing' || caseData?.status === 'closed');
  const opinionQuery = useCaseOpinion(opinionEnabled ? id : undefined, caseData?.status);
  const opinion = opinionQuery.data ?? null;

  const actions =
    caseData && id ? (
      <div className={styles.actions}>
        {role === 'prosecutor' &&
          caseData.status === 'assigned' &&
          caseData.assigned_prosecutor_id === userId && (
            <Button variant="primary" onClick={() => navigate(ROUTES.caseOpinion(id))}>
              <FileText size={14} /> Подать заключение
            </Button>
          )}
        {role === 'judge' && caseData.status === 'hearing' && (
          <Button variant="primary" onClick={() => navigate(ROUTES.caseVerdict(id))}>
            <Gavel size={14} /> Вынести приговор
          </Button>
        )}
      </div>
    ) : null;

  return (
    <div className={styles.wrap}>
      <Link to={ROUTES.cases} className={styles.back}>
        <ChevronLeft size={14} /> К картотеке
      </Link>

      {caseQuery.isLoading && <div className={styles.message}>Получение дела…</div>}
      {caseQuery.isError && (
        <div className={`${styles.message} ${styles.error}`}>
          Ошибка: {caseQuery.error instanceof Error ? caseQuery.error.message : 'неизвестно'}
        </div>
      )}
      {caseData && (
        <CaseView
          caseItem={caseData}
          {...(canSeeComplaints && complaintsQuery.data
            ? { complaints: complaintsQuery.data }
            : {})}
          opinion={opinion}
          actions={actions}
        />
      )}
    </div>
  );
};
